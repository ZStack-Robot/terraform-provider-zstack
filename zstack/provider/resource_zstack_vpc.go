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
	_ resource.Resource              = &vpcResource{}
	_ resource.ResourceWithConfigure = &vpcResource{}
)

type NetworkServiceProviderInventoryView struct {
	Type                string   `json:"type"`
	Uuid                string   `json:"uuid"`
	NetworkServiceTypes []string `json:"networkServiceTypes"`
}

type AttachNetworkServiceToL3NetworkParam struct {
	BaseParam struct{}
	Params    AttachNetworkServiceToL3NetworkDetailParam
}

type AttachNetworkServiceToL3NetworkDetailParam struct {
	NetworkServices map[string][]string
}

type vpcResource struct {
	client *client.ZSClient
}

type vpcModel struct {
	Uuid              types.String    `tfsdk:"uuid"`
	Name              types.String    `tfsdk:"name"`
	Description       types.String    `tfsdk:"description"`
	L2NetworkUuid     types.String    `tfsdk:"l2_network_uuid"`
	EnableIPAM        types.Bool      `tfsdk:"enable_ipam"`
	SubnetCidr        subnetCidrModel `tfsdk:"subnet_cidr"`
	Dns               types.String    `tfsdk:"dns"`
	VirtualRouterUuid types.String    `tfsdk:"virtual_router_uuid"`
}

type subnetCidrModel struct {
	Name        string `tfsdk:"name"`
	NetworkCidr string `tfsdk:"network_cidr"`
	Gateway     string `tfsdk:"gateway"`
}

func VpcResource() resource.Resource {
	return &vpcResource{}
}

func (r *vpcResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vpcResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc"
}

func (r *vpcResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Virtual Private Cloud (VPC) networks in ZStack. " +
			"A VPC network provides a logically isolated section of the cloud where you can launch resources such as virtual routers, subnets, and DNS services. " +
			"You can define the VPC's properties, such as its name, description, associated L2 network, and subnet CIDR.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VPC network.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC network.",
			},
			"description": schema.StringAttribute{
				Required:    true,
				Description: "A description for the VPC network.",
			},
			"l2_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L2 network associated with this VPC network.",
			},
			"enable_ipam": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable IP Address Management (IPAM) for this VPC network.",
			},
			"dns": schema.StringAttribute{
				Optional:    true,
				Description: "Attach Dns Server for this VPC network.",
			},
			"virtual_router_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "Attach virtual router  for this VPC network.",
			},
			"subnet_cidr": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the subnet CIDR.",
					},
					"network_cidr": schema.StringAttribute{
						Required:    true,
						Description: "The CIDR block for the subnet.",
					},
					"gateway": schema.StringAttribute{
						Required:    true,
						Description: "The gateway IP address for the subnet.",
					},
				},
				Description: "Details of the subnet CIDR to be configured in the VPC network.",
			},
		},
	}
}

func (r *vpcResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vpcModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	networkSvcs, _ := r.client.QueryNetworkServiceProvider(param.NewQueryParam())

	networkServices := make(map[string][]string)

	for _, svc := range networkSvcs {
		switch svc.Type {
		case "vrouter":
			networkServices[svc.Uuid] = []string{"IPsec", "VRouterRoute", "VipQos", "SNAT", "PortForwarding", "Eip", "DNS", "LoadBalancer", "CentralizedDNS"}
		case "Flat":
			networkServices[svc.Uuid] = []string{"DHCP", "Userdata"}
		case "SecurityGroup":
			networkServices[svc.Uuid] = []string{"SecurityGroup"}
		}
	}

	netSvcParam := param.AttachNetworkServiceToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachNetworkServiceToL3NetworkDetailParam{
			NetworkServices: networkServices,
		},
	}

	cidrParam := param.AddIpRangeByNetworkCidrParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpRangeByNetworkCidrDetailParam{
			Name:        plan.SubnetCidr.Name,
			NetworkCidr: plan.SubnetCidr.NetworkCidr,
			Gateway:     plan.SubnetCidr.Gateway,
		},
	}

	dnsParam := param.AddDnsToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AddDnsToL3NetworkDetailParam{
			Dns: plan.Dns.ValueString(),
		},
	}

	attachVRtoVPC := param.AttachL3NetworkToVmParam{
		BaseParam: param.BaseParam{},
		Params:    param.AttachL3NetworkToVmDetailParam{},
	}

	p := param.CreateL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL3NetworkDetailParam{
			Name:          plan.Name.ValueString(),
			Description:   plan.Description.ValueString(),
			L2NetworkUuid: plan.L2NetworkUuid.ValueString(),
			Type:          "L3VpcNetwork",
			EnableIPAM:    plan.EnableIPAM.ValueBool(),
		},
	}

	pvc, err := r.client.CreateL3Network(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create L3VpcNetwork",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(pvc.UUID)
	plan.Name = types.StringValue(pvc.Name)
	plan.Description = types.StringValue(pvc.Description)
	plan.L2NetworkUuid = types.StringValue(pvc.L2NetworkUuid)
	plan.EnableIPAM = types.BoolValue(pvc.EnableIPAM)

	r.client.AttachNetworkServiceToL3Network(pvc.UUID, netSvcParam)
	r.client.AddIpRangeByNetworkCidr(pvc.UUID, cidrParam)
	r.client.AddDnsToL3Network(pvc.UUID, dnsParam)
	r.client.AttachL3NetworkToVm(pvc.UUID, plan.VirtualRouterUuid.ValueString(), attachVRtoVPC)
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *vpcResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vpcModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vpcs, err := r.client.QueryL3Network(param.NewQueryParam())

	//Zql(fmt.Sprintf("query reservedIpRange where uuid='%s'", state.Uuid.ValueString()), &reservedIpRanges, "inventories")
	if err != nil {
		tflog.Warn(ctx, "cannot read vpcs, maybe it has been deleted, set uuid to 'empty'. vpcs was no longer managed by terraform. error: "+err.Error())
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vpc := range vpcs {
		if vpc.UUID == state.Uuid.ValueString() {
			// Update state with the matched subnet details
			state.Uuid = types.StringValue(vpc.UUID)
			state.Name = types.StringValue(vpc.Name)
			state.Description = types.StringValue(vpc.Description)
			state.L2NetworkUuid = types.StringValue(vpc.L2NetworkUuid)
			state.EnableIPAM = types.BoolValue(vpc.EnableIPAM)
			found = true
			break
		}
	}
	if !found {
		// If the subnet is not found, mark it as unmanaged
		tflog.Warn(ctx, "vpc not found. It might have been deleted outside of Terraform.")
		state = vpcModel{
			Uuid: types.StringValue(""),
		}
	}
	// update State
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *vpcResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vpcModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "reserved ip UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteL3Network(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete reserved ip range", ""+err.Error())
		return
	}

}
