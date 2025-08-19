// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

var (
	_ resource.Resource                = &securityGroupRuleResource{}
	_ resource.ResourceWithConfigure   = &securityGroupRuleResource{}
	_ resource.ResourceWithImportState = &securityGroupRuleResource{}
)

type securityGroupRuleResource struct {
	client *client.ZSClient
}

type securityGroupRuleModel struct {
	Uuid                    types.String `tfsdk:"uuid"`                       //Compute
	Name                    types.String `tfsdk:"name"`                       //require
	Description             types.String `tfsdk:"description"`                //Option
	SecurityGroupUuid       types.String `tfsdk:"security_group_uuid"`        //require
	Priority                types.Int32  `tfsdk:"priority"`                   //require
	Direction               types.String `tfsdk:"direction"`                  //require
	Action                  types.String `tfsdk:"action"`                     //require
	State                   types.String `tfsdk:"state"`                      //require. default is enabled
	IpVersion               types.Int32  `tfsdk:"ip_version"`                 //require. default is 4
	Protocol                types.String `tfsdk:"protocol"`                   //requere. TCP, UDP, ICMP
	IpRanges                types.String `tfsdk:"ip_ranges"`                  //requere
	DestinationPortRanges   types.String `tfsdk:"destination_port_ranges"`    //require
	RemoteSecurityGroupUuid types.String `tfsdk:"remote_security_group_uuid"` //Option 关联另一个rule
}

func SecurityGroupRuleResource() resource.Resource {
	return &securityGroupRuleResource{}
}

func (r *securityGroupRuleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *securityGroupRuleResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_networking_secgroup_rule"
}

func (r *securityGroupRuleResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource manages a security group rule in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the security group rule.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name for the security group rule.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the security group rule.",
			},
			"security_group_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the security group to which this rule belongs.",
			},
			"priority": schema.Int32Attribute{
				Required:    true,
				Description: "Priority of the rule. Must start from 1 and be consecutive (e.g., 1, 2, 3, 4). Skipping values (e.g., 1, 3, 10) is not allowed.",
			},
			"direction": schema.StringAttribute{
				Required:    true,
				Description: "The direction of the rule, either 'Ingress' or 'Egress'.",
				Validators: []validator.String{
					stringvalidator.OneOf("Ingress", "Egress"),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "The action to take when the rule matches, either 'ACCEPT' or 'DROP'.",
				Validators: []validator.String{
					stringvalidator.OneOf("ACCEPT", "DROP"),
				},
			},
			"state": schema.StringAttribute{
				Required:    true,
				Description: "The state of the rule, either 'Enabled' or 'Disabled'.",
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
			},
			"ip_version": schema.Int32Attribute{
				Required:    true,
				Description: "The IP version, either '4' or '6'.",
				Validators: []validator.Int32{
					int32validator.OneOf(4, 6),
				},
			},
			"protocol": schema.StringAttribute{
				Required:    true,
				Description: "The protocol, either 'TCP', 'UDP', or 'ICMP'.",
				Validators: []validator.String{
					stringvalidator.OneOf("TCP", "UDP", "ICMP"),
				},
			},
			"ip_ranges": schema.StringAttribute{
				Required:    true,
				Description: "The IP ranges for the rule, e.g., '10.0.0.0/24' or '10.0.0.1-10.0.0.10'. IP ranges. For Ingress rules, maps to `srcIpRange`. For Egress rules, maps to `dstIpRange`.",
			},
			"destination_port_ranges": schema.StringAttribute{
				Optional:    true,
				Description: "The destination port ranges, e.g., '80,443' or '8080-9090'.",
			},
			"remote_security_group_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the remote security group to associate with this rule.",
			},
		},
	}
}

