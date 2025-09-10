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
	_ datasource.DataSource              = &hookScriptsDataSource{}
	_ datasource.DataSourceWithConfigure = &hookScriptsDataSource{}
)

func ZStackHookScriptsDataSource() datasource.DataSource {
	return &hookScriptsDataSource{}
}

type hookScriptsDataSource struct {
	client *client.ZSClient
}

type hookScriptsDataSourceModel struct {
	Name        types.String       `tfsdk:"name"`
	NamePattern types.String       `tfsdk:"name_pattern"`
	Filter      []Filter           `tfsdk:"filter"`
	HostScripts []hookScriptsModel `tfsdk:"hook_scripts"`
}

type hookScriptsModel struct {
	Name       types.String `tfsdk:"name"`
	Uuid       types.String `tfsdk:"uuid"`
	Type       types.String `tfsdk:"type"`
	HookScript types.String `tfsdk:"hook_script"`
}

// todo modify mapping tools
// Configure implements datasource.DataSourceWithConfigure.
func (d *hookScriptsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *hookScriptsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hook_scripts"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *hookScriptsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state hookScriptsDataSourceModel
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

	hook_scripts, err := d.client.QueryVmUserDefinedXmlHookScript(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Hosts ",
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

	filterHostScripts, filterDiags := utils.FilterResource(ctx, hook_scripts, filters, "host_script")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, hostScripts := range filterHostScripts {
		HostScriptsState := hookScriptsModel{
			Name:       types.StringValue(hostScripts.Name),
			Uuid:       types.StringValue(hostScripts.UUID),
			Type:       types.StringValue(hostScripts.Type),
			HookScript: types.StringValue(hostScripts.HookScript),
		}

		state.HostScripts = append(state.HostScripts, HostScriptsState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *hookScriptsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of hook_scripts and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching hook_scripts",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"hook_scripts": schema.ListNestedAttribute{
				Description: "List of host entries matching the specified filters",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID Unique identifier of the host script",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the host script",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the host script (e.g., Customization)",
						},
						"hook_script": schema.StringAttribute{
							Computed:    true,
							Description: "content of the hook_script(Base64 encode)",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by type, use `name = \"type\"` and `values = [\"Customization\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., type).",
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
