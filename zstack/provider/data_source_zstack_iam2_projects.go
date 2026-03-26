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
	_ datasource.DataSource              = &iam2ProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &iam2ProjectDataSource{}
)

func ZStackIAM2ProjectDataSource() datasource.DataSource {
	return &iam2ProjectDataSource{}
}

type iam2ProjectItem struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.String `tfsdk:"state"`
}

type iam2ProjectDataSourceModel struct {
	Name         types.String      `tfsdk:"name"`
	NamePattern  types.String      `tfsdk:"name_pattern"`
	Filter       []Filter          `tfsdk:"filter"`
	IAM2Projects []iam2ProjectItem `tfsdk:"iam2_projects"`
}

type iam2ProjectDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *iam2ProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *iam2ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam2_projects"
}

// Read implements datasource.DataSource.
func (d *iam2ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state iam2ProjectDataSourceModel
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

	projects, err := d.client.QueryIAM2Project(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Query ZStack IAM2 Projects",
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

	filterProjects, filterDiags := utils.FilterResource(ctx, projects, filters, "iam2_project")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, proj := range filterProjects {
		state.IAM2Projects = append(state.IAM2Projects, iam2ProjectItem{
			Uuid:        types.StringValue(proj.UUID),
			Name:        types.StringValue(proj.Name),
			Description: types.StringValue(proj.Description),
			State:       types.StringValue(proj.State),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Schema implements datasource.DataSource.
func (d *iam2ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Query ZStack IAM2 Projects by name, name pattern, or additional filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for querying an IAM2 project.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching IAM2 project names. Use % or _ like SQL.",
				Optional:    true,
			},
			"iam2_projects": schema.ListNestedAttribute{
				Description: "List of matched IAM2 projects.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID of the IAM2 project.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the IAM2 project.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the IAM2 project.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the IAM2 project (Enabled, Disabled).",
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
