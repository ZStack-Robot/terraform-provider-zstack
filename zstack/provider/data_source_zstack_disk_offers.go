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
	_ datasource.DataSource              = &diskOfferingDataSource{}
	_ datasource.DataSourceWithConfigure = &diskOfferingDataSource{}
)

type diskOfferingDataSourceModel struct {
	Name         types.String        `tfsdk:"name"`
	NamePattern  types.String        `tfsdk:"name_pattern"`
	Filter       []Filter            `tfsdk:"filter"`
	DiskOffering []diskOfferingModel `tfsdk:"disk_offers"`
}

type diskOfferingModel struct {
	Name              types.String `tfsdk:"name"`
	Uuid              types.String `tfsdk:"uuid"`
	Description       types.String `tfsdk:"description"`
	DiskSize          types.Int64  `tfsdk:"disk_size"`
	Type              types.String `tfsdk:"type"`               // Type
	AllocatorStrategy types.String `tfsdk:"allocator_strategy"` // Allocation strategy
	SortKey           types.Int32  `tfsdk:"sort_key"`
	State             types.String `tfsdk:"state"` // State (Enabled, Disabled)
}

func ZStackDiskOfferingDataSource() datasource.DataSource {
	return &diskOfferingDataSource{}
}

type diskOfferingDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *diskOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *diskOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_disk_offers"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *diskOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state diskOfferingDataSourceModel

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

	diskOffers, err := d.client.QueryDiskOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read disk offers",
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

	filterDiskOffers, filterDiags := utils.FilterResource(ctx, diskOffers, filters, "disk_offer")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, diskOffer := range filterDiskOffers {
		diskOfferState := diskOfferingModel{
			Name:              types.StringValue(diskOffer.Name),
			Uuid:              types.StringValue(diskOffer.UUID),
			Description:       types.StringValue(diskOffer.Description),
			DiskSize:          types.Int64Value(utils.BytesToGB(int64(diskOffer.DiskSize))),
			Type:              types.StringValue(diskOffer.Type),
			AllocatorStrategy: types.StringValue(diskOffer.AllocatorStrategy),
			SortKey:           types.Int32Value(int32(diskOffer.SortKey)),
			State:             types.StringValue(diskOffer.State),
		}

		state.DiskOffering = append(state.DiskOffering, diskOfferState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *diskOfferingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of disk offers and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching  disk offer",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			/*
				"filter": schema.MapAttribute{
					Description: "Key-value pairs to filter disk offering. For example, to filter by State, use `State = \"Enabled\"`.",
					Optional:    true,
					ElementType: types.StringType,
				},*/
			"disk_offers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the disk offering.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the disk offering.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A brief description of the disk offering.",
						},
						"disk_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The disk size allocated to the disk offering, in gigabytes (GB).",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the disk offering.",
						},
						"allocator_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The strategy used for allocating resources to the disk.",
						},
						"sort_key": schema.Int32Attribute{
							Computed:    true,
							Description: "The sort key used for ordering disk offerings.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the disk offering (e.g., Enabled, Disabled).",
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
