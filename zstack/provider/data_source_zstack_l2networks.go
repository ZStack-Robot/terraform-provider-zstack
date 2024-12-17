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
	_ datasource.DataSource              = &l2NetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &l2NetworkDataSource{}
)

type l2NetworkDataSource struct {
	client *client.ZSClient
}

type l2NetworkDataSourceModel struct {
	Name        types.String      `tfsdk:"name"`
	NamePattern types.String      `tfsdk:"name_pattern"`
	L2networks  []l2networksModel `tfsdk:"l2networks"`
}

type l2networksModel struct {
	Name                 types.String   `tfsdk:"name"`
	Uuid                 types.String   `tfsdk:"uuid"`
	Vlan                 types.Int64    `tfsdk:"vlan"`
	ZoneUuid             types.String   `tfsdk:"zone_uuid"`
	PhysicalInterface    types.String   `tfsdk:"physical_interface"`
	Type                 types.String   `tfsdk:"type"`
	AttachedClusterUuids []types.String `tfsdk:"attached_cluster_uuids"`
}

func ZStackl2NetworkDataSource() datasource.DataSource {
	return &l2NetworkDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *l2NetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *l2NetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2networks"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *l2NetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state l2NetworkDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	//Create query parameters based on name

	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	//Query L2 networks with name filtering
	l2networks, err := d.client.QueryL2Network(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack L2Networks ",
			err.Error(),
		)
		return
	}

	state.L2networks = []l2networksModel{}

	// Process each L2 network and populate the state
	for _, l2network := range l2networks {
		attachedClusters := []types.String{}
		for _, clusterUuid := range l2network.AttachedClusterUuids {
			attachedClusters = append(attachedClusters, types.StringValue(clusterUuid))
		}

		l2networkState := l2networksModel{
			Name:                 types.StringValue(l2network.Name),
			Uuid:                 types.StringValue(l2network.UUID),
			Vlan:                 types.Int64Value(int64(l2network.Vlan)),
			ZoneUuid:             types.StringValue(l2network.ZoneUuid),
			PhysicalInterface:    types.StringValue(l2network.PhysicalInterface),
			Type:                 types.StringValue(l2network.Type),
			AttachedClusterUuids: attachedClusters,
		}
		state.L2networks = append(state.L2networks, l2networkState)
	}

	// Set the final state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *l2NetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of L2 networks and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching L2 Network.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"l2networks": schema.ListNestedAttribute{
				Description: "List of L2 networks matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the L2 network.",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Description: "UUID of the L2 network.",
							Computed:    true,
						},
						"vlan": schema.Int64Attribute{
							Description: "VLAN ID of the L2 network.",
							Computed:    true,
						},
						"zone_uuid": schema.StringAttribute{
							Description: "UUID of the zone where the L2 network resides.",
							Computed:    true,
						},
						"physical_interface": schema.StringAttribute{
							Description: "Physical interface associated with the L2 network.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the L2 network.",
							Computed:    true,
						},
						"attached_cluster_uuids": schema.ListAttribute{
							Description: "UUIDs of clusters attached to the L2 network.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
