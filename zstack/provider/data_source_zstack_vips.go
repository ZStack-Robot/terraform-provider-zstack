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
	_ datasource.DataSource              = &vipsDataSource{}
	_ datasource.DataSourceWithConfigure = &vipsDataSource{}
)

func ZStackVIPsDataSource() datasource.DataSource {
	return &vipsDataSource{}
}

type vipsDataSource struct {
	client *client.ZSClient
}

type vipsDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	NamePattern types.String `tfsdk:"name_pattern"`
	VIPs        []vipsModel  `tfsdk:"vips"`
}

type vipsModel struct {
	Name               types.String `tfsdk:"name"`
	Uuid               types.String `tfsdk:"uuid"`
	Description        types.String `tfsdk:"description"`
	PeerL3NetworkUuids types.String `tfsdk:"peer_l3_network_uuids"`
	L3NetworkUUID      types.String `tfsdk:"l3_network_uuid"`
	IP                 types.String `tfsdk:"ip"`
	State              types.String `tfsdk:"state"`
	Gateway            types.String `tfsdk:"gateway"`
	Netmask            types.String `tfsdk:"netmask"`
	UseFor             types.String `tfsdk:"use_for"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *vipsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vipsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vips"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *vipsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vipsDataSourceModel
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

	params.AddQ("system=" + "false") //Just return user VIPS, not include system vips

	vips, err := d.client.QueryVip(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack VIPS ",
			err.Error(),
		)
		return
	}

	for _, vip := range vips {
		VIPsState := vipsModel{
			Uuid:               types.StringValue(vip.UUID),
			Name:               types.StringValue(vip.Name),
			State:              types.StringValue(vip.State),
			Description:        types.StringValue(vip.Description),
			PeerL3NetworkUuids: types.StringValue(vip.PeerL3NetworkUuids),
			L3NetworkUUID:      types.StringValue(vip.L3NetworkUUID),
			IP:                 types.StringValue(vip.Ip),
			Gateway:            types.StringValue(vip.Gateway),
			Netmask:            types.StringValue(vip.Netmask),
			UseFor:             types.StringValue(vip.UseFor),
		}

		state.VIPs = append(state.VIPs, VIPsState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *vipsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of vips and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching VIPs",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"vips": schema.ListNestedAttribute{
				Description: "List of VIP entries matching the specified filters",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier of the VIP.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the VIP",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the VIP",
						},
						"l3_network_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the L3 network associated with the VIP.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "State of the VIP (e.g., Enabled, Disabled)",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "IP address of the VIP.",
						},
						"gateway": schema.StringAttribute{
							Computed:    true,
							Description: "Gateway address of the VIP.",
						},
						"netmask": schema.StringAttribute{
							Computed:    true,
							Description: "Netmask of the VIP.",
						},
						"use_for": schema.StringAttribute{
							Computed:    true,
							Description: "The purpose or usage of the VIP.",
						},
						"peer_l3_network_uuids": schema.StringAttribute{
							Computed:    true,
							Description: "The UUIDs of peer L3 networks associated with the VIP (e.g., related to EIP binding).",
						},
					},
				},
			},
		},
	}

}
