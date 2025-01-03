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
	resp.TypeName = req.ProviderTypeName + "_zones"
}

func (d *zoneDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of zones and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
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

	name := state.Name
	params := param.NewQueryParam()

	if !name.IsNull() {
		params.AddQ("name=" + name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	zones, err := d.client.QueryZone(params)
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
