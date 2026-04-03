// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vpcFirewallResource{}
	_ resource.ResourceWithConfigure   = &vpcFirewallResource{}
	_ resource.ResourceWithImportState = &vpcFirewallResource{}
)

type vpcFirewallResource struct {
	client *client.ZSClient
}

type vpcFirewallResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VpcUuid     types.String `tfsdk:"vpc_uuid"`
}

func VpcFirewallResource() resource.Resource {
	return &vpcFirewallResource{}
}

func (r *vpcFirewallResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vpcFirewallResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vpc_firewall"
}

func (r *vpcFirewallResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage VPC firewalls in ZStack. " +
			"A VPC firewall provides security and traffic filtering for a VPC.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VPC firewall.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VPC firewall.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VPC firewall.",
			},
			"vpc_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VPC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *vpcFirewallResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vpcFirewallResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateVpcFirewallParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVpcFirewallParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			VpcUuid:     plan.VpcUuid.ValueString(),
		},
	}

	vpcFirewall, err := r.client.CreateVpcFirewall(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create VPC firewall",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vpcFirewall.UUID)
	plan.Name = types.StringValue(vpcFirewall.Name)
	plan.Description = types.StringValue(vpcFirewall.Description)
	plan.VpcUuid = plan.VpcUuid

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcFirewallResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vpcFirewallResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	vpcFirewalls, err := r.client.QueryVpcFirewall(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query VPC firewalls. It may have been deleted.: "+err.Error())
		state = vpcFirewallResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, vpcFirewall := range vpcFirewalls {
		if vpcFirewall.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(vpcFirewall.UUID)
			state.Name = types.StringValue(vpcFirewall.Name)
			state.Description = types.StringValue(vpcFirewall.Description)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "VPC firewall not found. It might have been deleted outside of Terraform.")
		state = vpcFirewallResourceModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcFirewallResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan vpcFirewallResourceModel
	var state vpcFirewallResourceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateVpcFirewallParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVpcFirewallParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	vpcFirewall, err := r.client.UpdateVpcFirewall(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update VPC firewall",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vpcFirewall.UUID)
	plan.Name = types.StringValue(vpcFirewall.Name)
	plan.Description = types.StringValue(vpcFirewall.Description)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vpcFirewallResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vpcFirewallResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "VPC Firewall cannot be deleted via API, removing from state only.")
}

func (r *vpcFirewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
