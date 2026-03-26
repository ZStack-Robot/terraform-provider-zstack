// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &primaryStorageResource{}
	_ resource.ResourceWithConfigure   = &primaryStorageResource{}
	_ resource.ResourceWithImportState = &primaryStorageResource{}
)

type primaryStorageResource struct {
	client *client.ZSClient
}

type primaryStorageResourceModel struct {
	Uuid                      types.String `tfsdk:"uuid"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	ZoneUuid                  types.String `tfsdk:"zone_uuid"`
	Type                      types.String `tfsdk:"type"`
	Url                       types.String `tfsdk:"url"`
	State                     types.String `tfsdk:"state"`
	Status                    types.String `tfsdk:"status"`
	TotalCapacity             types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity         types.Int64  `tfsdk:"available_capacity"`
	TotalPhysicalCapacity     types.Int64  `tfsdk:"total_physical_capacity"`
	AvailablePhysicalCapacity types.Int64  `tfsdk:"available_physical_capacity"`
	MountPath                 types.String `tfsdk:"mount_path"`
	AttachedClusterUuids      types.List   `tfsdk:"attached_cluster_uuids"`
	MonUrls                   types.List   `tfsdk:"mon_urls"`
	RootVolumePoolName        types.String `tfsdk:"root_volume_pool_name"`
	DataVolumePoolName        types.String `tfsdk:"data_volume_pool_name"`
	ImageCachePoolName        types.String `tfsdk:"image_cache_pool_name"`
	DiskUuids                 types.List   `tfsdk:"disk_uuids"`
}

func PrimaryStorageResource() resource.Resource {
	return &primaryStorageResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *primaryStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *primaryStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_primary_storage"
}

// Schema implements resource.Resource.
func (r *primaryStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack primary storage. Supports LocalStorage, NFS, Ceph, and SharedBlock types.",
		MarkdownDescription: "Manage ZStack primary storage. Supports LocalStorage, NFS, Ceph, and SharedBlock types.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the primary storage.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the primary storage.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the primary storage.",
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the zone this primary storage belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of primary storage (LocalStorage, NFS, Ceph, SharedBlock).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The URL/path for the primary storage (required for LocalStorage and NFS types).",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the primary storage (Enabled, Disabled).",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the primary storage (Connected, Disconnected).",
			},
			"total_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The total capacity in bytes.",
			},
			"available_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The available capacity in bytes.",
			},
			"total_physical_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The total physical capacity in bytes.",
			},
			"available_physical_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The available physical capacity in bytes.",
			},
			"mount_path": schema.StringAttribute{
				Computed:    true,
				Description: "The mount path of the primary storage.",
			},
			"attached_cluster_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of cluster UUIDs to attach this primary storage to.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"mon_urls": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of Ceph monitor URLs (required for Ceph type).",
			},
			"root_volume_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The root volume pool name (for Ceph type).",
			},
			"data_volume_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The data volume pool name (for Ceph type).",
			},
			"image_cache_pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The image cache pool name (for Ceph type).",
			},
			"disk_uuids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of shared block disk UUIDs (required for SharedBlock type).",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *primaryStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan primaryStorageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating primary storage", map[string]any{"name": plan.Name.ValueString(), "type": plan.Type.ValueString()})

	var ps *view.PrimaryStorageInventoryView
	var err error

	switch plan.Type.ValueString() {
	case "LocalStorage":
		addLocalParam := param.AddLocalPrimaryStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddLocalPrimaryStorageParamDetail{
				Name:     plan.Name.ValueString(),
				Url:      plan.Url.ValueString(),
				ZoneUuid: plan.ZoneUuid.ValueString(),
			},
		}
		if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
			addLocalParam.Params.Description = stringPtr(plan.Description.ValueString())
		}
		// SDK bug: AddLocalPrimaryStorage() takes no parameters, use ZSHttpClient.Post directly
		ps = &view.PrimaryStorageInventoryView{}
		err = r.client.ZSHttpClient.Post("v1/primary-storage/local-storage", addLocalParam, ps)

	case "NFS":
		addNfsParam := param.AddNfsPrimaryStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddNfsPrimaryStorageParamDetail{
				Name:     plan.Name.ValueString(),
				Url:      plan.Url.ValueString(),
				ZoneUuid: plan.ZoneUuid.ValueString(),
			},
		}
		if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
			addNfsParam.Params.Description = stringPtr(plan.Description.ValueString())
		}
		// SDK bug: AddNfsPrimaryStorage() takes no parameters, use ZSHttpClient.Post directly
		ps = &view.PrimaryStorageInventoryView{}
		err = r.client.ZSHttpClient.Post("v1/primary-storage/nfs", addNfsParam, ps)

	case "Ceph":
		monUrls := listToStringSlice(plan.MonUrls)
		addCephParam := param.AddCephPrimaryStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddCephPrimaryStorageParamDetail{
				Name:     plan.Name.ValueString(),
				MonUrls:  monUrls,
				ZoneUuid: plan.ZoneUuid.ValueString(),
			},
		}
		if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
			addCephParam.Params.Description = stringPtr(plan.Description.ValueString())
		}
		if !plan.Url.IsNull() && plan.Url.ValueString() != "" {
			addCephParam.Params.Url = stringPtr(plan.Url.ValueString())
		}
		if !plan.RootVolumePoolName.IsNull() && plan.RootVolumePoolName.ValueString() != "" {
			addCephParam.Params.RootVolumePoolName = stringPtr(plan.RootVolumePoolName.ValueString())
		}
		if !plan.DataVolumePoolName.IsNull() && plan.DataVolumePoolName.ValueString() != "" {
			addCephParam.Params.DataVolumePoolName = stringPtr(plan.DataVolumePoolName.ValueString())
		}
		if !plan.ImageCachePoolName.IsNull() && plan.ImageCachePoolName.ValueString() != "" {
			addCephParam.Params.ImageCachePoolName = stringPtr(plan.ImageCachePoolName.ValueString())
		}
		ps, err = r.client.AddCephPrimaryStorage(addCephParam)

	case "SharedBlock":
		diskUuids := listToStringSlice(plan.DiskUuids)
		addSharedBlockParam := param.AddSharedBlockGroupPrimaryStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddSharedBlockGroupPrimaryStorageParamDetail{
				Name:      plan.Name.ValueString(),
				DiskUuids: diskUuids,
				ZoneUuid:  plan.ZoneUuid.ValueString(),
			},
		}
		if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
			addSharedBlockParam.Params.Description = stringPtr(plan.Description.ValueString())
		}
		if !plan.Url.IsNull() && plan.Url.ValueString() != "" {
			addSharedBlockParam.Params.Url = stringPtr(plan.Url.ValueString())
		}
		ps, err = r.client.AddSharedBlockGroupPrimaryStorage(addSharedBlockParam)

	default:
		resp.Diagnostics.AddError("Unsupported primary storage type", fmt.Sprintf("Type %q is not supported. Use LocalStorage, NFS, Ceph, or SharedBlock.", plan.Type.ValueString()))
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Could not create primary storage", err.Error())
		return
	}

	// Handle cluster attachment after creation
	desiredClusters := listToStringSlice(plan.AttachedClusterUuids)
	for _, clusterUuid := range desiredClusters {
		var attachResult view.PrimaryStorageInventoryView
		if err := r.client.ZSHttpClient.Post(
			fmt.Sprintf("v1/clusters/%s/primary-storage/%s", clusterUuid, ps.UUID),
			param.AttachPrimaryStorageToClusterParam{BaseParam: param.BaseParam{}},
			&attachResult,
		); err != nil {
			resp.Diagnostics.AddError("Could not attach primary storage to cluster",
				fmt.Sprintf("Failed to attach to cluster %s: %s", clusterUuid, err.Error()))
			return
		}
	}

	// Read back the final state
	psRead, err := r.client.GetPrimaryStorage(ps.UUID)
	if err != nil {
		resp.Diagnostics.AddError("Could not read created primary storage", err.Error())
		return
	}

	state := primaryStorageModelFromView(psRead, &plan)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *primaryStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state primaryStorageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ps, err := r.client.GetPrimaryStorage(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Could not read primary storage", err.Error())
		return
	}

	refreshedState := primaryStorageModelFromView(ps, &state)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *primaryStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan primaryStorageResourceModel
	var state primaryStorageResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update name, description, url if changed
	updateParam := param.UpdatePrimaryStorageParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdatePrimaryStorageParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}
	if !plan.Url.IsNull() && plan.Url.ValueString() != "" {
		updateParam.Params.Url = stringPtr(plan.Url.ValueString())
	}

	if _, err := r.client.UpdatePrimaryStorage(uuid, updateParam); err != nil {
		resp.Diagnostics.AddError("Could not update primary storage", err.Error())
		return
	}

	// Handle cluster attach/detach
	currentClusters := listToStringSlice(state.AttachedClusterUuids)
	desiredClusters := listToStringSlice(plan.AttachedClusterUuids)

	currentSet := make(map[string]bool)
	for _, c := range currentClusters {
		currentSet[c] = true
	}
	desiredSet := make(map[string]bool)
	for _, c := range desiredClusters {
		desiredSet[c] = true
	}

	// Attach new clusters
	for _, clusterUuid := range desiredClusters {
		if !currentSet[clusterUuid] {
			var attachResult view.PrimaryStorageInventoryView
			if err := r.client.ZSHttpClient.Post(
				fmt.Sprintf("v1/clusters/%s/primary-storage/%s", clusterUuid, uuid),
				param.AttachPrimaryStorageToClusterParam{BaseParam: param.BaseParam{}},
				&attachResult,
			); err != nil {
				resp.Diagnostics.AddError("Could not attach primary storage to cluster",
					fmt.Sprintf("Failed to attach to cluster %s: %s", clusterUuid, err.Error()))
				return
			}
		}
	}

	// Detach removed clusters
	for _, clusterUuid := range currentClusters {
		if !desiredSet[clusterUuid] {
			if err := r.client.DetachPrimaryStorageFromCluster(clusterUuid, uuid, param.DeleteModePermissive); err != nil {
				resp.Diagnostics.AddError("Could not detach primary storage from cluster",
					fmt.Sprintf("Failed to detach from cluster %s: %s", clusterUuid, err.Error()))
				return
			}
		}
	}

	// Read back the updated resource
	ps, err := r.client.GetPrimaryStorage(uuid)
	if err != nil {
		resp.Diagnostics.AddError("Could not read updated primary storage", err.Error())
		return
	}

	refreshedState := primaryStorageModelFromView(ps, &plan)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *primaryStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state primaryStorageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "primary storage uuid is empty, skip delete")
		return
	}

	uuid := state.Uuid.ValueString()

	// First detach from all clusters
	attachedClusters := listToStringSlice(state.AttachedClusterUuids)
	for _, clusterUuid := range attachedClusters {
		if err := r.client.DetachPrimaryStorageFromCluster(clusterUuid, uuid, param.DeleteModePermissive); err != nil {
			resp.Diagnostics.AddError("Could not detach primary storage from cluster",
				fmt.Sprintf("Failed to detach from cluster %s: %s", clusterUuid, err.Error()))
			return
		}
	}

	// Then delete the primary storage
	if err := r.client.DeletePrimaryStorage(uuid, param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Could not delete primary storage", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *primaryStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}


func primaryStorageModelFromView(ps *view.PrimaryStorageInventoryView, prior *primaryStorageResourceModel) primaryStorageResourceModel {
	model := primaryStorageResourceModel{
		Uuid:                      types.StringValue(ps.UUID),
		Name:                      types.StringValue(ps.Name),
		Description:               stringValueOrNull(ps.Description),
		ZoneUuid:                  types.StringValue(ps.ZoneUuid),
		Type:                      types.StringValue(ps.Type),
		Url:                       stringValueOrNull(ps.Url),
		State:                     stringValueOrNull(ps.State),
		Status:                    stringValueOrNull(ps.Status),
		TotalCapacity:             types.Int64Value(ps.TotalCapacity),
		AvailableCapacity:         types.Int64Value(ps.AvailableCapacity),
		TotalPhysicalCapacity:     types.Int64Value(ps.TotalPhysicalCapacity),
		AvailablePhysicalCapacity: types.Int64Value(ps.AvailablePhysicalCapacity),
		MountPath:                 stringValueOrNull(ps.MountPath),
		AttachedClusterUuids:      stringSliceToList(ps.AttachedClusterUuids),
	}

	// Preserve plan-only fields that are not returned by the API
	if prior != nil {
		model.MonUrls = prior.MonUrls
		model.RootVolumePoolName = prior.RootVolumePoolName
		model.DataVolumePoolName = prior.DataVolumePoolName
		model.ImageCachePoolName = prior.ImageCachePoolName
		model.DiskUuids = prior.DiskUuids
	}

	return model
}
