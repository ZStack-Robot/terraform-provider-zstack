// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vrouterRouteTableResource{}
	_ resource.ResourceWithConfigure   = &vrouterRouteTableResource{}
	_ resource.ResourceWithImportState = &vrouterRouteTableResource{}
)

type vrouterRouteTableResource struct {
	client *client.ZSClient
}

type vrouterRouteTableResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func VRouterRouteTableResource() resource.Resource {
	return &vrouterRouteTableResource{}
}

func (r *vrouterRouteTableResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vrouterRouteTableResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vrouter_route_table"
}

func (r *vrouterRouteTableResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage VRouter route tables in ZStack. " +
			"A VRouter route table contains routing entries for virtual routers.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VRouter route table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VRouter route table.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VRouter route table.",
			},
		},
	}
}

func (r *vrouterRouteTableResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vrouterRouteTableResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVRouterRouteTableParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVRouterRouteTableParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	vrouterRouteTable, err := r.client.CreateVRouterRouteTable(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VRouter Route Table",
			"Could not create vrouter route table, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vrouterRouteTable.UUID)
	plan.Name = types.StringValue(vrouterRouteTable.Name)
	plan.Description = types.StringValue(vrouterRouteTable.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vrouterRouteTableResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vrouterRouteTableResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vrouterRouteTable, err := findResourceByQuery(r.client.QueryVRouterRouteTable, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading VRouter Route Table",
			"Could not read VRouter Route Table, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(vrouterRouteTable.UUID)
	state.Name = types.StringValue(vrouterRouteTable.Name)
	state.Description = types.StringValue(vrouterRouteTable.Description)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vrouterRouteTableResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan vrouterRouteTableResourceModel
	var state vrouterRouteTableResourceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateVRouterRouteTableParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVRouterRouteTableParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	vrouterRouteTable, err := r.client.UpdateVRouterRouteTable(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating VRouter Route Table",
			"Could not update vrouter route table, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vrouterRouteTable.UUID)
	plan.Name = types.StringValue(vrouterRouteTable.Name)
	plan.Description = types.StringValue(vrouterRouteTable.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vrouterRouteTableResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vrouterRouteTableResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteVRouterRouteTable(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting VRouter Route Table",
			"Could not delete vrouter route table, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *vrouterRouteTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
