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
	_ datasource.DataSource              = &imgDataSource{}
	_ datasource.DataSourceWithConfigure = &imgDataSource{}
)

type imgDataSource struct {
	client *client.ZSClient
}

type imgModel struct {
	Name types.String `tfsdk:"name"`
	Uuid types.String `tfsdk:"uuid"`
}

func ZStackImgDataSource() datasource.DataSource {
	return &imgDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *imgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *imgDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_img"
}

// Read implements datasource.DataSource.
func (d *imgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	//var state imagesDataSourceModel

	var state []imgModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// query image
	params := param.NewQueryParam()
	params.AddQ("state=Enabled")
	params.AddQ("type=zstack")
	params.AddQ("format!=vmtx")
	params.AddQ("status=Ready")
	params.AddQ("system=false")
	params.AddQ("name=C790123newname")

	images, err := d.client.QueryImage(params)

	//var state []imagesModel
	// image, err := d.client.GetImage(state.Uuid.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Images ",
			err.Error(),
		)

		return
	}

	for r, image := range images {
		/*
			state.Name = types.StringValue(image.Name)
			state.State = types.StringValue(image.State)
			state.Status = types.StringValue(image.Status)
			state.Uuid = types.StringValue(image.UUID)
			state.GuestOsType = types.StringValue(image.GuestOsType)
			state.Format = types.StringValue(image.Format)
			state.Platform = types.StringValue(image.Platform)
			state.Architecture = types.StringValue(string(image.Architecture))
			state.Name_regx = types.StringValue(image.Name)
		*/
		state[r].Name = types.StringValue(image.Name)
		state[r].Uuid = types.StringValue(image.UUID)
		/*
			imageState := imagesModel{
				Name:         types.StringValue(image.Name),
				State:        types.StringValue(image.State),
				Status:       types.StringValue(image.Status),
				Uuid:         types.StringValue(image.UUID),
				GuestOsType:  types.StringValue(image.GuestOsType),
				Format:       types.StringValue(image.Format),
				Platform:     types.StringValue(image.Platform),
				Architecture: types.StringValue(string(image.Architecture)),
				Name_regx:    types.StringValue(image.Name),
			}

			state1 = append(state.Images, imageState)
		*/
	}

	diags = resp.State.Set(ctx, state)

	//nameRegex := state.Name_regx

	//diags = resp.State.Set(ctx, &state)

	//var name_regx types.String
	//diags1 := resp.State.GetAttribute(ctx, path.Root("name_regx"), &name_regx)
	//tflog.Info(ctx, diags1)

	//tflog.Info(ctx, resp.State.Schema.GetAttributes().GetOk)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *imgDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},

			"uuid": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
