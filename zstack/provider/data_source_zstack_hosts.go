// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"zstack.io/zstack-sdk-go/pkg/client"
)

var (
	_ datasource.DataSource              = &hostsDataSource{}
	_ datasource.DataSourceWithConfigure = &hostsDataSource{}
)

func NewHostDataSource() datasource.DataSource {
	return &hostsDataSource{}
}

type hostsDataSource struct {
	client *client.ZSClient
}

/*
type hostsDataSourceModel struct {
	Clusters []hostsModel `tfsdk:"clusters"`
}

type hostsModel struct {
	Name           types.String `tfsdk:"name"`
	HypervisorType types.String `tfsdk:"hypervisortype"`
	State          types.String `tfsdk:"state"`
	Type           types.String `tfsdk:"type"`
	Uuid           types.String `tfsdk:"uuid"`
	ZoneUuid       types.String `tfsdk:"zoneuuid"`
}
*/

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
func (t *hostsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hosts"
}

// Read implements datasource.DataSourceWithConfigure.
func (t *hostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements datasource.DataSourceWithConfigure.
func (t *hostsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	panic("unimplemented")
}
