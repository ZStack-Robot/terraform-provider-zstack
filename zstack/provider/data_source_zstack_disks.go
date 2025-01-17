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
	_ datasource.DataSource              = &disksDataSource{}
	_ datasource.DataSourceWithConfigure = &disksDataSource{}
)

type disksDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	NamePattern types.String `tfsdk:"name_pattern"`
	Filter      []Filter     `tfsdk:"filter"`
	Disks       []disksModel `tfsdk:"disks"`
}

type disksModel struct {
	Name               types.String `tfsdk:"name"`
	Uuid               types.String `tfsdk:"uuid"`
	Description        types.String `tfsdk:"description"`
	PrimaryStorageUUID types.String `tfsdk:"primary_storage_uuid"`
	VMInstanceUUID     types.String `tfsdk:"vm_instance_uuid"`
	DiskOfferingUUID   types.String `tfsdk:"disk_offering_uuid"`
	Type               types.String `tfsdk:"type"`
	Format             types.String `tfsdk:"format"`
	Size               types.Int64  `tfsdk:"size"`
	ActualSize         types.Int64  `tfsdk:"actual_size"`
	Status             types.String `tfsdk:"status"`
	State              types.String `tfsdk:"state"`
	IsShareable        types.Bool   `tfsdk:"is_shareable"`
}

func ZStackDisksDataSource() datasource.DataSource {
	return &disksDataSource{}
}

type disksDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *disksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *disksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disks"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *disksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state disksDataSourceModel

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

	disks, err := d.client.QueryVolume(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read disks",
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

	filterDisks, filterDiags := utils.FilterResource(ctx, disks, filters, "disks")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, disk := range filterDisks {
		diskState := disksModel{
			Name:               types.StringValue(disk.Name),
			Uuid:               types.StringValue(disk.UUID),
			Description:        types.StringValue(disk.Description),
			PrimaryStorageUUID: types.StringValue(disk.PrimaryStorageUUID),
			VMInstanceUUID:     types.StringValue(disk.VMInstanceUUID),
			DiskOfferingUUID:   types.StringValue(disk.DiskOfferingUUID),
			Type:               types.StringValue(disk.Type),
			Format:             types.StringValue(disk.Format),
			Size:               types.Int64Value(utils.BytesToGB(int64(disk.Size))),
			ActualSize:         types.Int64Value(utils.BytesToGB(int64(disk.ActualSize))),
			Status:             types.StringValue(disk.Status),
			State:              types.StringValue(disk.State),
			IsShareable:        types.BoolValue(disk.IsShareable),
		}
		state.Disks = append(state.Disks, diskState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *disksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of disks and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching  disks ",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"disks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the disk.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the disk.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A brief description of the disks.",
						},
						"primary_storage_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the primary storage.",
						},
						"vm_instance_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the VM instance.",
						},
						"disk_offering_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the disk offering.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the disk (e.g., Data, Root). ",
						},
						"format": schema.StringAttribute{
							Computed:    true,
							Description: "The format of the disk. ",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "The size of the disk in GB. ",
						},
						"actual_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The actual size of the disk in GB.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the disk (e.g., Enabled, Disabled).",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The current status of the disk. ",
						},
						"is_shareable": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the disk is shareable.",
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
