// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"zstack.io/zstack-sdk-go/pkg/client"
)

var (
	_ datasource.DataSource              = &virtualRouterImageDataSource{}
	_ datasource.DataSourceWithConfigure = &virtualRouterImageDataSource{}
)

type virtualRouterImageDataSource struct {
	client *client.ZSClient
}

type virtualRouterImagesModel struct {
	Name         types.String `tfsdk:"name"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
	Uuid         types.String `tfsdk:"uuid"`
	GuestOsType  types.String `tfsdk:"guest_os_type"`
	Format       types.String `tfsdk:"format"`
	Platform     types.String `tfsdk:"platform"`
	Architecture types.String `tfsdk:"architecture"`
}

type virtualRouterImagesDataSourceModel struct {
	Name        types.String               `tfsdk:"name"`
	NamePattern types.String               `tfsdk:"name_pattern"`
	Images      []virtualRouterImagesModel `tfsdk:"images"`
}

func ZStackVirtualRouterImageDataSource() datasource.DataSource {
	return &virtualRouterImageDataSource{}
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *virtualRouterImageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *virtualRouterImageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_router_images"
}

// Read implements datasource.DataSource.
func (d *virtualRouterImageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state virtualRouterImagesDataSourceModel
	//var state virtualRouterImagesModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var query string
	if !state.Name.IsNull() {
		query = fmt.Sprintf("query Image where __systemTag__='applianceType::vrouter' and name='%s'", state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		query = fmt.Sprintf("query Image where __systemTag__='applianceType::vrouter' and name like '%s'", state.NamePattern.ValueString())
	} else {
		query = "query Image where __systemTag__='applianceType::vrouter'"
	}

	//query := "query Image where __systemTag__='applianceType::vrouter'"
	var zqlResponse struct {
		Results []struct {
			Inventories []struct {
				Name         string `json:"name"`
				State        string `json:"state"`
				Status       string `json:"status"`
				Uuid         string `json:"uuid"`
				GuestOsType  string `json:"guestOsType"`
				Format       string `json:"format"`
				Platform     string `json:"platform"`
				Architecture string `json:"architecture"`
			} `json:"inventories"`
		} `json:"results"`
	}

	_, err := d.client.Zql(query, &zqlResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to execute ZQL query, got error: %s", err))
		return
	}

	// Map the results to the model
	for _, result := range zqlResponse.Results {
		for _, inventory := range result.Inventories {
			image := virtualRouterImagesModel{
				Name:         types.StringValue(inventory.Name),
				State:        types.StringValue(inventory.State),
				Status:       types.StringValue(inventory.Status),
				Uuid:         types.StringValue(inventory.Uuid),
				GuestOsType:  types.StringValue(inventory.GuestOsType),
				Format:       types.StringValue(inventory.Format),
				Platform:     types.StringValue(inventory.Platform),
				Architecture: types.StringValue(inventory.Architecture),
			}
			state.Images = append(state.Images, image)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSource.
func (d *virtualRouterImageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of virtual router images and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching virtual router images",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"images": schema.ListNestedAttribute{
				Description: "List of virtual router Images",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the virtual router image",
							Computed:    true,
						},

						"uuid": schema.StringAttribute{
							Description: "UUID identifier of the virtual router image",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "State of the virtual router image, indicating if it is Enabled or Disabled",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Readiness status of the virtual router image (e.g., Ready or Not Ready)",
							Computed:    true,
						},
						"guest_os_type": schema.StringAttribute{
							Description: "Operating system type of the virtual router image (e.g., Linux, Windows)",
							Computed:    true,
						},
						"format": schema.StringAttribute{
							Description: "Format of the virtual router image, such as qcow2, iso, vmdk, or raw",
							Computed:    true,
						},
						"platform": schema.StringAttribute{
							Description: "Platform of the virtual router image, such as Linux, Windows, or Other",
							Computed:    true,
						},
						"architecture": schema.StringAttribute{
							Description: "CPU architecture of the virtual router image, such as x86_64, aarch64, mips64, or longarch64",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
