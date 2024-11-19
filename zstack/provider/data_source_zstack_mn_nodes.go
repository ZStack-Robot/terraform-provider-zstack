// Copyright (c) ZStack.io, Inc.

package provider

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
	_ datasource.DataSource              = &mnNodeDataSource{}
	_ datasource.DataSourceWithConfigure = &mnNodeDataSource{}
)

func ZStackmnNodeDataSource() datasource.DataSource {
	return &mnNodeDataSource{}
}

type mnNodeDataSource struct {
	client *client.ZSClient
}

type mnNodeDataSourceModel struct {
	MN_Nodes []mnNodeModel `tfsdk:"mn_nodes"`
}

type mnNodeModel struct {
	Uuid     types.String `tfsdk:"uuid"`
	HostName types.String `tfsdk:"host_name"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *mnNodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mnNodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mnnodes"
}

func (d *mnNodeDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"mn_nodes": schema.ListNestedAttribute{
				Description: "Fetches a list of Management Nodes and their associated attributes from the ZStack environment",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the Management Node.",
						},
						"host_name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the Management Node",
						},
					},
				},
			},
		},
	}
}

func (d *mnNodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state mnNodeDataSourceModel
	//var state clusterModel
	mn_nodes, err := d.client.QueryManagementNode(param.NewQueryParam())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Management Nodes",
			err.Error(),
		)
		return
	}

	//map query Management Nodes
	for _, node := range mn_nodes {
		mnNodeState := mnNodeModel{
			Uuid:     types.StringValue(node.UUID),
			HostName: types.StringValue(node.HostName),
		}
		state.MN_Nodes = append(state.MN_Nodes, mnNodeState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
