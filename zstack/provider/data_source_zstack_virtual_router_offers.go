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
	_ datasource.DataSource              = &vrouterOfferingDataSource{}
	_ datasource.DataSourceWithConfigure = &vrouterOfferingDataSource{}
)

type vrouterOfferingDataSourceModel struct {
	Name            types.String           `tfsdk:"name"`
	NamePattern     types.String           `tfsdk:"name_pattern"`
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
	for _, vrouterOffer := range vrouterOffers {
		vrouterOfferState := vrouterOfferingModel{
			Name:              types.StringValue(vrouterOffer.Name),
			Uuid:              types.StringValue(vrouterOffer.UUID),
			Description:       types.StringValue(vrouterOffer.Description),
			CpuNum:            types.Int32Value(int32(vrouterOffer.CpuNum)),
			CpuSpeed:          types.Int32Value(int32(vrouterOffer.CpuSpeed)),
			MemorySize:        types.Int64Value(int64(vrouterOffer.MemorySize)),
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
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Exact name for searching virtual router offer",
				Optional:    true,
			},
			"name_pattern": schema.StringAttribute{
				Description: "Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.",
				Optional:    true,
			},
			"virtual_router_offers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The unique identifier (UUID) of the virtual router offer.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the virtual router offer.",
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: "The type of hypervisor on which the virtual router is running (e.g., KVM, VMware).",
						},
						"cpu_num": schema.Int32Attribute{
							Computed:    true,
							Description: "Specifies the type of appliance VM for the virtual router.",
						},
						"cpu_speed": schema.Int32Attribute{
							Computed:    true,
							Description: "The current state of the virtual router (e.g., Running, Stopped).",
						},
						"memory_size": schema.Int64Attribute{
							Computed:    true,
							Description: "Operational status of the virtual router (e.g., Connected, Disconnected)",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the virtual router (e.g., UserVm or SystemVm).",
						},
						"allocator_strategy": schema.StringAttribute{
							Computed:    true,
							Description: "The high-availability (HA) status of the virtual router.",
						},
						"zone_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the zone in which the virtual router is located.",
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
							Description: "The UUID of the host on which the virtual router is running.",
						},
						"sort_key": schema.Int32Attribute{
							Computed:    true,
							Description: "The UUID of the instance offering assigned to the virtual router.",
						},
						"state": schema.StringAttribute{
							Computed:    true,
							Description: "The platform (e.g., Linux, Windows) on which the virtual router is running.",
						},
						"reserved_memory_size": schema.StringAttribute{
							Computed:    true,
							Description: "The CPU architecture (e.g., x86_64, ARM) of the virtual router.",
						},
						"is_default": schema.BoolAttribute{
							Computed:    true,
							Description: "The CPU architecture (e.g., x86_64, ARM) of the virtual router.",
						},
					},
				},
			},
		},
	}
}
