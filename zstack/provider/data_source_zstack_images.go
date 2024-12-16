// Copyright (c) ZStack.io, Inc.

package provider

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
	GuestOsType  types.String `tfsdk:"guest_os_type"`
	Format       types.String `tfsdk:"format"`
	Platform     types.String `tfsdk:"platform"`
	Architecture types.String `tfsdk:"architecture"`
}

type imagesDataSourceModel struct {
	Name        types.String  `tfsdk:"name"`
	NamePattern types.String  `tfsdk:"name_pattern"`
	Images      []imagesModel `tfsdk:"images"`
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
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the ZStack Provider developer. ", req.ProviderData),
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

	//name_regex := state.Name
	params := param.NewQueryParam()

	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	images, err := d.client.QueryImage(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack Images",
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

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *imageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of images and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching images",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"images": schema.ListNestedAttribute{
				Description: "List of Images",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the image",
							Computed:    true,
						},

						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the image",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the image, indicating if it is Enabled or Disabled",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Readiness status of the image (e.g., Ready or Not Ready)",
							Computed:    true,
						},
						"guest_os_type": schema.StringAttribute{
							Description: "Operating system type of the image (e.g., Linux, Windows)",
							Computed:    true,
						},
						"format": schema.StringAttribute{
							Description: "Format of the image, such as qcow2, iso, vmdk, or raw",
							Computed:    true,
						},
						"platform": schema.StringAttribute{
							Description: "Platform of the image, such as Linux, Windows, or Other",
							Computed:    true,
						},
						"architecture": schema.StringAttribute{
							Description: "CPU architecture of the image, such as x86_64, aarch64, mips64, or longarch64",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
