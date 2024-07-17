// Copyright (c) HashiCorp, Inc.

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
	_ datasource.DataSource              = &hostsDataSource{}
	_ datasource.DataSourceWithConfigure = &hostsDataSource{}
)

func ZStackHostsDataSource() datasource.DataSource {
	return &hostsDataSource{}
}

type hostsDataSource struct {
	client *client.ZSClient
}

type hostsDataSourceModel struct {
	Name_regex types.String `tfsdk:"name_regex"`
	Hosts      []hostsModel `tfsdk:"hosts"`
}

type hostsModel struct {
	Name         types.String `tfsdk:"name"`
	Architecture types.String `tfsdk:"architecture"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
	Type         types.String `tfsdk:"type"`
	Uuid         types.String `tfsdk:"id"`
	ZoneUuid     types.String `tfsdk:"zoneuuid"`
	ClusterUuid  types.String `tfsdk:"clusteruuid"`
	ManagementIp types.String `tfsdk:"managementop"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *hostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hostsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	name_regex := state.Name_regex
	params := param.NewQueryParam()

	if !name_regex.IsNull() {
		params.AddQ("name=" + name_regex.ValueString())
	}

	hosts, err := d.client.QueryHost(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Hosts ",
			err.Error(),
		)
		return
	}

	for _, host := range hosts {
		HostsState := hostsModel{
			Name:         types.StringValue(host.Name),
			State:        types.StringValue(host.State),
			Status:       types.StringValue(host.Status),
			Uuid:         types.StringValue(host.UUID),
			Architecture: types.StringValue(host.Architecture),
			Type:         types.StringValue(host.HypervisorType),
			ZoneUuid:     types.StringValue(host.ZoneUuid),
			ClusterUuid:  types.StringValue(host.ClusterUuid),
			ManagementIp: types.StringValue(host.ManagementIp),
		}

		state.Hosts = append(state.Hosts, HostsState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *hostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"name_regex": schema.StringAttribute{
				Description: "name_regex for Search and filter clusters",
				Optional:    true,
			},
			"hosts": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"architecture": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"zoneuuid": schema.StringAttribute{
							Computed: true,
						},
						"clusteruuid": schema.StringAttribute{
							Computed: true,
						},
						"managementop": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}

}
