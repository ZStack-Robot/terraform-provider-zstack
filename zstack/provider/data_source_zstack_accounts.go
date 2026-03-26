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
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

func ZStackAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountItem struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

type accountDataSourceModel struct {
	Name        types.String  `tfsdk:"name"`
	NamePattern types.String  `tfsdk:"name_pattern"`
	Filter      []Filter      `tfsdk:"filter"`
	Accounts    []accountItem `tfsdk:"accounts"`
}

type accountDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *accountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

// Read implements datasource.DataSource.
func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountDataSourceModel
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

	accounts, err := d.client.QueryAccount(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack Accounts",
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

	filterAccounts, filterDiags := utils.FilterResource(ctx, accounts, filters, "account")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, acct := range filterAccounts {
		state.Accounts = append(state.Accounts, accountItem{
			Uuid:        types.StringValue(acct.UUID),
			Name:        types.StringValue(acct.Name),
			Description: types.StringValue(acct.Description),
			Type:        types.StringValue(acct.Type),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSource.
func (d *accountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack Accounts by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying an account.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching account names. Use % or _ like SQL.",
				Optional:    true,
			},
			"accounts": schema.ListNestedAttribute{
				Description: "List of matched accounts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the account.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the account.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the account.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the account (Normal, SystemAdmin).",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter results by field values.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by.",
							Required:    true,
						},
						"values": schema.SetAttribute{
							Description: "List of values to match. Treated as OR conditions.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}
