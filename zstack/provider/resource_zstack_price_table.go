// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &priceTableResource{}
	_ resource.ResourceWithConfigure   = &priceTableResource{}
	_ resource.ResourceWithImportState = &priceTableResource{}
)

type priceTableResource struct {
	client *client.ZSClient
}

type priceTableModel struct {
	Uuid        types.String           `tfsdk:"uuid"`
	Name        types.String           `tfsdk:"name"`
	Description types.String           `tfsdk:"description"`
	Prices      []priceTablePriceModel `tfsdk:"prices"`
}

type priceTablePriceModel struct {
	ResourceName types.String  `tfsdk:"resource_name"`
	ResourceUnit types.String  `tfsdk:"resource_unit"`
	TimeUnit     types.String  `tfsdk:"time_unit"`
	Price        types.Float64 `tfsdk:"price"`
	DateInLong   types.Int64   `tfsdk:"date_in_long"`
	SystemTags   types.List    `tfsdk:"system_tags"`
}

func PriceTableResource() resource.Resource {
	return &priceTableResource{}
}

func (r *priceTableResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)

		return
	}
	r.client = client
}

func (r *priceTableResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_price_table"
}

func (r *priceTableResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage price tables in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the price table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the price table.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the price table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"prices": schema.ListNestedAttribute{
				Required:    true,
				Description: "Price entries to create with the price table. ZStack does not return these entries from QueryPriceTable, so changes force replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_name": schema.StringAttribute{
							Required:    true,
							Description: "The billable resource name, for example cpu, memory, rootVolume, dataVolume, or publicIp.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"resource_unit": schema.StringAttribute{
							Required:    true,
							Description: "The billable resource unit.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"time_unit": schema.StringAttribute{
							Required:    true,
							Description: "The billing time unit.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"price": schema.Float64Attribute{
							Required:    true,
							Description: "The price for this resource and time unit.",
						},
						"date_in_long": schema.Int64Attribute{
							Optional:    true,
							Description: "Optional effective timestamp in milliseconds since the Unix epoch.",
						},
						"system_tags": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Optional system tags for this price entry.",
						},
					},
				},
			},
		},
	}
}

func (r *priceTableResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan priceTableModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	prices, priceDiags := buildPriceTablePriceParams(ctx, plan.Prices)
	response.Diagnostics.Append(priceDiags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreatePriceTableParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePriceTableParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Prices:      prices,
		},
	}

	priceTable, err := r.client.CreatePriceTable(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Price Table",
			"Could not create price table, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(priceTable.UUID)
	plan.Name = types.StringValue(priceTable.Name)
	plan.Description = stringValueOrNull(priceTable.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *priceTableResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state priceTableModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	priceTable, err := findResourceByQuery(r.client.QueryPriceTable, state.Uuid.ValueString())

	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Price Table",
			"Could not read Price Table, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(priceTable.UUID)
	state.Name = types.StringValue(priceTable.Name)
	state.Description = stringValueOrNull(priceTable.Description)
	// QueryPriceTable does not return price entries; keep the configured state.

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *priceTableResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan priceTableModel
	var state priceTableModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdatePriceTableParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdatePriceTableParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	priceTable, err := r.client.UpdatePriceTable(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Price Table",
			"Could not update price table UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(priceTable.UUID)
	plan.Name = types.StringValue(priceTable.Name)
	plan.Description = stringValueOrNull(priceTable.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *priceTableResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state priceTableModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePriceTable(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting Price Table", "Could not delete price table UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}
}

func buildPriceTablePriceParams(ctx context.Context, prices []priceTablePriceModel) ([]param.CreatePriceTable_PriceParam, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make([]param.CreatePriceTable_PriceParam, 0, len(prices))
	for _, price := range prices {
		var systemTags []string
		if !price.SystemTags.IsNull() && !price.SystemTags.IsUnknown() {
			diags.Append(price.SystemTags.ElementsAs(ctx, &systemTags, false)...)
			if diags.HasError() {
				return nil, diags
			}
		}

		entry := param.CreatePriceTable_PriceParam{
			ResourceName: stringPtr(price.ResourceName.ValueString()),
			ResourceUnit: stringPtr(price.ResourceUnit.ValueString()),
			TimeUnit:     stringPtr(price.TimeUnit.ValueString()),
			Price:        price.Price.ValueFloat64(),
			SystemTags:   systemTags,
		}
		if !price.DateInLong.IsNull() && !price.DateInLong.IsUnknown() {
			entry.DateInLong = price.DateInLong.ValueInt64Pointer()
		}
		result = append(result, entry)
	}
	return result, diags
}

func (r *priceTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
