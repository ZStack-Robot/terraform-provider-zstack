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
	_ datasource.DataSource              = &l2VlanNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &l2VlanNetworkDataSource{}
)

type l2VlanNetworkDataSource struct {
	client *client.ZSClient
}

type l2VlanNetworkDataSourceModel struct {
	Name            types.String           `tfsdk:"name"`
	NamePattern     types.String           `tfsdk:"name_pattern"`
	Filter          []Filter               `tfsdk:"filter"`
	L2VlanNetworks  []l2VlanNetworksModel  `tfsdk:"l2vlan_networks"`
}

type l2VlanNetworksModel struct {
	Uuid                 types.String   `tfsdk:"uuid"`
	Name                 types.String   `tfsdk:"name"`
	Description          types.String   `tfsdk:"description"`
	Vlan                 types.Int64    `tfsdk:"vlan"`
	ZoneUuid             types.String   `tfsdk:"zone_uuid"`
	PhysicalInterface    types.String   `tfsdk:"physical_interface"`
	Type                 types.String   `tfsdk:"type"`
	VSwitchType          types.String   `tfsdk:"vswitch_type"`
	AttachedClusterUuids []types.String `tfsdk:"attached_cluster_uuids"`
}

func ZStackL2VlanNetworkDataSource() datasource.DataSource {
	return &l2VlanNetworkDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *l2VlanNetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *l2VlanNetworkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vlan_networks"
}

// Schema implements datasource.DataSource.
func (d *l2VlanNetworkDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of L2 VLAN networks and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of L2 VLAN networks and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching L2 VLAN networks.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"l2vlan_networks": schema.ListNestedAttribute{
				Description: "List of L2 VLAN networks matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the L2 VLAN network.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the L2 VLAN network.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the L2 VLAN network.",
							Computed:    true,
						},
						"vlan": schema.Int64Attribute{
							Description: "VLAN ID of the L2 VLAN network.",
							Computed:    true,
						},
						"zone_uuid": schema.StringAttribute{
							Description: "UUID of the zone where the L2 VLAN network resides.",
							Computed:    true,
						},
						"physical_interface": schema.StringAttribute{
							Description: "Physical interface associated with the L2 VLAN network.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the L2 network.",
							Computed:    true,
						},
						"vswitch_type": schema.StringAttribute{
							Description: "Virtual switch type of the L2 VLAN network.",
							Computed:    true,
						},
						"attached_cluster_uuids": schema.ListAttribute{
							Description: "UUIDs of clusters attached to the L2 VLAN network.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by VLAN ID, use `name = \"vlan\"` and `values = [\"100\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., vlan, zone_uuid).",
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
func (d *l2VlanNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state l2VlanNetworkDataSourceModel

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

	l2VlanNetworks, err := d.client.QueryL2VlanNetwork(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack L2 VLAN Networks",
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

	filteredNetworks, filterDiags := utils.FilterResource(ctx, l2VlanNetworks, filters, "l2vlan_network")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.L2VlanNetworks = []l2VlanNetworksModel{}

	for _, network := range filteredNetworks {
		attachedClusters := make([]types.String, 0, len(network.AttachedClusterUuids))
		for _, clusterUuid := range network.AttachedClusterUuids {
			attachedClusters = append(attachedClusters, types.StringValue(clusterUuid))
		}

		networkState := l2VlanNetworksModel{
			Uuid:                 types.StringValue(network.UUID),
			Name:                 types.StringValue(network.Name),
			Description:          stringValueOrNull(network.Description),
			Vlan:                 types.Int64Value(int64(network.Vlan)),
			ZoneUuid:             types.StringValue(network.ZoneUuid),
			PhysicalInterface:    stringValueOrNull(network.PhysicalInterface),
			Type:                 stringValueOrNull(network.Type),
			VSwitchType:          stringValueOrNull(network.VSwitchType),
			AttachedClusterUuids: attachedClusters,
		}
		state.L2VlanNetworks = append(state.L2VlanNetworks, networkState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
