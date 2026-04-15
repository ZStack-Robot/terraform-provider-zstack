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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nic_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the vm instance NIC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Terraform resource ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

	attached, err := r.isNicAttached(secgroupUUID, nicUUID)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Security Group Attachment",
			"Could not query NIC membership: "+err.Error(),
		)
		return
	}

	if !attached {
		_, err = r.client.AddVmNicToSecurityGroup(secgroupUUID, param.AddVmNicToSecurityGroupParam{
			BaseParam: param.BaseParam{},
			Params: param.AddVmNicToSecurityGroupParamDetail{
				VmNicUuids: []string{nicUUID},
			},
		})
		if err != nil {
			response.Diagnostics.AddError(
				"Error creating Security Group Attachment",
				"Could not add VM NIC to security group: "+err.Error(),
			)
			return
		}
	}

	tflog.Info(ctx, "Security group attachment created successfully", map[string]interface{}{
		"id":            plan.ID.ValueString(),
		"secgroup_uuid": secgroupUUID,
		"nic_uuid":      nicUUID,
	})

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

	attached, err := r.isNicAttached(secgroupUUID, vmNicUUID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Security Group Attachment",
			"Could not query NIC membership: "+err.Error(),
		)
		return
	}

	if !attached {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s::%s", secgroupUUID, vmNicUUID))
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

	// The SDK's DeleteVmNicFromSecurityGroup is buggy — it sends
	// DELETE /v1/security-groups/{nicUUID} which deletes the security group
	// instead of detaching the NIC.  The correct endpoint is:
	//   DELETE /v1/security-groups/{sgUUID}/vm-instances/nics?vmNicUuids={nicUUID}
	// We call DeleteWithSpec directly to construct the correct URL.
	err := r.client.DeleteWithSpec(
		"v1/security-groups",
		secgroupUUID,
		"vm-instances/nics",
		fmt.Sprintf("vmNicUuids=%s", nicUUID),
		nil,
	)
	if err != nil {
		// Check whether the NIC is actually still attached.
		// If it is no longer attached, the detach error is harmless and we can remove state.
		attached, queryErr := r.isNicAttached(secgroupUUID, nicUUID)
		if queryErr != nil {
			response.Diagnostics.AddError(
				"Error detaching VM NIC from Security Group",
				fmt.Sprintf("Detach failed (%s) and could not verify attachment status: %s", err.Error(), queryErr.Error()),
			)
			return
		}
		if attached {
			response.Diagnostics.AddError(
				"Error detaching VM NIC from Security Group",
				"Could not detach VM NIC from security group: "+err.Error(),
			)
			return
		}
		// NIC is no longer attached — safe to remove state.
	}

	response.State.RemoveResource(ctx)
}

func (r *securityGroupAttachmentResource) isNicAttached(secgroupUUID, nicUUID string) (bool, error) {
	q := param.NewQueryParam()
	q.AddQ("securityGroupUuid=" + secgroupUUID)
	q.AddQ("vmNicUuid=" + nicUUID)
	refs, err := r.client.QueryVmNicInSecurityGroup(&q)
	if err != nil {
		return false, err
	}
	return len(refs) > 0, nil
}
