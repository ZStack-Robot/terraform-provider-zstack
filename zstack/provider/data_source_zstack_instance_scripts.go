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
	_ datasource.DataSource              = &instanceScriptDataSource{}
	_ datasource.DataSourceWithConfigure = &instanceScriptDataSource{}
)

func ZStackInstanceScriptDataSource() datasource.DataSource {
	return &instanceScriptDataSource{}
}

type instanceScriptDataSource struct {
	client *client.ZSClient
}

type instanceScriptDataSourceModel struct {
	Name        types.String          `tfsdk:"name"`
	NamePattern types.String          `tfsdk:"name_pattern"`
	Filter      []Filter              `tfsdk:"filter"`
	Scripts     []instanceScriptModel `tfsdk:"scripts"`
}

type instanceScriptModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	ScriptContent types.String `tfsdk:"script_content"`
	RenderParams  types.String `tfsdk:"render_params"`
	Platform      types.String `tfsdk:"platform"`
	ScriptType    types.String `tfsdk:"script_type"`
	ScriptTimeout types.Int64  `tfsdk:"script_timeout"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *instanceScriptDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *instanceScriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scripts"
}

func (d *instanceScriptDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of instance scripts and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching scripts",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"scripts": schema.ListNestedAttribute{
				Description: "Returned script objects.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the script.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the script.",
						},
						"script_content": schema.StringAttribute{
							Computed:    true,
							Description: "Script content.",
						},
						"render_params": schema.StringAttribute{
							Computed:    true,
							Description: "Render parameters used when executing the script.",
						},
						"platform": schema.StringAttribute{
							Computed:    true,
							Description: "Platform type, e.g., Linux.",
						},
						"script_type": schema.StringAttribute{
							Computed:    true,
							Description: "Script type (e.g., Shell, Python, Perl, Bat, Powershell).",
						},
						"script_timeout": schema.Int64Attribute{
							Computed:    true,
							Description: "Script timeout in seconds.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Additional filtering by field name and values.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Field name to filter by.",
							Required:    true,
						},
						"values": schema.SetAttribute{
							Description: "List of acceptable values for the field.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *instanceScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instanceScriptDataSourceModel
	//var state clusterModel

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

	scripts, err := d.client.QueryVmInstanceScript(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Scripts",
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

	filterScripts, filterDiags := utils.FilterResource(ctx, scripts, filters, "script")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, script := range filterScripts {
		scriptState := instanceScriptModel{
			Uuid:          types.StringValue(script.UUID),
			Name:          types.StringValue(script.Name),
			ScriptContent: types.StringValue(script.ScriptContent),
			RenderParams:  types.StringValue(script.RenderParams),
			Platform:      types.StringValue(script.Platform),
			ScriptType:    types.StringValue(script.ScriptType),
			ScriptTimeout: types.Int64Value(int64(script.ScriptTimeout)),
		}
		state.Scripts = append(state.Scripts, scriptState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
