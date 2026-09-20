// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ datasource.DataSource              = &slbOfferingDataSource{}
	_ datasource.DataSourceWithConfigure = &slbOfferingDataSource{}
)

type slbOfferingDataSourceModel struct {
	Uuid        types.String       `tfsdk:"uuid"`
	Name        types.String       `tfsdk:"name"`
	NamePattern types.String       `tfsdk:"name_pattern"`
	Filter      []Filter           `tfsdk:"filter"`
	SlbOffering []slbOfferingModel `tfsdk:"slb_offers"`
}

type slbOfferingModel struct {
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
	ZoneUuid              types.String `tfsdk:"zone_uuid"`
	ImageUuid             types.String `tfsdk:"image_uuid"`
	ReservedMemorySize    types.String `tfsdk:"reserved_memory_size"`
}

func ZStackSlbOfferingDataSource() datasource.DataSource {
	return &slbOfferingDataSource{}
}

type slbOfferingDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *slbOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *slbOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_slb_offerings"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *slbOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state slbOfferingDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	params := param.NewQueryParam()

	// 优先检查 `name` 精确查询
	applyUuidOrNameFilter(&params, state.Uuid, state.Name, state.NamePattern)

	slbOffers, err := d.client.QuerySlbOffering(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read SLB offers",
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

	filterSlbOffers, filterDiags := utils.FilterResource(ctx, slbOffers, filters, "slb_offering")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Return a known empty list when no offerings match.
	state.SlbOffering = make([]slbOfferingModel, 0, len(filterSlbOffers))
	for _, slbOffer := range filterSlbOffers {
		slbOfferState := slbOfferingModel{
			Name:              types.StringValue(slbOffer.Name),
			Uuid:              types.StringValue(slbOffer.UUID),
			Description:       types.StringValue(slbOffer.Description),
			CpuNum:            types.Int32Value(int32(slbOffer.CpuNum)),
			CpuSpeed:          types.Int32Value(int32(slbOffer.CpuSpeed)),
			MemorySize:        types.Int64Value(utils.BytesToMB(slbOffer.MemorySize)),
			Type:              types.StringValue(slbOffer.Type),
			AllocatorStrategy: types.StringValue(slbOffer.AllocatorStrategy),

			ZoneUuid:              types.StringValue(slbOffer.ZoneUuid),
			ManagementNetworkUuid: types.StringValue(slbOffer.ManagementNetworkUuid),
			ImageUuid:             types.StringValue(slbOffer.ImageUuid),

			SortKey:            types.Int32Value(int32(slbOffer.SortKey)),
			State:              types.StringValue(slbOffer.State),
			ReservedMemorySize: types.StringValue(fmt.Sprintf("%d", slbOffer.ReservedMemorySize)),
		}

		state.SlbOffering = append(state.SlbOffering, slbOfferState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *slbOfferingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of SLB offers and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "Exact UUID lookup. Recommended for automation: stable across renames, deterministic (0 or 1 match), idempotent. Mutually exclusive with `name` / `name_pattern`.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("name"),
						path.MatchRoot("name_pattern"),
					),
				},
			},
			"name": schema.StringAttribute{
				Description: "Exact name for searching SLB offer",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"slb_offers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the SLB offering.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the SLB offering.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A brief description of the SLB offering.",
						},
						"cpu_num": schema.Int32Attribute{
							Computed:    true,
							Description: "The number of CPUs allocated to the SLB offer.",
						},
						"cpu_speed": schema.Int32Attribute{
							Computed:    true,
							Description: "The speed of each CPU in MHz.",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The memory size allocated to the SLB, in mebibytes (MiB).",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the SLB offering (e.g., SLB).",
						},
						"allocator_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The strategy used for allocating resources to the SLB.",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the zone where the SLB is deployed.",
						},
						"management_network_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the management network connected to the SLB.",
						},

						"image_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the image used by the SLB offer.",
						},
						"sort_key": schema.Int32Attribute{
							Computed:    true,
							Description: "The sort key used for ordering SLB offerings.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the SLB offering (e.g., Enabled, Disabled).",
						},
						"reserved_memory_size": schema.StringAttribute{
							Computed:    true,
							Description: "The amount of memory reserved for the SLB, in bytes.",
						},
					},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				Description: "Filter resources based on any field in the schema. For example, to filter by state, use `name = \"state\"` and `values = [\"Enabled\"]`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the field to filter by (e.g., cpu_num, memory_size, state).",
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
