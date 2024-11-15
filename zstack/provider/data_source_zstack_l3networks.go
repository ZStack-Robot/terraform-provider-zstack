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
	_ datasource.DataSource              = &l3NetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &l3NetworkDataSource{}
)

type l3NetworkDataSource struct {
	client *client.ZSClient
}

type l3NetworkDataSourceModel struct {
	Name_regex types.String      `tfsdk:"name_regex"`
	L3networks []l3networksModel `tfsdk:"l3networks"`
}
type l3networksModel struct {
	Name     types.String   `tfsdk:"name"`
	Uuid     types.String   `tfsdk:"uuid"`
	Category types.String   `tfsdk:"category"`
	Dns      []dnsModel     `tfsdk:"dns"`
	Iprange  []ipRangeModel `tfsdk:"ip_range"`
	FreeIps  []freeIpModel  `tfsdk:"free_ips"`
}

type dnsModel struct {
	Dns types.String `tfsdk:"dns_model"`
}

type ipRangeModel struct {
	Name        types.String `tfsdk:"ip_range_name"`
	StartIp     types.String `tfsdk:"start_ip"`
	EndIp       types.String `tfsdk:"end_ip"`
	Netmask     types.String `tfsdk:"netmask"`
	Gateway     types.String `tfsdk:"gateway"`
	NetworkCidr types.String `tfsdk:"cidr"`
}
type freeIpModel struct {
	IpRangeUuid string `tfsdk:"ip_range_uuid"`
	Ip          string `tfsdk:"ip"`
	Netmask     string `tfsdk:"netmask"`
	Gateway     string `tfsdk:"gateway"`
}

func ZStackl3NetworkDataSource() datasource.DataSource {
	return &l3NetworkDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *l3NetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *l3NetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l3networks"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *l3NetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state l3NetworkDataSourceModel
	//var L3state l3networksModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	//Create query parameters based on name_regex
	//uuid := L3state.Uuid.ValueString()
	//name_regex := state.Name_regex
	params := param.NewQueryParam()

	if !state.Name_regex.IsNull() {
		params.AddQ("name=" + state.Name_regex.ValueString())
	}

	//Query L3 networks with name filtering
	l3networks, err := d.client.QueryL3Network(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack L3Networks ",
			err.Error(),
		)
		return
	}

	// Process each L3 network in the result
	//l3freeIps, err := d.client.GetFreeIp(uuid, param.QueryParam{})
	for _, l3network := range l3networks {
		// Query free IPs for the current L3 network UUID
		l3freeIps, err := d.client.GetFreeIp(l3network.UUID, param.QueryParam{})
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Fetch Free IPs for L3 Network",
				fmt.Sprintf("Failed to retrieve free IPs for L3 network with UUID %s: %s", l3network.UUID, err.Error()),
			)
			return
		}

		// Build the L3 network model with nested attributes
		l3networkState := l3networksModel{
			Name:     types.StringValue(l3network.Name),
			Uuid:     types.StringValue(l3network.UUID),
			Category: types.StringValue(l3network.Category),
			Dns:      make([]dnsModel, len(l3network.Dns)),
			Iprange:  make([]ipRangeModel, len(l3network.IpRanges)),
			FreeIps:  make([]freeIpModel, len(l3freeIps)),
		}

		// Populate DNS information
		for i, dns := range l3network.Dns {
			l3networkState.Dns[i] = dnsModel{
				Dns: types.StringValue(dns),
			}
		}

		// Populate IP range information
		for i, iprange := range l3network.IpRanges {
			l3networkState.Iprange[i] = ipRangeModel{
				Name:        types.StringValue(iprange.Name),
				StartIp:     types.StringValue(iprange.StartIp),
				EndIp:       types.StringValue(iprange.EndIp),
				Netmask:     types.StringValue(iprange.Netmask),
				Gateway:     types.StringValue(iprange.Gateway),
				NetworkCidr: types.StringValue(iprange.NetworkCidr),
			}
		}

		// Populate free IP information
		for i, freeIp := range l3freeIps {
			l3networkState.FreeIps[i] = freeIpModel{
				IpRangeUuid: freeIp.IpRangeUuid,
				Ip:          freeIp.Ip,
				Netmask:     freeIp.Netmask,
				Gateway:     freeIp.Gateway,
			}
		}

		state.L3networks = append(state.L3networks, l3networkState)
	}

	// Set the final state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *l3NetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name_regex": schema.StringAttribute{
				Description: "Regex pattern to filter L3 networks by name.",
				Optional:    true,
			},
			"l3networks": schema.ListNestedAttribute{
				Description: "List of L3 networks matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the L3 network",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the L3 network.",
						},
						"category": schema.StringAttribute{
							Computed:    true,
							Description: "Category of the L3 network.",
						},
						"dns": schema.ListNestedAttribute{
							Description: "List of DNS servers for the L3 network.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"dns_model": schema.StringAttribute{
										Description: "DNS server address.",
										Computed:    true,
									},
								},
							},
						},
						"ip_range": schema.ListNestedAttribute{
							Description: "List of IP ranges in the L3 network.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip_range_name": schema.StringAttribute{
										Description: "Name of the IP range.",
										Computed:    true,
									},
									"start_ip": schema.StringAttribute{
										Description: "Starting IP address in the range.",
										Computed:    true,
									},
									"end_ip": schema.StringAttribute{
										Description: "Ending IP address in the range.",
										Computed:    true,
									},
									"netmask": schema.StringAttribute{
										Description: "Netmask of the IP range.",
										Computed:    true,
									},
									"gateway": schema.StringAttribute{
										Description: "Gateway for the IP range.",
										Computed:    true,
									},
									"cidr": schema.StringAttribute{
										Description: "CIDR notation for the IP range.",
										Computed:    true,
									},
								},
							},
						},
						"free_ips": schema.ListNestedAttribute{
							Description: "List of free IPs available in the L3 network.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip_range_uuid": schema.StringAttribute{
										Description: "UUID of the IP range containing the free IP.",
										Computed:    true,
									},
									"ip": schema.StringAttribute{
										Description: "Free IP address.",
										Computed:    true,
									},
									"netmask": schema.StringAttribute{
										Description: "Netmask for the free IP.",
										Computed:    true,
									},
									"gateway": schema.StringAttribute{
										Description: "Gateway for the free IP.",
										Computed:    true,
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
