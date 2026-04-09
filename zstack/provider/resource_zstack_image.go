// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &imageResource{}
	_ resource.ResourceWithConfigure   = &imageResource{}
	_ resource.ResourceWithImportState = &imageResource{}
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
	Virtio             types.String `tfsdk:"virtio"`
	//	Type               types.String `tfsdk:"type"`
	//Marketplace        types.Bool   `tfsdk:"marketplace"`
	BootMode types.String `tfsdk:"boot_mode"`
	Expunge  types.Bool   `tfsdk:"expunge"`
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
		storage, err := r.client.QueryBackupStorage(&param.QueryParam{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating Image",
				fmt.Sprintf("Could not create image because backup storage could not be read: %v", err),
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
			systemTags = append(systemTags, "bootMode::UEFI")
		} else {
			systemTags = append(systemTags, "bootMode::Legacy")
		}
	} else {
		bootMode := strings.ToLower(imagePlan.BootMode.ValueString())

		switch bootMode {
		case "uefi":
			systemTags = append(systemTags, "bootMode::UEFI")
		case "legacy":
			systemTags = append(systemTags, "bootMode::Legacy")
		default:
			resp.Diagnostics.AddError(
				"Error creating Image",
				fmt.Sprintf("Could not create image, invalid boot mode: %s", bootMode),
			)
			return
		}
	}

	if imagePlan.Description.IsNull() {
		imagePlan.Description = types.StringValue("")
	}
	if imagePlan.GuestOsType.IsNull() {
		imagePlan.GuestOsType = types.StringValue("Linux")
	}
	if imagePlan.Platform.IsNull() {
		imagePlan.Platform = types.StringValue("Linux")
	}

	tflog.Info(ctx, "Configuring ZStack client")
	imageParam := param.AddImageParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
		},
		Params: param.AddImageParamDetail{
			Name:               imagePlan.Name.ValueString(),
			Description:        stringPtr(imagePlan.Description.ValueString()),
			Url:                imagePlan.Url.ValueString(),
			MediaType:          stringPtr(imagePlan.MediaType.ValueString()),
			GuestOsType:        stringPtr(imagePlan.GuestOsType.ValueString()),
			System:             false,
			Format:             stringPtr(imagePlan.Format.ValueString()),
			Platform:           stringPtr(imagePlan.Platform.ValueString()),
			BackupStorageUuids: backupStorageUuids,
			//Type:               imagePlan.Type.ValueString(),
			//ResourceUuid: stringPtr(""),
			Architecture: stringPtr(imagePlan.Architecture.ValueString()),
			Virtio:       boolPtr(imagePlan.Virtio.ValueString() == "true"),
		},
	}

	ctx = tflog.SetField(ctx, "url", imagePlan.Url)
	image, err := r.client.AddImage(imageParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Image", "Could not create image, unexpected error: "+err.Error(),
		)
		return
	}

	imagePlan.Uuid = types.StringValue(image.UUID)
	imagePlan.Name = types.StringValue(image.Name)
	imagePlan.Description = types.StringValue(image.Description)
	imagePlan.Url = types.StringValue(image.Url)
	imagePlan.GuestOsType = types.StringValue(image.GuestOsType)
	imagePlan.System = types.StringValue(fmt.Sprintf("%t", image.System))
	imagePlan.Platform = types.StringValue(image.Platform)
	//imagePlan.Type = types.StringValue(image.Type)
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

	expunge := false
	if !state.Expunge.IsNull() && !state.Expunge.IsUnknown() {
		expunge = state.Expunge.ValueBool()
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "image uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteImage(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("Error deleting Image", "Could not delete image, unexpected error: "+err.Error())
		return
	}

	if expunge {
		tflog.Info(ctx, fmt.Sprintf("expunge image %s", state.Uuid.ValueString()))

		err = r.client.ExpungeImage(state.Uuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error expunging Image", "Could not expunge image, unexpected error: "+err.Error(),
			)
			return
		}
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
	image, err := findResourceByGet(r.client.GetImage, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Image", "Could not read image UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(image.UUID)
	state.Name = types.StringValue(image.Name)
	state.Url = types.StringValue(image.Url)
	state.LastUpdated = types.StringValue(image.LastOpDate.GoString())
	//state.Description = types.StringValue(image.Description)

	if !state.Description.IsNull() {
		state.Description = types.StringValue(image.Description)
	}
	if !state.GuestOsType.IsNull() {
		state.GuestOsType = types.StringValue(image.GuestOsType)
	}
	if !state.Platform.IsNull() {
		state.Platform = types.StringValue(image.Platform)
	}

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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The timestamp of the last update to the image resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the image. This is a mandatory field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the image, providing additional context or details.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL where the image is located. This can be a file path or an HTTP link.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"media_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of media for the image. Examples include 'ISO' or 'RootVolumeTemplate' or DataVolumeTemplate.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ISO", "RootVolumeTemplate", "DataVolumeTemplate"),
				},
			},
			"guest_os_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The guest operating system type that the image is optimized for.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"system": schema.StringAttribute{
				Computed:    true,
				Description: "Indicates if the image is a system image. Set automatically by ZStack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"platform": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The platform that the image is intended for, such as 'Linux', 'Windows', or others.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Linux", "Windows", "Other"),
				},
			},
			"format": schema.StringAttribute{
				Required:    true,
				Description: "The format of the image file, such as 'qcow2', 'raw', or 'vmdk'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("qcow2", "iso", "raw", "vmdk"),
				},
			},
			"backup_storage_uuids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "A list of UUIDs for the backup storages where the image is stored.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"architecture": schema.StringAttribute{
				Optional:    true,
				Description: "The architecture of the image, such as 'x86_64' or 'aarch64'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("x86_64", "aarch64", "mips64el", "loongarch64"),
				},
			},
			/*
				"type": schema.StringAttribute{
					Computed:    true,
					Description: "The type of the image, for example, 'ISO' or 'RootVolumeTemplate'.",
				},
			*/
			"virtio": schema.StringAttribute{
				Optional:    true,
				Description: "Indicates if the VirtIO drivers are required for the image.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expunge": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the image should be expunged after deletion.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"boot_mode": schema.StringAttribute{
				Optional:    true,
				Description: "The boot mode supported by the image, such as 'Legacy' or 'UEFI'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("Legacy", "UEFI"),
				},
			},
		},
	}
}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Image resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *imageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

/*

func removeStringFromSlice(slice []string, s string) []string {
	for i, v := range slice {
		if v == s {
			slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return slice
}
*/
