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
	_ resource.Resource                = &autoScalingGroupResource{}
	_ resource.ResourceWithConfigure   = &autoScalingGroupResource{}
	_ resource.ResourceWithImportState = &autoScalingGroupResource{}
)

type autoScalingGroupResource struct {
	client *client.ZSClient
}

type autoScalingGroupResourceModel struct {
	Uuid                types.String `tfsdk:"uuid"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	ScalingResourceType types.String `tfsdk:"scaling_resource_type"`
	State               types.String `tfsdk:"state"`
	DefaultCooldown     types.Int64  `tfsdk:"default_cooldown"`
	MinResourceSize     types.Int64  `tfsdk:"min_resource_size"`
	MaxResourceSize     types.Int64  `tfsdk:"max_resource_size"`
	RemovalPolicy       types.String `tfsdk:"removal_policy"`
}

func AutoScalingGroupResource() resource.Resource {
	return &autoScalingGroupResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *autoScalingGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *autoScalingGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_scaling_group"
}

// Schema implements resource.Resource.
func (r *autoScalingGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack auto scaling groups. An auto scaling group automatically adjusts the number of VM instances based on scaling rules and policies.",
		MarkdownDescription: "Manage ZStack auto scaling groups. An auto scaling group automatically adjusts the number of VM instances based on scaling rules and policies.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the auto scaling group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the auto scaling group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the auto scaling group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scaling_resource_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of resource to scale (e.g., VmInstance).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the auto scaling group (Enabled, Disabled).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_cooldown": schema.Int64Attribute{
				Required:    true,
				Description: "The default cooldown period in seconds between two scaling activities.",
			},
			"min_resource_size": schema.Int64Attribute{
				Required:    true,
				Description: "The minimum number of instances in the scaling group.",
			},
			"max_resource_size": schema.Int64Attribute{
				Required:    true,
				Description: "The maximum number of instances in the scaling group.",
			},
			"removal_policy": schema.StringAttribute{
				Required:    true,
				Description: "The policy for removing instances when scaling in (e.g., OldestInstance, NewestInstance, OldestScalingConfiguration).",
				Validators: []validator.String{
					stringvalidator.OneOf("OldestInstance", "NewestInstance", "OldestScalingConfiguration"),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *autoScalingGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan autoScalingGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating auto scaling group", map[string]any{"name": plan.Name.ValueString()})

	createParam := param.CreateAutoScalingGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAutoScalingGroupParamDetail{
			Name:                plan.Name.ValueString(),
			ScalingResourceType: plan.ScalingResourceType.ValueString(),
			MinResourceSize:     int(plan.MinResourceSize.ValueInt64()),
			MaxResourceSize:     int(plan.MaxResourceSize.ValueInt64()),
			DefaultCooldown:     plan.DefaultCooldown.ValueInt64(),
			RemovalPolicy:       plan.RemovalPolicy.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	group, err := r.client.CreateAutoScalingGroup(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Auto Scaling Group", "Could not create auto scaling group, unexpected error: "+err.Error())
		return
	}

	state := autoScalingGroupModelFromView(group)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *autoScalingGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state autoScalingGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := findResourceByGet(r.client.GetAutoScalingGroup, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Auto Scaling Group", "Could not read auto scaling group, unexpected error: "+err.Error())
		return
	}

	refreshedState := autoScalingGroupModelFromView(group)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *autoScalingGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan autoScalingGroupResourceModel
	var state autoScalingGroupResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	updateParam := param.UpdateAutoScalingGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAutoScalingGroupParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			MinResourceSize: intPtr(int(plan.MinResourceSize.ValueInt64())),
			MaxResourceSize: intPtr(int(plan.MaxResourceSize.ValueInt64())),
			RemovalPolicy:   stringPtr(plan.RemovalPolicy.ValueString()),
		},
	}

	if _, err := r.client.UpdateAutoScalingGroup(uuid, updateParam); err != nil {
		resp.Diagnostics.AddError("Error updating Auto Scaling Group", "Could not update auto scaling group, unexpected error: "+err.Error())
		return
	}

	// Read back the updated resource
	group, err := r.client.GetAutoScalingGroup(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Auto Scaling Group", "Could not read updated auto scaling group: "+err.Error())
		return
	}

	refreshedState := autoScalingGroupModelFromView(group)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *autoScalingGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state autoScalingGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "auto scaling group uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteAutoScalingGroup(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting Auto Scaling Group", "Could not delete auto scaling group, unexpected error: "+err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *autoScalingGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func autoScalingGroupModelFromView(g *view.AutoScalingGroupInventoryView) autoScalingGroupResourceModel {
	return autoScalingGroupResourceModel{
		Uuid:                types.StringValue(g.UUID),
		Name:                types.StringValue(g.Name),
		Description:         stringValueOrNull(g.Description),
		ScalingResourceType: types.StringValue(g.ScalingResourceType),
		State:               stringValueOrNull(g.State),
		DefaultCooldown:     types.Int64Value(g.DefaultCooldown),
		MinResourceSize:     types.Int64Value(int64(g.MinResourceSize)),
		MaxResourceSize:     types.Int64Value(int64(g.MaxResourceSize)),
		RemovalPolicy:       stringValueOrNull(g.RemovalPolicy),
	}
}
