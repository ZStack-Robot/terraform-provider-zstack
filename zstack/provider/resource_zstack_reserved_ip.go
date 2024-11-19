// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

var (
	_ resource.Resource              = &reservedIpResource{}
	_ resource.ResourceWithConfigure = &reservedIpResource{}
)

type reservedIpResource struct {
	client *client.ZSClient
}

type reservedIpModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	StartIp       types.String `tfsdk:"start_ip"`
	EndIp         types.String `tfsdk:"end_ip"`
	IpVersion     types.Int64  `tfsdk:"ip_version"`
}

func ReservedIpResource() resource.Resource {
	return &reservedIpResource{}
}

func (r *reservedIpResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *reservedIpResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_reserved_ip"
}

func (r *reservedIpResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the reserved ip range.",
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network where the ip range will be reserved.",
			},
			"start_ip": schema.StringAttribute{
				Required:    true,
				Description: "The start IP address of the reserved range.",
			},
			"end_ip": schema.StringAttribute{
				Required:    true,
				Description: "The end IP address of the reserved range.",
			},
			"ip_version": schema.Int64Attribute{
				Computed:    true,
				Description: "The IP version (e.g., 4 for IPv4 or 6 for IPv6) of the reserved range.",
			},
		},
	}
}

func (r *reservedIpResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var reservedIpPlan reservedIpModel
	diags := request.Plan.Get(ctx, &reservedIpPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddReservedIpRangeParam{
		BaseParam: param.BaseParam{},
		Params: param.AddReservedIpRangeDetailParam{
			StartIp: reservedIpPlan.StartIp.ValueString(),
			EndIp:   reservedIpPlan.EndIp.ValueString(),
		},
	}

	ipRange, err := r.client.AddReservedIpRange(reservedIpPlan.L3NetworkUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to add reserved ip range to L3 network",
			"Error "+err.Error(),
		)
		return
	}

	reservedIpPlan.Uuid = types.StringValue(ipRange.Uuid)
	reservedIpPlan.IpVersion = types.Int64Value(int64(ipRange.IpVersion))
	diags = response.State.Set(ctx, reservedIpPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *reservedIpResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state reservedIpModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var reservedIpRanges []view.ReservedIpRangeInventoryView
	_, err := r.client.Zql(fmt.Sprintf("query reservedIpRange where uuid='%s'", state.Uuid.ValueString()), &reservedIpRanges, "inventories")
	if err != nil {
		return
	}

	if len(reservedIpRanges) == 0 {
		state.Uuid = types.StringValue("")
	} else {
		state.Uuid = types.StringValue(reservedIpRanges[0].Uuid)
		state.StartIp = types.StringValue(reservedIpRanges[0].StartIp)
		state.EndIp = types.StringValue(reservedIpRanges[0].EndIp)
		state.IpVersion = types.Int64Value(int64(reservedIpRanges[0].IpVersion))
		state.L3NetworkUuid = types.StringValue(reservedIpRanges[0].L3NetworkUuid)
	}
}

func (r *reservedIpResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *reservedIpResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state reservedIpModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "reserved ip UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteReservedIpRange(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete reserved ip range", ""+err.Error())
		return
	}

}
