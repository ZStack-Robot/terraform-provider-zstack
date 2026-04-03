// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &zboxBackupResource{}
	_ resource.ResourceWithConfigure   = &zboxBackupResource{}
	_ resource.ResourceWithImportState = &zboxBackupResource{}
)

type zboxBackupResource struct {
	client *client.ZSClient
}

type zboxBackupModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ZBoxUuid    types.String `tfsdk:"zbox_uuid"`
	DryRun      types.Bool   `tfsdk:"dry_run"`
	State       types.String `tfsdk:"state"`
	InstallPath types.String `tfsdk:"install_path"`
	TotalSize   types.Int64  `tfsdk:"total_size"`
	Version     types.String `tfsdk:"version"`
	Type        types.String `tfsdk:"type"`
}

func ZBoxBackupResource() resource.Resource {
	return &zboxBackupResource{}
}

func (r *zboxBackupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *zboxBackupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_zbox_backup"
}

func (r *zboxBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage ZBox backups in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the ZBox backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ZBox backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the ZBox backup.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zbox_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the ZBox.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dry_run": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to run in dry-run mode.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the ZBox backup.",
			},
			"install_path": schema.StringAttribute{
				Computed:    true,
				Description: "The install path of the ZBox backup.",
			},
			"total_size": schema.Int64Attribute{
				Computed:    true,
				Description: "The total size of the ZBox backup in bytes.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the ZBox backup.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the ZBox backup.",
			},
		},
	}
}

func (r *zboxBackupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan zboxBackupModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateZBoxBackupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateZBoxBackupParamDetail{
			ZBoxUuid:    plan.ZBoxUuid.ValueString(),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if !plan.DryRun.IsNull() && !plan.DryRun.IsUnknown() {
		p.Params.DryRun = boolPtr(plan.DryRun.ValueBool())
	}

	result, err := r.client.CreateZBoxBackup(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create zbox backup",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = stringValueOrNull(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.State = stringValueOrNull(result.State)
	plan.InstallPath = stringValueOrNull(result.InstallPath)
	plan.TotalSize = types.Int64Value(result.TotalSize)
	plan.Version = stringValueOrNull(result.Version)
	plan.Type = stringValueOrNull(result.Type)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *zboxBackupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state zboxBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	zboxBackups, err := r.client.QueryZBoxBackup(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query zbox backups. It may have been deleted.: "+err.Error())
		state = zboxBackupModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, zboxBackup := range zboxBackups {
		if zboxBackup.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(zboxBackup.UUID)
			state.Name = stringValueOrNull(zboxBackup.Name)
			state.ZBoxUuid = stringValueOrNull(zboxBackup.ZBoxUuid)
			state.Description = stringValueOrNull(zboxBackup.Description)
			state.State = stringValueOrNull(zboxBackup.State)
			state.InstallPath = stringValueOrNull(zboxBackup.InstallPath)
			state.TotalSize = types.Int64Value(zboxBackup.TotalSize)
			state.Version = stringValueOrNull(zboxBackup.Version)
			state.Type = stringValueOrNull(zboxBackup.Type)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "ZBox backup not found. It might have been deleted outside of Terraform.")
		state = zboxBackupModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *zboxBackupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r *zboxBackupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state zboxBackupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "ZBox backup UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteExternalBackup(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete zbox backup", err.Error())
		return
	}
}

func (r *zboxBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
