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
	_ resource.Resource                = &licenseResource{}
	_ resource.ResourceWithConfigure   = &licenseResource{}
	_ resource.ResourceWithImportState = &licenseResource{}
)

type licenseResource struct {
	client *client.ZSClient
}

type licenseModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	License            types.String `tfsdk:"license"`
	ManagementNodeUuid types.String `tfsdk:"management_node_uuid"`
	User               types.String `tfsdk:"user"`
	ProdInfo           types.String `tfsdk:"prod_info"`
	CpuNum             types.Int64  `tfsdk:"cpu_num"`
	HostNum            types.Int64  `tfsdk:"host_num"`
	VmNum              types.Int64  `tfsdk:"vm_num"`
	Capacity           types.Int64  `tfsdk:"capacity"`
	LicenseType        types.String `tfsdk:"license_type"`
	QuotaType          types.String `tfsdk:"quota_type"`
	ExpiredDate        types.String `tfsdk:"expired_date"`
	IssuedDate         types.String `tfsdk:"issued_date"`
	UploadDate         types.String `tfsdk:"upload_date"`
	Expired            types.Bool   `tfsdk:"expired"`
	Source             types.String `tfsdk:"source"`
	PlatformId         types.String `tfsdk:"platform_id"`
	AvailableHostNum   types.Int64  `tfsdk:"available_host_num"`
	AvailableCpuNum    types.Int64  `tfsdk:"available_cpu_num"`
	AvailableVmNum     types.Int64  `tfsdk:"available_vm_num"`
}

func populateLicenseModelFromInventory(state *licenseModel, license *view.LicenseInventoryView) {
	state.Uuid = stringValueOrNull(license.UUID)
	state.ManagementNodeUuid = stringValueOrNull(license.ManagementNodeUuid)
	state.User = stringValueOrNull(license.User)
	state.ProdInfo = stringValueOrNull(license.ProdInfo)
	state.CpuNum = types.Int64Value(int64(license.CpuNum))
	state.HostNum = types.Int64Value(int64(license.HostNum))
	state.VmNum = types.Int64Value(int64(license.VmNum))
	state.Capacity = types.Int64Value(int64(license.Capacity))
	state.LicenseType = stringValueOrNull(license.LicenseType)
	state.QuotaType = stringValueOrNull(license.QuotaType)
	state.ExpiredDate = stringValueOrNull(license.ExpiredDate)
	state.IssuedDate = stringValueOrNull(license.IssuedDate)
	state.UploadDate = stringValueOrNull(license.UploadDate)
	state.Expired = types.BoolValue(license.Expired)
	state.Source = stringValueOrNull(license.Source)
	state.PlatformId = stringValueOrNull(license.PlatformId)
	state.AvailableHostNum = types.Int64Value(int64(license.AvailableHostNum))
	state.AvailableCpuNum = types.Int64Value(int64(license.AvailableCpuNum))
	state.AvailableVmNum = types.Int64Value(int64(license.AvailableVmNum))
}

func LicenseResource() resource.Resource {
	return &licenseResource{}
}

func (r *licenseResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *licenseResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_license"
}

func (r *licenseResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a ZStack license.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the license.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"license": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The license text to upload.",
			},
			"management_node_uuid": schema.StringAttribute{
				Required:    true,
				Computed:    true,
				Description: "The management node UUID used for license upload.",
			},
			"user": schema.StringAttribute{
				Computed:    true,
				Description: "The license user.",
			},
			"prod_info": schema.StringAttribute{
				Computed:    true,
				Description: "The product information from the license.",
			},
			"cpu_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The licensed CPU quota.",
			},
			"host_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The licensed host quota.",
			},
			"vm_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The licensed VM quota.",
			},
			"capacity": schema.Int64Attribute{
				Computed:    true,
				Description: "The licensed capacity quota.",
			},
			"license_type": schema.StringAttribute{
				Computed:    true,
				Description: "The license type.",
			},
			"quota_type": schema.StringAttribute{
				Computed:    true,
				Description: "The quota type.",
			},
			"expired_date": schema.StringAttribute{
				Computed:    true,
				Description: "The expiration date of the license.",
			},
			"issued_date": schema.StringAttribute{
				Computed:    true,
				Description: "The issued date of the license.",
			},
			"upload_date": schema.StringAttribute{
				Computed:    true,
				Description: "The upload date of the license.",
			},
			"expired": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the license is expired.",
			},
			"source": schema.StringAttribute{
				Computed:    true,
				Description: "The source of the license.",
			},
			"platform_id": schema.StringAttribute{
				Computed:    true,
				Description: "The platform ID from the license.",
			},
			"available_host_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The available host quota.",
			},
			"available_cpu_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The available CPU quota.",
			},
			"available_vm_num": schema.Int64Attribute{
				Computed:    true,
				Description: "The available VM quota.",
			},
		},
	}
}

func (r *licenseResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan licenseModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateLicenseParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateLicenseParamDetail{
			License: plan.License.ValueString(),
		},
	}

	license, err := r.client.UpdateLicense(plan.ManagementNodeUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Error creating license", err.Error())
		return
	}

	populateLicenseModelFromInventory(&plan, license)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
}

func (r *licenseResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state licenseModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	existingLicense := state.License
	license, err := r.client.GetLicenseInfo()
	if err != nil {
		if isZStackNotFoundError(err) {
			tflog.Warn(ctx, "License not found. It may have been deleted outside of Terraform: "+err.Error())
			response.State.RemoveResource(ctx)
			return
		}

		response.Diagnostics.AddError("Unable to read license info", err.Error())
		return
	}

	existingManagementNodeUuid := state.ManagementNodeUuid
	populateLicenseModelFromInventory(&state, license)
	state.License = existingLicense

	// Preserve management_node_uuid from state if API returns null/empty
	if state.ManagementNodeUuid.IsNull() || state.ManagementNodeUuid.IsUnknown() {
		if !existingManagementNodeUuid.IsNull() && !existingManagementNodeUuid.IsUnknown() {
			state.ManagementNodeUuid = existingManagementNodeUuid
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *licenseResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan licenseModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateLicenseParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateLicenseParamDetail{
			License: plan.License.ValueString(),
		},
	}

	license, err := r.client.UpdateLicense(plan.ManagementNodeUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Error updating license", err.Error())
		return
	}

	populateLicenseModelFromInventory(&plan, license)

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
}

func (r *licenseResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state licenseModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLicense(state.ManagementNodeUuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting license", err.Error())
		return
	}
}

func (r *licenseResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
