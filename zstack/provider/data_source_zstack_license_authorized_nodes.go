// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &licenseAuthorizedNodeDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseAuthorizedNodeDataSource{}
)

func ZStackLicenseAuthorizedNodeDataSource() datasource.DataSource {
	return &licenseAuthorizedNodeDataSource{}
}

type licenseAuthorizedNodeItem struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	AppId        types.String `tfsdk:"app_id"`
	Ip           types.String `tfsdk:"ip"`
	LastSyncDate types.String `tfsdk:"last_sync_date"`
	Status       types.String `tfsdk:"status"`
	Type         types.String `tfsdk:"type"`
}

type licenseAuthorizedNodeDataSourceModel struct {
	Name        types.String                `tfsdk:"name"`
	NamePattern types.String                `tfsdk:"name_pattern"`
	Filter      []Filter                    `tfsdk:"filter"`
	Nodes       []licenseAuthorizedNodeItem `tfsdk:"nodes"`
}

type licenseAuthorizedNodeDataSource struct {
	client *client.ZSClient
}

func (d *licenseAuthorizedNodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *licenseAuthorizedNodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_authorized_nodes"
}

func (d *licenseAuthorizedNodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state licenseAuthorizedNodeDataSourceModel
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

	nodes, err := d.client.QueryLicenseAuthorizedNode(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack License Authorized Nodes",
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

	filterNodes, filterDiags := utils.FilterResource(ctx, nodes, filters, "license_authorized_node")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, node := range filterNodes {
		state.Nodes = append(state.Nodes, licenseAuthorizedNodeItem{
			Uuid:         types.StringValue(node.UUID),
			Name:         types.StringValue(node.Name),
			AppId:        types.StringValue(node.AppId),
			Ip:           types.StringValue(node.Ip),
			LastSyncDate: types.StringValue(node.LastSyncDate.Format(time.RFC3339)),
			Status:       types.StringValue(node.Status),
			Type:         types.StringValue(node.Type),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *licenseAuthorizedNodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack License Authorized Nodes by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching names.",
				Optional:    true,
			},
			"nodes": schema.ListNestedAttribute{
				Description: "List of matched license authorized nodes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":           schema.StringAttribute{Computed: true, Description: "UUID of the license authorized node."},
						"name":           schema.StringAttribute{Computed: true, Description: "Name of the license authorized node."},
						"app_id":         schema.StringAttribute{Computed: true, Description: "Application ID of the node."},
						"ip":             schema.StringAttribute{Computed: true, Description: "IP address of the node."},
						"last_sync_date": schema.StringAttribute{Computed: true, Description: "Last synchronization timestamp in RFC3339 format."},
						"status":         schema.StringAttribute{Computed: true, Description: "Status of the node."},
						"type":           schema.StringAttribute{Computed: true, Description: "Type of the node."},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by field values.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":   schema.StringAttribute{Required: true, Description: "Name of the field to filter by."},
						"values": schema.SetAttribute{Required: true, ElementType: types.StringType, Description: "List of values to match."},
					},
				},
			},
		},
	}
}
