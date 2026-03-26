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
	_ datasource.DataSource              = &autoScalingGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &autoScalingGroupDataSource{}
)

type autoScalingGroupDataSource struct {
	client *client.ZSClient
}

type autoScalingGroupDataSourceModel struct {
	Name              types.String             `tfsdk:"name"`
	NamePattern       types.String             `tfsdk:"name_pattern"`
	Filter            []Filter                 `tfsdk:"filter"`
	AutoScalingGroups []autoScalingGroupsModel `tfsdk:"auto_scaling_groups"`
}

type autoScalingGroupsModel struct {
	Uuid                types.String `tfsdk:"uuid"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	ScalingResourceType types.String `tfsdk:"scaling_resource_type"`
	State               types.String `tfsdk:"state"`
	DefaultCooldown     types.Int64  `tfsdk:"default_cooldown"`
	MinResourceSize     types.Int64  `tfsdk:"min_resource_size"`
	MaxResourceSize     types.Int64  `tfsdk:"max_resource_size"`
	RemovalPolicy       types.String `tfsdk:"removal_policy"`
}

func ZStackAutoScalingGroupDataSource() datasource.DataSource {
	return &autoScalingGroupDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *autoScalingGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *autoScalingGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_auto_scaling_groups"
}

// Schema implements datasource.DataSource.
func (d *autoScalingGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of auto scaling groups and their associated attributes from the ZStack environment.",
		MarkdownDescription: "Fetches a list of auto scaling groups and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching auto scaling groups.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"auto_scaling_groups": schema.ListNestedAttribute{
				Description: "List of auto scaling groups matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the auto scaling group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the auto scaling group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the auto scaling group.",
							Computed:    true,
						},
						"scaling_resource_type": schema.StringAttribute{
							Description: "Type of resource being scaled (e.g., VmInstance).",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the auto scaling group (Enabled, Disabled).",
							Computed:    true,
						},
						"default_cooldown": schema.Int64Attribute{
							Description: "Default cooldown period in seconds.",
							Computed:    true,
						},
						"min_resource_size": schema.Int64Attribute{
							Description: "Minimum number of instances.",
							Computed:    true,
						},
						"max_resource_size": schema.Int64Attribute{
							Description: "Maximum number of instances.",
							Computed:    true,
						},
						"removal_policy": schema.StringAttribute{
							Description: "Policy for removing instances when scaling in.",
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
							Description: "Name of the field to filter by (e.g., state, scaling_resource_type).",
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
func (d *autoScalingGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state autoScalingGroupDataSourceModel

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

	groups, err := d.client.QueryAutoScalingGroup(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Auto Scaling Groups",
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

	filteredGroups, filterDiags := utils.FilterResource(ctx, groups, filters, "auto_scaling_group")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AutoScalingGroups = []autoScalingGroupsModel{}

	for _, g := range filteredGroups {
		groupState := autoScalingGroupsModel{
			Uuid:                types.StringValue(g.UUID),
			Name:                types.StringValue(g.Name),
			Description:         stringValueOrNull(g.Description),
			ScalingResourceType: types.StringValue(g.ScalingResourceType),
			State:               stringValueOrNull(g.State),
			DefaultCooldown:     types.Int64Value(g.DefaultCooldown),
			MinResourceSize:     types.Int64Value(int64(g.MinResourceSize)),
			MaxResourceSize:     types.Int64Value(int64(g.MaxResourceSize)),
			RemovalPolicy:       stringValueOrNull(g.RemovalPolicy),
		}
		state.AutoScalingGroups = append(state.AutoScalingGroups, groupState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
