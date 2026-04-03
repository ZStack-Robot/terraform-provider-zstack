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
	_ resource.Resource                = &volumeBackupResource{}
	_ resource.ResourceWithConfigure   = &volumeBackupResource{}
	_ resource.ResourceWithImportState = &volumeBackupResource{}
)

type volumeBackupResource struct {
	client *client.ZSClient
}

type volumeBackupModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	VolumeUuid        types.String `tfsdk:"volume_uuid"`
	BackupStorageUuid types.String `tfsdk:"backup_storage_uuid"`
	Type              types.String `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	Size              types.Int64  `tfsdk:"size"`
}

func VolumeBackupResource() resource.Resource {
	return &volumeBackupResource{}
}

func (r *volumeBackupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *volumeBackupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_volume_backup"
}

func (r *volumeBackupResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage volume backups in ZStack. " +
			"A volume backup is a point-in-time copy of a volume stored in backup storage.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the volume backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the volume backup.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the volume backup.",
			},
			"volume_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the volume to backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"backup_storage_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the backup storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the volume backup.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the volume backup.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the volume backup.",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "The size of the volume backup in bytes.",
			},
		},
	}
}

func (r *volumeBackupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan volumeBackupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVolumeBackupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVolumeBackupParamDetail{
			BackupStorageUuid: plan.BackupStorageUuid.ValueString(),
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.CreateVolumeBackup(plan.VolumeUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create volume backup",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.VolumeUuid = types.StringValue(result.VolumeUuid)
	plan.BackupStorageUuid = plan.BackupStorageUuid
	if len(result.BackupStorageRefs) > 0 {
		plan.BackupStorageUuid = types.StringValue(result.BackupStorageRefs[0].BackupStorageUuid)
	}
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)
	plan.Status = types.StringValue(result.Status)
	plan.Size = types.Int64Value(result.Size)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *volumeBackupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state volumeBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	volumeBackups, err := r.client.QueryVolumeBackup(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query volume backups. It may have been deleted.: "+err.Error())
		state = volumeBackupModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, volumeBackup := range volumeBackups {
		if volumeBackup.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(volumeBackup.UUID)
			state.Name = types.StringValue(volumeBackup.Name)
			state.Description = stringValueOrNull(volumeBackup.Description)
			state.VolumeUuid = types.StringValue(volumeBackup.VolumeUuid)
			state.BackupStorageUuid = state.BackupStorageUuid
			if len(volumeBackup.BackupStorageRefs) > 0 {
				state.BackupStorageUuid = types.StringValue(volumeBackup.BackupStorageRefs[0].BackupStorageUuid)
			}
			state.Type = types.StringValue(volumeBackup.Type)
			state.State = types.StringValue(volumeBackup.State)
			state.Status = types.StringValue(volumeBackup.Status)
			state.Size = types.Int64Value(volumeBackup.Size)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "Volume backup not found. It might have been deleted outside of Terraform.")
		state = volumeBackupModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *volumeBackupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *volumeBackupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state volumeBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Volume backup UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVolumeBackup(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete volume backup", ""+err.Error())
		return
	}
}

func (r *volumeBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
