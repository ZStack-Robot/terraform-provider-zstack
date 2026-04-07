// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
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
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the price table.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the price table.",
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

	p := param.CreatePriceTableParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePriceTableParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Prices:      []param.CreatePriceTable_PriceParam{},
		},
	}

	priceTable, err := r.client.CreatePriceTable(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create price table",
			"Error "+err.Error(),
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

	queryParam := param.NewQueryParam()
	priceTables, err := r.client.QueryPriceTable(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query price tables. It may have been deleted.: "+err.Error())
		state = priceTableModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, priceTable := range priceTables {
		if priceTable.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(priceTable.UUID)
			state.Name = types.StringValue(priceTable.Name)
			state.Description = stringValueOrNull(priceTable.Description)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Price table not found. It might have been deleted outside of Terraform.")
		state = priceTableModel{
			Uuid: types.StringValue(""),
		}
	}

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
			"Fail to update price table",
			"Error "+err.Error(),
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

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Price table UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeletePriceTable(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Fail to delete price table", "Error "+err.Error())
		return
	}
}

func (r *priceTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