func (r *securityGroupRuleResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var rulePlan securityGroupRuleModel
	diags := request.Plan.Get(ctx, &rulePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	// Prepare the rule to be added
	rule := param.AddSecurityGroupRule{
		RuleType:  rulePlan.Direction.ValueString(), // Maps to "type" in the API
		State:     rulePlan.State.ValueString(),
		Protocol:  rulePlan.Protocol.ValueString(),
		Action:    rulePlan.Action.ValueString(),
		IpVersion: int(rulePlan.IpVersion.ValueInt32()),
		//SrcIpRange: rulePlan.IpRanges.ValueString(), // Maps to srcIpRange in the API
	}

	// Set IP ranges based on direction
	if rulePlan.Direction.ValueString() == "Ingress" {
		rule.SrcIpRange = rulePlan.IpRanges.ValueString()
	} else {
		rule.DstIpRange = rulePlan.IpRanges.ValueString()
	}

	// Add optional fields if provided
	if !rulePlan.Description.IsNull() {
		rule.Description = rulePlan.Description.ValueString()
	}

	if !rulePlan.DestinationPortRanges.IsNull() {
		rule.DstPortRange = rulePlan.DestinationPortRanges.ValueString()
	}

	if !rulePlan.RemoteSecurityGroupUuid.IsNull() {
		rule.RemoteSecurityGroupUuid = rulePlan.RemoteSecurityGroupUuid.ValueString()
	}

	// Create the parameter for adding security group rule
	params := param.AddSecurityGroupRuleParam{
		BaseParam: param.BaseParam{},
		Params: param.AddSecurityGroupRuleDetailParam{
			Rules:    []param.AddSecurityGroupRule{rule},
			Priority: int(rulePlan.Priority.ValueInt32()),
		},
	}

	// Call the API to add the security group rule
	resp, err := r.client.AddSecurityGroupRule(rulePlan.SecurityGroupUuid.ValueString(), params)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to create security group rule",
			fmt.Sprintf("Error creating security group rule: %s", err.Error()),
		)
		return
	}
	if resp == nil || len(resp.Rules) == 0 {
		response.Diagnostics.AddError("Empty response", "AddSecurityGroupRule returned no rules")
		return
	}

	var created *view.SecurityGroupRuleInventoryView
	for i := range resp.Rules {
		rv := resp.Rules[i]
		if rv.Type != rule.RuleType {
			continue
		}
		if rv.Protocol != rule.Protocol {
			continue
		}
		if rv.Action != rule.Action {
			continue
		}
		if rv.State != rule.State {
			continue
		}
		if rv.IpVersion != rule.IpVersion {
			continue
		}
		if rv.Priority != params.Params.Priority {
			continue
		}
		if rule.RuleType == "Ingress" {
			if rv.SrcIpRange != rule.SrcIpRange {
				continue
			}
		} else {
			if rv.DstIpRange != rule.DstIpRange {
				continue
			}
		}
		if rv.DstPortRange != rule.DstPortRange {
			continue
		}
		if rule.Description != "" && rv.Description != rule.Description {
			continue
		}
		created = &rv
		break
	}

	if created == nil {
		response.Diagnostics.AddError(
			"Failed to locate created rule",
			"Could not uniquely identify the created rule in API response; check matching logic or API behavior.")
		return
	}

	rulePlan.Uuid = types.StringValue(created.UUID)
	rulePlan.SecurityGroupUuid = types.StringValue(created.SecurityGroupUuid)
	rulePlan.Direction = types.StringValue(created.Type)
	rulePlan.Action = types.StringValue(created.Action)
	rulePlan.State = types.StringValue(created.State)
	rulePlan.Protocol = types.StringValue(created.Protocol)
	rulePlan.Priority = types.Int32Value(int32(created.Priority))
	rulePlan.IpVersion = types.Int32Value(int32(created.IpVersion))

	if created.Description != "" {
		rulePlan.Description = types.StringValue(created.Description)
	} else {
		rulePlan.Description = types.StringNull()
	}

	if created.Type == "Ingress" {
		rulePlan.IpRanges = types.StringValue(created.SrcIpRange)
	} else {
		rulePlan.IpRanges = types.StringValue(created.DstIpRange)
	}

	if created.DstPortRange != "" {
		rulePlan.DestinationPortRanges = types.StringValue(created.DstPortRange)
	} else {
		rulePlan.DestinationPortRanges = types.StringNull()
	}

	if !rulePlan.Name.IsNull() && rulePlan.Name.ValueString() != "" {
		// keep user provided name
	} else {
		rulePlan.Name = types.StringValue(fmt.Sprintf("%s-%s-Rule", created.Type, created.Protocol))
	}

	// Set the state
	diags = response.State.Set(ctx, &rulePlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Created security group rule", map[string]interface{}{
		"uuid": rulePlan.Uuid.ValueString(),
	})
}

