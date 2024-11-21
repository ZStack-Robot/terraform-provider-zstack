// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ datasource.DataSource              = &vmsDataSource{}
	_ datasource.DataSourceWithConfigure = &vmsDataSource{}
)

type vmsDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	NamePattern types.String `tfsdk:"name_pattern"`
	VmInstances []vmsModel   `tfsdk:"vminstances"`
}

type vmsModel struct {
	Name           types.String      `tfsdk:"name"`
	HypervisorType types.String      `tfsdk:"hypervisor_type"`
	State          types.String      `tfsdk:"state"`
	Type           types.String      `tfsdk:"type"`
	Uuid           types.String      `tfsdk:"uuid"`
	ZoneUuid       types.String      `tfsdk:"zone_uuid"`
	ClusterUuid    types.String      `tfsdk:"cluster_uuid"`
	ImageUuid      types.String      `tfsdk:"image_uuid"`
	HostUuid       types.String      `tfsdk:"host_uuid"`
	Platform       types.String      `tfsdk:"platform"`
	Architecture   types.String      `tfsdk:"architecture"`
	CPUNum         types.Int64       `tfsdk:"cup_num"`
	MemorySize     types.Int64       `tfsdk:"memory_size"`
	VmNics         []vmNicsModel     `tfsdk:"vm_nics"`
	AllVolumes     []allVolumesModel `tfsdk:"all_volumes"`
}

type vmNicsModel struct {
	IP      types.String `tfsdk:"ip"`
	Mac     types.String `tfsdk:"mac"`
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
}

type allVolumesModel struct {
	VolumeUuid        types.String `tfsdk:"volume_uuid"`
	VolumeDescription types.String `tfsdk:"volume_description"`
	VolumeType        types.String `tfsdk:"volume_type"`
	VolumeFormat      types.String `tfsdk:"volume_format"`
	VolumeSize        types.Int64  `tfsdk:"volume_size"`
	VolumeActualSize  types.Int64  `tfsdk:"volume_actual_size"`
	VolumeState       types.String `tfsdk:"volume_state"`
	VolumeStatus      types.String `tfsdk:"volume_status"`
}

func ZStackvmsDataSource() datasource.DataSource {
	return &vmsDataSource{}
}

type vmsDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *vmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Metadata implements datasource.DataSourceWithConfigure.
func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instances"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vmsDataSourceModel

	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	vminstances, err := d.client.QueryVmInstance(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read vm instances",
			err.Error(),
		)
		return
	}
	for _, vminstance := range vminstances {
		vminstanceState := vmsModel{
			Name:           types.StringValue(vminstance.Name),
			HypervisorType: types.StringValue(vminstance.HypervisorType),
			State:          types.StringValue(vminstance.State),
			Type:           types.StringValue(vminstance.Type),
			Uuid:           types.StringValue(vminstance.UUID),
			ZoneUuid:       types.StringValue(vminstance.ZoneUUID),
			ClusterUuid:    types.StringValue(vminstance.ClusterUUID),
			ImageUuid:      types.StringValue(vminstance.ImageUUID),
			HostUuid:       types.StringValue(vminstance.HostUUID),
			Platform:       types.StringValue(vminstance.Platform),
			Architecture:   types.StringValue(vminstance.Architecture),
			CPUNum:         types.Int64Value(int64(vminstance.CPUNum)),
			MemorySize:     types.Int64Value(int64(vminstance.MemorySize)),
		}

		for _, vmnics := range vminstance.VMNics {
			vminstanceState.VmNics = append(vminstanceState.VmNics, vmNicsModel{
				IP:      types.StringValue(vmnics.IP),
				Mac:     types.StringValue(vmnics.Mac),
				Netmask: types.StringValue(vmnics.Netmask),
				Gateway: types.StringValue(vmnics.Gateway),
			})
		}

		for _, allvolumes := range vminstance.AllVolumes {
			vminstanceState.AllVolumes = append(vminstanceState.AllVolumes, allVolumesModel{
				VolumeUuid:        types.StringValue(allvolumes.UUID),
				VolumeDescription: types.StringValue(allvolumes.Description),
				VolumeType:        types.StringValue(allvolumes.Type),
				VolumeFormat:      types.StringValue(allvolumes.Format),
				VolumeSize:        types.Int64Value(int64(allvolumes.Size)),
				VolumeActualSize:  types.Int64Value(int64(allvolumes.ActualSize)),
				VolumeState:       types.StringValue(allvolumes.State),
				VolumeStatus:      types.StringValue(allvolumes.Status),
			})
		}

		state.VmInstances = append(state.VmInstances, vminstanceState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *vmsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching VM instance",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"vminstances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the VM instance.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the VM instance.",
						},
						"hypervisor_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of hypervisor on which the VM is running (e.g., KVM, VMware).",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the VM (e.g., Running, Stopped).",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the VM (e.g., UserVm or SystemVm).",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the zone in which the VM is located.",
						},
						"cluster_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the cluster in which the VM is located.",
						},
						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the image used to create the VM.",
						},
						"host_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the host on which the VM is running.",
						},
						"platform": schema.StringAttribute{
							Computed:    true,
							Description: "The platform (e.g., Linux, Windows) on which the VM is running.",
						},
						"architecture": schema.StringAttribute{
							Computed:    true,
							Description: "The CPU architecture (e.g., x86_64, ARM) of the VM.",
						},
						"cup_num": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of CPUs allocated to the VM.",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The amount of memory (in bytes) allocated to the VM.",
						},
						"vm_nics": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{
										Computed:    true,
										Description: "The IP address assigned to the VM NIC.",
									},
									"mac": schema.StringAttribute{
										Computed:    true,
										Description: "The MAC address of the VM NIC.",
									},
									"netmask": schema.StringAttribute{
										Computed:    true,
										Description: "The network mask of the VM NIC.",
									},
									"gateway": schema.StringAttribute{
										Computed:    true,
										Description: "The gateway IP address for the VM NIC.",
									},
								},
							},
						},
						"all_volumes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volume_uuid": schema.StringAttribute{
										Computed:    true,
										Description: "The UUID of the volume attached to the VM.",
									},
									"volume_description": schema.StringAttribute{
										Computed:    true,
										Description: "The description of the volume attached to the VM.",
									},
									"volume_type": schema.StringAttribute{
										Computed:    true,
										Description: "The type of the volume (e.g., root, data).",
									},
									"volume_format": schema.StringAttribute{
										Computed:    true,
										Description: "The format of the volume (e.g., RAW, QCOW2).",
									},
									"volume_size": schema.Int64Attribute{
										Computed:    true,
										Description: "The size of the volume (in bytes).",
									},
									"volume_actual_size": schema.Int64Attribute{
										Computed:    true,
										Description: "The actual size of the volume (in bytes), which might differ from the requested size.",
									},
									"volume_state": schema.StringAttribute{
										Computed:    true,
										Description: "The state of the volume (e.g., Enabled, Disabled).",
									},
									"volume_status": schema.StringAttribute{
										Computed:    true,
										Description: "The status of the volume (e.g., Ready, NoReady).",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
