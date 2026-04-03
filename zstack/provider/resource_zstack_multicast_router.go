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
	_ resource.Resource                = &multicastRouterResource{}
	_ resource.ResourceWithConfigure   = &multicastRouterResource{}
	_ resource.ResourceWithImportState = &multicastRouterResource{}
)

type multicastRouterResource struct {
	client *client.ZSClient
}

type multicastRouterModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	VpcRouterVmUuid types.String `tfsdk:"vpc_router_vm_uuid"`
	State           types.String `tfsdk:"state"`
}

func MulticastRouterResource() resource.Resource {
	return &multicastRouterResource{}
}

func (r *multicastRouterResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *multicastRouterResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_multicast_router"
}

func (r *multicastRouterResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a multicast router in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the multicast router.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the multicast router.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the multicast router.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_router_vm_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VPC router VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the multicast router.",
			},
		},
	}
}

func (r *multicastRouterResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan multicastRouterModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateMulticastRouterParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateMulticastRouterParamDetail{
			VpcRouterVmUuid: plan.VpcRouterVmUuid.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	multicastRouter, err := r.client.CreateMulticastRouter(p)
	if err != nil {
		response.Diagnostics.AddError("Fail to create multicast router", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(multicastRouter.UUID)
	plan.Name = stringValueOrNull(multicastRouter.Name)
	plan.Description = stringValueOrNull(multicastRouter.Description)
	plan.State = stringValueOrNull(multicastRouter.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *multicastRouterResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state multicastRouterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	multicastRouters, err := r.client.QueryMulticastRouter(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query multicast routers. It may have been deleted.: "+err.Error())
		state = multicastRouterModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, multicastRouter := range multicastRouters {
		if multicastRouter.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(multicastRouter.UUID)
			state.Name = stringValueOrNull(multicastRouter.Name)
			state.Description = stringValueOrNull(multicastRouter.Description)
			state.State = stringValueOrNull(multicastRouter.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Multicast router not found. It might have been deleted outside of Terraform.")
		state = multicastRouterModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *multicastRouterResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r *multicastRouterResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state multicastRouterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Multicast router UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteMulticastRouter(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete multicast router", ""+err.Error())
		return
	}
}

func (r *multicastRouterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
