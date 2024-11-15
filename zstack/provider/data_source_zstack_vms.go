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
	VmInstances []vmsModel `tfsdk:"vminstances"`
}

type vmsModel struct {
	Name           types.String      `tfsdk:"name"`
	HypervisorType types.String      `tfsdk:"hypervisortype"`
	State          types.String      `tfsdk:"state"`
	Type           types.String      `tfsdk:"type"`
	Uuid           types.String      `tfsdk:"uuid"`
	ZoneUuid       types.String      `tfsdk:"zoneuuid"`
	ClusterUuid    types.String      `tfsdk:"clusteruuid"`
	ImageUuid      types.String      `tfsdk:"imageuuid"`
	HostUuid       types.String      `tfsdk:"hostuuid"`
	Platform       types.String      `tfsdk:"platform"`
	Architecture   types.String      `tfsdk:"architecture"`
	CPUNum         types.Int64       `tfsdk:"cupnum"`
	MemorySize     types.Int64       `tfsdk:"memorysize"`
	VmNics         []vmNicsModel     `tfsdk:"vmnics"`
	AllVolumes     []allVolumesModel `tfsdk:"allvolumes"`
}

type vmNicsModel struct {
	IP      types.String `tfsdk:"ip"`
	Mac     types.String `tfsdk:"mac"`
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
}

type allVolumesModel struct {
	VolumeUuid        types.String `tfsdk:"volumeuuid"`
	VolumeDescription types.String `tfsdk:"volumedescription"`
	VolumeType        types.String `tfsdk:"volumetype"`
	VolumeFormat      types.String `tfsdk:"volumeformat"`
	VolumeSize        types.Int64  `tfsdk:"volumesize"`
	VolumeActualSize  types.Int64  `tfsdk:"volumeactualsize"`
	VolumeState       types.String `tfsdk:"volumestate"`
	VolumeStatus      types.String `tfsdk:"volumestatus"`
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
	resp.TypeName = req.ProviderTypeName + "_vminstances"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vmsDataSourceModel

	vminstances, err := d.client.QueryVmInstance(param.NewQueryParam())
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
			"vminstances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"hypervisortype": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"zoneuuid": schema.StringAttribute{
							Computed: true,
						},
						"clusteruuid": schema.StringAttribute{
							Computed: true,
						},
						"imageuuid": schema.StringAttribute{
							Computed: true,
						},
						"hostuuid": schema.StringAttribute{
							Computed: true,
						},
						"platform": schema.StringAttribute{
							Computed: true,
						},
						"architecture": schema.StringAttribute{
							Computed: true,
						},
						"cupnum": schema.Int64Attribute{
							Computed: true,
						},
						"memorysize": schema.Int64Attribute{
							Computed: true,
						},
						"vmnics": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{
										Computed: true,
									},
									"mac": schema.StringAttribute{
										Computed: true,
									},
									"netmask": schema.StringAttribute{
										Computed: true,
									},
									"gateway": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"allvolumes": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"volumeuuid": schema.StringAttribute{
										Computed: true,
									},
									"volumedescription": schema.StringAttribute{
										Computed: true,
									},
									"volumetype": schema.StringAttribute{
										Computed: true,
									},
									"volumeformat": schema.StringAttribute{
										Computed: true,
									},
									"volumesize": schema.Int64Attribute{
										Computed: true,
									},
									"volumeactualsize": schema.Int64Attribute{
										Computed: true,
									},
									"volumestate": schema.StringAttribute{
										Computed: true,
									},
									"volumestatus": schema.StringAttribute{
										Computed: true,
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
