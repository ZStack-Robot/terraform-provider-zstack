// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &zoneDataSource{}
	_ datasource.DataSourceWithConfigure = &zoneDataSource{}
)

func ZStackZoneDataSource() datasource.DataSource {
	return &zoneDataSource{}
}

type zoneDataSource struct {
	client *client.ZSClient
}

type zoneDataSourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	NamePattern types.String `tfsdk:"name_pattern"`
	Filter      []Filter     `tfsdk:"filter"`
	Zones       []zoneModel  `tfsdk:"zones"`
}

type zoneModel struct {
	Uuid  types.String `tfsdk:"uuid"`
	Name  types.String `tfsdk:"name"`
	State types.String `tfsdk:"state"`
	Type  types.String `tfsdk:"type"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *zoneDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *zoneDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (d *zoneDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of zones and their associated attributes from the ZStack environment. " +
			"For automation / AI-generated configurations, prefer `uuid` for stable, deterministic lookups; " +
			"`name` / `name_pattern` remain ergonomic for human-authored configs.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "Exact UUID lookup. Recommended for automation: stable across renames, " +
					"deterministic (0 or 1 match), idempotent. Mutually exclusive with `name` / `name_pattern`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("name"),
						path.MatchRoot("name_pattern"),
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Exact name for Searching  zones",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			/*
				"filter": schema.MapAttribute{
					Description: "Key-value pairs to filter Zones . For example, to filter by State, use `State = \"Enabled\"`.",
					Optional:    true,
					ElementType: types.StringType,
				},
			*/
			"zones": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
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

func (d *zoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zoneDataSourceModel
	//var state zoneModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	params := param.NewQueryParam()
	applyUuidOrNameFilter(&params, state.Uuid, state.Name, state.NamePattern)

	zones, err := d.client.QueryZone(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack zones",
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

	filterZones, filterDiags := utils.FilterResource(ctx, zones, filters, "zone")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, zone := range filterZones {
		zoneState := zoneModel{
			Name:  types.StringValue(zone.Name),
			State: types.StringValue(zone.State),
			Type:  types.StringValue(zone.Type),
			Uuid:  types.StringValue(zone.UUID),
		}

		state.Zones = append(state.Zones, zoneState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
