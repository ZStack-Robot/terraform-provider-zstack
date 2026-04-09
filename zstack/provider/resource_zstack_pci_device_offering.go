// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
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
	_ resource.Resource                = &pciDeviceOfferingResource{}
	_ resource.ResourceWithConfigure   = &pciDeviceOfferingResource{}
	_ resource.ResourceWithImportState = &pciDeviceOfferingResource{}
)

type pciDeviceOfferingResource struct {
	client *client.ZSClient
}

type pciDeviceOfferingModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VendorId    types.String `tfsdk:"vendor_id"`
	DeviceId    types.String `tfsdk:"device_id"`
	SubvendorId types.String `tfsdk:"subvendor_id"`
	SubdeviceId types.String `tfsdk:"subdevice_id"`
	RamSize     types.String `tfsdk:"ram_size"`
	Type        types.String `tfsdk:"type"`
}

func PciDeviceOfferingResource() resource.Resource {
	return &pciDeviceOfferingResource{}
}

func (r *pciDeviceOfferingResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *pciDeviceOfferingResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_pci_device_offering"
}

func (r *pciDeviceOfferingResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage PCI device offerings in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the PCI device offering.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "The offering name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vendor_id": schema.StringAttribute{
				Required:    true,
				Description: "Vendor ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_id": schema.StringAttribute{
				Required:    true,
				Description: "Device ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subvendor_id": schema.StringAttribute{
				Optional:    true,
				Description: "Subvendor ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subdevice_id": schema.StringAttribute{
				Optional:    true,
				Description: "Subdevice ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ram_size": schema.StringAttribute{
				Optional:    true,
				Description: "RAM size.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Offering type.",
			},
		},
	}
}

func (r *pciDeviceOfferingResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan pciDeviceOfferingModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreatePciDeviceOfferingParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePciDeviceOfferingParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			VendorId:    plan.VendorId.ValueString(),
			DeviceId:    plan.DeviceId.ValueString(),
			SubvendorId: stringPtrOrNil(plan.SubvendorId.ValueString()),
			SubdeviceId: stringPtrOrNil(plan.SubdeviceId.ValueString()),
			RamSize:     stringPtrOrNil(plan.RamSize.ValueString()),
		},
	}

	offering, err := r.client.CreatePciDeviceOffering(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating PCI Device Offering",
			"Could not create PCI device offering, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(offering.UUID)
	plan.Name = stringValueOrNull(offering.Name)
	plan.Description = stringValueOrNull(offering.Description)
	plan.VendorId = stringValueOrNull(offering.VendorId)
	plan.DeviceId = stringValueOrNull(offering.DeviceId)
	plan.SubvendorId = stringValueOrNull(offering.SubvendorId)
	plan.SubdeviceId = stringValueOrNull(offering.SubdeviceId)
	plan.RamSize = stringValueOrNull(offering.RamSize)
	plan.Type = stringValueOrNull(offering.Type)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *pciDeviceOfferingResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state pciDeviceOfferingModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	offering, err := findResourceByQuery(r.client.QueryPciDeviceOffering, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query PCI device offerings. It may have been deleted.: "+err.Error())
		state = pciDeviceOfferingModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(offering.UUID)
	state.Name = stringValueOrNull(offering.Name)
	state.Description = stringValueOrNull(offering.Description)
	state.VendorId = stringValueOrNull(offering.VendorId)
	state.DeviceId = stringValueOrNull(offering.DeviceId)
	state.SubvendorId = stringValueOrNull(offering.SubvendorId)
	state.SubdeviceId = stringValueOrNull(offering.SubdeviceId)
	state.RamSize = stringValueOrNull(offering.RamSize)
	state.Type = stringValueOrNull(offering.Type)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *pciDeviceOfferingResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"PCI Device Offering resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *pciDeviceOfferingResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state pciDeviceOfferingModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "PCI device offering UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeletePciDeviceOffering(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting PCI Device Offering", "Could not delete PCI device offering UUID "+state.Uuid.ValueString()+": "+err.Error())
		return
	}
}

func (r *pciDeviceOfferingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
