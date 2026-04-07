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
	_ datasource.DataSource              = &subnetIpRangeDataSource{}
	_ datasource.DataSourceWithConfigure = &subnetIpRangeDataSource{}
)

type subnetIpRangeDataSource struct {
	client *client.ZSClient
}

type subnetIpRangeItemModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	Description   types.String `tfsdk:"description"`
	StartIp       types.String `tfsdk:"start_ip"`
	EndIp         types.String `tfsdk:"end_ip"`
	Netmask       types.String `tfsdk:"netmask"`
	Gateway       types.String `tfsdk:"gateway"`
	NetworkCidr   types.String `tfsdk:"network_cidr"`
	IpVersion     types.Int64  `tfsdk:"ip_version"`
	AddressMode   types.String `tfsdk:"address_mode"`
	PrefixLen     types.Int64  `tfsdk:"prefix_len"`
	IpRangeType   types.String `tfsdk:"ip_range_type"`
}

type subnetIpRangeDataSourceModel struct {
	Name           types.String             `tfsdk:"name"`
	NamePattern    types.String             `tfsdk:"name_pattern"`
	SubnetIpRanges []subnetIpRangeItemModel `tfsdk:"subnet_ip_ranges"`
	Filter         []Filter                 `tfsdk:"filter"`
}

func ZStackSubnetIpRangeDataSource() datasource.DataSource {
	return &subnetIpRangeDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *subnetIpRangeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the ZStack Provider developer. ", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Metadata implements datasource.DataSource.
func (d *subnetIpRangeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet_ip_ranges"
}

// Read implements datasource.DataSource.
func (d *subnetIpRangeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state subnetIpRangeDataSourceModel
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

	ipRanges, err := d.client.QueryIpRange(&params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Subnet IP Ranges",
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

	filterIpRanges, filterDiags := utils.FilterResource(ctx, ipRanges, filters, "subnet_ip_range")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, ipRange := range filterIpRanges {
		ipRangeState := subnetIpRangeItemModel{
			Uuid:          types.StringValue(ipRange.UUID),
			Name:          types.StringValue(ipRange.Name),
			L3NetworkUuid: types.StringValue(ipRange.L3NetworkUuid),
			Description:   types.StringValue(ipRange.Description),
			StartIp:       types.StringValue(ipRange.StartIp),
			EndIp:         types.StringValue(ipRange.EndIp),
			Netmask:       types.StringValue(ipRange.Netmask),
			Gateway:       types.StringValue(ipRange.Gateway),
			NetworkCidr:   types.StringValue(ipRange.NetworkCidr),
			IpVersion:     types.Int64Value(int64(ipRange.IpVersion)),
			AddressMode:   types.StringValue(ipRange.AddressMode),
			PrefixLen:     types.Int64Value(int64(ipRange.PrefixLen)),
			IpRangeType:   types.StringValue(ipRange.IpRangeType),
		}
		state.SubnetIpRanges = append(state.SubnetIpRanges, ipRangeState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *subnetIpRangeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of subnet IP ranges and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching subnet IP ranges",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"subnet_ip_ranges": schema.ListNestedAttribute{
				Description: "List of Subnet IP Ranges",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the subnet IP range",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the subnet IP range",
							Computed:    true,
						},
						"l3_network_uuid": schema.StringAttribute{
							Description: "UUID of the L3 network",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the subnet IP range",
							Computed:    true,
						},
						"start_ip": schema.StringAttribute{
							Description: "Start IP address of the range",
							Computed:    true,
						},
						"end_ip": schema.StringAttribute{
							Description: "End IP address of the range",
							Computed:    true,
						},
						"netmask": schema.StringAttribute{
							Description: "Netmask for the IP range",
							Computed:    true,
						},
						"gateway": schema.StringAttribute{
							Description: "Gateway IP address",
							Computed:    true,
						},
						"network_cidr": schema.StringAttribute{
							Description: "Network CIDR notation",
							Computed:    true,
						},
						"ip_version": schema.Int64Attribute{
							Description: "IP version (4 for IPv4 or 6 for IPv6)",
							Computed:    true,
						},
						"address_mode": schema.StringAttribute{
							Description: "Address mode",
							Computed:    true,
						},
						"prefix_len": schema.Int64Attribute{
							Description: "Prefix length for IPv6",
							Computed:    true,
						},
						"ip_range_type": schema.StringAttribute{
							Description: "Type of IP range",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by IP version, use `name = \"ipVersion\"` and `values = [\"4\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., ipVersion, l3NetworkUuid).",
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
