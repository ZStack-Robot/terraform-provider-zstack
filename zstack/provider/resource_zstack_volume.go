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
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

type volumeResource struct {
	client *client.ZSClient
}

type volumeResourceModel struct {
	Uuid               types.String   `tfsdk:"uuid"`
	Name               types.String   `tfsdk:"name"`
	Description        types.String   `tfsdk:"description"`
	DiskOfferingUuid   types.String   `tfsdk:"disk_offering_uuid"`
	DiskSize           types.Int64    `tfsdk:"disk_size"`
	PrimaryStorageUuid types.String   `tfsdk:"primary_storage_uuid"`
	ResourceUuid       types.String   `tfsdk:"resource_uuid"`
	TagUuids           types.List     `tfsdk:"tag_uuids"`
	VmInstanceUuid     types.String   `tfsdk:"vm_instance_uuid"`
	Type               types.String   `tfsdk:"type"`
	Format             types.String   `tfsdk:"format"`
	State              types.String   `tfsdk:"state"`
	Status             types.String   `tfsdk:"status"`
	ActualSize         types.Int64    `tfsdk:"actual_size"`
	IsShareable        types.Bool     `tfsdk:"is_shareable"`
}

func VolumeResource() resource.Resource {
	return &volumeResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *volumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema implements resource.Resource.
func (r *volumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage ZStack data volumes and optionally attach them to virtual machines.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the volume.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the volume.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the volume.",
			},
			"disk_offering_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the disk offering used to create the volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"disk_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The size of the volume in bytes.",
			},
			"primary_storage_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the primary storage where the volume is created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The custom UUID requested at creation time.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tag_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The tag UUIDs attached during creation.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"vm_instance_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the VM instance that the volume is attached to.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The volume type reported by ZStack.",
			},
			"format": schema.StringAttribute{
				Computed:    true,
				Description: "The volume format reported by ZStack.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The administrative state of the volume.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The operational status of the volume.",
			},
			"actual_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The actual size of the volume in bytes.",
			},
			"is_shareable": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the volume is shareable.",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *volumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.DiskOfferingUuid.IsNull() && (plan.DiskSize.IsNull() || plan.DiskSize.ValueInt64() <= 0) {
		resp.Diagnostics.AddError(
			"Missing Volume Size Configuration",
			"Set either disk_offering_uuid or a positive disk_size when creating a volume.",
		)
		return
	}

	createParam := param.CreateDataVolumeParam{
		Params: param.CreateDataVolumeParamDetail{
			Name:             plan.Name.ValueString(),
			Description:      stringPtr(plan.Description.ValueString()),
			DiskOfferingUuid: stringPtr(plan.DiskOfferingUuid.ValueString()),
			TagUuids:         listToStringSlice(plan.TagUuids),
		},
	}

	if !plan.DiskSize.IsNull() && plan.DiskSize.ValueInt64() > 0 {
		createParam.Params.DiskSize = int64Ptr(plan.DiskSize.ValueInt64())
	}
	if !plan.PrimaryStorageUuid.IsNull() && plan.PrimaryStorageUuid.ValueString() != "" {
		createParam.Params.PrimaryStorageUuid = stringPtr(plan.PrimaryStorageUuid.ValueString())
	}
	if !plan.ResourceUuid.IsNull() && plan.ResourceUuid.ValueString() != "" {
		createParam.Params.ResourceUuid = stringPtr(plan.ResourceUuid.ValueString())
	}

	volume, err := r.client.CreateDataVolume(createParam)
	if err != nil {
		resp.Diagnostics.AddError("Could not create volume", err.Error())
		return
	}

	if !plan.VmInstanceUuid.IsNull() && plan.VmInstanceUuid.ValueString() != "" {
		if _, err := r.client.AttachDataVolumeToVm(param.AttachDataVolumeToVmParam{
			BaseParam: param.BaseParam{},
		}); err != nil {
			resp.Diagnostics.AddError("Volume created but attach failed", err.Error())
			return
		}
	}

	state, err := r.readVolume(volume.UUID, plan)
	if err != nil {
		resp.Diagnostics.AddError("Could not read created volume", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *volumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := r.readVolume(state.Uuid.ValueString(), state)
	if err != nil {
		resp.Diagnostics.AddError("Could not read volume", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan volumeResourceModel
	var state volumeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		updateParam := param.UpdateVolumeParam{
			Params: param.UpdateVolumeParamDetail{
				Name:        plan.Name.ValueString(),
				Description: plan.Description.ValueStringPointer(),
			},
		}

		if _, err := r.client.UpdateVolume(state.Uuid.ValueString(), updateParam); err != nil {
			resp.Diagnostics.AddError("Could not update volume metadata", err.Error())
			return
		}
	}

	if !plan.DiskSize.IsNull() && plan.DiskSize.ValueInt64() != state.DiskSize.ValueInt64() {
		if plan.DiskSize.ValueInt64() < state.DiskSize.ValueInt64() {
			resp.Diagnostics.AddError(
				"Shrinking a volume is not supported",
				"disk_size can only stay the same or increase for an existing volume.",
			)
			return
		}

		if _, err := r.client.ResizeDataVolume(state.Uuid.ValueString(), param.ResizeDataVolumeParam{
			Params: param.ResizeDataVolumeParamDetail{
				Size: plan.DiskSize.ValueInt64(),
			},
		}); err != nil {
			resp.Diagnostics.AddError("Could not resize volume", err.Error())
			return
		}
	}

	if err := r.reconcileAttachment(state, plan); err != nil {
		resp.Diagnostics.AddError("Could not update volume attachment", err.Error())
		return
	}

	refreshedState, err := r.readVolume(state.Uuid.ValueString(), plan)
	if err != nil {
		resp.Diagnostics.AddError("Could not read updated volume", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *volumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "volume uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteDataVolume(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Could not delete volume", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *volumeResource) readVolume(uuid string, prior volumeResourceModel) (volumeResourceModel, error) {
	volume, err := r.client.GetVolume(uuid)
	if err != nil {
		return volumeResourceModel{}, err
	}

	return volumeModelFromView(volume, prior), nil
}

func (r *volumeResource) reconcileAttachment(state, plan volumeResourceModel) error {
	currentVm := state.VmInstanceUuid.ValueString()
	desiredVm := plan.VmInstanceUuid.ValueString()

	if currentVm == desiredVm {
		return nil
	}

	if currentVm != "" {
		if err := r.client.DetachDataVolumeFromVm(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
			return err
		}
	}

	if desiredVm != "" {
		if _, err := r.client.AttachDataVolumeToVm(param.AttachDataVolumeToVmParam{
			BaseParam: param.BaseParam{},
		}); err != nil {
			return err
		}
	}

	return nil
}

func volumeModelFromView(volume *view.VolumeInventoryView, prior volumeResourceModel) volumeResourceModel {
	// Ensure resource_uuid and tag_uuids are never unknown after apply.
	resourceUuid := prior.ResourceUuid
	if resourceUuid.IsUnknown() {
		resourceUuid = types.StringNull()
	}
	tagUuids := prior.TagUuids
	if tagUuids.IsUnknown() {
		tagUuids = types.ListNull(types.StringType)
	}

	return volumeResourceModel{
		Uuid:               types.StringValue(volume.UUID),
		Name:               types.StringValue(volume.Name),
		Description:        stringValueOrNull(volume.Description),
		DiskOfferingUuid:   stringValueOrNull(volume.DiskOfferingUuid),
		DiskSize:           types.Int64Value(int64(volume.Size)),
		PrimaryStorageUuid: stringValueOrNull(volume.PrimaryStorageUuid),
		ResourceUuid:       resourceUuid,
		TagUuids:           tagUuids,
		VmInstanceUuid:     stringValueOrNull(volume.VmInstanceUuid),
		Type:               stringValueOrNull(volume.Type),
		Format:             stringValueOrNull(volume.Format),
		State:              stringValueOrNull(volume.State),
		Status:             stringValueOrNull(volume.Status),
		ActualSize:         types.Int64Value(int64(volume.ActualSize)),
		IsShareable:        types.BoolValue(volume.IsShareable),
	}
}
