// Copyright (c) HashiCorp, Inc.

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
	Iprange  []iprangeModel `tfsdk:"iprange"`
}

type dnsModel struct {
	Dns types.String `tfsdk:"dnsmodel"`
}

type iprangeModel struct {
	Name        types.String `tfsdk:"iprangename"`
	StartIp     types.String `tfsdk:"startip"`
	EndIp       types.String `tfsdk:"endip"`
	Netmask     types.String `tfsdk:"netmask"`
	Gateway     types.String `tfsdk:"gateway"`
	NetworkCidr types.String `tfsdk:"cidr"`
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
	resp.TypeName = req.ProviderTypeName + "_l3network"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *l3NetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state l3NetworkDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	name_regex := state.Name_regex
	params := param.NewQueryParam()

	if !name_regex.IsNull() {
		params.AddQ("name=" + name_regex.ValueString())
	}

	l3networks, err := d.client.QueryL3Network(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack L3Networks ",
			err.Error(),
		)

		return
	}
	for _, l3network := range l3networks {
		l3networkState := l3networksModel{
			Name:     types.StringValue(l3network.Name),
			Uuid:     types.StringValue(l3network.UUID),
			Category: types.StringValue(l3network.Category),
			Dns:      make([]dnsModel, len(l3network.Dns)),
			Iprange:  make([]iprangeModel, len(l3network.IpRanges)),
		}

		for i, dns := range l3network.Dns {
			l3networkState.Dns[i] = dnsModel{
				Dns: types.StringValue(dns),
			}
		}
		for i, iprange := range l3network.IpRanges {
			l3networkState.Iprange[i] = iprangeModel{
				Name:        types.StringValue(iprange.Name),
				StartIp:     types.StringValue(iprange.StartIp),
				EndIp:       types.StringValue(iprange.EndIp),
				Netmask:     types.StringValue(iprange.Netmask),
				Gateway:     types.StringValue(iprange.Gateway),
				NetworkCidr: types.StringValue(iprange.NetworkCidr),
			}
		}

		state.L3networks = append(state.L3networks, l3networkState)
	}

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
				Description: "name_regex for Search and filter L3 network",
				Optional:    true,
			},
			"l3networks": schema.ListNestedAttribute{
				Description: "List of L3 networks",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the L3 network",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"category": schema.StringAttribute{
							Computed: true,
						},
						"dns": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"dnsmodel": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"iprange": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"iprangename": schema.StringAttribute{
										Computed: true,
									},
									"startip": schema.StringAttribute{
										Computed: true,
									},
									"endip": schema.StringAttribute{
										Computed: true,
									},
									"netmask": schema.StringAttribute{
										Computed: true,
									},
									"gateway": schema.StringAttribute{
										Computed: true,
									},
									"cidr": schema.StringAttribute{
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
