// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource              = &guestToolsResource{}
	_ resource.ResourceWithConfigure = &guestToolsResource{}
)

type guestToolsResource struct {
	client *client.ZSClient
}

type qgaModel struct {
	ID            types.String `tfsdk:"id"`
	Instance_Uuid types.String `tfsdk:"instance_uuid"`
	Version       types.String `tfsdk:"guest_tools_version"`
	Status        types.String `tfsdk:"guest_tools_status"`
}

func GuestToolsResource() resource.Resource {
	return &guestToolsResource{}
}

func (r *guestToolsResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *guestToolsResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_guest_tools_attachment"
}

func (r *guestToolsResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Attaches guest tools ISO to a ZStack VM instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Same as the vm_instance_uuid. Used for Terraform tracking.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the ZStack VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"guest_tools_version": schema.StringAttribute{
				Computed:    true,
				Description: "Version of the ZStack guest tools installed on the VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guest_tools_status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of the ZStack guest tools on the VM (e.g., 'Connected', 'Disconnected').",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *guestToolsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan qgaModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.Instance_Uuid.IsNull() || plan.Instance_Uuid.IsUnknown() {
		response.Diagnostics.AddError(
			"Error creating Guest Tool Attachment",
			"Could not create guest tool attachment, missing instance_uuid.",
		)
		return
	}
	instance_uuid := plan.Instance_Uuid.ValueString()

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	_, err := r.client.AttachGuestToolsIsoToVm(instance_uuid, param.AttachGuestToolsIsoToVmParam{})

	if err != nil {
		response.Diagnostics.AddError(
			"Error attaching Guest Tools to VM Instance",
			"Could not attach guest tools to VM instance, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Guest tools ISO attached to VM instance successfully.")
	guest_tools, err := r.client.GetVmGuestToolsInfo(instance_uuid)
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading Guest Tool Attachment",
			"Could not read guest tool attachment after attach: "+err.Error(),
		)
		return
	}

	plan.ID = plan.Instance_Uuid
	plan.Status = types.StringValue(guest_tools.Status)
	plan.Version = types.StringValue(guest_tools.Version)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *guestToolsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state qgaModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	guest_tools, err := r.client.GetVmGuestToolsInfo(state.Instance_Uuid.ValueString())

	//Zql(fmt.Sprintf("query reservedIpRange where uuid='%s'", state.Uuid.ValueString()), &reservedIpRanges, "inventories")
	if err != nil {
		if isZStackNotFoundError(err) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query guest tools info. It may have been detached or the VM is powered off: "+err.Error())
		state = qgaModel{
			ID: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	// Update state with the matched subnet details
	state.Status = types.StringValue(guest_tools.Status)
	state.Version = types.StringValue(guest_tools.Version)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *guestToolsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update not supported",
		"Guest Tool Attachment resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *guestToolsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete is a no-op. Guest tools ISO will remain attached.")
}
