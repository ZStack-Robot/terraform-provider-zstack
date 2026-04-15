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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &policyRouteRuleResource{}
	_ resource.ResourceWithConfigure   = &policyRouteRuleResource{}
	_ resource.ResourceWithImportState = &policyRouteRuleResource{}
)

type policyRouteRuleResource struct {
	client *client.ZSClient
}

type policyRouteRuleModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	RuleSetUuid types.String `tfsdk:"rule_set_uuid"`
	TableUuid   types.String `tfsdk:"table_uuid"`
	RuleNumber  types.Int64  `tfsdk:"rule_number"`
	DestIp      types.String `tfsdk:"dest_ip"`
	SourceIp    types.String `tfsdk:"source_ip"`
	DestPort    types.String `tfsdk:"dest_port"`
	SourcePort  types.String `tfsdk:"source_port"`
	Protocol    types.String `tfsdk:"protocol"`
	State       types.String `tfsdk:"state"`
}

func PolicyRouteRuleResource() resource.Resource {
	return &policyRouteRuleResource{}
}

func (r *policyRouteRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_route_rule"
}

func (r *policyRouteRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Policy Route Rule resources in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Policy Route Rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rule_set_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the rule set this rule belongs to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"table_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the routing table.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_number": schema.Int64Attribute{
				Required:    true,
				Description: "The rule number (priority).",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"dest_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The destination IP address or CIDR.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The source IP address or CIDR.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dest_port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The destination port or port range.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_port": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The source port or port range.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The protocol (TCP, UDP, ICMP, etc.).",
				Validators: []validator.String{
					stringvalidator.OneOf("TCP", "UDP", "ICMP"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the rule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *policyRouteRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *policyRouteRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyRouteRuleModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleNumber := int(plan.RuleNumber.ValueInt64())
	createParam := param.CreatePolicyRouteRuleParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePolicyRouteRuleParamDetail{
			RuleSetUuid: plan.RuleSetUuid.ValueString(),
			TableUuid:   plan.TableUuid.ValueString(),
			RuleNumber:  ruleNumber,
		},
	}

	if !plan.DestIp.IsNull() {
		createParam.Params.DestIp = stringPtr(plan.DestIp.ValueString())
	}
	if !plan.SourceIp.IsNull() {
		createParam.Params.SourceIp = stringPtr(plan.SourceIp.ValueString())
	}
	if !plan.DestPort.IsNull() {
		createParam.Params.DestPort = stringPtr(plan.DestPort.ValueString())
	}
	if !plan.SourcePort.IsNull() {
		createParam.Params.SourcePort = stringPtr(plan.SourcePort.ValueString())
	}
	if !plan.Protocol.IsNull() {
		createParam.Params.Protocol = stringPtr(plan.Protocol.ValueString())
	}

	result, err := r.client.CreatePolicyRouteRule(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Policy Route Rule",
			"Could not create policy route rule, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.RuleSetUuid = types.StringValue(result.RuleSetUuid)
	plan.TableUuid = types.StringValue(result.TableUuid)
	plan.RuleNumber = types.Int64Value(int64(result.RuleNumber))
	plan.DestIp = stringValueOrNull(result.DestIp)
	plan.SourceIp = stringValueOrNull(result.SourceIp)
	plan.DestPort = stringValueOrNull(result.DestPort)
	plan.SourcePort = stringValueOrNull(result.SourcePort)
	plan.Protocol = stringValueOrNull(result.Protocol)
	plan.State = stringValueOrNull(result.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Policy Route Rule created", map[string]interface{}{
		"uuid":        result.UUID,
		"rule_number": result.RuleNumber,
	})
}

func (r *policyRouteRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyRouteRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := findResourceByQuery(r.client.QueryPolicyRouteRule, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "Policy Route Rule not found, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Policy Route Rule",
			"Could not read policy route rule, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(rule.UUID)
	state.RuleSetUuid = types.StringValue(rule.RuleSetUuid)
	state.TableUuid = types.StringValue(rule.TableUuid)
	state.RuleNumber = types.Int64Value(int64(rule.RuleNumber))
	state.DestIp = stringValueOrNull(rule.DestIp)
	state.SourceIp = stringValueOrNull(rule.SourceIp)
	state.DestPort = stringValueOrNull(rule.DestPort)
	state.SourcePort = stringValueOrNull(rule.SourcePort)
	state.Protocol = stringValueOrNull(rule.Protocol)
	state.State = stringValueOrNull(rule.State)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyRouteRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Policy Route Rule has no Update method - all fields are ForceNew
	// This method is required by the interface but should not be called
	resp.Diagnostics.AddError(
		"Update not supported",
		"Policy Route Rule does not support updates. All fields require replacement.",
	)
}

func (r *policyRouteRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyRouteRuleModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicyRouteRule(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Policy Route Rule",
			"Could not delete policy route rule UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Policy Route Rule deleted", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *policyRouteRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
