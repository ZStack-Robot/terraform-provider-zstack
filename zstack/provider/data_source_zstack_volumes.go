// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &volumesDataSource{}
	_ datasource.DataSourceWithConfigure = &volumesDataSource{}
)

type volumesDataSource struct {
	client *client.ZSClient
}

type volumesDataSourceModel struct {
	Name        types.String      `tfsdk:"name"`
	NamePattern types.String      `tfsdk:"name_pattern"`
	Filter      []Filter          `tfsdk:"filter"`
	Volumes     []volumeDataModel `tfsdk:"volumes"`
}

type volumeDataModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	DiskOfferingUuid   types.String `tfsdk:"disk_offering_uuid"`
	DiskSize           types.Int64  `tfsdk:"disk_size"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
	VmInstanceUuid     types.String `tfsdk:"vm_instance_uuid"`
	LastVmInstanceUuid types.String `tfsdk:"last_vm_instance_uuid"`
	Type               types.String `tfsdk:"type"`
	Format             types.String `tfsdk:"format"`
	ActualSize         types.Int64  `tfsdk:"actual_size"`
	State              types.String `tfsdk:"state"`
	Status             types.String `tfsdk:"status"`
	IsShareable        types.Bool   `tfsdk:"is_shareable"`
}

func ZStackVolumesDataSource() datasource.DataSource {
	return &volumesDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *volumesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Metadata implements datasource.DataSource.
func (d *volumesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes"
}

// Schema implements datasource.DataSource.
func (d *volumesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack data volumes.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Exact name to match.",
			},
			"name_pattern": schema.StringAttribute{
				Optional:    true,
				Description: "Pattern for fuzzy name search. Use % for multiple characters and _ for exactly one character.",
			},
			"volumes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The matching data volumes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":                  schema.StringAttribute{Computed: true},
						"name":                  schema.StringAttribute{Computed: true},
						"description":           schema.StringAttribute{Computed: true},
						"disk_offering_uuid":    schema.StringAttribute{Computed: true},
						"disk_size":             schema.Int64Attribute{Computed: true},
						"primary_storage_uuid":  schema.StringAttribute{Computed: true},
						"vm_instance_uuid":      schema.StringAttribute{Computed: true},
						"last_vm_instance_uuid": schema.StringAttribute{Computed: true},
						"type":                  schema.StringAttribute{Computed: true},
						"format":                schema.StringAttribute{Computed: true},
						"actual_size":           schema.Int64Attribute{Computed: true},
						"state":                 schema.StringAttribute{Computed: true},
						"status":                schema.StringAttribute{Computed: true},
						"is_shareable":          schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any returned field.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by.",
						},
						"values": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Accepted values for the field. Multiple values are treated as OR.",
						},
					},
				},
			},
		},
	}
}

// Read implements datasource.DataSource.
func (d *volumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state volumesDataSourceModel
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

	volumes, err := d.client.QueryVolume(&params)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read volumes", err.Error())
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

	filteredVolumes, filterDiags := utils.FilterResource(ctx, volumes, filters, "volume")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Volumes = []volumeDataModel{}
	for _, volume := range filteredVolumes {
		state.Volumes = append(state.Volumes, volumeDataModel{
			Uuid:               types.StringValue(volume.UUID),
			Name:               types.StringValue(volume.Name),
			Description:        stringValueOrNull(volume.Description),
			DiskOfferingUuid:   stringValueOrNull(volume.DiskOfferingUuid),
			DiskSize:           types.Int64Value(int64(volume.Size)),
			PrimaryStorageUuid: stringValueOrNull(volume.PrimaryStorageUuid),
			VmInstanceUuid:     stringValueOrNull(volume.VmInstanceUuid),
			LastVmInstanceUuid: stringValueOrNull(volume.LastVmInstanceUuid),
			Type:               stringValueOrNull(volume.Type),
			Format:             stringValueOrNull(volume.Format),
			ActualSize:         types.Int64Value(int64(volume.ActualSize)),
			State:              stringValueOrNull(volume.State),
			Status:             stringValueOrNull(volume.Status),
			IsShareable:        types.BoolValue(volume.IsShareable),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