func (r *securityGroupRuleResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state securityGroupRuleModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		response.Diagnostics.AddError("Missing UUID", "Cannot read SecurityGroupRule without UUID")
		return
	}

	//rules, err := r.client.GetSecurityGroupRule(state.Uuid.ValueString())
	rule, err := r.client.GetSecurityGroupRule(state.Uuid.ValueString())

	if err != nil {
		response.Diagnostics.AddError("Fail to get security gourp and rules ", err.Error())
		return
	}

	if rule == nil || rule.UUID == "" {
		response.State.RemoveResource(ctx)
		return
	}

	state.Uuid = types.StringValue(rule.UUID)
	state.SecurityGroupUuid = types.StringValue(rule.SecurityGroupUuid)
	state.Direction = types.StringValue(rule.Type)
	state.Action = types.StringValue(rule.Action)
	state.State = types.StringValue(rule.State)
	state.Protocol = types.StringValue(rule.Protocol)
	state.Priority = types.Int32Value(int32(rule.Priority))
	state.IpVersion = types.Int32Value(int32(rule.IpVersion))

	if rule.Description != "" {
		state.Description = types.StringValue(rule.Description)
	} else {
		state.Description = types.StringNull()
	}

	if rule.DstPortRange != "" {
		state.DestinationPortRanges = types.StringValue(rule.DstPortRange)
	} else {
		state.DestinationPortRanges = types.StringNull()
	}

	if rule.Type == "Ingress" {
		state.IpRanges = types.StringValue(rule.SrcIpRange)
	} else {
		state.IpRanges = types.StringValue(rule.DstIpRange)
	}

	if rule.RemoteSecurityGroupUuid != "" {
		state.RemoteSecurityGroupUuid = types.StringValue(rule.RemoteSecurityGroupUuid)
	} else {
		state.RemoteSecurityGroupUuid = types.StringNull()
	}

	// Update State
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupRuleResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan, state securityGroupRuleModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if r.client == nil {
		response.Diagnostics.AddError("Client Not Configured", "The Client was not properly configured")
		return
	}

	ruleUUID := state.Uuid.ValueString()
	if ruleUUID == "" {
		response.Diagnostics.AddError("Missing Security Group Rule UUID", "UUID is required to update the rule.")
		return
	}

	change := param.UpdateSecurityGroupRuleDetailParam{
		Priority: int(plan.Priority.ValueInt32()),
	}

	if !plan.State.Equal(state.State) {
		change.State = plan.State.ValueString()
	}
	if !plan.Action.Equal(state.Action) {
		change.Action = plan.Action.ValueString()
	}
	if !plan.Protocol.Equal(state.Protocol) {
		change.Protocol = plan.Protocol.ValueString()
	}
	if !plan.Description.Equal(state.Description) {
		if !plan.Description.IsNull() {
			change.Description = plan.Description.ValueString()
		} else {
			change.Description = ""
		}
	}
	if !plan.RemoteSecurityGroupUuid.Equal(state.RemoteSecurityGroupUuid) {
		change.RemoteSecurityGroupUuid = plan.RemoteSecurityGroupUuid.ValueString()
	}
	if !plan.IpRanges.Equal(state.IpRanges) {
		if plan.Direction.ValueString() == "Ingress" {
			change.SrcIpRange = plan.IpRanges.ValueString()
		} else {
			change.DstIpRange = plan.IpRanges.ValueString()
		}
	}
	if !plan.DestinationPortRanges.Equal(state.DestinationPortRanges) {
		if !plan.DestinationPortRanges.IsNull() {
			change.DstPortRange = plan.DestinationPortRanges.ValueString()
		} else {
			change.DstPortRange = ""
		}
	}

	needUpdate := true
	_, err := r.client.UpdateSecurityGroupRule(ruleUUID, param.UpdateSecurityGroupRuleParam{
		BaseParam:               param.BaseParam{},
		ChangeSecurityGroupRule: change,
	})
	if err != nil {
		response.Diagnostics.AddError("Failed to update security group rule", fmt.Sprintf("Error: %s", err))
		return
	}

	rule, err := r.client.GetSecurityGroupRule(ruleUUID)
	if err != nil {
		response.Diagnostics.AddError("Failed to read security group rule after update", fmt.Sprintf("Error: %s", err))
		return
	}
	_ = needUpdate

	plan.Uuid = types.StringValue(rule.UUID)
	plan.SecurityGroupUuid = types.StringValue(rule.SecurityGroupUuid)
	plan.Direction = types.StringValue(rule.Type)
	plan.Action = types.StringValue(rule.Action)
	plan.State = types.StringValue(rule.State)
	plan.Protocol = types.StringValue(rule.Protocol)
	plan.Priority = types.Int32Value(int32(rule.Priority))
	plan.IpVersion = types.Int32Value(int32(rule.IpVersion))
	if rule.Description != "" {
		plan.Description = types.StringValue(rule.Description)
	} else {
		plan.Description = types.StringNull()
	}
	if rule.Type == "Ingress" {
		plan.IpRanges = types.StringValue(rule.SrcIpRange)
	} else {
		plan.IpRanges = types.StringValue(rule.DstIpRange)
	}
	if rule.DstPortRange != "" {
		plan.DestinationPortRanges = types.StringValue(rule.DstPortRange)
	} else {
		plan.DestinationPortRanges = types.StringNull()
	}
	if rule.RemoteSecurityGroupUuid != "" {
		plan.RemoteSecurityGroupUuid = types.StringValue(rule.RemoteSecurityGroupUuid)
	}

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	tflog.Info(ctx, "Updated security group rule", map[string]interface{}{"uuid": plan.Uuid.ValueString()})
}

