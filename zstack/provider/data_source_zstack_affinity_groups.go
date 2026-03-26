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
	_ datasource.DataSource              = &affinityGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &affinityGroupDataSource{}
)

func ZStackAffinityGroupDataSource() datasource.DataSource {
	return &affinityGroupDataSource{}
}

type affinityGroupItem struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Policy      types.String `tfsdk:"policy"`
	Type        types.String `tfsdk:"type"`
	ZoneUuid    types.String `tfsdk:"zone_uuid"`
	State       types.String `tfsdk:"state"`
}

type affinityGroupDataSourceModel struct {
	Name           types.String        `tfsdk:"name"`
	NamePattern    types.String        `tfsdk:"name_pattern"`
	Filter         []Filter            `tfsdk:"filter"`
	AffinityGroups []affinityGroupItem `tfsdk:"affinity_groups"`
}

type affinityGroupDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *affinityGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *affinityGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_affinity_groups"
}

// Read implements datasource.DataSource.
func (d *affinityGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state affinityGroupDataSourceModel
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

	affinityGroups, err := d.client.QueryAffinityGroup(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack Affinity Groups",
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

	filterAffinityGroups, filterDiags := utils.FilterResource(ctx, affinityGroups, filters, "affinity_group")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, ag := range filterAffinityGroups {
		state.AffinityGroups = append(state.AffinityGroups, affinityGroupItem{
			Uuid:        types.StringValue(ag.UUID),
			Name:        types.StringValue(ag.Name),
			Description: types.StringValue(ag.Description),
			Policy:      types.StringValue(ag.Policy),
			Type:        types.StringValue(ag.Type),
			ZoneUuid:    types.StringValue(ag.ZoneUuid),
			State:       types.StringValue(ag.State),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSource.
func (d *affinityGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack Affinity Groups by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying an affinity group.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching affinity group names. Use % or _ like SQL.",
				Optional:    true,
			},
			"affinity_groups": schema.ListNestedAttribute{
				Description: "List of matched affinity groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the affinity group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the affinity group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the affinity group.",
							Computed:    true,
						},
						"policy": schema.StringAttribute{
							Description: "Placement policy (antiSoft, antiHard, proSoft, proHard).",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the affinity group.",
							Computed:    true,
						},
						"zone_uuid": schema.StringAttribute{
							Description: "UUID of the zone.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the affinity group (Enabled, Disabled).",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by field values.",
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
