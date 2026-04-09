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
	_ datasource.DataSource              = &userTagDataSource{}
	_ datasource.DataSourceWithConfigure = &userTagDataSource{}
)

type userTagDataSource struct {
	client *client.ZSClient
}

type userTagItemModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	TagPatternUuid types.String `tfsdk:"tag_pattern_uuid"`
	ResourceUuid   types.String `tfsdk:"resource_uuid"`
	ResourceType   types.String `tfsdk:"resource_type"`
	Tag            types.String `tfsdk:"tag"`
	Type           types.String `tfsdk:"type"`
}

type userTagDataSourceModel struct {
	Name        types.String       `tfsdk:"name"`
	NamePattern types.String       `tfsdk:"name_pattern"`
	UserTags    []userTagItemModel `tfsdk:"user_tags"`
	Filter      []Filter           `tfsdk:"filter"`
}

func ZStackUserTagDataSource() datasource.DataSource {
	return &userTagDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *userTagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the ZStack Provider developer. ", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Metadata implements datasource.DataSource.
func (d *userTagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_tags"
}

// Read implements datasource.DataSource.
func (d *userTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state userTagDataSourceModel
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

	userTags, err := d.client.QueryUserTag(&params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack User Tags",
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

	filterUserTags, filterDiags := utils.FilterResource(ctx, userTags, filters, "user_tag")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, userTag := range filterUserTags {
		userTagState := userTagItemModel{
			Uuid:           types.StringValue(userTag.UUID),
			Name:           types.StringValue(userTag.Name),
			TagPatternUuid: types.StringValue(userTag.TagPatternUuid),
			ResourceUuid:   types.StringValue(userTag.ResourceUuid),
			ResourceType:   types.StringValue(userTag.ResourceType),
			Tag:            types.StringValue(userTag.Tag),
			Type:           types.StringValue(userTag.Type),
		}
		state.UserTags = append(state.UserTags, userTagState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *userTagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of user tags and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching user tags",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"user_tags": schema.ListNestedAttribute{
				Description: "List of User Tags",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the user tag",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the user tag",
							Computed:    true,
						},
						"tag_pattern_uuid": schema.StringAttribute{
							Description: "UUID of the tag pattern",
							Computed:    true,
						},
						"resource_uuid": schema.StringAttribute{
							Description: "UUID of the resource",
							Computed:    true,
						},
						"resource_type": schema.StringAttribute{
							Description: "Type of the resource",
							Computed:    true,
						},
						"tag": schema.StringAttribute{
							Description: "Tag value",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the user tag",
							Computed:    true,
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by resource type, use `name = \"resourceType\"` and `values = [\"VmInstance\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., resourceType, tag).",
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
