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
)

var (
	_ resource.Resource              = &vipResource{}
	_ resource.ResourceWithConfigure = &vipResource{}
)

type vipResource struct {
	client *client.ZSClient
}

type vipModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
	IpRangeUuid   types.String `tfsdk:"ip_range_uuid"`
	VIP           types.String `tfsdk:"vip"`
}

func VipResource() resource.Resource {
	return &vipResource{}
}

func (r *vipResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vipResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vip"
}

func (r *vipResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VIP network service.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VIP network service.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the VIP network service.",
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network associated with this VIP network service.",
			},
			"ip_range_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The type of IP range. Possible values depend on the ZStack configuration (e.g., 'Normal' or 'Reserved').",
			},
			"vip": schema.StringAttribute{
				Optional:    true,
				Description: "create vip ip address  for this VPC network.",
			},
		},
	}
}

func (r *vipResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vipModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVipDetailParam{
			Name:          plan.Name.ValueString(),
			Description:   plan.Description.ValueString(), //.EndIp.ValueString(),
			L3NetworkUUID: plan.L3NetworkUuid.ValueString(),
			IpRangeUUID:   plan.IpRangeUuid.ValueString(),
			RequiredIp:    plan.VIP.ValueString(),
		},
	}

	vip, err := r.client.CreateVip(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to add reserved ip range to L3 network",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vip.UUID)
	plan.Name = types.StringValue(vip.Name)
	plan.Description = types.StringValue(vip.Description)
	plan.L3NetworkUuid = types.StringValue(vip.L3NetworkUUID)
	plan.VIP = types.StringValue(vip.Ip)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *vipResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vips, err := r.client.QueryVip(param.NewQueryParam())

	//Zql(fmt.Sprintf("query reservedIpRange where uuid='%s'", state.Uuid.ValueString()), &reservedIpRanges, "inventories")
	if err != nil {
		tflog.Warn(ctx, "cannot read vpcs, maybe it has been deleted, set uuid to 'empty'. vpcs was no longer managed by terraform. error: "+err.Error())
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vip := range vips {
		if vip.UUID == state.Uuid.ValueString() {
			// Update state with the matched subnet details
			state.Uuid = types.StringValue(vip.UUID)
			state.Name = types.StringValue(vip.Name)
			state.Description = types.StringValue(vip.Description)
			state.L3NetworkUuid = types.StringValue(vip.L3NetworkUUID)
			state.VIP = types.StringValue(vip.Ip)
			found = true
			break
		}
	}
	if !found {
		// If the subnet is not found, mark it as unmanaged
		tflog.Warn(ctx, "vpc not found. It might have been deleted outside of Terraform.")
		state = vipModel{
			Uuid: types.StringValue(""),
		}
	}
	// 更新 State
	/*   # 报状态不一致，需修复
		│ When applying changes to zstack_vip.test, provider
	│ "provider[\"zstack.io/terraform-provider-zstack/zstack\"]" produced an
	│ unexpected new value: .vip: was null, but now
	│ cty.StringVal("192.168.110.229").
	*/
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *vipResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *vipResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "reserved ip UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteVip(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete reserved ip range", ""+err.Error())
		return
	}

}
