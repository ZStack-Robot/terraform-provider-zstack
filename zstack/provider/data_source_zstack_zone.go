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
	Name_regex types.String `tfsdk:"name_regex"`
	Zones      []zoneModel  `tfsdk:"zones"`
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
		Description: "Fetches the list of zones. ",
		Attributes: map[string]schema.Attribute{
			"name_regex": schema.StringAttribute{
				Description: "name_regex for Search and filter zones",
				Optional:    true,
			},
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
	}
}

func (d *zoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state zoneDataSourceModel
	//var state zoneModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	name_regex := state.Name_regex
	params := param.NewQueryParam()

	if !name_regex.IsNull() {
		params.AddQ("name=" + name_regex.ValueString())
	}

	//images, err := d.client.QueryImage(params)

	zones, err := d.client.QueryZone(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack zones",
			err.Error(),
		)

		return
	}

	//map query zones body to mode
	for _, zone := range zones {
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
