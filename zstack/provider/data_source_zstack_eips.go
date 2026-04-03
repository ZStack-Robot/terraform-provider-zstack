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
	_ datasource.DataSource              = &eipDataSource{}
	_ datasource.DataSourceWithConfigure = &eipDataSource{}
)

type eipDataSource struct {
	client *client.ZSClient
}

type eipItemModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VmNicUuid   types.String `tfsdk:"vm_nic_uuid"`
	VipUuid     types.String `tfsdk:"vip_uuid"`
	State       types.String `tfsdk:"state"`
	VipIp       types.String `tfsdk:"vip_ip"`
	GuestIp     types.String `tfsdk:"guest_ip"`
}

type eipDataSourceModel struct {
	Name        types.String   `tfsdk:"name"`
	NamePattern types.String   `tfsdk:"name_pattern"`
	Eips        []eipItemModel `tfsdk:"eips"`
	Filter      []Filter       `tfsdk:"filter"`
}

func ZStackEipDataSource() datasource.DataSource {
	return &eipDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *eipDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *eipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eips"
}

// Read implements datasource.DataSource.
func (d *eipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state eipDataSourceModel
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

	eips, err := d.client.QueryEip(&params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack EIPs",
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

	filterEips, filterDiags := utils.FilterResource(ctx, eips, filters, "eip")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, eip := range filterEips {
		eipState := eipItemModel{
			Uuid:        types.StringValue(eip.UUID),
			Name:        types.StringValue(eip.Name),
			Description: types.StringValue(eip.Description),
			VmNicUuid:   types.StringValue(eip.VmNicUuid),
			VipUuid:     types.StringValue(eip.VipUuid),
			State:       types.StringValue(eip.State),
			VipIp:       types.StringValue(eip.VipIp),
			GuestIp:     types.StringValue(eip.GuestIp),
		}
		state.Eips = append(state.Eips, eipState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *eipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of elastic IPs and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching elastic IPs",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"eips": schema.ListNestedAttribute{
				Description: "List of Elastic IPs",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the elastic IP",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the elastic IP",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the elastic IP",
							Computed:    true,
						},
						"vm_nic_uuid": schema.StringAttribute{
							Description: "UUID of the VM NIC attached to this elastic IP",
							Computed:    true,
						},
						"vip_uuid": schema.StringAttribute{
							Description: "UUID of the VIP associated with this elastic IP",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the elastic IP",
							Computed:    true,
						},
						"vip_ip": schema.StringAttribute{
							Description: "Virtual IP address",
							Computed:    true,
						},
						"guest_ip": schema.StringAttribute{
							Description: "Guest IP address",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by state, use `name = \"state\"` and `values = [\"Enabled\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., state, name).",
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
