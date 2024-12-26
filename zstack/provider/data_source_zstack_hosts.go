// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ datasource.DataSource              = &hostsDataSource{}
	_ datasource.DataSourceWithConfigure = &hostsDataSource{}
)

func ZStackHostsDataSource() datasource.DataSource {
	return &hostsDataSource{}
}

type hostsDataSource struct {
	client *client.ZSClient
}

type hostsDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	NamePattern types.String `tfsdk:"name_pattern"`
	Filter      types.Map    `tfsdk:"filter"`
	Hosts       []hostsModel `tfsdk:"hosts"`
}

type hostsModel struct {
	Name         types.String `tfsdk:"name"`
	Architecture types.String `tfsdk:"architecture"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
	Type         types.String `tfsdk:"type"`
	Uuid         types.String `tfsdk:"uuid"`
	ZoneUuid     types.String `tfsdk:"zone_uuid"`
	ClusterUuid  types.String `tfsdk:"cluster_uuid"`
	ManagementIp types.String `tfsdk:"managementip"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *hostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	//name_regex := state.Name
	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	hosts, err := d.client.QueryHost(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Hosts ",
			err.Error(),
		)
		return
	}

	filters := make(map[string]string)
	if !state.Filter.IsNull() {
		diags := state.Filter.ElementsAs(ctx, &filters, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	filterHosts, filterDiags := utils.FilterResource(ctx, hosts, filters)
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, host := range filterHosts {
		HostsState := hostsModel{
			Name:         types.StringValue(host.Name),
			State:        types.StringValue(host.State),
			Status:       types.StringValue(host.Status),
			Uuid:         types.StringValue(host.UUID),
			Architecture: types.StringValue(host.Architecture),
			Type:         types.StringValue(host.HypervisorType),
			ZoneUuid:     types.StringValue(host.ZoneUuid),
			ClusterUuid:  types.StringValue(host.ClusterUuid),
			ManagementIp: types.StringValue(host.ManagementIp),
		}

		state.Hosts = append(state.Hosts, HostsState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of hosts and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching hosts",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"filter": schema.MapAttribute{
				Description: "Key-value pairs to filter hosts. For example, to filter by State, use `State = \"Enabled\"`.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"hosts": schema.ListNestedAttribute{
				Description: "List of host entries matching the specified filters",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID Unique identifier of the host",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the host",
						},
						"architecture": schema.StringAttribute{
							Computed:    true,
							Description: "CPU architecture of the host (e.g., x86_64, arm64)",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "State of the host (e.g., Enabled, Disabled)",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Operational status of the host (e.g., Connected, Disconnected)",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the host (e.g., bare metal, virtualized)",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the zone to which the host belongs",
						},
						"cluster_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the cluster to which the host belongs",
						},
						"managementip": schema.StringAttribute{
							Computed:    true,
							Description: "Current management operation status on the host (e.g., Pending, Completed)",
						},
					},
				},
			},
		},
	}

}
