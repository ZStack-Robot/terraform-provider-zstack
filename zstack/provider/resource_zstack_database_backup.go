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
	_ resource.Resource                = &databaseBackupResource{}
	_ resource.ResourceWithConfigure   = &databaseBackupResource{}
	_ resource.ResourceWithImportState = &databaseBackupResource{}
)

type databaseBackupResource struct {
	client *client.ZSClient
}

type databaseBackupModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	BackupStorageUuid types.String `tfsdk:"backup_storage_uuid"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	Size              types.Int64  `tfsdk:"size"`
	Metadata          types.String `tfsdk:"metadata"`
}

func DatabaseBackupResource() resource.Resource {
	return &databaseBackupResource{}
}

func (r *databaseBackupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *databaseBackupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_database_backup"
}

func (r *databaseBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage database backups in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the database backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the database backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the database backup.",
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
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the database backup.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the database backup.",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "The size of the database backup in bytes.",
			},
			"metadata": schema.StringAttribute{
				Computed:    true,
				Description: "The metadata of the database backup.",
			},
		},
	}
}

func (r *databaseBackupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan databaseBackupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateDatabaseBackupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateDatabaseBackupParamDetail{
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
			BackupStorageUuid: plan.BackupStorageUuid.ValueString(),
		},
	}

	result, err := r.client.CreateDatabaseBackup(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create database backup",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = stringValueOrNull(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.BackupStorageUuid = plan.BackupStorageUuid
	plan.State = stringValueOrNull(result.State)
	plan.Status = stringValueOrNull(result.Status)
	plan.Size = types.Int64Value(result.Size)
	plan.Metadata = stringValueOrNull(result.Metadata)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *databaseBackupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state databaseBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	databaseBackups, err := r.client.QueryDatabaseBackup(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query database backups. It may have been deleted.: "+err.Error())
		state = databaseBackupModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, databaseBackup := range databaseBackups {
		if databaseBackup.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(databaseBackup.UUID)
			state.Name = stringValueOrNull(databaseBackup.Name)
			state.Description = stringValueOrNull(databaseBackup.Description)
			state.BackupStorageUuid = state.BackupStorageUuid
			state.State = stringValueOrNull(databaseBackup.State)
			state.Status = stringValueOrNull(databaseBackup.Status)
			state.Size = types.Int64Value(databaseBackup.Size)
			state.Metadata = stringValueOrNull(databaseBackup.Metadata)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Database backup not found. It might have been deleted outside of Terraform.")
		state = databaseBackupModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *databaseBackupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r *databaseBackupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state databaseBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Database backup UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteDatabaseBackup(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete database backup", err.Error())
		return
	}
}

func (r *databaseBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
