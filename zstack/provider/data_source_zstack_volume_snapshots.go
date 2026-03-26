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
	_ datasource.DataSource              = &volumeSnapshotsDataSource{}
	_ datasource.DataSourceWithConfigure = &volumeSnapshotsDataSource{}
)

type volumeSnapshotsDataSource struct {
	client *client.ZSClient
}

type volumeSnapshotsDataSourceModel struct {
	Name        types.String              `tfsdk:"name"`
	NamePattern types.String              `tfsdk:"name_pattern"`
	Filter      []Filter                  `tfsdk:"filter"`
	Snapshots   []volumeSnapshotDataModel `tfsdk:"snapshots"`
}

type volumeSnapshotDataModel struct {
	Uuid               types.String `tfsdk:"uuid"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	VolumeUuid         types.String `tfsdk:"volume_uuid"`
	TreeUuid           types.String `tfsdk:"tree_uuid"`
	ParentUuid         types.String `tfsdk:"parent_uuid"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
	VolumeType         types.String `tfsdk:"volume_type"`
	Format             types.String `tfsdk:"format"`
	Latest             types.Bool   `tfsdk:"latest"`
	Size               types.Int64  `tfsdk:"size"`
	State              types.String `tfsdk:"state"`
	Status             types.String `tfsdk:"status"`
	Distance           types.Int64  `tfsdk:"distance"`
	GroupUuid          types.String `tfsdk:"group_uuid"`
}

func ZStackVolumeSnapshotsDataSource() datasource.DataSource {
	return &volumeSnapshotsDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *volumeSnapshotsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *volumeSnapshotsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_snapshots"
}

// Schema implements datasource.DataSource.
func (d *volumeSnapshotsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack volume snapshots.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Exact name to match.",
			},
			"name_pattern": schema.StringAttribute{
				Optional:    true,
				Description: "Pattern for fuzzy name search. Use % for multiple characters and _ for exactly one character.",
			},
			"snapshots": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The matching volume snapshots.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":                 schema.StringAttribute{Computed: true},
						"name":                 schema.StringAttribute{Computed: true},
						"description":          schema.StringAttribute{Computed: true},
						"volume_uuid":          schema.StringAttribute{Computed: true},
						"tree_uuid":            schema.StringAttribute{Computed: true},
						"parent_uuid":          schema.StringAttribute{Computed: true},
						"primary_storage_uuid": schema.StringAttribute{Computed: true},
						"volume_type":          schema.StringAttribute{Computed: true},
						"format":               schema.StringAttribute{Computed: true},
						"latest":               schema.BoolAttribute{Computed: true},
						"size":                 schema.Int64Attribute{Computed: true},
						"state":                schema.StringAttribute{Computed: true},
						"status":               schema.StringAttribute{Computed: true},
						"distance":             schema.Int64Attribute{Computed: true},
						"group_uuid":           schema.StringAttribute{Computed: true},
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
func (d *volumeSnapshotsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state volumeSnapshotsDataSourceModel
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

	snapshots, err := d.client.QueryVolumeSnapshot(&params)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read volume snapshots", err.Error())
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

	filteredSnapshots, filterDiags := utils.FilterResource(ctx, snapshots, filters, "volume_snapshot")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Snapshots = []volumeSnapshotDataModel{}
	for _, snapshot := range filteredSnapshots {
		state.Snapshots = append(state.Snapshots, volumeSnapshotDataModel{
			Uuid:               types.StringValue(snapshot.UUID),
			Name:               types.StringValue(snapshot.Name),
			Description:        stringValueOrNull(snapshot.Description),
			VolumeUuid:         stringValueOrNull(snapshot.VolumeUuid),
			TreeUuid:           stringValueOrNull(snapshot.TreeUuid),
			ParentUuid:         stringValueOrNull(snapshot.ParentUuid),
			PrimaryStorageUuid: stringValueOrNull(snapshot.PrimaryStorageUuid),
			VolumeType:         stringValueOrNull(snapshot.VolumeType),
			Format:             stringValueOrNull(snapshot.Format),
			Latest:             types.BoolValue(snapshot.Latest),
			Size:               types.Int64Value(snapshot.Size),
			State:              stringValueOrNull(snapshot.State),
			Status:             stringValueOrNull(snapshot.Status),
			Distance:           types.Int64Value(int64(snapshot.Distance)),
			GroupUuid:          stringValueOrNull(snapshot.GroupUuid),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
