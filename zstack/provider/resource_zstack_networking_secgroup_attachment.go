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
	_ resource.Resource              = &securityGroupAttachmentResource{}
	_ resource.ResourceWithConfigure = &securityGroupAttachmentResource{}
)

type securityGroupAttachmentResource struct {
	client *client.ZSClient
}

type securityGroupAttachmentModel struct {
	SecurityGroupUuid types.String `tfsdk:"secgroup_uuid"`
	VmNicUuid         types.String `tfsdk:"nic_uuid"`
	ID                types.String `tfsdk:"id"`
}

func SecurityGroupAttachmentResource() resource.Resource {
	return &securityGroupAttachmentResource{}
}

func (r *securityGroupAttachmentResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *securityGroupAttachmentResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_networking_secgroup_attachment"
}

func (r *securityGroupAttachmentResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Attach VM instance NICs to security groups in ZStack.",
		Attributes: map[string]schema.Attribute{
			"secgroup_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the security group.",
			},
			"nic_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the vm instance NIC.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform resource ID.",
			},
		},
	}
}

func (r *securityGroupAttachmentResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan securityGroupAttachmentModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	secgroupUUID := plan.SecurityGroupUuid.ValueString()
	nicUUID := plan.VmNicUuid.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf("%s::%s", secgroupUUID, nicUUID))

	tflog.Debug(ctx, "Creating security group attachment", map[string]interface{}{
		"secgroup_uuid": secgroupUUID,
		"nic_uuid":      nicUUID,
	})

	// Check if the attachment already exists. This makes the Create operation idempotent.
	attachedNics, err := r.client.GetSecurityGroup(secgroupUUID)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to Query ZStack Security Groups",
			err.Error(),
		)
		return
	}

	for _, nic := range attachedNics {
		if nic.UUID == nicUUID {
			tflog.Info(ctx, "Security group attachment already exists, importing into state.", map[string]interface{}{
				"id":            plan.ID.ValueString(),
				"secgroup_uuid": secgroupUUID,
				"nic_uuid":      nicUUID,
			})
			// The resource already exists, just set the state.
			diags := response.State.Set(ctx, plan)
			response.Diagnostics.Append(diags...)
			return
		}
	}

	// Step 1: get candidates vm nics
	candidates, err := r.client.GetCandidateVmNicForSecurityGroup(secgroupUUID)
	if err != nil {
		response.Diagnostics.AddError("GetCandidateVmNicForSecurityGroup failed", err.Error())
		return
	}

	// Step 2: Set for query
	allowed := make(map[string]struct{}, len(candidates))
	for _, nic := range candidates {
		allowed[nic.UUID] = struct{}{}
	}

	if _, ok := allowed[nicUUID]; !ok {
		response.Diagnostics.AddError(
			"Invalid NIC UUID",
			fmt.Sprintf("VM NIC UUID %s is not a valid candidate for security group %s", nicUUID, secgroupUUID),
		)
		return
	}

	// Step 3: Add VM NIC to security group
	err = r.client.AddVmNicToSecurityGroup(secgroupUUID, param.AddVmNicToSecurityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.AddVmNicToSecurityGroupDetailParam{
			VmNicUuids: []string{nicUUID},
		},
	})
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to add VM NIC to security group",
			"Error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Security group attachment created successfully", map[string]interface{}{
		"id":            plan.ID.ValueString(),
		"secgroup_uuid": secgroupUUID,
		"nic_uuid":      nicUUID,
	})

	// Set the final state
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *securityGroupAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state securityGroupAttachmentModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	secgroupUUID := state.SecurityGroupUuid.ValueString()
	vmNicUUID := state.VmNicUuid.ValueString()

	candidates, err := r.client.GetCandidateVmNicForSecurityGroup(secgroupUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Security Group Attachment",
			fmt.Sprintf("GetCandidateVmNicForSecurityGroup failed: %v", err),
		)
		return
	}

	for _, candidate := range candidates {
		if candidate.UUID == vmNicUUID {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	state.ID = types.StringValue(fmt.Sprintf("%s-%s", secgroupUUID, vmNicUUID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *securityGroupAttachmentResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// This resource is a simple link, so it doesn't support updates.
	// Any change should be handled by Terraform as a replacement (destroy and create).
	response.Diagnostics.AddError(
		"Update Not Supported",
		"This resource does not support updates. Any changes will require replacement.",
	)
}

func (r *securityGroupAttachmentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state securityGroupAttachmentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// During delete, the state is populated from the last known state.
	// We can directly use the values.
	secgroupUUID := state.SecurityGroupUuid.ValueString()
	nicUUID := state.VmNicUuid.ValueString()

	tflog.Debug(ctx, "Deleting security group attachment", map[string]interface{}{
		"id":            state.ID.ValueString(),
		"secgroup_uuid": secgroupUUID,
		"nic_uuid":      nicUUID,
	})

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	err := r.client.DeleteVmNicFromSecurityGroup(secgroupUUID, []string{nicUUID})
	if err != nil {
		tflog.Warn(ctx, "Failed to delete VM NIC from security group, it might already be gone.", map[string]interface{}{"error": err.Error()})
	}

	tflog.Info(ctx, "Successfully deleted security group attachment.")
	response.State.RemoveResource(ctx)
}
