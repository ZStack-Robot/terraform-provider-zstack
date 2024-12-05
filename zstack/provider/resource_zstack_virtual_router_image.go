// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &virtualRouterImageResource{}
	_ resource.ResourceWithConfigure = &virtualRouterImageResource{}
)

type virtualRouterImageResource struct {
	client *client.ZSClient
}

type virtualRouterImageResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Url         types.String `tfsdk:"url"`
	MediaType   types.String `tfsdk:"media_type"`
	GuestOsType types.String `tfsdk:"guest_os_type"`
	//System             types.String `tfsdk:"system"`
	Platform types.String `tfsdk:"platform"`
	//Format             types.String `tfsdk:"format"`
	BackupStorageUuids types.List   `tfsdk:"backup_storage_uuids"`
	Architecture       types.String `tfsdk:"architecture"`
	Virtio             types.Bool   `tfsdk:"virtio"`
	//Type               types.String `tfsdk:"type"`
	BootMode types.String `tfsdk:"boot_mode"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *virtualRouterImageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

func VirtualRouterImageResource() resource.Resource {
	return &virtualRouterImageResource{}
}

// Create implements resource.Resource.
func (r *virtualRouterImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualRouterImageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupStorageUuids []string
	if plan.BackupStorageUuids.IsNull() {
		storage, err := r.client.QueryBackupStorage(param.QueryParam{})
		if err != nil {
			resp.Diagnostics.AddError(
				"fail to get backup storage",
				fmt.Sprintf("fail to get backup storage, err: %v", err),
			)
			return
		}
		backupStorageUuids = []string{storage[0].UUID}
	} else {
		plan.BackupStorageUuids.ElementsAs(ctx, &backupStorageUuids, false)
	}

	var systemTags []string

	if plan.BootMode.IsNull() || plan.BootMode.ValueString() == "" {
		// if boot mode not set, use uefi in aarch64 and legacy in x86_64
		if plan.Architecture.ValueString() == "aarch64" {
			systemTags = append(systemTags, param.SystemTagBootModeUEFI)
		} else {
			systemTags = append(systemTags, param.SystemTagBootModeLegacy)
		}
	} else {
		bootMode := strings.ToLower(plan.BootMode.ValueString())

		switch bootMode {
		case "uefi":
			systemTags = append(systemTags, param.SystemTagBootModeUEFI)
		case "legacy":
			systemTags = append(systemTags, param.SystemTagBootModeLegacy)
		default:
			resp.Diagnostics.AddError(
				"invalid boot mode",
				fmt.Sprintf("invalid boot mode: %s", bootMode),
			)
			return
		}
	}

	systemTags = append(systemTags, param.SystemTagApplianceTypeVRouter)

	//systemTags = append(systemTags, param.vr)
	tflog.Info(ctx, "Configuring ZStack client")
	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.AddImageDetailParam{
			Name:               plan.Name.ValueString(),
			Description:        plan.Description.ValueString(),
			Url:                plan.Url.ValueString(),
			MediaType:          param.RootVolumeTemplate,
			GuestOsType:        plan.GuestOsType.ValueString(),
			System:             true,
			Format:             param.Qcow2,                 //param.ImageFormat(plan.Format.ValueString()), // param.Qcow2,
			Platform:           plan.Platform.ValueString(), //plan.Platform.ValueString(),
			BackupStorageUuids: backupStorageUuids,
			//Type:               string(param.ApplianceVm), //"ApplianceVm",
			ResourceUuid: "",
			Architecture: param.Architecture(plan.Architecture.ValueString()),
			Virtio:       plan.Virtio.ValueBool(),
		},
	}

	ctx = tflog.SetField(ctx, "url", plan.Url)
	image, err := r.client.AddImage(imageParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add image to ZStack Image storage"+image.Name, "Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(image.UUID)
	plan.Name = types.StringValue(image.Name)
	plan.Description = types.StringValue(image.Description)
	plan.Url = types.StringValue(image.Url)
	//	plan.GuestOsType = types.StringValue(image.GuestOsType)
	//plan.System = types.StringValue(image.System)
	plan.Platform = types.StringValue(image.Platform)
	//plan.Type = types.StringValue(image.Type)
	//plan.LastUpdated = types.StringValue(image.LastOpDate.String())
	ctx = tflog.SetField(ctx, "url", image.Url)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *virtualRouterImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualRouterImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "virtual router image uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteImage(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete virtual router image", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *virtualRouterImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_router_image"
}

// Read implements resource.Resource.
func (r *virtualRouterImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualRouterImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	image, err := r.client.GetImage(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting virtual router image uuid", "Could not read image uuid "+err.Error(),
		)
		return
	}

	state.Platform = types.StringValue(image.Platform)
	if image.Platform == "" {
		state.Platform = types.StringValue("Linux")
	} else {
		state.Platform = types.StringValue(image.Platform)
	}

	state.Uuid = types.StringValue(image.UUID)
	state.Name = types.StringValue(image.Name)
	state.Url = types.StringValue(image.Url)

	//state.LastUpdated = types.StringValue(image.LastOpDate.GoString())
	state.Description = types.StringValue(image.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *virtualRouterImageResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the virtual router image. Automatically generated by ZStack.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual router image. This is a mandatory field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the virtual router image, providing additional context or details.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL where the virtual router image is located. This can be a file path or an HTTP link.",
			},
			"media_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of media for the image. Examples include 'ISO' or 'Template'.",
			},
			"guest_os_type": schema.StringAttribute{
				Optional:    true,
				Description: "The virtual router OS type that the image is optimized for. The Value must be one of: ['VyOS 1.1.7' 'openEuler 22.03']",

				Validators: []validator.String{
					stringvalidator.OneOf("VyOS 1.1.7", "openEuler 22.03"),
				},
			},
			"platform": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The platform type of the image. Defaults to Linux.",
			},
			"backup_storage_uuids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "A list of UUIDs for the backup storages where the image is stored.",
			},
			"architecture": schema.StringAttribute{
				Optional:    true,
				Description: "The architecture of the image, such as 'x86_64' or 'arm64'.",
			},
			"virtio": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the VirtIO drivers are required for the image.",
			},
			"boot_mode": schema.StringAttribute{
				Optional:    true,
				Description: "The boot mode supported by the image, such as 'Legacy' or 'UEFI'.",
			},
		},
	}
}

func (r *virtualRouterImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
