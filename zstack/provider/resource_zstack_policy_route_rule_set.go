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
	_ resource.Resource                = &policyRouteRuleSetResource{}
	_ resource.ResourceWithConfigure   = &policyRouteRuleSetResource{}
	_ resource.ResourceWithImportState = &policyRouteRuleSetResource{}
)

type policyRouteRuleSetResource struct {
	client *client.ZSClient
}

type policyRouteRuleSetModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VrouterUuid types.String `tfsdk:"vrouter_uuid"`
	Type        types.String `tfsdk:"type"`
}

func PolicyRouteRuleSetResource() resource.Resource {
	return &policyRouteRuleSetResource{}
}

func (r *policyRouteRuleSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_route_rule_set"
}

func (r *policyRouteRuleSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Policy Route Rule Set resources in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the Policy Route Rule Set.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Policy Route Rule Set.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the Policy Route Rule Set.",
			},
			"vrouter_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the virtual router.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the Policy Route Rule Set.",
			},
		},
	}
}

func (r *policyRouteRuleSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyRouteRuleSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyRouteRuleSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreatePolicyRouteRuleSetParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePolicyRouteRuleSetParamDetail{
			Name:        plan.Name.ValueString(),
			VRouterUuid: plan.VrouterUuid.ValueString(),
		},
	}

	if !plan.Description.IsNull() {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}
	if !plan.Type.IsNull() {
		createParam.Params.Type = stringPtr(plan.Type.ValueString())
	}

	result, err := r.client.CreatePolicyRouteRuleSet(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Policy Route Rule Set",
			"Could not create Policy Route Rule Set, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = stringValueOrNull(result.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Policy Route Rule Set created", map[string]interface{}{
		"uuid": result.UUID,
		"name": result.Name,
	})
}

func (r *policyRouteRuleSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyRouteRuleSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	results, err := r.client.QueryPolicyRouteRuleSet(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query Policy Route Rule Set: "+err.Error())
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, result := range results {
		if result.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(result.UUID)
			state.Name = types.StringValue(result.Name)
			state.Description = stringValueOrNull(result.Description)
			// VrouterUuid is not in the result, keep from state
			state.Type = stringValueOrNull(result.Type)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Policy Route Rule Set not found, removing from state")
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *policyRouteRuleSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyRouteRuleSetModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdatePolicyRouteRuleSetParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdatePolicyRouteRuleSetParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdatePolicyRouteRuleSet(plan.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Policy Route Rule Set",
			"Could not update Policy Route Rule Set, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	// VrouterUuid is not in the result, keep from plan
	plan.Type = stringValueOrNull(result.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Policy Route Rule Set updated", map[string]interface{}{
		"uuid": result.UUID,
		"name": result.Name,
	})
}

func (r *policyRouteRuleSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyRouteRuleSetModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicyRouteRuleSet(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Policy Route Rule Set",
			"Could not delete Policy Route Rule Set, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Policy Route Rule Set deleted", map[string]interface{}{
		"uuid": state.Uuid.ValueString(),
	})
}

func (r *policyRouteRuleSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
