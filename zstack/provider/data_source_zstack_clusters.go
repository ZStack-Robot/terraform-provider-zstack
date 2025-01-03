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
	_ datasource.DataSource              = &clusterDataSource{}
	_ datasource.DataSourceWithConfigure = &clusterDataSource{}
)

func ZStackClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	client *client.ZSClient
}

type clusterDataSourceModel struct {
	Name        types.String   `tfsdk:"name"`
	NamePattern types.String   `tfsdk:"name_pattern"`
	Filter      []Filter       `tfsdk:"filter"`
	Clusters    []clusterModel `tfsdk:"clusters"`
}

type clusterModel struct {
	Name           types.String `tfsdk:"name"`
	HypervisorType types.String `tfsdk:"hypervisor_type"`
	State          types.String `tfsdk:"state"`
	Type           types.String `tfsdk:"type"`
	Uuid           types.String `tfsdk:"uuid"`
	ZoneUuid       types.String `tfsdk:"zone_uuid"`
	Architecture   types.String `tfsdk:"architecture"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *clusterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clusters"
}

func (d *clusterDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of clusters and their associated attributes.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching Cluster",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			/*
				"filter": schema.MapAttribute{
					Description: "Key-value pairs to filter Clusters. For example, to filter by CPU Architecture, use `Architecture = \"x86_64\"`.",
					Optional:    true,
					ElementType: types.StringType,
				},
			*/
			"clusters": schema.ListNestedAttribute{
				Description: "List of clusters matching the specified filters",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the cluster",
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID identifier of the cluster",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the zone to which the cluster belongs",
						},
						"hypervisor_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of hypervisor used by the cluster (e.g., KVM, ESXi)",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "ype of the cluster",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "State of the cluster (e.g., Enabled, Disabled)",
						},
						"architecture": schema.StringAttribute{
							Computed:    true,
							Description: "Architecture of the cluster",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by status, use `name = \"status\"` and `values = [\"Ready\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., status, state).",
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

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clusterDataSourceModel
	//var state clusterModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	//name_regex := state.Name_pattern
	params := param.NewQueryParam()

	// 优先检查 `name` 精确查询
	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	//images, err := d.client.QueryImage(params)

	clusters, err := d.client.QueryCluster(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Clusters",
			err.Error(),
		)
		return
	}
	/*
	   filters := make(map[string]string)

	   	if !state.Filter.IsNull() {
	   		diags := state.Filter.ElementsAs(ctx, &filters, false)
	   		resp.Diagnostics.Append(diags...)
	   		if resp.Diagnostics.HasError() {
	   			return
	   		}
	   	}

	   filterClusters, filterDiags := utils.FilterResource(ctx, clusters, filters)
	   resp.Diagnostics.Append(filterDiags...)

	   	if resp.Diagnostics.HasError() {
	   		return
	   	}
	*/
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

	// 过滤资源
	filterClusters, filterDiags := utils.FilterResource(ctx, clusters, filters, "cluster")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//map query clusters body to mode
	for _, cluster := range filterClusters {
		clusterState := clusterModel{
			HypervisorType: types.StringValue(cluster.HypervisorType),
			State:          types.StringValue(cluster.State),
			Type:           types.StringValue(cluster.Type),
			Uuid:           types.StringValue(cluster.Uuid),
			ZoneUuid:       types.StringValue(cluster.ZoneUuid),
			Name:           types.StringValue(cluster.Name),
			Architecture:   types.StringValue(cluster.Architecture),
		}

		state.Clusters = append(state.Clusters, clusterState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
