// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ datasource.DataSource              = &tagDataSource{}
	_ datasource.DataSourceWithConfigure = &tagDataSource{}
)

type tagDataSource struct {
	client *client.ZSClient
}

type tagDataSourceModel struct {
	Name        types.String     `tfsdk:"name"`
	NamePattern types.String     `tfsdk:"name_pattern"`
	TagType     types.String     `tfsdk:"tag_type"`
	Filter      []Filter         `tfsdk:"filter"`
	Tags        []tagModel       `tfsdk:"tags"`
	UserTags    []userTagModel   `tfsdk:"user_tags"`   // Optional, if you want to fetch User Tags separately
	SystemTags  []systemTagModel `tfsdk:"system_tags"` // Optional, if you want to fetch System Tags separately
}
type userTagModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceUuid types.String `tfsdk:"resource_uuid"`
	Tag          types.String `tfsdk:"tag"`
	Type         types.String `tfsdk:"type"`
}

type systemTagModel struct {
	Uuid types.String `tfsdk:"uuid"`
	//	Name         types.String `tfsdk:"name"`
	//	Description  types.String `tfsdk:"description"`
	Inherent     types.Bool   `tfsdk:"inherent"`
	ResourceUuid types.String `tfsdk:"resource_uuid"`
	ResourceType types.String `tfsdk:"resource_type"`
	Tag          types.String `tfsdk:"tag"`
}

type tagModel struct {
	Name        types.String `tfsdk:"name"`
	Uuid        types.String `tfsdk:"uuid"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
	Type        types.String `tfsdk:"type"`
}

func ZStackTagDataSource() datasource.DataSource {
	return &tagDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *tagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *tagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tags"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *tagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state tagDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.TagType.IsNull() {
		resp.Diagnostics.AddError("Missing tag_type", "You must specify 'tag_type': 'user', 'system', or 'tag'.")
		return
	}

	tagType := state.TagType.ValueString()
	params := param.NewQueryParam()

	if tagType == "tag" {
		if !state.Name.IsNull() && state.Name.ValueString() != "" {
			params.AddQ("name=" + state.Name.ValueString())
		} else if !state.NamePattern.IsNull() && state.NamePattern.ValueString() != "" {
			params.AddQ("name~=" + state.NamePattern.ValueString())
		}
	}

	/*
		if !state.Name.IsNull() && state.Name.ValueString() != "" {
			params.AddQ("name=" + state.Name.ValueString())
		} else if !state.NamePattern.IsNull() && state.NamePattern.ValueString() != "" {
			params.AddQ("name~=" + state.NamePattern.ValueString())
		}
	*/

	// Apply filters
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

	switch tagType {
	case "user":
		userTags, err := d.client.QueryUserTag(params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Fetch User Tags from ZStack",
				fmt.Sprintf("Error while querying user tags: %s", err.Error()),
			)
			return
		}
		filteredTags, filterDiags := utils.FilterResource(ctx, userTags, filters, "user_tags")
		resp.Diagnostics.Append(filterDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, tag := range filteredTags {
			state.UserTags = append(state.UserTags, userTagModel{
				Uuid:         types.StringValue(tag.Uuid),
				ResourceType: types.StringValue(tag.ResourceType),
				ResourceUuid: types.StringValue(tag.ResourceUuid),
				Tag:          types.StringValue(tag.Tag),
				Type:         types.StringValue("User"),
			})
		}
	case "system":
		systemTags, err := d.client.QuerySystemTags(params)
		if err != nil {
			resp.Diagnostics.AddError("Unable to query system tags", err.Error())
			return
		}
		filteredTags, filterDiags := utils.FilterResource(ctx, systemTags, filters, "system_tags")
		resp.Diagnostics.Append(filterDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, tag := range filteredTags {
			state.SystemTags = append(state.SystemTags, systemTagModel{
				Uuid: types.StringValue(tag.UUID),
				//		Name:         types.StringValue(tag.Name),
				//		Description:  types.StringValue(tag.Description),
				Inherent:     types.BoolValue(tag.Inherent),
				ResourceUuid: types.StringValue(tag.ResourceUuid),
				ResourceType: types.StringValue(tag.ResourceType),
				Tag:          types.StringValue(tag.Tag),
			})
		}
	case "tag":
		tags, err := d.client.QueryTag(params)
		if err != nil {
			resp.Diagnostics.AddError("Unable to query tags", err.Error())
			return
		}
		filteredTags, filterDiags := utils.FilterResource(ctx, tags, filters, "tag")
		resp.Diagnostics.Append(filterDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, tag := range filteredTags {
			state.Tags = append(state.Tags, tagModel{
				Uuid:        types.StringValue(tag.UUID),
				Name:        types.StringValue(tag.Name),
				Description: types.StringValue(tag.Description),
				Color:       types.StringValue(tag.Color),
				Type:        types.StringValue(tag.Type),
			})
		}

	default:
		resp.Diagnostics.AddError("Invalid tag_type", fmt.Sprintf("Unsupported tag_type: %s. Use 'user', 'system', or 'tag'.", tagType))
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *tagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve User, System, or regular Tags from the ZStack environment based on name, pattern, or filters.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name of the tag to match.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy matching the tag name (supports % as wildcard, _ as single character).",
				Optional:    true,
			},
			"tag_type": schema.StringAttribute{
				Description: "Specifies which type of tag to query: 'user', 'system', or 'tag'.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("user", "system", "tag"),
				},
			},
			"tags": schema.ListNestedAttribute{
				Description: "List of regular tags matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "Unique identifier of the tag.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the tag.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the tag.",
							Computed:    true,
						},
						"color": schema.StringAttribute{
							Description: "Color label assigned to the tag.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the tag.",
							Computed:    true,
						},
					},
				},
			},
			"user_tags": schema.ListNestedAttribute{
				Description: "List of user tags matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "Unique identifier of the user tag.",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Type of the resource the tag is attached to.",
							Computed:    true,
						},
						"resource_uuid": schema.StringAttribute{
							Description: "UUID of the resource the tag is attached to.",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "Tag value.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Tag category, typically 'User'.",
							Computed:    true,
						},
					},
				},
			},
			"system_tags": schema.ListNestedAttribute{
				Description: "List of system tags matching the query.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "Unique identifier of the system tag.",
							Computed:    true,
						},
						/*
							"name": schema.StringAttribute{
								Description: "Name of the system tag.",
								Computed:    true,
							},
							"description": schema.StringAttribute{
								Description: "Description of the system tag.",
								Computed:    true,
							},*/
						"inherent": schema.BoolAttribute{
							Description: "Indicates if the tag is inherent (built-in) or user-defined.",
							Computed:    true,
						},
						"resource_uuid": schema.StringAttribute{
							Description: "UUID of the resource the system tag is attached to.",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Type of the resource the system tag is attached to.",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "System tag value.",
							Computed:    true,
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
