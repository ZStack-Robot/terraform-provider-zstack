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
	_ datasource.DataSource              = &loadBalancerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerDataSource{}
)

type loadBalancerDataSource struct {
	client *client.ZSClient
}

type loadBalancerDataSourceModel struct {
	Name          types.String        `tfsdk:"name"`
	NamePattern   types.String        `tfsdk:"name_pattern"`
	Filter        []Filter            `tfsdk:"filter"`
	LoadBalancers []loadBalancersModel `tfsdk:"load_balancers"`
}

type loadBalancersModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	VipUuid         types.String `tfsdk:"vip_uuid"`
	State           types.String `tfsdk:"state"`
	Type            types.String `tfsdk:"type"`
	ServerGroupUuid types.String `tfsdk:"server_group_uuid"`
}

func ZStackLoadBalancerDataSource() datasource.DataSource {
	return &loadBalancerDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *loadBalancerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadBalancerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancers"
}

// Schema implements datasource.DataSource.
func (d *loadBalancerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of load balancers and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of load balancers and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching load balancers.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"load_balancers": schema.ListNestedAttribute{
				Description: "List of load balancers matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the load balancer.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the load balancer.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the load balancer.",
							Computed:    true,
						},
						"vip_uuid": schema.StringAttribute{
							Description: "UUID of the VIP bound to the load balancer.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the load balancer (Enabled, Disabled).",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the load balancer.",
							Computed:    true,
						},
						"server_group_uuid": schema.StringAttribute{
							Description: "UUID of the default server group.",
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
							Description: "Name of the field to filter by (e.g., state, type).",
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
func (d *loadBalancerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadBalancerDataSourceModel

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

	lbs, err := d.client.QueryLoadBalancer(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Load Balancers",
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

	filteredLbs, filterDiags := utils.FilterResource(ctx, lbs, filters, "load_balancer")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.LoadBalancers = []loadBalancersModel{}

	for _, lb := range filteredLbs {
		lbState := loadBalancersModel{
			Uuid:            types.StringValue(lb.UUID),
			Name:            types.StringValue(lb.Name),
			Description:     stringValueOrNull(lb.Description),
			VipUuid:         types.StringValue(lb.VipUuid),
			State:           stringValueOrNull(lb.State),
			Type:            stringValueOrNull(lb.Type),
			ServerGroupUuid: stringValueOrNull(lb.ServerGroupUuid),
		}
		state.LoadBalancers = append(state.LoadBalancers, lbState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
