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
	_ datasource.DataSource              = &vrouterOfferingDataSource{}
	_ datasource.DataSourceWithConfigure = &vrouterOfferingDataSource{}
)

type vrouterOfferingDataSourceModel struct {
	Name            types.String           `tfsdk:"name"`
	NamePattern     types.String           `tfsdk:"name_pattern"`
	Filter          []Filter               `tfsdk:"filter"`
	VRouterOffering []vrouterOfferingModel `tfsdk:"virtual_router_offers"`
}

type vrouterOfferingModel struct {
	Name                  types.String `tfsdk:"name"`
	Uuid                  types.String `tfsdk:"uuid"`
	Description           types.String `tfsdk:"description"`
	CpuNum                types.Int32  `tfsdk:"cpu_num"`            // Number of CPUs
	CpuSpeed              types.Int32  `tfsdk:"cpu_speed"`          // CPU speed
	MemorySize            types.Int64  `tfsdk:"memory_size"`        // Memory size
	Type                  types.String `tfsdk:"type"`               // Type
	AllocatorStrategy     types.String `tfsdk:"allocator_strategy"` // Allocation strategy
	SortKey               types.Int32  `tfsdk:"sort_key"`
	State                 types.String `tfsdk:"state"` // State (Enabled, Disabled)
	ManagementNetworkUuid types.String `tfsdk:"management_network_uuid"`
	PublicNetworkUuid     types.String `tfsdk:"public_network_uuid"`
	ZoneUuid              types.String `tfsdk:"zone_uuid"`
	ImageUuid             types.String `tfsdk:"image_uuid"`
	IsDefault             types.Bool   `tfsdk:"is_default"`
	ReservedMemorySize    types.String `tfsdk:"reserved_memory_size"`
}

func ZStackVRouterOfferingDataSource() datasource.DataSource {
	return &vrouterOfferingDataSource{}
}

type vrouterOfferingDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *vrouterOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *vrouterOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_router_offers"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *vrouterOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vrouterOfferingDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	// 优先检查 `name` 精确查询
	if !state.Name.IsNull() {
		params.AddQ("name=" + state.Name.ValueString())
	} else if !state.NamePattern.IsNull() {
		params.AddQ("name~=" + state.NamePattern.ValueString())
	}

	vrouterOffers, err := d.client.QueryVirtualRouterOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read virtual router offers",
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

	filterVrouterOffers, filterDiags := utils.FilterResource(ctx, vrouterOffers, filters, "virtual_router_offer")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, vrouterOffer := range filterVrouterOffers {
		vrouterOfferState := vrouterOfferingModel{
			Name:              types.StringValue(vrouterOffer.Name),
			Uuid:              types.StringValue(vrouterOffer.UUID),
			Description:       types.StringValue(vrouterOffer.Description),
			CpuNum:            types.Int32Value(int32(vrouterOffer.CpuNum)),
			CpuSpeed:          types.Int32Value(int32(vrouterOffer.CpuSpeed)),
			MemorySize:        types.Int64Value(utils.BytesToMB(vrouterOffer.MemorySize)),
			Type:              types.StringValue(vrouterOffer.Type),
			AllocatorStrategy: types.StringValue(vrouterOffer.AllocatorStrategy),

			ZoneUuid:              types.StringValue(vrouterOffer.ZoneUuid),
			ManagementNetworkUuid: types.StringValue(vrouterOffer.ManagementNetworkUuid),
			PublicNetworkUuid:     types.StringValue(vrouterOffer.PublicNetworkUuid),
			ImageUuid:             types.StringValue(vrouterOffer.ImageUuid),

			SortKey:            types.Int32Value(int32(vrouterOffer.SortKey)),
			State:              types.StringValue(vrouterOffer.State),
			IsDefault:          types.BoolValue(vrouterOffer.IsDefault),
			ReservedMemorySize: types.StringValue(vrouterOffer.ReservedMemorySize),
		}

		state.VRouterOffering = append(state.VRouterOffering, vrouterOfferState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *vrouterOfferingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of virtual router offers and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching virtual router offer",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			/*
				"filter": schema.MapAttribute{
					Description: "Key-value pairs to filter virtual router offering . For example, to filter by State, use `State = \"Enabled\"`.",
					Optional:    true,
					ElementType: types.StringType,
				},
			*/
			"virtual_router_offers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the virtual router offering.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the virtual router offering.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A brief description of the virtual router offering.",
						},
						"cpu_num": schema.Int32Attribute{
							Computed:    true,
							Description: "The number of CPUs allocated to the virtual router offer.",
						},
						"cpu_speed": schema.Int32Attribute{
							Computed:    true,
							Description: "The speed of each CPU in MHz.",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The memory size allocated to the virtual router, in megabytes (MB).",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the virtual router offering (e.g., VirtualRouter).",
						},
						"allocator_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The strategy used for allocating resources to the virtual router.",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the zone where the virtual router is deployed.",
						},
						"management_network_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the management network connected to the virtual router.",
						},
						"public_network_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the public network connected to the virtual router.",
						},
						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the image used by the virtual router offer.",
						},
						"sort_key": schema.Int32Attribute{
							Computed:    true,
							Description: "The sort key used for ordering virtual router offerings.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the virtual router offering (e.g., Enabled, Disabled).",
						},
						"reserved_memory_size": schema.StringAttribute{
							Computed:    true,
							Description: "The amount of memory reserved for the virtual router, in bytes.",
						},
						"is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "Indicates whether this virtual router offering is the default configuration.",
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
