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
	_ datasource.DataSource              = &backupStorageDataSource{}
	_ datasource.DataSourceWithConfigure = &backupStorageDataSource{}
)

func ZStackBackupStorageDataSource() datasource.DataSource {
	return &backupStorageDataSource{}
}

type backupStorage struct {
	Name              types.String `tfsdk:"name"`
	Uuid              types.String `tfsdk:"uuid"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
	TotalCapacity     types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity types.Int64  `tfsdk:"available_capacity"`
}

type backupStorageDataSourceModel struct {
	Name          types.String    `tfsdk:"name"`
	NamePattern   types.String    `tfsdk:"name_pattern"`
	BackupStorges []backupStorage `tfsdk:"backup_storages"`
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
	resp.TypeName = req.ProviderTypeName + "_backupstorages"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *backupStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state backupStorageDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	// 优先检查 `name` 精确查询
	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
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
		Description: "List all backup storages, or query backup storages by exact name match, or query backup storages by name pattern fuzzy match.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching backup storage.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"backup_storages": schema.ListNestedAttribute{
				Description: "List of backup storage entries",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the backup storage",
							Computed:    true,
						},

						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the backup storage",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the backup storage (Enabled or Disabled)",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Readiness status of the backup storage",
							Computed:    true,
						},
						"total_capacity": schema.Int64Attribute{
							Description: "Total capacity of the backup storage in bytes",
							Computed:    true,
						},
						"available_capacity": schema.Int64Attribute{
							Description: "Available capacity of the backup storage in bytes",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
