// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &gpuDeviceDataSource{}
	_ datasource.DataSourceWithConfigure = &gpuDeviceDataSource{}
)

type gpuDeviceDataSource struct {
	client *client.ZSClient
}

type gpuDeviceDataSourceModel struct {
	Name       types.String     `tfsdk:"name"`
	NamePattern types.String    `tfsdk:"name_pattern"`
	Filter     []Filter         `tfsdk:"filter"`
	GpuDevices []gpuDevicesModel `tfsdk:"gpu_devices"`
}

type gpuDevicesModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	GpuType          types.String `tfsdk:"gpu_type"`
	GpuStatus        types.String `tfsdk:"gpu_status"`
	AllocateStatus   types.String `tfsdk:"allocate_status"`
	HostUuid         types.String `tfsdk:"host_uuid"`
	VmInstanceUuid   types.String `tfsdk:"vm_instance_uuid"`
	VendorId         types.String `tfsdk:"vendor_id"`
	Vendor           types.String `tfsdk:"vendor"`
	DeviceId         types.String `tfsdk:"device_id"`
	Device           types.String `tfsdk:"device"`
	Type             types.String `tfsdk:"type"`
	State            types.String `tfsdk:"state"`
	Status           types.String `tfsdk:"status"`
	Memory           types.Int64  `tfsdk:"memory"`
	PciDeviceAddress types.String `tfsdk:"pci_device_address"`
}

func ZStackGpuDeviceDataSource() datasource.DataSource {
	return &gpuDeviceDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *gpuDeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = client
}

// Metadata implements datasource.DataSource.
func (d *gpuDeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpu_devices"
}

// Schema implements datasource.DataSource.
func (d *gpuDeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of GPU devices and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of GPU devices and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching GPU devices.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"gpu_devices": schema.ListNestedAttribute{
				Description: "List of GPU devices matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the GPU device.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the GPU device.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the GPU device.",
							Computed:    true,
						},
						"gpu_type": schema.StringAttribute{
							Description: "Type of the GPU (e.g., vGPU, passthrough).",
							Computed:    true,
						},
						"gpu_status": schema.StringAttribute{
							Description: "Status of the GPU device.",
							Computed:    true,
						},
						"allocate_status": schema.StringAttribute{
							Description: "Allocation status of the GPU (e.g., allocated, available).",
							Computed:    true,
						},
						"host_uuid": schema.StringAttribute{
							Description: "UUID of the host where the GPU is installed.",
							Computed:    true,
						},
						"vm_instance_uuid": schema.StringAttribute{
							Description: "UUID of the VM instance the GPU is attached to, if any.",
							Computed:    true,
						},
						"vendor_id": schema.StringAttribute{
							Description: "PCI vendor ID.",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "GPU vendor name (e.g., NVIDIA).",
							Computed:    true,
						},
						"device_id": schema.StringAttribute{
							Description: "PCI device ID.",
							Computed:    true,
						},
						"device": schema.StringAttribute{
							Description: "GPU device model name.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Device type.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the GPU device (Enabled, Disabled).",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Connection status of the GPU device.",
							Computed:    true,
						},
						"memory": schema.Int64Attribute{
							Description: "GPU memory in bytes.",
							Computed:    true,
						},
						"pci_device_address": schema.StringAttribute{
							Description: "PCI device address.",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by vendor, use `name = \"vendor\"` and `values = [\"NVIDIA\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., vendor, gpu_type, state).",
							Required:    true,
						},
						"values": schema.SetAttribute{
							Description: "Values to filter by. Multiple values will be treated as an OR condition.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Read implements datasource.DataSource.
func (d *gpuDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state gpuDeviceDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	gpus, err := d.client.QueryGpuDevice(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack GPU Devices",
			err.Error(),
		)
		return
	}

	filters := make(map[string][]string)
	for _, filter := range state.Filter {
		values := make([]string, 0, len(filter.Values.Elements()))
		diags := filter.Values.ElementsAs(ctx, &values, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		filters[filter.Name.ValueString()] = values
	}

	filteredGpus, filterDiags := utils.FilterResource(ctx, gpus, filters, "gpu_device")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.GpuDevices = []gpuDevicesModel{}

	for _, gpu := range filteredGpus {
		gpuState := gpuDevicesModel{
			Uuid:             types.StringValue(gpu.UUID),
			Name:             types.StringValue(gpu.Name),
			Description:      stringValueOrNull(gpu.Description),
			GpuType:          stringValueOrNull(gpu.GpuType),
			GpuStatus:        stringValueOrNull(gpu.GpuStatus),
			AllocateStatus:   stringValueOrNull(gpu.AllocateStatus),
			HostUuid:         stringValueOrNull(gpu.HostUuid),
			VmInstanceUuid:   stringValueOrNull(gpu.VmInstanceUuid),
			VendorId:         stringValueOrNull(gpu.VendorId),
			Vendor:           stringValueOrNull(gpu.Vendor),
			DeviceId:         stringValueOrNull(gpu.DeviceId),
			Device:           stringValueOrNull(gpu.Device),
			Type:             stringValueOrNull(gpu.Type),
			State:            stringValueOrNull(gpu.State),
			Status:           stringValueOrNull(gpu.Status),
			Memory:           types.Int64Value(gpu.Memory),
			PciDeviceAddress: stringValueOrNull(gpu.PciDeviceAddress),
		}
		state.GpuDevices = append(state.GpuDevices, gpuState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
