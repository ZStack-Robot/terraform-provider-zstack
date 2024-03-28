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
	_ datasource.DataSource              = &backupStorageDataSource{}
	_ datasource.DataSourceWithConfigure = &backupStorageDataSource{}
)

func ZstackBackupStorageDataSource() datasource.DataSource {
	return &backupStorageDataSource{}
}

type backupStorage struct {
	Name              types.String `tfsdk:"name"`
	Uuid              types.String `tfsdk:"uuid"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	TotalCapacity     types.Int64  `tfsdk:"totalcapacity"`
	AvailableCapacity types.Int64  `tfsdk:"availablecapacity"`
}

type backupStorageDataSourceModel struct {
	Name_regex    types.String    `tfsdk:"name_regex"`
	BackupStorges []backupStorage `tfsdk:"backupstorages"`
}

type backupStorageDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *backupStorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *backupStorageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backupstorage"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *backupStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backupStorageDataSourceModel
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

	backupstorages, err := d.client.QueryBackupStorage(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Backup Storages",
			err.Error(),
		)

		return
	}

	for _, backupstorage := range backupstorages {
		backupStorageState := backupStorage{
			TotalCapacity:     types.Int64Value(backupstorage.TotalCapacity),
			State:             types.StringValue(backupstorage.State),
			Status:            types.StringValue(backupstorage.Status),
			Uuid:              types.StringValue(backupstorage.UUID),
			AvailableCapacity: types.Int64Value(backupstorage.AvailableCapacity),
			Name:              types.StringValue(backupstorage.Name),
		}

		state.BackupStorges = append(state.BackupStorges, backupStorageState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *backupStorageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name_regex": schema.StringAttribute{
				Description: "name_regex for Search and filter backup storage",
				Optional:    true,
			},
			"backupstorages": schema.ListNestedAttribute{
				Description: "List of backup storages",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "name of backup storage ",
							Computed:    true,
						},

						"uuid": schema.StringAttribute{
							Description: "uuid identifier of the backup storage",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "Enabled or Disabled state of the backup storage",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Ready or not of the  backup storage",
							Computed:    true,
						},
						"totalcapacity": schema.Int64Attribute{
							Description: "Total capacity of the backup storage",
							Computed:    true,
						},
						"availablecapacity": schema.Int64Attribute{
							Description: "AvailableCapacity of the backup storage",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
