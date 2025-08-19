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
	_ datasource.DataSource              = &networkingSecGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &networkingSecGroupDataSource{}
)

func ZStackNetworkingSecGroupDataSource() datasource.DataSource {
	return &networkingSecGroupDataSource{}
}

type networkingSecGroup struct {
	Name                   types.String `tfsdk:"name"`
	Uuid                   types.String `tfsdk:"uuid"`
	Description            types.String `tfsdk:"description"`
	State                  types.String `tfsdk:"state"`
	AttachedL3NetworkUuids types.Set    `tfsdk:"attached_l3network_uuids"`
	VSwitchType            types.String `tfsdk:"vswitch_type"` // LinuxBridge, OvnDpdk
	Rules                  []rules      `tfsdk:"rules"`
}

type networkingSecGroupDataSourceModel struct {
	Name                types.String         `tfsdk:"name"`
	NamePattern         types.String         `tfsdk:"name_pattern"`
	Filter              []Filter             `tfsdk:"filter"`
	NetworkingSecGroups []networkingSecGroup `tfsdk:"networking_secgroups"`
}

type rules struct {
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

type networkingSecGroupDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *networkingSecGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networking_secgroups"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state networkingSecGroupDataSourceModel
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

	securityGroups, err := d.client.QuerySecurityGroup(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack Security Groups",
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

	filterSecurityGroups, filterDiags := utils.FilterResource(ctx, securityGroups, filters, "security_group")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, securitygroups := range filterSecurityGroups {
		networkingSecGroupState := networkingSecGroup{
			Name:        types.StringValue(securitygroups.Name),
			Uuid:        types.StringValue(securitygroups.UUID),
			Description: types.StringValue(securitygroups.Description),
			VSwitchType: types.StringValue(securitygroups.VSwitchType),
			State:       types.StringValue(securitygroups.State),
		}

		l3uuidSet, diags := types.SetValueFrom(ctx, types.StringType, securitygroups.AttachedL3NetworkUuids)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		networkingSecGroupState.AttachedL3NetworkUuids = l3uuidSet

		for _, rule := range securitygroups.Rules {
			networkingSecGroupState.Rules = append(networkingSecGroupState.Rules,
				rules{
					Uuid:              types.StringValue(rule.UUID),
					Type:              types.StringValue(rule.Type),
					IpVersion:         types.Int32Value(int32(rule.IpVersion)),
					Description:       types.StringValue(rule.Description),
					Priority:          types.Int32Value(int32(rule.Priority)),
					DstPortRange:      types.StringValue(rule.DstPortRange),
					SecurityGroupUuid: types.StringValue(rule.SecurityGroupUuid),
					SrcIpRange:        types.StringValue(rule.SrcIpRange),
					DstIpRange:        types.StringValue(rule.DstIpRange),
					Action:            types.StringValue(rule.Action),
					Protocol:          types.StringValue(rule.Protocol),
					State:             types.StringValue(rule.State),
				})
		}
		state.NetworkingSecGroups = append(state.NetworkingSecGroups, networkingSecGroupState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *networkingSecGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack Security Groups by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying a security group.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching security group names. Use % or _ like SQL.",
				Optional:    true,
			},
			"networking_secgroups": schema.ListNestedAttribute{
				Description: "List of matched security groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the security group.",
							Computed:    true,
						},
						"uuid": schema.StringAttribute{
							Description: "UUID of the security group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the security group.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the security group (Enabled, Disabled).",
							Computed:    true,
						},
						"attached_l3network_uuids": schema.SetAttribute{
							Description: "Set of L3 network UUIDs attached to the security group.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"vswitch_type": schema.StringAttribute{
							Description: "Type of the virtual switch (LinuxBridge, OvnDpdk).",
							Computed:    true,
						},
						"rules": schema.SetNestedAttribute{
							Description: "List of security group rules.",
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
										Description: "IP version (IPv4 or IPv6).",
										Computed:    true,
									},
									"priority": schema.Int32Attribute{
										Description: "Priority of the rule, default is 0.",
										Computed:    true,
									},
									"dst_port_range": schema.StringAttribute{
										Description: "Destination port range, e.g., '21, 80-443'",
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
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by fields in the security group, such as state or IP version.",
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
