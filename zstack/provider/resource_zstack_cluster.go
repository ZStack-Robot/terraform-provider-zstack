// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	_ resource.Resource                = &clusterResource{}
	_ resource.ResourceWithConfigure   = &clusterResource{}
	_ resource.ResourceWithImportState = &clusterResource{}
)

type clusterResource struct {
	client *client.ZSClient
}

type clusterResourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	State          types.String `tfsdk:"state"`
	HypervisorType types.String `tfsdk:"hypervisor_type"`
	ZoneUuid       types.String `tfsdk:"zone_uuid"`
	Type           types.String `tfsdk:"type"`
	Architecture   types.String `tfsdk:"architecture"`
}

func ClusterResource() resource.Resource {
	return &clusterResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

// Schema implements resource.Resource.
func (r *clusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack clusters. A cluster is a logical grouping of compute hosts within a zone.",
		MarkdownDescription: "Manage ZStack clusters. A cluster is a logical grouping of compute hosts within a zone.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the cluster.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone to which the cluster belongs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hypervisor_type": schema.StringAttribute{
				Required:    true,
				Description: "The type of hypervisor used by the cluster (e.g., KVM).",
				Validators: []validator.String{
					stringvalidator.OneOf("KVM"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The state of the cluster (Enabled, Disabled).",
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the cluster.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"architecture": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The architecture of the cluster (e.g., x86_64, aarch64).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating cluster", map[string]any{"name": plan.Name.ValueString()})

	createParam := param.CreateClusterParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateClusterParamDetail{
			Name:           plan.Name.ValueString(),
			ZoneUuid:       plan.ZoneUuid.ValueString(),
			HypervisorType: plan.HypervisorType.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	if !plan.Type.IsNull() && plan.Type.ValueString() != "" {
		createParam.Params.Type = stringPtr(plan.Type.ValueString())
	}

	if !plan.Architecture.IsNull() && plan.Architecture.ValueString() != "" {
		createParam.Params.Architecture = stringPtr(plan.Architecture.ValueString())
	}

	cluster, err := r.client.CreateCluster(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cluster",
			"Could not create cluster, unexpected error: "+err.Error(),
		)
		return
	}

	// If the desired state is Disabled, change the cluster state after creation
	if !plan.State.IsNull() && strings.EqualFold(plan.State.ValueString(), "Disabled") {
		stateParam := param.ChangeClusterStateParam{
			BaseParam: param.BaseParam{},
			Params: param.ChangeClusterStateParamDetail{
				StateEvent: "disable",
			},
		}
		cluster, err = r.client.ChangeClusterState(cluster.UUID, stateParam)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error changing Cluster state",
				"Could not change cluster state to Disabled: "+err.Error(),
			)
			return
		}
	}

	state := clusterModelFromView(cluster)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := findResourceByGet(r.client.GetCluster, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Cluster",
			"Could not read cluster, unexpected error: "+err.Error(),
		)
		return
	}

	refreshedState := clusterModelFromView(cluster)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan clusterResourceModel
	var state clusterResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update name and description
	updateParam := param.UpdateClusterParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateClusterParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if _, err := r.client.UpdateCluster(uuid, updateParam); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Cluster",
			"Could not update cluster, unexpected error: "+err.Error(),
		)
		return
	}

	// Handle state change if needed
	if !plan.State.IsNull() && plan.State.ValueString() != state.State.ValueString() {
		var stateEvent string
		if strings.EqualFold(plan.State.ValueString(), "Enabled") {
			stateEvent = "enable"
		} else if strings.EqualFold(plan.State.ValueString(), "Disabled") {
			stateEvent = "disable"
		}

		if stateEvent != "" {
			stateParam := param.ChangeClusterStateParam{
				BaseParam: param.BaseParam{},
				Params: param.ChangeClusterStateParamDetail{
					StateEvent: stateEvent,
				},
			}
			if _, err := r.client.ChangeClusterState(uuid, stateParam); err != nil {
				resp.Diagnostics.AddError(
					"Error changing Cluster state",
					"Could not change cluster state: "+err.Error(),
				)
				return
			}
		}
	}

	// Read back the updated resource
	cluster, err := r.client.GetCluster(uuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Cluster",
			"Could not read updated cluster: "+err.Error(),
		)
		return
	}

	refreshedState := clusterModelFromView(cluster)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "cluster uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteCluster(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Cluster",
			"Could not delete cluster, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func clusterModelFromView(c *view.ClusterInventoryView) clusterResourceModel {
	return clusterResourceModel{
		Uuid:           types.StringValue(c.UUID),
		Name:           types.StringValue(c.Name),
		Description:    stringValueOrNull(c.Description),
		State:          stringValueOrNull(c.State),
		HypervisorType: stringValueOrNull(c.HypervisorType),
		ZoneUuid:       types.StringValue(c.ZoneUuid),
		Type:           stringValueOrNull(c.Type),
		Architecture:   stringValueOrNull(c.Architecture),
	}
}
