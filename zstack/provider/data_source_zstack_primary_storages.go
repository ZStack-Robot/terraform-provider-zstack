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
	_ datasource.DataSource              = &primaryStorageDataSource{}
	_ datasource.DataSourceWithConfigure = &primaryStorageDataSource{}
)

func ZStackPrimaryStorageDataSource() datasource.DataSource {
	return &primaryStorageDataSource{}
}

type primaryStorage struct {
	Name                      types.String `tfsdk:"name"`
	Uuid                      types.String `tfsdk:"uuid"`
	State                     types.String `tfsdk:"state"`
	Status                    types.String `tfsdk:"status"`
	TotalCapacity             types.Int64  `tfsdk:"total_capacity"`
	AvailableCapacity         types.Int64  `tfsdk:"available_capacity"`
	TotalPhysicalCapacity     types.Int64  `tfsdk:"total_physical_capacity"`
	AvailablePhysicalCapacity types.Int64  `tfsdk:"available_physical_capacity"`
	SystemUsedCapacity        types.Int64  `tfsdk:"system_used_capacity"`
}

type primaryStorageDataSourceModel struct {
	Name           types.String     `tfsdk:"name"`
	NamePattern    types.String     `tfsdk:"name_pattern"`
	Filter         []Filter         `tfsdk:"filter"`
	PrimaryStorges []primaryStorage `tfsdk:"primary_storages"`
}

type primaryStorageDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *primaryStorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *primaryStorageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_primary_storages"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *primaryStorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state primaryStorageDataSourceModel
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

	primaryStorages, err := d.client.QueryPrimaryStorage(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack primary Storages",
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

	filterPrimaryStorage, filterDiags := utils.FilterResource(ctx, primaryStorages, filters, "primary_storage")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, primarystorage := range filterPrimaryStorage {
		primaryStorageState := primaryStorage{
			TotalCapacity:             types.Int64Value(primarystorage.TotalCapacity),
			State:                     types.StringValue(primarystorage.State),
			Status:                    types.StringValue(primarystorage.Status),
			Uuid:                      types.StringValue(primarystorage.UUID),
			AvailableCapacity:         types.Int64Value(primarystorage.AvailableCapacity),
			Name:                      types.StringValue(primarystorage.Name),
			TotalPhysicalCapacity:     types.Int64Value(primarystorage.TotalPhysicalCapacity),
			AvailablePhysicalCapacity: types.Int64Value(primarystorage.AvailablePhysicalCapacity),
			SystemUsedCapacity:        types.Int64Value(primarystorage.SystemUsedCapacity),
		}

		state.PrimaryStorges = append(state.PrimaryStorges, primaryStorageState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *primaryStorageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all primary storages, or query primary storages by exact name match, or query primary storages by name pattern fuzzy match.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching primary storage.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"primary_storages": schema.ListNestedAttribute{
				Description: "List of primary storage entries",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the primary storage",
							Computed:    true,
						},

						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the primary storage",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the primary storage (Enabled or Disabled)",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Readiness status of the primary storage",
							Computed:    true,
						},
						"total_capacity": schema.Int64Attribute{
							Description: "Total capacity of the primary storage in bytes",
							Computed:    true,
						},
						"available_capacity": schema.Int64Attribute{
							Description: "Available capacity of the primary storage in bytes",
							Computed:    true,
						},
						"total_physical_capacity": schema.Int64Attribute{
							Description: "Total physical capacity of the primary storage in bytes",
							Computed:    true,
						},
						"available_physical_capacity": schema.Int64Attribute{
							Description: "Available physical capacity of the primary storage in bytes",
							Computed:    true,
						},
						"system_used_capacity": schema.Int64Attribute{
							Description: "System used capacity of the primary storage in bytes",
							Computed:    true,
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
