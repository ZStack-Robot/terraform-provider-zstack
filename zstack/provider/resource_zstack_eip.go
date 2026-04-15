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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &eipResource{}
	_ resource.ResourceWithConfigure   = &eipResource{}
	_ resource.ResourceWithImportState = &eipResource{}
)

type eipResource struct {
	client *client.ZSClient
}

type eipModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VipUuid     types.String `tfsdk:"vip_uuid"`
	VmNicUuid   types.String `tfsdk:"vm_nic_uuid"`
}

func EIPResource() resource.Resource {
	return &eipResource{}
}

func (r *eipResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *eipResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_eip"
}

func (r *eipResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage Elastic IP (EIP) services in ZStack. " +
			"An EIP is a virtual IP address that can be associated with a virtual machine NIC. " +
			"It provides external network access for the virtual machine.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VIP network service.",
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
				Description: "A description for the VIP network service.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vip_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the  VIP IP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_nic_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the virtual machine NIC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *eipResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan eipModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateEipParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateEipParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtr(plan.Description.ValueString()),
			VipUuid:     plan.VipUuid.ValueString(),
			VmNicUuid:   stringPtr(plan.VmNicUuid.ValueString()),
		},
	}

	eip, err := r.client.CreateEip(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating EIP",
			"Could not create EIP, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(eip.UUID)
	plan.Name = types.StringValue(eip.Name)
	plan.Description = types.StringValue(eip.Description)
	plan.VipUuid = types.StringValue(eip.VipUuid)
	plan.VmNicUuid = types.StringValue(eip.VmNicUuid)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *eipResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state eipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	eip, err := findResourceByQuery(r.client.QueryEip, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading EIP",
			"Could not read EIP UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(eip.UUID)
	state.Name = types.StringValue(eip.Name)
	state.Description = types.StringValue(eip.Description)
	state.VipUuid = types.StringValue(eip.VipUuid)
	state.VmNicUuid = types.StringValue(eip.VmNicUuid)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *eipResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"EIP resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *eipResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state eipModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteEip(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting EIP",
			"Could not delete EIP, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *eipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
