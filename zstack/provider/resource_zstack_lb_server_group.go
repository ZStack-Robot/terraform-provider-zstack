// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &lbServerGroupResource{}
	_ resource.ResourceWithConfigure   = &lbServerGroupResource{}
	_ resource.ResourceWithImportState = &lbServerGroupResource{}
)

type lbServerGroupResource struct {
	client *client.ZSClient
}

type lbServerGroupModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	LoadBalancerUuid types.String `tfsdk:"load_balancer_uuid"`
	IpVersion        types.Int64  `tfsdk:"ip_version"`
}

func LBServerGroupResource() resource.Resource {
	return &lbServerGroupResource{}
}

func (r *lbServerGroupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *lbServerGroupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_lb_server_group"
}

func (r *lbServerGroupResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Load Balancer Server Groups in ZStack. " +
			"A Load Balancer Server Group is a collection of backend servers for a load balancer.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Load Balancer Server Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Load Balancer Server Group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the Load Balancer Server Group.",
			},
			"load_balancer_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the Load Balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_version": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The IP version (4 for IPv4, 6 for IPv6).",
				Validators: []validator.Int64{
					int64validator.OneOf(4, 6),
				},
			},
		},
	}
}

func (r *lbServerGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan lbServerGroupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateLoadBalancerServerGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateLoadBalancerServerGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			IpVersion:   intPtrFromInt64OrNil(plan.IpVersion),
		},
	}

	lbServerGroup, err := r.client.CreateLoadBalancerServerGroup(plan.LoadBalancerUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Load Balancer Server Group",
			"Could not create load balancer server group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(lbServerGroup.UUID)
	plan.Name = types.StringValue(lbServerGroup.Name)
	plan.Description = stringValueOrNull(lbServerGroup.Description)
	plan.LoadBalancerUuid = types.StringValue(lbServerGroup.LoadBalancerUuid)
	plan.IpVersion = types.Int64Value(int64(lbServerGroup.IpVersion))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *lbServerGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state lbServerGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	lbServerGroup, err := findResourceByQuery(r.client.QueryLoadBalancerServerGroup, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query Load Balancer Server Groups. It may have been deleted.: "+err.Error())
		state = lbServerGroupModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(lbServerGroup.UUID)
	state.Name = types.StringValue(lbServerGroup.Name)
	state.Description = stringValueOrNull(lbServerGroup.Description)
	state.LoadBalancerUuid = types.StringValue(lbServerGroup.LoadBalancerUuid)
	state.IpVersion = types.Int64Value(int64(lbServerGroup.IpVersion))

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *lbServerGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan lbServerGroupModel
	var state lbServerGroupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateLoadBalancerServerGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateLoadBalancerServerGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	lbServerGroup, err := r.client.UpdateLoadBalancerServerGroup(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Load Balancer Server Group",
			"Could not update load balancer server group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(lbServerGroup.UUID)
	plan.Name = types.StringValue(lbServerGroup.Name)
	plan.Description = stringValueOrNull(lbServerGroup.Description)
	plan.LoadBalancerUuid = types.StringValue(lbServerGroup.LoadBalancerUuid)
	plan.IpVersion = types.Int64Value(int64(lbServerGroup.IpVersion))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *lbServerGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state lbServerGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Load Balancer Server Group UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteLoadBalancerServerGroup(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Error deleting Load Balancer Server Group", "Could not delete load balancer server group, unexpected error: "+err.Error())
		return
	}
}

func (r *lbServerGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
