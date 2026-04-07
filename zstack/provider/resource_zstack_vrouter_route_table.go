// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			"Fail to create VRouter route table",
			"Error "+err.Error(),
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

	queryParam := param.NewQueryParam()
	vrouterRouteTables, err := r.client.QueryVRouterRouteTable(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query VRouter route tables. It may have been deleted.: "+err.Error())
		state = vrouterRouteTableResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vrouterRouteTable := range vrouterRouteTables {
		if vrouterRouteTable.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(vrouterRouteTable.UUID)
			state.Name = types.StringValue(vrouterRouteTable.Name)
			state.Description = types.StringValue(vrouterRouteTable.Description)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "VRouter route table not found. It might have been deleted outside of Terraform.")
		state = vrouterRouteTableResourceModel{
			Uuid: types.StringValue(""),
		}
	}

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
			"Fail to update VRouter route table",
			"Error "+err.Error(),
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

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "VRouter route table UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVRouterRouteTable(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete VRouter route table", ""+err.Error())
		return
	}
}

func (r *vrouterRouteTableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
