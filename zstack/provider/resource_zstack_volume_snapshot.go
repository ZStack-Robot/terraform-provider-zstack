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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &volumeSnapshotResource{}
	_ resource.ResourceWithConfigure   = &volumeSnapshotResource{}
	_ resource.ResourceWithImportState = &volumeSnapshotResource{}
)

type volumeSnapshotResource struct {
	client *client.ZSClient
}

type volumeSnapshotResourceModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	VolumeUuid         types.String `tfsdk:"volume_uuid"`
	TreeUuid           types.String `tfsdk:"tree_uuid"`
	ParentUuid         types.String `tfsdk:"parent_uuid"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
	VolumeType         types.String `tfsdk:"volume_type"`
	Format             types.String `tfsdk:"format"`
	Latest             types.Bool   `tfsdk:"latest"`
	Size               types.Int64  `tfsdk:"size"`
	State              types.String `tfsdk:"state"`
	Status             types.String `tfsdk:"status"`
	Distance           types.Int64  `tfsdk:"distance"`
	GroupUuid          types.String `tfsdk:"group_uuid"`
	Revert             types.Bool   `tfsdk:"revert"`
}

func VolumeSnapshotResource() resource.Resource {
	return &volumeSnapshotResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *volumeSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *volumeSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_snapshot"
}

// Schema implements resource.Resource.
func (r *volumeSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage ZStack data volume snapshots.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the snapshot.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"volume_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the source volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tree_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the snapshot tree.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the parent snapshot, if any.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_storage_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the primary storage holding the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"volume_type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the source volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"format": schema.StringAttribute{
				Computed:    true,
				Description: "The snapshot format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"latest": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is the latest snapshot in its tree.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "The size of the snapshot in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The administrative state of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The operational status of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"distance": schema.Int64Attribute{
				Computed:    true,
				Description: "The distance from the root snapshot.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"group_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The snapshot group UUID, if any.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"revert": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Set to true to revert the source volume to this snapshot. The revert is triggered during update when this field changes from false to true. After the revert completes, the value resets to false.",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *volumeSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan volumeSnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := r.client.CreateVolumeSnapshot(plan.VolumeUuid.ValueString(), param.CreateVolumeSnapshotParam{
		Params: param.CreateVolumeSnapshotParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtr(plan.Description.ValueString()),
		},
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Volume Snapshot",
			"Could not create volume snapshot, unexpected error: "+err.Error(),
		)
		return
	}

	state := volumeSnapshotModelFromView(snapshot)
	state.Revert = types.BoolValue(false)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *volumeSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state volumeSnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := findResourceByGet(r.client.GetVolumeSnapshot, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Volume Snapshot",
			"Could not read volume snapshot UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state = volumeSnapshotModelFromView(snapshot)
	state.Revert = types.BoolValue(false)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *volumeSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan volumeSnapshotResourceModel
	var state volumeSnapshotResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() || plan.Description.ValueString() != state.Description.ValueString() {
		_, err := r.client.UpdateVolumeSnapshot(state.Uuid.ValueString(), param.UpdateVolumeSnapshotParam{
			Params: param.UpdateVolumeSnapshotParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtr(plan.Description.ValueString()),
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating Volume Snapshot",
				"Could not update volume snapshot, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// One-shot trigger: only revert on a false→true edge (not when state already holds true)
	if plan.Revert.ValueBool() && !state.Revert.ValueBool() {
		tflog.Info(ctx, "Reverting volume from snapshot", map[string]any{"snapshot_uuid": state.Uuid.ValueString()})
		if _, err := r.client.RevertVolumeFromSnapshot(state.Uuid.ValueString(), param.RevertVolumeFromSnapshotParam{}); err != nil {
			resp.Diagnostics.AddError(
				"Error reverting Volume Snapshot",
				"Could not revert volume from volume snapshot UUID "+state.Uuid.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	snapshot, err := r.client.GetVolumeSnapshot(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Volume Snapshot",
			"Could not read volume snapshot after update: "+err.Error(),
		)
		return
	}

	updatedState := volumeSnapshotModelFromView(snapshot)
	// Always reset revert to false so it doesn't re-trigger on subsequent updates
	updatedState.Revert = types.BoolValue(false)
	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *volumeSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state volumeSnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "volume snapshot uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteVolumeSnapshot(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Volume Snapshot",
			"Could not delete volume snapshot, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *volumeSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func volumeSnapshotModelFromView(snapshot *view.VolumeSnapshotInventoryView) volumeSnapshotResourceModel {
	return volumeSnapshotResourceModel{
		Uuid:               types.StringValue(snapshot.UUID),
		Name:               types.StringValue(snapshot.Name),
		Description:        stringValueOrNull(snapshot.Description),
		VolumeUuid:         stringValueOrNull(snapshot.VolumeUuid),
		TreeUuid:           stringValueOrNull(snapshot.TreeUuid),
		ParentUuid:         stringValueOrNull(snapshot.ParentUuid),
		PrimaryStorageUuid: stringValueOrNull(snapshot.PrimaryStorageUuid),
		VolumeType:         stringValueOrNull(snapshot.VolumeType),
		Format:             stringValueOrNull(snapshot.Format),
		Latest:             types.BoolValue(snapshot.Latest),
		Size:               types.Int64Value(snapshot.Size),
		State:              stringValueOrNull(snapshot.State),
		Status:             stringValueOrNull(snapshot.Status),
		Distance:           types.Int64Value(int64(snapshot.Distance)),
		GroupUuid:          stringValueOrNull(snapshot.GroupUuid),
		Revert:             types.BoolNull(),
	}
}
