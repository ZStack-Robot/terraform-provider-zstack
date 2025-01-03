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
	_ datasource.DataSource              = &instanceOfferingDataSource{}
	_ datasource.DataSourceWithConfigure = &instanceOfferingDataSource{}
)

type instanceOfferingDataSourceModel struct {
	Name             types.String            `tfsdk:"name"`
	NamePattern      types.String            `tfsdk:"name_pattern"`
	Filter           []Filter                `tfsdk:"filter"`
	InstanceOffering []instanceOfferingModel `tfsdk:"instance_offers"`
}

type instanceOfferingModel struct {
	Name              types.String `tfsdk:"name"`
	Uuid              types.String `tfsdk:"uuid"`
	Description       types.String `tfsdk:"description"`
	CpuNum            types.Int32  `tfsdk:"cpu_num"`            // Number of CPUs
	CpuSpeed          types.Int32  `tfsdk:"cpu_speed"`          // CPU speed
	MemorySize        types.Int64  `tfsdk:"memory_size"`        // Memory size
	Type              types.String `tfsdk:"type"`               // Type
	AllocatorStrategy types.String `tfsdk:"allocator_strategy"` // Allocation strategy
	SortKey           types.Int32  `tfsdk:"sort_key"`
	State             types.String `tfsdk:"state"` // State (Enabled, Disabled)
}

func ZStackInstanceOfferingDataSource() datasource.DataSource {
	return &instanceOfferingDataSource{}
}

type instanceOfferingDataSource struct {
	client *client.ZSClient
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *instanceOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *instanceOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance_offers"
}

// Read implements datasource.DataSourceWithConfigure.
func (d *instanceOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state instanceOfferingDataSourceModel

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

	params.AddQ("type=UserVm")

	instanceOffers, err := d.client.QueryInstaceOffering(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read instance offers",
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

	filterInstanceOffers, filterDiags := utils.FilterResource(ctx, instanceOffers, filters, "instance_offer")
	resp.Diagnostics.Append(filterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, instanceOffer := range filterInstanceOffers {
		instanceOfferState := instanceOfferingModel{
			Name:              types.StringValue(instanceOffer.Name),
			Uuid:              types.StringValue(instanceOffer.UUID),
			Description:       types.StringValue(instanceOffer.Description),
			CpuNum:            types.Int32Value(int32(instanceOffer.CpuNum)),
			CpuSpeed:          types.Int32Value(int32(instanceOffer.CpuSpeed)),
			MemorySize:        types.Int64Value(int64(instanceOffer.MemorySize)),
			Type:              types.StringValue(instanceOffer.Type),
			AllocatorStrategy: types.StringValue(instanceOffer.AllocatorStrategy),
			SortKey:           types.Int32Value(int32(instanceOffer.SortKey)),
			State:             types.StringValue(instanceOffer.State),
		}

		state.InstanceOffering = append(state.InstanceOffering, instanceOfferState)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements datasource.DataSourceWithConfigure.
func (d *instanceOfferingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a list of instance offers and their associated attributes from the ZStack environment.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching  instance offer",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			/*
				"filter": schema.MapAttribute{
					Description: "Key-value pairs to filter instance offering. For example, to filter by State, use `State = \"Enabled\"`.",
					Optional:    true,
					ElementType: types.StringType,
				},
			*/
			"instance_offers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the instance offering.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the instance offering.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "A brief description of the instance offering.",
						},
						"cpu_num": schema.Int32Attribute{
							Computed:    true,
							Description: "The number of CPUs allocated to the instance offer.",
						},
						"cpu_speed": schema.Int32Attribute{
							Computed:    true,
							Description: "The speed of each CPU in MHz.",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "The memory size allocated to the instance, in bytes.",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the instance offering.",
						},
						"allocator_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The strategy used for allocating resources to the instance.",
						},
						"sort_key": schema.Int32Attribute{
							Computed:    true,
							Description: "The sort key used for ordering instance offerings.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The current state of the instance offering (e.g., Enabled, Disabled).",
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
