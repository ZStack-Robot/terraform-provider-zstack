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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &imageStoreBackupStorageResource{}
	_ resource.ResourceWithConfigure   = &imageStoreBackupStorageResource{}
	_ resource.ResourceWithImportState = &imageStoreBackupStorageResource{}
)

type imageStoreBackupStorageResource struct {
	client *client.ZSClient
}

type imageStoreBackupStorageModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Hostname          types.String `tfsdk:"hostname"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	SshPort           types.Int64  `tfsdk:"ssh_port"`
	Url               types.String `tfsdk:"url"`
	Type              types.String `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	TotalCapacity     types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity types.Int64  `tfsdk:"available_capacity"`
}

func ImageStoreBackupStorageResource() resource.Resource {
	return &imageStoreBackupStorageResource{}
}

func (r *imageStoreBackupStorageResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T.", request.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *imageStoreBackupStorageResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_image_store_backup_storage"
}

func (r *imageStoreBackupStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages an Image Store Backup Storage in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the image store backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the image store backup storage",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the image store backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The hostname or IP address of the image store backup storage",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The username for SSH access to the image store backup storage",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for SSH access to the image store backup storage",
			},
			"ssh_port": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The SSH port of the image store backup storage",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the image store backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the backup storage",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"total_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The total capacity of the backup storage in bytes",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"available_capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The available capacity of the backup storage in bytes",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *imageStoreBackupStorageResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan imageStoreBackupStorageModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	createParam := param.AddImageStoreBackupStorageParam{}
	createParam.Params.Name = plan.Name.ValueString()
	createParam.Params.Url = plan.Url.ValueString()
	createParam.Params.Hostname = plan.Hostname.ValueString()
	createParam.Params.Username = plan.Username.ValueString()
	createParam.Params.Password = plan.Password.ValueString()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createParam.Params.Description = stringPtrOrNil(plan.Description.ValueString())
	}

	if !plan.SshPort.IsNull() && !plan.SshPort.IsUnknown() {
		sshPort := int(plan.SshPort.ValueInt64())
		createParam.Params.SshPort = &sshPort
	}

	result, err := r.client.AddImageStoreBackupStorage(createParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Image Store Backup Storage",
			"Could not create image store backup storage, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Hostname = types.StringValue(result.Hostname)
	plan.Username = types.StringValue(result.Username)
	plan.SshPort = types.Int64Value(int64(result.SshPort))
	plan.Url = types.StringValue(result.Url)
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)
	plan.Status = types.StringValue(result.Status)
	plan.TotalCapacity = types.Int64Value(result.TotalCapacity)
	plan.AvailableCapacity = types.Int64Value(result.AvailableCapacity)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "created an image store backup storage")

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *imageStoreBackupStorageResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state imageStoreBackupStorageModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	backupStorages, err := r.client.QueryImageStoreBackupStorage(&queryParam)
	if err != nil {
		if isZStackNotFoundError(err) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Image Store Backup Storage",
			"Could not read image store backup storage UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	if len(backupStorages) == 0 {
		tflog.Warn(ctx, "image store backup storage not found, removing from state", map[string]interface{}{
			"uuid": state.Uuid.ValueString(),
		})
		response.State.RemoveResource(ctx)
		return
	}

	backupStorage := backupStorages[0]

	state.Name = types.StringValue(backupStorage.Name)
	state.Description = stringValueOrNull(backupStorage.Description)
	state.Hostname = types.StringValue(backupStorage.Hostname)
	state.Username = types.StringValue(backupStorage.Username)
	state.SshPort = types.Int64Value(int64(backupStorage.SshPort))
	state.Url = types.StringValue(backupStorage.Url)
	state.Type = types.StringValue(backupStorage.Type)
	state.State = types.StringValue(backupStorage.State)
	state.Status = types.StringValue(backupStorage.Status)
	state.TotalCapacity = types.Int64Value(backupStorage.TotalCapacity)
	state.AvailableCapacity = types.Int64Value(backupStorage.AvailableCapacity)
	// Password is not returned by query, preserve from state

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *imageStoreBackupStorageResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan imageStoreBackupStorageModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state imageStoreBackupStorageModel
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateImageStoreBackupStorageParam{}
	updateParam.Params.Name = plan.Name.ValueString()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateParam.Params.Description = stringPtrOrNil(plan.Description.ValueString())
	}

	if !plan.Hostname.Equal(state.Hostname) {
		updateParam.Params.Hostname = stringPtr(plan.Hostname.ValueString())
	}

	if !plan.Username.Equal(state.Username) {
		updateParam.Params.Username = stringPtr(plan.Username.ValueString())
	}

	if !plan.Password.Equal(state.Password) {
		updateParam.Params.Password = stringPtr(plan.Password.ValueString())
	}

	if !plan.SshPort.Equal(state.SshPort) {
		sshPort := int(plan.SshPort.ValueInt64())
		updateParam.Params.SshPort = &sshPort
	}

	_, err := r.client.UpdateImageStoreBackupStorage(state.Uuid.ValueString(), updateParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Image Store Backup Storage",
			"Could not update image store backup storage, unexpected error: "+err.Error(),
		)
		return
	}

	// After update, do a fresh read to get full details
	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	backupStorages, err := r.client.QueryImageStoreBackupStorage(&queryParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading Image Store Backup Storage",
			"Could not read image store backup storage UUID "+state.Uuid.ValueString()+" after update: "+err.Error(),
		)
		return
	}

	if len(backupStorages) == 0 {
		response.Diagnostics.AddError(
			"Error reading Image Store Backup Storage",
			"Could not read image store backup storage after update: not found",
		)
		return
	}

	backupStorage := backupStorages[0]

	plan.Uuid = types.StringValue(backupStorage.UUID)
	plan.Name = types.StringValue(backupStorage.Name)
	plan.Description = stringValueOrNull(backupStorage.Description)
	plan.Hostname = types.StringValue(backupStorage.Hostname)
	plan.Username = types.StringValue(backupStorage.Username)
	plan.SshPort = types.Int64Value(int64(backupStorage.SshPort))
	plan.Url = types.StringValue(backupStorage.Url)
	plan.Type = types.StringValue(backupStorage.Type)
	plan.State = types.StringValue(backupStorage.State)
	plan.Status = types.StringValue(backupStorage.Status)
	plan.TotalCapacity = types.Int64Value(backupStorage.TotalCapacity)
	plan.AvailableCapacity = types.Int64Value(backupStorage.AvailableCapacity)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "updated an image store backup storage")

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *imageStoreBackupStorageResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state imageStoreBackupStorageModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBackupStorage(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Image Store Backup Storage",
			"Could not delete image store backup storage, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted an image store backup storage")
}

func (r *imageStoreBackupStorageResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
