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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ datasource.DataSource              = &reservedIpDataSource{}
	_ datasource.DataSourceWithConfigure = &reservedIpDataSource{}
)

type reservedIpDataSource struct {
	client *client.ZSClient
}

type reservedIpItemModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	Description   types.String `tfsdk:"description"`
	StartIp       types.String `tfsdk:"start_ip"`
	EndIp         types.String `tfsdk:"end_ip"`
	IpVersion     types.Int64  `tfsdk:"ip_version"`
}

type reservedIpDataSourceModel struct {
	Name        types.String          `tfsdk:"name"`
	NamePattern types.String          `tfsdk:"name_pattern"`
	ReservedIps []reservedIpItemModel `tfsdk:"reserved_ips"`
	Filter      []Filter              `tfsdk:"filter"`
}

func ZStackReservedIpDataSource() datasource.DataSource {
	return &reservedIpDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *reservedIpDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *reservedIpDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reserved_ips"
}

// Read implements datasource.DataSource.
func (d *reservedIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state reservedIpDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	queryStr := "query reservedIpRange"
	if !state.Name.IsNull() {
		queryStr += " where name='" + state.Name.ValueString() + "'"
	} else if !state.NamePattern.IsNull() {
		queryStr += " where name like '%" + state.NamePattern.ValueString() + "%'"
	}

	var reservedIps []view.ReservedIpRangeInventoryView
	_, err := d.client.Zql(ctx, queryStr, &reservedIps, "inventories")

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Reserved IPs",
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

	filterReservedIps, filterDiags := utils.FilterResource(ctx, reservedIps, filters, "reserved_ip")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, reservedIp := range filterReservedIps {
		reservedIpState := reservedIpItemModel{
			Uuid:          types.StringValue(reservedIp.UUID),
			Name:          types.StringValue(reservedIp.Name),
			L3NetworkUuid: types.StringValue(reservedIp.L3NetworkUuid),
			Description:   types.StringValue(reservedIp.Description),
			StartIp:       types.StringValue(reservedIp.StartIp),
			EndIp:         types.StringValue(reservedIp.EndIp),
			IpVersion:     types.Int64Value(int64(reservedIp.IpVersion)),
		}
		state.ReservedIps = append(state.ReservedIps, reservedIpState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *reservedIpDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of reserved IP ranges and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching reserved IP ranges",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"reserved_ips": schema.ListNestedAttribute{
				Description: "List of Reserved IP Ranges",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the reserved IP range",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the reserved IP range",
							Computed:    true,
						},
						"l3_network_uuid": schema.StringAttribute{
							Description: "UUID of the L3 network",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the reserved IP range",
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
						"ip_version": schema.Int64Attribute{
							Description: "IP version (4 for IPv4 or 6 for IPv6)",
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
