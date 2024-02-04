// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package zstack

import (
	"context"
	"fmt"

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

func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	client *client.ZSClient
}

type clusterDataSourceModel struct {
	Clusters []clusterModel `tfsdk:"clusters"`
}

type clusterModel struct {
	Name           types.String `tfsdk:"name"`
	HypervisorType types.String `tfsdk:"hypervisortype"`
	State          types.String `tfsdk:"state"`
	Type           types.String `tfsdk:"type"`
	Uuid           types.String `tfsdk:"uuid"`
	ZoneUuid       types.String `tfsdk:"zoneuuid"`
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
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. jiajian.chi@zstack.io", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *clusterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zsclusters"
}

func (d *clusterDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"clusters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"zoneuuid": schema.StringAttribute{
							Computed: true,
						},
						"hypervisortype": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
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
	clusters, err := d.client.QueryCluster(param.NewQueryParam())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Clusters",
			err.Error(),
		)

		return
	}

	//map query clusters body to mode
	for _, cluster := range clusters {
		clusterState := clusterModel{
			HypervisorType: types.StringValue(cluster.HypervisorType),
			State:          types.StringValue(cluster.State),
			Type:           types.StringValue(cluster.Type),
			Uuid:           types.StringValue(cluster.Uuid),
			ZoneUuid:       types.StringValue(cluster.ZoneUuid),
			Name:           types.StringValue(cluster.Name),
		}

		state.Clusters = append(state.Clusters, clusterState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
