// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &imageResource{}
	_ resource.ResourceWithConfigure = &imageResource{}
)

type imageResource struct {
	client *client.ZSClient
}

type imageResourceModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Url                types.String `tfsdk:"url"`
	MediaType          types.String `tfsdk:"media_type"`
	GuestOsType        types.String `tfsdk:"guest_os_type"`
	System             types.String `tfsdk:"system"`
	Platform           types.String `tfsdk:"platform"`
	Format             types.String `tfsdk:"format"`
	BackupStorageUuids types.List   `tfsdk:"backup_storage_uuids"`
	Architecture       types.String `tfsdk:"architecture"`
	Virtio             types.Bool   `tfsdk:"virtio"`
	Type               types.String `tfsdk:"type"`
	Marketplace        types.Bool   `tfsdk:"marketplace"`
	BootMode           types.String `tfsdk:"boot_mode"`
}

// Configure implements resource.ResourceWithConfigure.
func (r *imageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func ImageResource() resource.Resource {
	return &imageResource{}
}

// Create implements resource.Resource.
func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var imagePlan imageResourceModel
	diags := req.Plan.Get(ctx, &imagePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var backupStorageUuids []string
	if imagePlan.BackupStorageUuids.IsNull() {
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
		imagePlan.BackupStorageUuids.ElementsAs(ctx, &backupStorageUuids, false)
	}

	var systemTags []string

	if imagePlan.BootMode.IsNull() || imagePlan.BootMode.ValueString() == "" {
		// if boot mode not set, use uefi in aarch64 and legacy in x86_64
		if imagePlan.Architecture.ValueString() == "aarch64" {
			systemTags = append(systemTags, param.SystemTagBootModeUEFI)
		} else {
			systemTags = append(systemTags, param.SystemTagBootModeLegacy)
		}
	} else {
		bootMode := strings.ToLower(imagePlan.BootMode.ValueString())

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

	if imagePlan.Marketplace.ValueBool() {

		// for marketplace image, if the image with same name status is Ready and state is Enabled, then skip add it.
		qparam := param.NewQueryParam()
		qparam.AddQ("name=" + imagePlan.Name.ValueString())
		//qparam.AddQ("url=" + imagePlan.Url.ValueString())
		qparam.AddQ("architecture=" + imagePlan.Architecture.ValueString())
		qparam.AddQ("status=Ready")
		qparam.AddQ("state=Enabled")

		images, err := r.client.QueryImage(qparam)
		if err != nil {
			resp.Diagnostics.AddError(
				"fail to get image",
				fmt.Sprintf("fail to get image: %v", err),
			)
			return
		}

		tflog.Info(ctx, fmt.Sprintf("find %d images", len(images)))
		for _, image := range images {
			for _, backupStorageRef := range image.BackupStorageRefs {
				backupStorageUuids = removeStringFromSlice(backupStorageUuids, backupStorageRef.BackupStorageUuid)
			}

			if len(backupStorageUuids) == 0 {
				tflog.Info(ctx, "image has been imported to all backup storage")

				imagePlan.Uuid = types.StringValue(image.UUID)
				imagePlan.Name = types.StringValue(image.Name)
				imagePlan.Description = types.StringValue(image.Description)
				imagePlan.Url = types.StringValue(image.Url)
				imagePlan.GuestOsType = types.StringValue(image.GuestOsType)
				imagePlan.System = types.StringValue(image.System)
				imagePlan.Platform = types.StringValue(image.Platform)
				imagePlan.Type = types.StringValue(image.Type)
				imagePlan.LastUpdated = types.StringValue(image.LastOpDate.String())
				ctx = tflog.SetField(ctx, "url", image.Url)
				diags = resp.State.Set(ctx, imagePlan)
				resp.Diagnostics.Append(diags...)
				return
			}
		}

		tflog.Info(ctx, fmt.Sprintf("unimported backupStorageUuids: %v", backupStorageUuids))
		systemTags = append(systemTags, "marketplace::true")
	}

	tflog.Info(ctx, "Configuring ZStack client")
	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.AddImageDetailParam{
			Name:               imagePlan.Name.ValueString(),
			Description:        imagePlan.Description.ValueString(),
			Url:                imagePlan.Url.ValueString(),
			MediaType:          param.RootVolumeTemplate,
			GuestOsType:        imagePlan.GuestOsType.ValueString(),
			System:             false,
			Format:             param.ImageFormat(imagePlan.Format.ValueString()), // param.Qcow2,
			Platform:           imagePlan.Platform.ValueString(),
			BackupStorageUuids: backupStorageUuids,
			Type:               imagePlan.Type.ValueString(),
			ResourceUuid:       "",
			Architecture:       param.Architecture(imagePlan.Architecture.ValueString()),
			Virtio:             imagePlan.Virtio.ValueBool(),
		},
	}

	ctx = tflog.SetField(ctx, "url", imagePlan.Url)
	image, err := r.client.AddImage(imageParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Add image to ZStack Image storage"+image.Name, "Error "+err.Error(),
		)
		return
	}

	imagePlan.Uuid = types.StringValue(image.UUID)
	imagePlan.Name = types.StringValue(image.Name)
	imagePlan.Description = types.StringValue(image.Description)
	imagePlan.Url = types.StringValue(image.Url)
	imagePlan.GuestOsType = types.StringValue(image.GuestOsType)
	imagePlan.System = types.StringValue(image.System)
	imagePlan.Platform = types.StringValue(image.Platform)
	imagePlan.Type = types.StringValue(image.Type)
	imagePlan.LastUpdated = types.StringValue(image.LastOpDate.String())
	ctx = tflog.SetField(ctx, "url", image.Url)
	diags = resp.State.Set(ctx, imagePlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "image uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteImage(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete image", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Read implements resource.Resource.
func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	image, err := r.client.GetImage(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting ZStack Image uuid", "Could not read image uuid"+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(image.UUID)
	state.Name = types.StringValue(image.Name)
	state.Url = types.StringValue(image.Url)
	state.LastUpdated = types.StringValue(image.LastOpDate.GoString())
	//state.Description = types.StringValue(image.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *imageResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage images in ZStack. " +
			"An image represents a virtual machine image format qcow2, raw, vmdk or an ISO file that can be used to create or boot virtual machines. " +
			"You can define the image's properties, such as its URL, format, architecture, and backup storage locations.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the image. Automatically generated by ZStack.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of the last update to the image resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the image. This is a mandatory field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the image, providing additional context or details.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL where the image is located. This can be a file path or an HTTP link.",
			},
			"media_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of media for the image. Examples include 'ISO' or 'Template' or DataVolumeTemplate.",
			},
			"guest_os_type": schema.StringAttribute{
				Optional:    true,
				Description: "The guest operating system type that the image is optimized for.",
			},
			"system": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates if the image is a system image. Set automatically by ZStack.",
			},
			"platform": schema.StringAttribute{
				Optional:    true,
				Description: "The platform that the image is intended for, such as 'Linux', 'Windows', or others.",
			},
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The format of the image file, such as 'qcow2', 'raw', or 'vmdk'.",
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
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the image, for example, 'ISO' or 'RootVolumeTemplate'.",
			},
			"virtio": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the VirtIO drivers are required for the image.",
			},
			"marketplace": schema.BoolAttribute{
				Optional:    true,
				Description: "Specifies whether the image is from a marketplace.",
			},
			"boot_mode": schema.StringAttribute{
				Optional:    true,
				Description: "The boot mode supported by the image, such as 'Legacy' or 'UEFI'.",
			},
		},
	}
}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	/*
		var imageplan imageResourceModel
		diags := req.Plan.Get(ctx, &imageplan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		//uuid := imageplan.UUID.ValueString()
		//params := param.UpdateImageParam{
		//	UpdateImage: param.UpdateImageDetailParam{
		//		Name:        imageplan.Name.ValueString(), // "image4chenjtest",
		//		Description: &desc,                        //imageplan.Description.ValueStringPointer(),
		//	},
		//}
		//_, err := r.client.UpdateImage(uuid, params)
		//if err != nil {
		//	resp.Diagnostics.AddError(
		//		"", ""+err.Error(),
		//	)
		//	return
		//}
		//
		//image, err := r.client.GetImage(imageplan.UUID.ValueString())
		//if err != nil {
		//	resp.Diagnostics.AddError("", ""+err.Error())
		//	return
		//}
		//imageplan.UUID = types.StringValue(image.UUID)
		//imageplan.Name = types.StringValue(image.Name)
		//imageplan.Description = types.StringValue(image.Description)
		//imageplan.url = types.StringValue(image.Url)
		//imageplan.MediaType = types.StringValue(image.MediaType)
		//imageplan.GuestOsType = types.StringValue(image.GuestOsType)
		//imageplan.System = types.StringValue(image.System)
		////in ceph backup storage, the image format is always raw
		////imageplan.Format = types.StringValue(image.Format)
		//imageplan.Platform = types.StringValue(image.Platform)
		//imageplan.Type = types.StringValue(image.Type)
		//imageplan.Architecture = types.StringValue(string(image.Architecture))
		//imageplan.Virtio = types.BoolValue(image.Virtio)
		//imageplan.BackupStorageUuid = types.StringValue(string(image.BackupStorageRefs[0].BackupStorageUuid))
		//imageplan.LastUpdated = types.StringValue(image.LastOpDate.String())
		//
		//diags = resp.State.Set(ctx, imageplan)
		//resp.Diagnostics.Append(diags...)
		//if resp.Diagnostics.HasError() {
		//	return
		//}
	*/
}

func removeStringFromSlice(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}
