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
	_ datasource.DataSource              = &vrouterDataSource{}
	_ datasource.DataSourceWithConfigure = &vrouterDataSource{}
)

type vrouterDataSourceModel struct {
	Name        types.String   `tfsdk:"name"`
	NamePattern types.String   `tfsdk:"name_pattern"`
	Vrouter     []vrouterModel `tfsdk:"virtual_router"`
}

type vrouterModel struct {
	Name            types.String `tfsdk:"name"`
	Uuid            types.String `tfsdk:"uuid"`
	HypervisorType  types.String `tfsdk:"hypervisor_type"`
	ApplianceVmType types.String `tfsdk:"appliance_vm_type"`
	State           types.String `tfsdk:"state"`
	Status          types.String `tfsdk:"status"`
	AgentPort       types.Int64  `tfsdk:"agent_port"`
	Type            types.String `tfsdk:"type"`
	HaStatus        types.String `tfsdk:"ha_status"`

	ZoneUuid              types.String       `tfsdk:"zone_uuid"`
	ClusterUuid           types.String       `tfsdk:"cluster_uuid"`
	ManagementNetworkUuid types.String       `tfsdk:"management_network_uuid"`
	ImageUuid             types.String       `tfsdk:"image_uuid"`
	HostUuid              types.String       `tfsdk:"host_uuid"`
	InstanceOfferingUUID  types.String       `tfsdk:"instance_offering_uuid"`
	Platform              types.String       `tfsdk:"platform"`
	Architecture          types.String       `tfsdk:"architecture"`
	CPUNum                types.Int64        `tfsdk:"cup_num"`
	MemorySize            types.Int64        `tfsdk:"memory_size"`
	VmNics                []vrouterNicsModel `tfsdk:"vm_nics"`
}

type vrouterNicsModel struct {
	IP      types.String `tfsdk:"ip"`
	Mac     types.String `tfsdk:"mac"`
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
}

func ZStackvrouterDataSource() datasource.DataSource {
	return &vrouterDataSource{}
}

type vrouterDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *vrouterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vrouterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_routers"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *vrouterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vrouterDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	// 优先检查 `name` 精确查询
	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	vrouters, err := d.client.QueryVirtualRouterVm(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read virtual router instances",
			err.Error(),
		)
		return
	}
	for _, vrouter := range vrouters {
		vrouterState := vrouterModel{
			Name:            types.StringValue(vrouter.Name),
			Uuid:            types.StringValue(vrouter.UUID),
			HypervisorType:  types.StringValue(vrouter.HypervisorType),
			ApplianceVmType: types.StringValue(vrouter.ApplianceVmType),
			State:           types.StringValue(vrouter.State),
			Status:          types.StringValue(vrouter.Status),
			AgentPort:       types.Int64Value(int64(vrouter.AgentPort)),
			Type:            types.StringValue(vrouter.Type),
			HaStatus:        types.StringValue(vrouter.HaStatus),

			ZoneUuid:              types.StringValue(vrouter.ZoneUuid),
			ClusterUuid:           types.StringValue(vrouter.ClusterUUID),
			ManagementNetworkUuid: types.StringValue(vrouter.ManagementNetworkUuid),
			ImageUuid:             types.StringValue(vrouter.ImageUUID),
			HostUuid:              types.StringValue(vrouter.HostUuid),
			InstanceOfferingUUID:  types.StringValue(vrouter.InstanceOfferingUUID),

			Platform:     types.StringValue(vrouter.Platform),
			Architecture: types.StringValue(vrouter.Architecture),
			CPUNum:       types.Int64Value(int64(vrouter.CPUNum)),
			MemorySize:   types.Int64Value(int64(vrouter.MemorySize)),
		}

		for _, vmnics := range vrouter.VMNics {
			vrouterState.VmNics = append(vrouterState.VmNics, vrouterNicsModel{
				IP:      types.StringValue(vmnics.IP),
				Mac:     types.StringValue(vmnics.Mac),
				Netmask: types.StringValue(vmnics.Netmask),
				Gateway: types.StringValue(vmnics.Gateway),
			})
		}

		state.Vrouter = append(state.Vrouter, vrouterState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *vrouterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching virtual router instance",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"virtual_router": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the virtual router instance.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the virtual router instance.",
						},
						"hypervisor_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of hypervisor on which the virtual router is running (e.g., KVM, VMware).",
						},
						"appliance_vm_type": schema.StringAttribute{
							Computed:    true,
							Description: "Specifies the type of appliance VM for the virtual router.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the virtual router (e.g., Running, Stopped).",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Operational status of the virtual router (e.g., Connected, Disconnected)",
						},
						"agent_port": schema.Int64Attribute{
							Computed:    true,
							Description: "The agent's listening port on the virtual router.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the virtual router (e.g., UserVm or SystemVm).",
						},
						"ha_status": schema.StringAttribute{
							Computed:    true,
							Description: "The high-availability (HA) status of the virtual router.",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the zone in which the virtual router is located.",
						},
						"cluster_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the cluster in which the virtual router is located.",
						},
						"management_network_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the management network connected to the virtual router.",
						},
						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the image used to create the virtual router.",
						},
						"host_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the host on which the virtual router is running.",
						},
						"instance_offering_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the instance offering assigned to the virtual router.",
						},
						"platform": schema.StringAttribute{
							Computed:    true,
							Description: "The platform (e.g., Linux, Windows) on which the virtual router is running.",
						},
						"architecture": schema.StringAttribute{
							Computed:    true,
							Description: "The CPU architecture (e.g., x86_64, ARM) of the virtual router.",
						},
						"cup_num": schema.Int64Attribute{
							Computed:    true,
							Description: "The number of CPUs allocated to the virtual router.",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The amount of memory (in bytes) allocated to the virtual router.",
						},
						"vm_nics": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{
										Computed:    true,
										Description: "The IP address assigned to the virtual router NIC.",
									},
									"mac": schema.StringAttribute{
										Computed:    true,
										Description: "The MAC address of the virtual router NIC.",
									},
									"netmask": schema.StringAttribute{
										Computed:    true,
										Description: "The network mask of the virtual router NIC.",
									},
									"gateway": schema.StringAttribute{
										Computed:    true,
										Description: "The gateway IP address for the virtual router NIC.",
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
