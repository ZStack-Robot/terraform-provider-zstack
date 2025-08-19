// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ datasource.DataSource              = &networkingSecGroupRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &networkingSecGroupRuleDataSource{}
)

func ZStackNetworkingSecGroupRuleDataSource() datasource.DataSource {
	return &networkingSecGroupRuleDataSource{}
}

type networkingSecGroupRuleDataSourceModel struct {
	Priority types.Int32  `tfsdk:"priority"`
	Filter   []Filter     `tfsdk:"filter"`
	Rules    []rulesModel `tfsdk:"rules"`
}

type rulesModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Action            types.String `tfsdk:"action"` // Allow, Deny
	Description       types.String `tfsdk:"description"`
	Protocol          types.String `tfsdk:"protocol"`
	IpVersion         types.Int32  `tfsdk:"ip_version"`
	Priority          types.Int32  `tfsdk:"priority"`
	DstPortRange      types.String `tfsdk:"dst_port_range"`
	SecurityGroupUuid types.String `tfsdk:"security_group_uuid"`
	SrcIpRange        types.String `tfsdk:"src_ip_range"` // CIDR
	DstIpRange        types.String `tfsdk:"dst_ip_range"` // CIDR
	Type              types.String `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
}

type networkingSecGroupRuleDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *networkingSecGroupRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networking_secgroup_rules"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networkingSecGroupRuleDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	if !state.Priority.IsNull() {
		params.AddQ("priority=" + fmt.Sprint(state.Priority.ValueInt32()))
	}

	securityGroupRules, err := d.client.QuerySecurityGroupRule(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack Security Groups Rules",
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

	securityGroupRules, filterDiags := utils.FilterResource(ctx, securityGroupRules, filters, "security_group_rule")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secGroupMap := make(map[string][]rulesModel)
	for _, rule := range securityGroupRules {
		ruleModel := rulesModel{
			Uuid:              types.StringValue(rule.UUID),
			Action:            types.StringValue(rule.Action),
			Protocol:          types.StringValue(rule.Protocol),
			IpVersion:         types.Int32Value(int32(rule.IpVersion)),
			Priority:          types.Int32Value(int32(rule.Priority)),
			SecurityGroupUuid: types.StringValue(rule.SecurityGroupUuid),
			Type:              types.StringValue(rule.Type),
			State:             types.StringValue(rule.State),
		}
		if rule.Description != "" {
			ruleModel.Description = types.StringValue(rule.Description)
		}
		if rule.DstPortRange != "" {
			ruleModel.DstPortRange = types.StringValue(rule.DstPortRange)
		}
		if rule.SrcIpRange != "" {
			ruleModel.SrcIpRange = types.StringValue(rule.SrcIpRange)
		}
		if rule.DstIpRange != "" {
			ruleModel.DstIpRange = types.StringValue(rule.DstIpRange)
		}
		secGroupMap[rule.SecurityGroupUuid] = append(secGroupMap[rule.SecurityGroupUuid], ruleModel)
	}

	state.Rules = state.Rules[:0]
	for _, rules := range secGroupMap {
		state.Rules = append(state.Rules, rules...)
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack Security Group Rules by exact name, name pattern, and optional filters.",
		Attributes: map[string]schema.Attribute{
			"priority": schema.Int32Attribute{
				Description: "Exact priority for querying security group rules.",
				Optional:    true,
			},
			"rules": schema.ListNestedAttribute{
				Description: "List of matched security group rules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the rule.",
							Computed:    true,
						},
						"action": schema.StringAttribute{
							Description: "Action of the rule (Allow, Deny).",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the rule.",
							Computed:    true,
						},
						"protocol": schema.StringAttribute{
							Description: "Protocol of the rule (TCP, UDP, ICMP, ALL).",
							Computed:    true,
						},
						"ip_version": schema.Int32Attribute{
							Description: "IP version (4 for IPv4, 6 for IPv6).",
							Computed:    true,
						},
						"priority": schema.Int32Attribute{
							Description: "Priority of the rule, default is 0.",
							Computed:    true,
						},
						"dst_port_range": schema.StringAttribute{
							Description: "Destination port range, e.g., '21, 80-443'.",
							Computed:    true,
						},
						"security_group_uuid": schema.StringAttribute{
							Description: "UUID of the security group this rule belongs to.",
							Computed:    true,
						},
						"src_ip_range": schema.StringAttribute{
							Description: "Source IP range in CIDR format, e.g., '192.168.1.0/24'.",
							Computed:    true,
						},
						"dst_ip_range": schema.StringAttribute{
							Description: "Destination IP range in CIDR format, e.g., '192.168.1.0/24'.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the rule (Ingress, Egress).",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the rule (Enabled, Disabled).",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by specific rule fields.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by.",
							Required:    true,
						},
						"values": schema.SetAttribute{
							Description: "List of values to match. Treated as OR conditions.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}
