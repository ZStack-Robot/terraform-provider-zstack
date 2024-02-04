// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package zstack

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
	_ datasource.DataSource              = &imageDataSource{}
	_ datasource.DataSourceWithConfigure = &imageDataSource{}
)

type imageDataSource struct {
	client *client.ZSClient
}

type imagesModel struct {
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
	Uuid         types.String `tfsdk:"uuid"`
	GuestOsType  types.String `tfsdk:"guestostype"`
	Format       types.String `tfsdk:"format"`
	Platform     types.String `tfsdk:"platform"`
	Architecture types.String `tfsdk:"architecture"`
}

type imagesDataSourceModel struct {
	Name_regx types.String  `tfsdk:"name_regx"`
	Images    []imagesModel `tfsdk:"images"`
}

func ZStackImageDataSource() datasource.DataSource {
	return &imageDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *imageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. jiajian.chi@zstack.io", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata implements datasource.DataSource.
func (d *imageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

// Read implements datasource.DataSource.
func (d *imageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state imagesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	name_regex := state.Name_regx

	if !name_regex.IsNull() {
		params := param.NewQueryParam()
		params.AddQ("name=" + name_regex.ValueString())
		images, err := d.client.QueryImage(params)

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read ZStack Images ",
				err.Error(),
			)

			return
		}
		for _, image := range images {
			imageState := imagesModel{
				Name:         types.StringValue(image.Name),
				State:        types.StringValue(image.State),
				Status:       types.StringValue(image.Status),
				Uuid:         types.StringValue(image.UUID),
				GuestOsType:  types.StringValue(image.GuestOsType),
				Format:       types.StringValue(image.Format),
				Platform:     types.StringValue(image.Platform),
				Architecture: types.StringValue(string(image.Architecture)),
			}

			state.Images = append(state.Images, imageState)
		}

	}

	if name_regex.IsNull() {
		images, err := d.client.QueryImage(param.NewQueryParam())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read ZStack Images ",
				err.Error(),
			)

			return
		}
		for _, image := range images {

			imageState := imagesModel{
				Name:         types.StringValue(image.Name),
				State:        types.StringValue(image.State),
				Status:       types.StringValue(image.Status),
				Uuid:         types.StringValue(image.UUID),
				GuestOsType:  types.StringValue(image.GuestOsType),
				Format:       types.StringValue(image.Format),
				Platform:     types.StringValue(image.Platform),
				Architecture: types.StringValue(string(image.Architecture)),
			}

			state.Images = append(state.Images, imageState)
		}

	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *imageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name_regx": schema.StringAttribute{
				Optional: true,
			},
			"images": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},

						"uuid": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"guestostype": schema.StringAttribute{
							Computed: true,
						},
						"format": schema.StringAttribute{
							Computed: true,
						},
						"platform": schema.StringAttribute{
							Computed: true,
						},
						"architecture": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
