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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &loadBalancerResource{}
	_ resource.ResourceWithConfigure   = &loadBalancerResource{}
	_ resource.ResourceWithImportState = &loadBalancerResource{}
)

type loadBalancerResource struct {
	client *client.ZSClient
}

type loadBalancerResourceModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	VipUuid         types.String `tfsdk:"vip_uuid"`
	State           types.String `tfsdk:"state"`
	Type            types.String `tfsdk:"type"`
	ServerGroupUuid types.String `tfsdk:"server_group_uuid"`
}

func LoadBalancerResource() resource.Resource {
	return &loadBalancerResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *loadBalancerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *loadBalancerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer"
}

// Schema implements resource.Resource.
func (r *loadBalancerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack load balancers. A load balancer distributes incoming traffic across multiple backend servers using a VIP address.",
		MarkdownDescription: "Manage ZStack load balancers. A load balancer distributes incoming traffic across multiple backend servers using a VIP address.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the load balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the load balancer.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the load balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VIP bound to the load balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the load balancer (Enabled, Disabled).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the load balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_group_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the default server group associated with the load balancer.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *loadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loadBalancerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating load balancer", map[string]any{"name": plan.Name.ValueString()})

	createParam := param.CreateLoadBalancerParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateLoadBalancerParamDetail{
			Name:    plan.Name.ValueString(),
			VipUuid: stringPtr(plan.VipUuid.ValueString()),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	lb, err := r.client.CreateLoadBalancer(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Load Balancer", "Could not create load balancer, unexpected error: "+err.Error())
		return
	}

	state := loadBalancerModelFromView(lb)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *loadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loadBalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	lb, err := findResourceByGet(r.client.GetLoadBalancer, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Load Balancer", "Could not read load balancer, unexpected error: "+err.Error())
		return
	}

	refreshedState := loadBalancerModelFromView(lb)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *loadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan loadBalancerResourceModel
	var state loadBalancerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		updateParam := param.UpdateLoadBalancerParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateLoadBalancerParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		if _, err := r.client.UpdateLoadBalancer(uuid, updateParam); err != nil {
			resp.Diagnostics.AddError("Error updating Load Balancer", "Could not update load balancer, unexpected error: "+err.Error())
			return
		}
	}

	// Read back the updated resource
	lb, err := r.client.GetLoadBalancer(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Load Balancer", "Could not read load balancer after update: "+err.Error())
		return
	}

	refreshedState := loadBalancerModelFromView(lb)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *loadBalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loadBalancerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "load balancer uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteLoadBalancer(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting Load Balancer", "Could not delete load balancer, unexpected error: "+err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *loadBalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func loadBalancerModelFromView(lb *view.LoadBalancerInventoryView) loadBalancerResourceModel {
	return loadBalancerResourceModel{
		Uuid:            types.StringValue(lb.UUID),
		Name:            types.StringValue(lb.Name),
		Description:     stringValueOrNull(lb.Description),
		VipUuid:         types.StringValue(lb.VipUuid),
		State:           stringValueOrNull(lb.State),
		Type:            stringValueOrNull(lb.Type),
		ServerGroupUuid: stringValueOrNull(lb.ServerGroupUuid),
	}
}
