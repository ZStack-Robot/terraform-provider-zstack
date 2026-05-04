// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

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
	_ datasource.DataSource              = &globalConfigsDataSource{}
	_ datasource.DataSourceWithConfigure = &globalConfigsDataSource{}
)

func ZStackGlobalConfigsDataSource() datasource.DataSource {
	return &globalConfigsDataSource{}
}

type globalConfigsDataSource struct {
	client *client.ZSClient
}

type globalConfigsDataSourceModel struct {
	Category      types.String       `tfsdk:"category"`
	Name          types.String       `tfsdk:"name"`
	NamePattern   types.String       `tfsdk:"name_pattern"`
	GlobalConfigs []globalConfigItem `tfsdk:"global_configs"`
}

type globalConfigItem struct {
	Category     types.String `tfsdk:"category"`
	Name         types.String `tfsdk:"name"`
	Value        types.String `tfsdk:"value"`
	DefaultValue types.String `tfsdk:"default_value"`
	Description  types.String `tfsdk:"description"`
}

func (d *globalConfigsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *globalConfigsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_configs"
}

func (d *globalConfigsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state globalConfigsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()
	if !state.Category.IsNull() && !state.Category.IsUnknown() && state.Category.ValueString() != "" {
		params.AddQ("category=" + state.Category.ValueString())
	}
	if !state.Name.IsNull() && !state.Name.IsUnknown() && state.Name.ValueString() != "" {
		params.AddQ("name=" + state.Name.ValueString())
	}
	if !state.NamePattern.IsNull() && !state.NamePattern.IsUnknown() && state.NamePattern.ValueString() != "" {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	configs, err := d.client.QueryGlobalConfig(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack Global Configs",
			err.Error(),
		)
		return
	}

	if len(configs) == 0 && (!state.Category.IsNull() || !state.Name.IsNull() || !state.NamePattern.IsNull()) {
		resp.Diagnostics.AddError(
			"No Matching Global Configs",
			"No global configs matched the supplied filters. Use data.zstack_global_configs without filters to inspect valid category/name pairs for the target ZStack environment.",
		)
		return
	}

	state.GlobalConfigs = make([]globalConfigItem, 0, len(configs))
	for _, config := range configs {
		state.GlobalConfigs = append(state.GlobalConfigs, globalConfigItem{
			Category:     types.StringValue(config.Category),
			Name:         types.StringValue(config.Name),
			Value:        types.StringValue(config.Value),
			DefaultValue: stringValueOrNull(config.DefaultValue),
			Description:  stringValueOrNull(config.Description),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (d *globalConfigsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack global configuration keys and values. Use this data source to discover valid category/name pairs before managing a zstack_global_config resource.",
		Attributes: map[string]schema.Attribute{
			"category": schema.StringAttribute{
				Description: "Exact global config category to query.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Exact global config name to query. Mutually exclusive with `name_pattern`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("name_pattern")),
				},
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching global config names. Use % for multiple characters and _ for exactly one character. Mutually exclusive with `name`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("name")),
				},
			},
			"global_configs": schema.ListNestedAttribute{
				Description: "List of matched global configs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{
							Description: "Global config category.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Global config name.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Current global config value.",
							Computed:    true,
						},
						"default_value": schema.StringAttribute{
							Description: "Default global config value.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Global config description.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
