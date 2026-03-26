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
	_ datasource.DataSource              = &loadBalancerListenerDataSource{}
	_ datasource.DataSourceWithConfigure = &loadBalancerListenerDataSource{}
)

type loadBalancerListenerDataSource struct {
	client *client.ZSClient
}

type loadBalancerListenerDataSourceModel struct {
	Name                   types.String                     `tfsdk:"name"`
	NamePattern            types.String                     `tfsdk:"name_pattern"`
	Filter                 []Filter                         `tfsdk:"filter"`
	LoadBalancerListeners  []loadBalancerListenersModel     `tfsdk:"load_balancer_listeners"`
}

type loadBalancerListenersModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	LoadBalancerUuid   types.String `tfsdk:"load_balancer_uuid"`
	Protocol           types.String `tfsdk:"protocol"`
	LoadBalancerPort   types.Int64  `tfsdk:"load_balancer_port"`
	InstancePort       types.Int64  `tfsdk:"instance_port"`
	SecurityPolicyType types.String `tfsdk:"security_policy_type"`
	ServerGroupUuid    types.String `tfsdk:"server_group_uuid"`
}

func ZStackLoadBalancerListenerDataSource() datasource.DataSource {
	return &loadBalancerListenerDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *loadBalancerListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *loadBalancerListenerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_listeners"
}

// Schema implements datasource.DataSource.
func (d *loadBalancerListenerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of load balancer listeners and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of load balancer listeners and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching load balancer listeners.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"load_balancer_listeners": schema.ListNestedAttribute{
				Description: "List of load balancer listeners matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the load balancer listener.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the load balancer listener.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the load balancer listener.",
							Computed:    true,
						},
						"load_balancer_uuid": schema.StringAttribute{
							Description: "UUID of the parent load balancer.",
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "Protocol of the listener (tcp, udp, http, https).",
							Computed:    true,
						},
						"load_balancer_port": schema.Int64Attribute{
							Description: "Frontend port on the load balancer.",
							Computed:    true,
						},
						"instance_port": schema.Int64Attribute{
							Description: "Backend port on the instances.",
							Computed:    true,
						},
						"security_policy_type": schema.StringAttribute{
							Description: "Security policy type for HTTPS listeners.",
							Computed:    true,
						},
						"server_group_uuid": schema.StringAttribute{
							Description: "UUID of the associated server group.",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by protocol, use `name = \"protocol\"` and `values = [\"tcp\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., protocol, load_balancer_uuid).",
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
func (d *loadBalancerListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loadBalancerListenerDataSourceModel

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

	listeners, err := d.client.QueryLoadBalancerListener(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Load Balancer Listeners",
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

	filteredListeners, filterDiags := utils.FilterResource(ctx, listeners, filters, "load_balancer_listener")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.LoadBalancerListeners = []loadBalancerListenersModel{}

	for _, l := range filteredListeners {
		listenerState := loadBalancerListenersModel{
			Uuid:               types.StringValue(l.UUID),
			Name:               types.StringValue(l.Name),
			Description:        stringValueOrNull(l.Description),
			LoadBalancerUuid:   types.StringValue(l.LoadBalancerUuid),
			Protocol:           stringValueOrNull(l.Protocol),
			LoadBalancerPort:   types.Int64Value(int64(l.LoadBalancerPort)),
			InstancePort:       types.Int64Value(int64(l.InstancePort)),
			SecurityPolicyType: stringValueOrNull(l.SecurityPolicyType),
			ServerGroupUuid:    stringValueOrNull(l.ServerGroupUuid),
		}
		state.LoadBalancerListeners = append(state.LoadBalancerListeners, listenerState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
