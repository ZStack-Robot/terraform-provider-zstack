// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vpcResource{}
	_ resource.ResourceWithConfigure   = &vpcResource{}
	_ resource.ResourceWithImportState = &vpcResource{}
)

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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC network.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VPC network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"l2_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the L2 network associated with this VPC network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable_ipam": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable IP Address Management (IPAM) for this VPC network.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"dns": schema.StringAttribute{
				Optional:    true,
				Description: "Attach Dns Server for this VPC network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"virtual_router_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "Attach virtual router  for this VPC network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet_cidr": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:    true,
						Description: "The name of the subnet CIDR.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"network_cidr": schema.StringAttribute{
						Required:    true,
						Description: "The CIDR block for the subnet.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"gateway": schema.StringAttribute{
						Required:    true,
						Description: "The gateway IP address for the subnet.",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
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

	q := param.NewQueryParam()
	networkSvcs, _ := r.client.QueryNetworkServiceProvider(&q)

	networkServices := make(map[string]interface{})

	for _, svc := range networkSvcs {
		switch svc.Type {
		case "vrouter":
			networkServices[svc.UUID] = []string{"IPsec", "VRouterRoute", "VipQos", "SNAT", "PortForwarding", "Eip", "DNS", "LoadBalancer", "CentralizedDNS"}
		case "Flat":
			networkServices[svc.UUID] = []string{"DHCP", "Userdata"}
		case "SecurityGroup":
			networkServices[svc.UUID] = []string{"SecurityGroup"}
		}
	}

	netSvcParam := param.AttachNetworkServiceToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AttachNetworkServiceToL3NetworkParamDetail{
			NetworkServices: networkServices,
		},
	}

	cidrParam := param.AddIpRangeByNetworkCidrParam{
		BaseParam: param.BaseParam{},
		Params: param.AddIpRangeByNetworkCidrParamDetail{
			Name:        plan.SubnetCidr.Name,
			NetworkCidr: plan.SubnetCidr.NetworkCidr,
			Gateway:     stringPtr(plan.SubnetCidr.Gateway),
		},
	}

	dnsParam := param.AddDnsToL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.AddDnsToL3NetworkParamDetail{
			Dns: plan.Dns.ValueString(),
		},
	}

	attachVRtoVPC := param.AttachL3NetworkToVmParam{
		BaseParam: param.BaseParam{},
		Params:    param.AttachL3NetworkToVmParamDetail{},
	}

	p := param.CreateL3NetworkParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateL3NetworkParamDetail{
			Name:          plan.Name.ValueString(),
			Description:   stringPtr(plan.Description.ValueString()),
			L2NetworkUuid: plan.L2NetworkUuid.ValueString(),
			Type:          stringPtr("L3VpcNetwork"),
			EnableIPAM:    boolPtr(plan.EnableIPAM.ValueBool()),
		},
	}

	pvc, err := r.client.CreateL3Network(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating VPC",
			"Could not create vpc, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(pvc.UUID)
	plan.Name = types.StringValue(pvc.Name)
	plan.Description = types.StringValue(pvc.Description)
	plan.L2NetworkUuid = types.StringValue(pvc.L2NetworkUuid)
	plan.EnableIPAM = types.BoolValue(pvc.EnableIPAM)

	// Save partial state so the L3 network UUID is tracked even if follow-up steps fail
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.AttachNetworkServiceToL3Network(pvc.UUID, netSvcParam); err != nil {
		response.Diagnostics.AddError(
			"Error attaching network services to VPC",
			"Could not attach network services to vpc UUID "+pvc.UUID+": "+err.Error(),
		)
		return
	}

	if _, err := r.client.AddIpRangeByNetworkCidr(pvc.UUID, cidrParam); err != nil {
		response.Diagnostics.AddError(
			"Error adding IP range to VPC",
			"Could not add IP range to vpc UUID "+pvc.UUID+": "+err.Error(),
		)
		return
	}

	if !plan.Dns.IsNull() && plan.Dns.ValueString() != "" {
		if _, err := r.client.AddDnsToL3Network(pvc.UUID, dnsParam); err != nil {
			response.Diagnostics.AddError(
				"Error adding DNS to VPC",
				"Could not add DNS to vpc UUID "+pvc.UUID+": "+err.Error(),
			)
			return
		}
	}

	if !plan.VirtualRouterUuid.IsNull() && plan.VirtualRouterUuid.ValueString() != "" {
		if _, err := r.client.AttachL3NetworkToVm(plan.VirtualRouterUuid.ValueString(), pvc.UUID, attachVRtoVPC); err != nil {
			response.Diagnostics.AddError(
				"Error attaching VPC to Virtual Router",
				"Could not attach vpc UUID "+pvc.UUID+" to virtual router UUID "+plan.VirtualRouterUuid.ValueString()+": "+err.Error(),
			)
			return
		}
	}
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

	vpc, err := findResourceByQuery(r.client.QueryL3Network, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "cannot read vpcs, maybe it has been deleted, set uuid to 'empty'. vpcs was no longer managed by terraform. error: "+err.Error())
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(vpc.UUID)
	state.Name = types.StringValue(vpc.Name)
	state.Description = types.StringValue(vpc.Description)
	state.L2NetworkUuid = types.StringValue(vpc.L2NetworkUuid)
	state.EnableIPAM = types.BoolValue(vpc.EnableIPAM)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"VPC resource does not support updates. Please recreate the resource instead.",
	)
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
		response.Diagnostics.AddError(
			"Error deleting VPC",
			"Could not delete vpc, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *vpcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
