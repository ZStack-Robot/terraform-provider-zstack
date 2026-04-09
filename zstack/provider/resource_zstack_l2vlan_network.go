// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
	_ resource.Resource                = &l2VlanNetworkResource{}
	_ resource.ResourceWithConfigure   = &l2VlanNetworkResource{}
	_ resource.ResourceWithImportState = &l2VlanNetworkResource{}
)

type l2VlanNetworkResource struct {
	client *client.ZSClient
}

type l2VlanNetworkResourceModel struct {
	Uuid                 types.String   `tfsdk:"uuid"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	Vlan                 types.Int64    `tfsdk:"vlan"`
	ZoneUuid             types.String   `tfsdk:"zone_uuid"`
	PhysicalInterface    types.String   `tfsdk:"physical_interface"`
	Type                 types.String   `tfsdk:"type"`
	VSwitchType          types.String   `tfsdk:"vswitch_type"`
	AttachedClusterUuids []types.String `tfsdk:"attached_cluster_uuids"`
}

func L2VlanNetworkResource() resource.Resource {
	return &l2VlanNetworkResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *l2VlanNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *l2VlanNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vlan_network"
}

// Schema implements resource.Resource.
func (r *l2VlanNetworkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack L2 VLAN networks. L2 VLAN networks provide VLAN-based network isolation and can be attached to clusters.",
		MarkdownDescription: "Manage ZStack L2 VLAN networks. L2 VLAN networks provide VLAN-based network isolation and can be attached to clusters.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the L2 VLAN network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the L2 VLAN network.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the L2 VLAN network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan": schema.Int64Attribute{
				Required:    true,
				Description: "The VLAN ID for this L2 network (1-4094).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone where the L2 VLAN network resides.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"physical_interface": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The physical network interface (e.g., eth0, bond0) used by this L2 VLAN network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the L2 network (e.g., L2VlanNetwork).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vswitch_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The virtual switch type (e.g., LinuxBridge, OvsDpdk).",
				Validators: []validator.String{
					stringvalidator.OneOf("LinuxBridge", "OvsDpdk"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"attached_cluster_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The list of cluster UUIDs to which this L2 VLAN network is attached.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *l2VlanNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan l2VlanNetworkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating L2 VLAN network", map[string]any{"name": plan.Name.ValueString(), "vlan": plan.Vlan.ValueInt64()})

	createParam := param.CreateL2VlanNetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL2VlanNetworkParamDetail{
			Vlan:     int(plan.Vlan.ValueInt64()),
			Name:     plan.Name.ValueString(),
			ZoneUuid: plan.ZoneUuid.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.PhysicalInterface.IsNull() && plan.PhysicalInterface.ValueString() != "" {
		createParam.Params.PhysicalInterface = stringPtr(plan.PhysicalInterface.ValueString())
	}
	if !plan.VSwitchType.IsNull() && plan.VSwitchType.ValueString() != "" {
		createParam.Params.VSwitchType = stringPtr(plan.VSwitchType.ValueString())
	}

	l2Network, err := r.client.CreateL2VlanNetwork(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Error creating L2 VLAN Network", "Could not create L2 VLAN network, unexpected error: "+err.Error())
		return
	}

	// Save partial state so the L2 VLAN network UUID is tracked even if cluster attachment fails
	partialState, err := r.readL2VlanNetwork(l2Network.UUID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading L2 VLAN Network", "Could not read L2 VLAN network after create: "+err.Error())
		return
	}
	diags = resp.State.Set(ctx, &partialState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Attach to clusters if specified
	desiredClusters := terraformStringsToSlice(plan.AttachedClusterUuids)
	for _, clusterUuid := range desiredClusters {
		if err := r.attachCluster(l2Network.UUID, clusterUuid); err != nil {
			resp.Diagnostics.AddError(
				"Error attaching Cluster to L2 VLAN Network",
				fmt.Sprintf("Could not attach cluster %s to L2 VLAN network: %s", clusterUuid, err.Error()),
			)
			return
		}
	}

	// Read back the created resource
	state, err := r.readL2VlanNetwork(l2Network.UUID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading L2 VLAN Network", "Could not read L2 VLAN network after create: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *l2VlanNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state l2VlanNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.readL2VlanNetwork(state.Uuid.ValueString())
	if err != nil {
		if isZStackNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading L2 VLAN Network", "Could not read L2 VLAN network, unexpected error: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *l2VlanNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan l2VlanNetworkResourceModel
	var state l2VlanNetworkResourceModel

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
		updateParam := param.UpdateL2NetworkParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateL2NetworkParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		if _, err := r.client.UpdateL2Network(uuid, updateParam); err != nil {
			resp.Diagnostics.AddError("Error updating L2 VLAN Network", "Could not update L2 VLAN network, unexpected error: "+err.Error())
			return
		}
	}

	// Reconcile cluster attachments
	if err := r.reconcileClusterAttachments(uuid, state.AttachedClusterUuids, plan.AttachedClusterUuids); err != nil {
		resp.Diagnostics.AddError("Error updating L2 VLAN Network Cluster Attachments", "Could not update L2 VLAN network cluster attachments, unexpected error: "+err.Error())
		return
	}

	// Read back the updated resource
	refreshedState, err := r.readL2VlanNetwork(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading L2 VLAN Network", "Could not read L2 VLAN network after update: "+err.Error())
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *l2VlanNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state l2VlanNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "L2 VLAN network uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteL2Network(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting L2 VLAN Network", "Could not delete L2 VLAN network, unexpected error: "+err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *l2VlanNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *l2VlanNetworkResource) readL2VlanNetwork(uuid string) (l2VlanNetworkResourceModel, error) {
	l2Network, err := r.client.GetL2VlanNetwork(uuid)
	if err != nil {
		return l2VlanNetworkResourceModel{}, err
	}

	return l2VlanNetworkModelFromView(l2Network), nil
}

func (r *l2VlanNetworkResource) attachCluster(l2NetworkUuid, clusterUuid string) error {
	attachParam := param.AttachL2NetworkToClusterParam{
		BaseParam: param.BaseParam{},
		Params:    param.AttachL2NetworkToClusterParamDetail{},
	}

	if _, err := r.client.AttachL2NetworkToCluster(l2NetworkUuid, clusterUuid, attachParam); err != nil {
		return err
	}
	return nil
}

func (r *l2VlanNetworkResource) reconcileClusterAttachments(uuid string, current, desired []types.String) error {
	currentSet := make(map[string]bool)
	for _, c := range current {
		if !c.IsNull() && c.ValueString() != "" {
			currentSet[c.ValueString()] = true
		}
	}

	desiredSet := make(map[string]bool)
	for _, d := range desired {
		if !d.IsNull() && d.ValueString() != "" {
			desiredSet[d.ValueString()] = true
		}
	}

	// Detach clusters that are no longer desired
	for clusterUuid := range currentSet {
		if !desiredSet[clusterUuid] {
			if err := r.client.DetachL2NetworkFromCluster(uuid, clusterUuid, param.DeleteModePermissive); err != nil {
				return fmt.Errorf("failed to detach cluster %s: %w", clusterUuid, err)
			}
		}
	}

	// Attach clusters that are newly desired
	for clusterUuid := range desiredSet {
		if !currentSet[clusterUuid] {
			if err := r.attachCluster(uuid, clusterUuid); err != nil {
				return fmt.Errorf("failed to attach cluster %s: %w", clusterUuid, err)
			}
		}
	}

	return nil
}

func l2VlanNetworkModelFromView(l2Network *view.L2VlanNetworkInventoryView) l2VlanNetworkResourceModel {
	attachedClusters := make([]types.String, 0, len(l2Network.AttachedClusterUuids))
	for _, clusterUuid := range l2Network.AttachedClusterUuids {
		attachedClusters = append(attachedClusters, types.StringValue(clusterUuid))
	}

	return l2VlanNetworkResourceModel{
		Uuid:                 types.StringValue(l2Network.UUID),
		Name:                 types.StringValue(l2Network.Name),
		Description:          stringValueOrNull(l2Network.Description),
		Vlan:                 types.Int64Value(int64(l2Network.Vlan)),
		ZoneUuid:             types.StringValue(l2Network.ZoneUuid),
		PhysicalInterface:    stringValueOrNull(l2Network.PhysicalInterface),
		Type:                 stringValueOrNull(l2Network.Type),
		VSwitchType:          stringValueOrNull(l2Network.VSwitchType),
		AttachedClusterUuids: attachedClusters,
	}
}
