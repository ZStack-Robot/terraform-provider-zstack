// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &loadBalancerListenerResource{}
	_ resource.ResourceWithConfigure   = &loadBalancerListenerResource{}
	_ resource.ResourceWithImportState = &loadBalancerListenerResource{}
)

type loadBalancerListenerResource struct {
	client *client.ZSClient
}

type loadBalancerListenerResourceModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	LoadBalancerUuid   types.String `tfsdk:"load_balancer_uuid"`
	Protocol           types.String `tfsdk:"protocol"`
	LoadBalancerPort   types.Int64  `tfsdk:"load_balancer_port"`
	InstancePort       types.Int64  `tfsdk:"instance_port"`
	SecurityPolicyType types.String `tfsdk:"security_policy_type"`
	ServerGroupUuid    types.String `tfsdk:"server_group_uuid"`
}

func LoadBalancerListenerResource() resource.Resource {
	return &loadBalancerListenerResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *loadBalancerListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *loadBalancerListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_listener"
}

// Schema implements resource.Resource.
func (r *loadBalancerListenerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack load balancer listeners. A listener defines the protocol and port configuration for a load balancer, routing traffic from a frontend port to backend instance ports.",
		MarkdownDescription: "Manage ZStack load balancer listeners. A listener defines the protocol and port configuration for a load balancer, routing traffic from a frontend port to backend instance ports.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the load balancer listener.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the load balancer listener.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the load balancer listener.",
			},
			"load_balancer_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the load balancer this listener belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "The protocol for the listener: tcp, udp, http, or https.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"load_balancer_port": schema.Int64Attribute{
				Required:    true,
				Description: "The port on the load balancer (frontend port) that receives traffic.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"instance_port": schema.Int64Attribute{
				Required:    true,
				Description: "The port on the backend instances that receives forwarded traffic.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"security_policy_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The security policy type for HTTPS listeners (e.g., tls_cipher_policy_default).",
			},
			"server_group_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the server group associated with this listener.",
			},
		},
	}
}

// Create implements resource.Resource.
//
// SDK Workaround: CreateLoadBalancerListener uses cli.Post("v1/load-balancers/{loadBalancerUuid}/listeners")
// which has the URL template bug (ZSClient.Post shadows ZSHttpClient.Post, templates are not resolved).
// We call ZSHttpClient.Post directly with the resolved URL.
// See docs/SDK_URL_TEMPLATE_BUG.md for details.
func (r *loadBalancerListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loadBalancerListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating load balancer listener", map[string]any{
		"name":              plan.Name.ValueString(),
		"load_balancer_uuid": plan.LoadBalancerUuid.ValueString(),
	})

	createParam := param.CreateLoadBalancerListenerParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateLoadBalancerListenerParamDetail{
			Name:             plan.Name.ValueString(),
			LoadBalancerPort: int(plan.LoadBalancerPort.ValueInt64()),
			Protocol:         stringPtr(plan.Protocol.ValueString()),
			InstancePort:     intPtr(int(plan.InstancePort.ValueInt64())),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.SecurityPolicyType.IsNull() && plan.SecurityPolicyType.ValueString() != "" {
		createParam.Params.SecurityPolicyType = stringPtr(plan.SecurityPolicyType.ValueString())
	}

	listener, err := r.client.CreateLoadBalancerListener(plan.LoadBalancerUuid.ValueString(), createParam)
	if err != nil {
		resp.Diagnostics.AddError("Could not create load balancer listener", err.Error())
		return
	}

	state := loadBalancerListenerModelFromView(listener)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *loadBalancerListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loadBalancerListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listener, err := r.client.GetLoadBalancerListener(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not read load balancer listener", err.Error())
		return
	}

	refreshedState := loadBalancerListenerModelFromView(listener)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *loadBalancerListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan loadBalancerListenerResourceModel
	var state loadBalancerListenerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update name/description if changed
	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		updateParam := param.UpdateLoadBalancerListenerParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateLoadBalancerListenerParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		if _, err := r.client.UpdateLoadBalancerListener(uuid, updateParam); err != nil {
			resp.Diagnostics.AddError("Could not update load balancer listener", err.Error())
			return
		}
	}

	// Read back the updated resource
	listener, err := r.client.GetLoadBalancerListener(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Could not read updated load balancer listener", err.Error())
		return
	}

	refreshedState := loadBalancerListenerModelFromView(listener)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *loadBalancerListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loadBalancerListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "load balancer listener uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteLoadBalancerListener(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Could not delete load balancer listener", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *loadBalancerListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func loadBalancerListenerModelFromView(l *view.LoadBalancerListenerInventoryView) loadBalancerListenerResourceModel {
	return loadBalancerListenerResourceModel{
		Uuid:               types.StringValue(l.UUID),
		Name:               types.StringValue(l.Name),
		Description:        stringValueOrNull(l.Description),
		LoadBalancerUuid:   types.StringValue(l.LoadBalancerUuid),
		Protocol:           stringValueOrNull(l.Protocol),
		LoadBalancerPort:   types.Int64Value(int64(l.LoadBalancerPort)),
		InstancePort:       types.Int64Value(int64(l.InstancePort)),
		SecurityPolicyType: stringValueOrNull(l.SecurityPolicyType),
		ServerGroupUuid:    stringValueOrNull(l.ServerGroupUuid),
	}
}
