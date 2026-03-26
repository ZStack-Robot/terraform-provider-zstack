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
	_ datasource.DataSource              = &portForwardingRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &portForwardingRuleDataSource{}
)

type portForwardingRuleDataSource struct {
	client *client.ZSClient
}

type portForwardingRuleDataSourceModel struct {
	Name                types.String                    `tfsdk:"name"`
	NamePattern         types.String                    `tfsdk:"name_pattern"`
	Filter              []Filter                        `tfsdk:"filter"`
	PortForwardingRules []portForwardingRulesModel      `tfsdk:"port_forwarding_rules"`
}

type portForwardingRulesModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	VipUuid          types.String `tfsdk:"vip_uuid"`
	VipIp            types.String `tfsdk:"vip_ip"`
	GuestIp          types.String `tfsdk:"guest_ip"`
	VipPortStart     types.Int64  `tfsdk:"vip_port_start"`
	VipPortEnd       types.Int64  `tfsdk:"vip_port_end"`
	PrivatePortStart types.Int64  `tfsdk:"private_port_start"`
	PrivatePortEnd   types.Int64  `tfsdk:"private_port_end"`
	VmNicUuid        types.String `tfsdk:"vm_nic_uuid"`
	ProtocolType     types.String `tfsdk:"protocol_type"`
	State            types.String `tfsdk:"state"`
	AllowedCidr      types.String `tfsdk:"allowed_cidr"`
}

func ZStackPortForwardingRuleDataSource() datasource.DataSource {
	return &portForwardingRuleDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *portForwardingRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *portForwardingRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forwarding_rules"
}

// Schema implements datasource.DataSource.
func (d *portForwardingRuleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of port forwarding rules and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of port forwarding rules and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching port forwarding rules.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"port_forwarding_rules": schema.ListNestedAttribute{
				Description: "List of port forwarding rules matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the port forwarding rule.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the port forwarding rule.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the port forwarding rule.",
							Computed:    true,
						},
						"vip_uuid": schema.StringAttribute{
							Description: "UUID of the VIP used for port forwarding.",
							Computed:    true,
						},
						"vip_ip": schema.StringAttribute{
							Description: "IP address of the VIP.",
							Computed:    true,
						},
						"guest_ip": schema.StringAttribute{
							Description: "IP address of the guest VM.",
							Computed:    true,
						},
						"vip_port_start": schema.Int64Attribute{
							Description: "Start port on the VIP side.",
							Computed:    true,
						},
						"vip_port_end": schema.Int64Attribute{
							Description: "End port on the VIP side.",
							Computed:    true,
						},
						"private_port_start": schema.Int64Attribute{
							Description: "Start port on the private (VM) side.",
							Computed:    true,
						},
						"private_port_end": schema.Int64Attribute{
							Description: "End port on the private (VM) side.",
							Computed:    true,
						},
						"vm_nic_uuid": schema.StringAttribute{
							Description: "UUID of the VM NIC attached to this rule.",
							Computed:    true,
						},
						"protocol_type": schema.StringAttribute{
							Description: "Protocol type (TCP or UDP).",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the port forwarding rule.",
							Computed:    true,
						},
						"allowed_cidr": schema.StringAttribute{
							Description: "CIDR block allowed to access this rule.",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by protocol, use `name = \"protocol_type\"` and `values = [\"TCP\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., protocol_type, state).",
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
func (d *portForwardingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state portForwardingRuleDataSourceModel

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

	rules, err := d.client.QueryPortForwardingRule(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Port Forwarding Rules",
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

	filteredRules, filterDiags := utils.FilterResource(ctx, rules, filters, "port_forwarding_rule")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.PortForwardingRules = []portForwardingRulesModel{}

	for _, rule := range filteredRules {
		ruleState := portForwardingRulesModel{
			Uuid:             types.StringValue(rule.UUID),
			Name:             types.StringValue(rule.Name),
			Description:      stringValueOrNull(rule.Description),
			VipUuid:          types.StringValue(rule.VipUuid),
			VipIp:            stringValueOrNull(rule.VipIp),
			GuestIp:          stringValueOrNull(rule.GuestIp),
			VipPortStart:     types.Int64Value(int64(rule.VipPortStart)),
			VipPortEnd:       types.Int64Value(int64(rule.VipPortEnd)),
			PrivatePortStart: types.Int64Value(int64(rule.PrivatePortStart)),
			PrivatePortEnd:   types.Int64Value(int64(rule.PrivatePortEnd)),
			VmNicUuid:        stringValueOrNull(rule.VmNicUuid),
			ProtocolType:     types.StringValue(rule.ProtocolType),
			State:            stringValueOrNull(rule.State),
			AllowedCidr:      stringValueOrNull(rule.AllowedCidr),
		}
		state.PortForwardingRules = append(state.PortForwardingRules, ruleState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
