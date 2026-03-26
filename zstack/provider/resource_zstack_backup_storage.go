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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &backupStorageResource{}
	_ resource.ResourceWithConfigure   = &backupStorageResource{}
	_ resource.ResourceWithImportState = &backupStorageResource{}
)

type backupStorageResource struct {
	client *client.ZSClient
}

type backupStorageResourceModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Type              types.String `tfsdk:"type"`
	Url               types.String `tfsdk:"url"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	TotalCapacity     types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity types.Int64  `tfsdk:"available_capacity"`
	AttachedZoneUuids types.List   `tfsdk:"attached_zone_uuids"`
	Hostname          types.String `tfsdk:"hostname"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	SshPort           types.Int64  `tfsdk:"ssh_port"`
	MonUrls           types.List   `tfsdk:"mon_urls"`
	PoolName          types.String `tfsdk:"pool_name"`
}

func BackupStorageResource() resource.Resource {
	return &backupStorageResource{}
}

func (r *backupStorageResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *backupStorageResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_backup_storage"
}

func (r *backupStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage backup storage in ZStack. " +
			"Supported types include ImageStoreBackupStorage, CephBackupStorage, and SftpBackupStorage.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the backup storage.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the backup storage.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the backup storage.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of backup storage: ImageStoreBackupStorage, CephBackupStorage, or SftpBackupStorage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The URL/path for the backup storage (required for ImageStore and Sftp types).",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the backup storage.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the backup storage.",
			},
			"total_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The total capacity of the backup storage in bytes.",
			},
			"available_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The available capacity of the backup storage in bytes.",
			},
			"attached_zone_uuids": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of zone UUIDs to attach this backup storage to.",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Description: "The hostname or IP address of the backup storage server (required for ImageStore and Sftp types).",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The SSH username for the backup storage server (required for ImageStore and Sftp types).",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The SSH password for the backup storage server (required for ImageStore and Sftp types).",
			},
			"ssh_port": schema.Int64Attribute{
				Optional:    true,
				Description: "The SSH port for the backup storage server (for ImageStore and Sftp types).",
			},
			"mon_urls": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of Ceph monitor URLs (required for Ceph type).",
			},
			"pool_name": schema.StringAttribute{
				Optional:    true,
				Description: "The Ceph pool name (for Ceph type).",
			},
		},
	}
}

func (r *backupStorageResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan backupStorageResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	var bsUuid string

	switch plan.Type.ValueString() {
	case "ImageStoreBackupStorage":
		p := param.AddImageStoreBackupStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddImageStoreBackupStorageParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
				Hostname:    plan.Hostname.ValueString(),
				Username:    plan.Username.ValueString(),
				Password:    plan.Password.ValueString(),
				Url:         plan.Url.ValueString(),
			},
		}
		if !plan.SshPort.IsNull() && !plan.SshPort.IsUnknown() {
			port := int(plan.SshPort.ValueInt64())
			p.Params.SshPort = &port
		}

		result, err := r.client.AddImageStoreBackupStorage(p)
		if err != nil {
			response.Diagnostics.AddError("Failed to create ImageStore backup storage", err.Error())
			return
		}
		bsUuid = result.UUID

	case "CephBackupStorage":
		monUrls := listToStringSlice(plan.MonUrls)
		p := param.AddCephBackupStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddCephBackupStorageParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
				MonUrls:     monUrls,
			},
		}
		if !plan.PoolName.IsNull() && !plan.PoolName.IsUnknown() {
			p.Params.PoolName = stringPtr(plan.PoolName.ValueString())
		}
		if !plan.Url.IsNull() && !plan.Url.IsUnknown() {
			p.Params.Url = stringPtr(plan.Url.ValueString())
		}

		result, err := r.client.AddCephBackupStorage(p)
		if err != nil {
			response.Diagnostics.AddError("Failed to create Ceph backup storage", err.Error())
			return
		}
		bsUuid = result.UUID

	case "SftpBackupStorage":
		p := param.AddSftpBackupStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.AddSftpBackupStorageParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
				Hostname:    plan.Hostname.ValueString(),
				Username:    plan.Username.ValueString(),
				Password:    plan.Password.ValueString(),
				Url:         plan.Url.ValueString(),
			},
		}
		if !plan.SshPort.IsNull() && !plan.SshPort.IsUnknown() {
			port := int(plan.SshPort.ValueInt64())
			p.Params.SshPort = &port
		}

		result, err := r.client.AddSftpBackupStorage(p)
		if err != nil {
			response.Diagnostics.AddError("Failed to create Sftp backup storage", err.Error())
			return
		}
		bsUuid = result.UUID

	default:
		response.Diagnostics.AddError("Unsupported backup storage type", fmt.Sprintf("Type %q is not supported. Use ImageStoreBackupStorage, CephBackupStorage, or SftpBackupStorage.", plan.Type.ValueString()))
		return
	}

	// Handle zone attachments
	zoneUuids := listToStringSlice(plan.AttachedZoneUuids)
	for _, zoneUuid := range zoneUuids {
		attachParam := param.AttachBackupStorageToZoneParam{
			BaseParam: param.BaseParam{},
		}
		var result view.BackupStorageInventoryView
		// Use direct Post to work around URL template bug in AttachBackupStorageToZone
		if err := r.client.Post(fmt.Sprintf("v1/zones/%s/backup-storage/%s", zoneUuid, bsUuid), attachParam, &result); err != nil {
			response.Diagnostics.AddError("Failed to attach backup storage to zone",
				fmt.Sprintf("Error attaching backup storage %s to zone %s: %s", bsUuid, zoneUuid, err.Error()))
			return
		}
	}

	// Read back the created resource to get full state
	bs, err := r.client.GetBackupStorage(bsUuid)
	if err != nil {
		response.Diagnostics.AddError("Failed to read backup storage after creation", err.Error())
		return
	}

	model := backupStorageModelFromView(bs, plan)

	diags = response.State.Set(ctx, model)
	response.Diagnostics.Append(diags...)
}