func (r *securityGroupRuleResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state securityGroupRuleModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "Security group UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteSecurityGroupRule(state.Uuid.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to delete security group",
			"Error "+err.Error(),
		)
		return
	}

}

func (r *securityGroupRuleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// Extract the resource ID (UUID) from the import request
	resourceID := request.ID

	// Validate the import ID format
	if resourceID == "" {
		response.Diagnostics.AddError(
			"Invalid Import ID",
			"The import ID must be the UUID of the security group rule. Format: <rule-uuid>",
		)
		return
	}

	// Check if the client is configured
	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	// Call the API to get the security group rule details
	rule, err := r.client.GetSecurityGroupRule(resourceID)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to import security group rule",
			fmt.Sprintf("Error retrieving security group rule with UUID %s: %s", resourceID, err.Error()),
		)
		return
	}
	if rule == nil || rule.UUID == "" {
		response.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Security group rule with UUID %s does not exist", resourceID),
		)
		return
	}

	// Create a new state model
	state := securityGroupRuleModel{
		Uuid:              types.StringValue(rule.UUID),
		SecurityGroupUuid: types.StringValue(rule.SecurityGroupUuid),
		Direction:         types.StringValue(rule.Type),
		Action:            types.StringValue(rule.Action),
		State:             types.StringValue(rule.State),
		Protocol:          types.StringValue(rule.Protocol),
		Priority:          types.Int32Value(int32(rule.Priority)),
		IpVersion:         types.Int32Value(int32(rule.IpVersion)),
	}

	// Set optional fields
	if rule.Description != "" {
		state.Description = types.StringValue(rule.Description)
	} else {
		state.Description = types.StringNull()
	}

	// Set IP ranges based on direction
	if rule.Type == "Ingress" {
		state.IpRanges = types.StringValue(rule.SrcIpRange)
	} else {
		state.IpRanges = types.StringValue(rule.DstIpRange)
	}

	// Set destination port ranges
	if rule.DstPortRange != "" {
		state.DestinationPortRanges = types.StringValue(rule.DstPortRange)
	} else {
		state.DestinationPortRanges = types.StringNull()
	}

	state.Name = types.StringNull()

	// Set the imported state
	diags := response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Imported security group rule", map[string]interface{}{
		"uuid": resourceID,
	})
}
