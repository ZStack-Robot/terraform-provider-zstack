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
	_ datasource.DataSource              = &sdnControllerDataSource{}
	_ datasource.DataSourceWithConfigure = &sdnControllerDataSource{}
)

func ZStackSdnControllerDataSource() datasource.DataSource {
	return &sdnControllerDataSource{}
}

type sdnControllerDataSource struct {
	client *client.ZSClient
}

type sdnControllerDataSourceModel struct {
	Name           types.String         `tfsdk:"name"`
	NamePattern    types.String         `tfsdk:"name_pattern"`
	Filter         []Filter             `tfsdk:"filter"`
	SdnControllers []sdnControllerModel `tfsdk:"sdn_controllers"`
}

type sdnControllerModel struct {
	Name        types.String `tfsdk:"name"`
	Uuid        types.String `tfsdk:"uuid"`
	Description types.String `tfsdk:"description"`
	Ip          types.String `tfsdk:"ip"`
	Status      types.String `tfsdk:"status"`
	//UserName    types.String `tfsdk:"username"`
	//Passwordd   types.String `tfsdk:"password"`
	VendorType types.String `tfsdk:"vendor_type"`
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *sdnControllerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *sdnControllerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networking_sdn_controllers"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *sdnControllerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sdnControllerDataSourceModel
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

	sdnControllers, err := d.client.QuerySdnController(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read ZStack SDN Controllers ",
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

	filterControllers, filterDiags := utils.FilterResource(ctx, sdnControllers, filters, "sdn_controller")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, sdn := range filterControllers {
		SdnsState := sdnControllerModel{
			Uuid:        types.StringValue(sdn.UUID),
			Name:        types.StringValue(sdn.Name),
			Description: types.StringValue(sdn.Description),
			Ip:          types.StringValue(sdn.Ip),
			Status:      types.StringValue(sdn.Status),
			VendorType:  types.StringValue(sdn.VendorType),
		}

		state.SdnControllers = append(state.SdnControllers, SdnsState)
	}

	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements datasource.DataSourceWithConfigure.
func (d *sdnControllerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of SDN Controllers and their attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching SDN Controllers.",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"sdn_controllers": schema.ListNestedAttribute{
				Description: "List of SDN Controller entries matching the specified filters.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "UUID of the SDN Controller.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the SDN Controller.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "Description of the SDN Controller.",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "IP address of the SDN Controller.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "Operational status of the SDN Controller. (e.g.,Connecting, Connected, Disconnected )",
						},
						"vendor_type": schema.StringAttribute{
							Computed:    true,
							Description: "Vendor type of the SDN Controller (e.g., Ovn).",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter SDN controllers based on any field in the schema. For example, to filter by status, use `name = \"status\"` and `values = [\"Ready\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., status, ip).",
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