func (r *backupStorageResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state backupStorageResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "Backup storage UUID is empty, skipping read.")
		return
	}

	bs, err := r.client.GetBackupStorage(state.Uuid.ValueString())
	if err != nil {
		tflog.Warn(ctx, "Failed to read backup storage, it may have been deleted: "+err.Error())
		response.State.RemoveResource(ctx)
		return
	}

	model := backupStorageModelFromView(bs, state)

	diags = response.State.Set(ctx, model)
	response.Diagnostics.Append(diags...)
}

func (r *backupStorageResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan backupStorageResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state backupStorageResourceModel
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	// Update name/description if changed
	if !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) {
		updateParam := param.UpdateBackupStorageParam{
			BaseParam: param.BaseParam{},
			Params: param.UpdateBackupStorageParamDetail{
				Name:        plan.Name.ValueString(),
				Description: stringPtrOrNil(plan.Description.ValueString()),
			},
		}

		_, err := r.client.UpdateBackupStorage(uuid, updateParam)
		if err != nil {
			response.Diagnostics.AddError("Failed to update backup storage", err.Error())
			return
		}
	}

	// Handle zone attach/detach
	currentZones := listToStringSlice(state.AttachedZoneUuids)
	desiredZones := listToStringSlice(plan.AttachedZoneUuids)

	currentZoneMap := make(map[string]bool)
	for _, z := range currentZones {
		currentZoneMap[z] = true
	}
	desiredZoneMap := make(map[string]bool)
	for _, z := range desiredZones {
		desiredZoneMap[z] = true
	}

	// Attach new zones
	for _, zoneUuid := range desiredZones {
		if !currentZoneMap[zoneUuid] {
			attachParam := param.AttachBackupStorageToZoneParam{
				BaseParam: param.BaseParam{},
			}
			var result view.BackupStorageInventoryView
			// Use direct Post to work around URL template bug
			if err := r.client.Post(fmt.Sprintf("v1/zones/%s/backup-storage/%s", zoneUuid, uuid), attachParam, &result); err != nil {
				response.Diagnostics.AddError("Failed to attach backup storage to zone",
					fmt.Sprintf("Error attaching backup storage %s to zone %s: %s", uuid, zoneUuid, err.Error()))
				return
			}
		}
	}

	// Detach removed zones
	for _, zoneUuid := range currentZones {
		if !desiredZoneMap[zoneUuid] {
			if err := r.client.DetachBackupStorageFromZone(zoneUuid, uuid, param.DeleteModePermissive); err != nil {
				response.Diagnostics.AddError("Failed to detach backup storage from zone",
					fmt.Sprintf("Error detaching backup storage %s from zone %s: %s", uuid, zoneUuid, err.Error()))
				return
			}
		}
	}

	// Read back the updated resource
	bs, err := r.client.GetBackupStorage(uuid)
	if err != nil {
		response.Diagnostics.AddError("Failed to read backup storage after update", err.Error())
		return
	}

	model := backupStorageModelFromView(bs, plan)

	diags = response.State.Set(ctx, model)
	response.Diagnostics.Append(diags...)
}

func (r *backupStorageResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state backupStorageResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()
	if uuid == "" {
		tflog.Warn(ctx, "Backup storage UUID is empty, skipping delete.")
		return
	}

	// Detach from all zones first
	zoneUuids := listToStringSlice(state.AttachedZoneUuids)
	for _, zoneUuid := range zoneUuids {
		if err := r.client.DetachBackupStorageFromZone(zoneUuid, uuid, param.DeleteModePermissive); err != nil {
			response.Diagnostics.AddError("Failed to detach backup storage from zone",
				fmt.Sprintf("Error detaching backup storage %s from zone %s: %s", uuid, zoneUuid, err.Error()))
			return
		}
	}

	// Delete the backup storage
	if err := r.client.DeleteBackupStorage(uuid, param.DeleteModePermissive); err != nil {
		response.Diagnostics.AddError("Failed to delete backup storage", err.Error())
		return
	}
}

func (r *backupStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}


func backupStorageModelFromView(bs *view.BackupStorageInventoryView, plan backupStorageResourceModel) backupStorageResourceModel {
	model := backupStorageResourceModel{
		Uuid:              types.StringValue(bs.UUID),
		Name:              types.StringValue(bs.Name),
		Description:       stringValueOrNull(bs.Description),
		Type:              types.StringValue(bs.Type),
		Url:               stringValueOrNull(bs.Url),
		State:             stringValueOrNull(bs.State),
		Status:            stringValueOrNull(bs.Status),
		TotalCapacity:     types.Int64Value(bs.TotalCapacity),
		AvailableCapacity: types.Int64Value(bs.AvailableCapacity),
		AttachedZoneUuids: stringSliceToList(bs.AttachedZoneUuids),
		// Preserve sensitive fields from plan
		Hostname: plan.Hostname,
		Username: plan.Username,
		Password: plan.Password,
		SshPort:  plan.SshPort,
		MonUrls:  plan.MonUrls,
		PoolName: plan.PoolName,
	}
	return model
}
