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
	_ resource.Resource              = &subnetResource{}
	_ resource.ResourceWithConfigure = &subnetResource{}
)

type subnetResource struct {
	client *client.ZSClient
}

type subnetModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	StartIp       types.String `tfsdk:"start_ip"`
	EndIp         types.String `tfsdk:"end_ip"`
	Netmask       types.String `tfsdk:"netmask"`
	Gateway       types.String `tfsdk:"gateway"`
	IpRangeType   types.String `tfsdk:"ip_range_type"`
	L3NetworkUuid types.String `tfsdk:"l3_network_uuid"`
}

func SubnetResource() resource.Resource {
	return &subnetResource{}
}

func (r *subnetResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *subnetResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_subnet"
}

func (r *subnetResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the subnet.",
			},
			"l3_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L3 network to which the subnet belongs.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the subnet. This is a user-defined identifier for the subnet.",
			},
			"start_ip": schema.StringAttribute{
				Required:    true,
				Description: "The starting IP address of the subnet range.",
			},
			"end_ip": schema.StringAttribute{
				Required:    true,
				Description: "The ending IP address of the subnet range.",
			},
			"netmask": schema.StringAttribute{
				Required:    true,
				Description: "The subnet mask, used to define the network portion of an IP address.",
			},
			"gateway": schema.StringAttribute{
				Required:    true,
				Description: "The default gateway for the subnet.",
			},
			"ip_range_type": schema.StringAttribute{
				Optional:    true,
				Description: "The type of IP range. Possible values depend on the ZStack configuration (e.g., 'Normal' or 'Reserved').",
			},
		},
	}
}

func (r *subnetResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan subnetModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddIpRangeParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpRangeDetailParam{
			Name:        plan.Name.ValueString(),
			StartIp:     plan.StartIp.ValueString(),
			EndIp:       plan.EndIp.ValueString(),
			Netmask:     plan.Netmask.ValueString(),
			Gateway:     plan.Gateway.ValueString(),
			IpRangeType: plan.IpRangeType.ValueString(),
		},
	}

	subnet, err := r.client.AddIpRange(plan.L3NetworkUuid.ValueString(), p)

	//AddReservedIpRange(plan.Uuid.String().L3NetworkUuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to add subnet(ip range) to L3 network",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(subnet.Uuid)
	plan.Name = types.StringValue(subnet.Name)
	plan.StartIp = types.StringValue(subnet.StartIp)
	plan.EndIp = types.StringValue(subnet.EndIp)
	plan.Netmask = types.StringValue(subnet.Netmask)
	plan.Gateway = types.StringValue(subnet.Gateway)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *subnetResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state subnetModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	subnets, err := r.client.QueryIpRange(param.NewQueryParam())

	//Zql(fmt.Sprintf("query reservedIpRange where uuid='%s'", state.Uuid.ValueString()), &reservedIpRanges, "inventories")
	if err != nil {
		tflog.Warn(ctx, "cannot read subnet, maybe it has been deleted, set uuid to 'empty'. subnet was no longer managed by terraform. error: "+err.Error())
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	// Iterate over the fetched subnets and match the UUID
	found := false
	for _, subnet := range subnets {
		if subnet.Uuid == state.Uuid.ValueString() {
			// Update state with the matched subnet details
			state = subnetModel{
				Uuid:          types.StringValue(subnet.Uuid),
				Name:          types.StringValue(subnet.Name),
				StartIp:       types.StringValue(subnet.StartIp),
				EndIp:         types.StringValue(subnet.EndIp),
				Netmask:       types.StringValue(subnet.Netmask),
				Gateway:       types.StringValue(subnet.Gateway),
				L3NetworkUuid: types.StringValue(subnet.L3NetworkUuid),
			}
			found = true
			break
		}
	}

	if !found {
		// If the subnet is not found, mark it as unmanaged
		tflog.Warn(ctx, "Subnet not found. It might have been deleted outside of Terraform.")
		state = subnetModel{
			Uuid: types.StringValue(""),
		}
	}

	// 更新 State
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *subnetResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *subnetResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state subnetModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Subnet name is empty, skipping delete.")
		return
	}

	err := r.client.DeleteIpRange(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("Failed to delete subnet(ip range)", ""+err.Error())
		return
	}

}
