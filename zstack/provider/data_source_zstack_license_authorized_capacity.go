// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
)

var (
	_ datasource.DataSource              = &licenseAuthorizedCapacityDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseAuthorizedCapacityDataSource{}
)

func ZStackLicenseAuthorizedCapacityDataSource() datasource.DataSource {
	return &licenseAuthorizedCapacityDataSource{}
}

type licenseAuthorizedCapacityDataSource struct {
	client *client.ZSClient
}

type licenseAuthorizedCapacityDataSourceModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Id           types.Int64  `tfsdk:"id"`
	NodeUuid     types.String `tfsdk:"node_uuid"`
	ResourceUuid types.String `tfsdk:"resource_uuid"`
	ResourceInfo types.String `tfsdk:"resource_info"`
	QuotaType    types.String `tfsdk:"quota_type"`
	Quota        types.Int64  `tfsdk:"quota"`
	LicenseType  types.String `tfsdk:"license_type"`
	Type         types.String `tfsdk:"type"`
}

func (d *licenseAuthorizedCapacityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *licenseAuthorizedCapacityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_authorized_capacity"
}

func (d *licenseAuthorizedCapacityDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches license authorized capacity from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"uuid":          schema.StringAttribute{Computed: true, Description: "UUID of the license authorized capacity."},
			"name":          schema.StringAttribute{Computed: true, Description: "Name of the license authorized capacity."},
			"id":            schema.Int64Attribute{Computed: true, Description: "ID of the license authorized capacity."},
			"node_uuid":     schema.StringAttribute{Computed: true, Description: "Management node UUID for this capacity record."},
			"resource_uuid": schema.StringAttribute{Computed: true, Description: "Resource UUID for this capacity record."},
			"resource_info": schema.StringAttribute{Computed: true, Description: "Resource information for this capacity record."},
			"quota_type":    schema.StringAttribute{Computed: true, Description: "Quota type."},
			"quota":         schema.Int64Attribute{Computed: true, Description: "Authorized quota value."},
			"license_type":  schema.StringAttribute{Computed: true, Description: "License type."},
			"type":          schema.StringAttribute{Computed: true, Description: "Capacity type."},
		},
	}
}

func (d *licenseAuthorizedCapacityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseAuthorizedCapacityDataSourceModel

	capacity, err := d.client.GetLicenseAuthorizedCapacity()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack License Authorized Capacity",
			err.Error(),
		)
		return
	}

	if capacity == nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack License Authorized Capacity",
			"API returned empty license authorized capacity response.",
		)
		return
	}

	state.Uuid = types.StringValue(capacity.UUID)
	state.Name = types.StringValue(capacity.Name)
	state.Id = types.Int64Value(capacity.Id)
	state.NodeUuid = types.StringValue(capacity.NodeUuid)
	state.ResourceUuid = types.StringValue(capacity.ResourceUuid)
	state.ResourceInfo = types.StringValue(capacity.ResourceInfo)
	state.QuotaType = types.StringValue(capacity.QuotaType)
	state.Quota = types.Int64Value(capacity.Quota)
	state.LicenseType = types.StringValue(capacity.LicenseType)
	state.Type = types.StringValue(capacity.Type)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
