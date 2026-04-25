// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &policyResource{}
	_ resource.ResourceWithConfigure   = &policyResource{}
	_ resource.ResourceWithImportState = &policyResource{}
)

type policyResource struct {
	client *client.ZSClient
}

type policyResourceModel struct {
	Uuid        types.String          `tfsdk:"uuid"`
	Name        types.String          `tfsdk:"name"`
	Description types.String          `tfsdk:"description"`
	AccountUuid types.String          `tfsdk:"account_uuid"`
	Statements  []policyStatementModel `tfsdk:"statements"`
}

type policyStatementModel struct {
	Name       types.String `tfsdk:"name"`
	Effect     types.String `tfsdk:"effect"`
	Principals types.List   `tfsdk:"principals"`
	Actions    types.List   `tfsdk:"actions"`
	Resources  types.List   `tfsdk:"resources"`
}

func PolicyResource() resource.Resource {
	return &policyResource{}
}

func (r *policyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Policy in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the policy",
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
				Description: "The description of the policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the account the policy belongs to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"statements": schema.ListNestedAttribute{
				Required:    true,
				Description: "Policy statements defining Allow/Deny rules. Each statement specifies effect, actions, and resources. ZStack Policy is immutable — changes here force resource replacement.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Optional human-readable name for the statement.",
						},
						"effect": schema.StringAttribute{
							Required:    true,
							Description: "Effect of the statement. One of: Allow, Deny.",
							Validators: []validator.String{
								stringvalidator.OneOf("Allow", "Deny"),
							},
						},
						"principals": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Principals this statement applies to (optional).",
						},
						"actions": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "API actions this statement covers, e.g. [\"vm:Start\", \"vm:Stop\"].",
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
						"resources": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "Resource UUIDs or patterns this statement applies to (optional).",
						},
					},
				},
			},
		},
	}
}

func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// BUG-054: build real statements from plan instead of sending empty array.
	statements, sdiags := buildPolicyStatementParams(ctx, plan.Statements)
	resp.Diagnostics.Append(sdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreatePolicyParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePolicyParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Statements:  statements,
		},
	}

	result, err := r.client.CreatePolicy(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Policy",
			"Could not create policy, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	// Description is not returned by CreatePolicy
	plan.AccountUuid = types.StringValue(result.AccountUuid)

	tflog.Trace(ctx, "created a policy")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// buildPolicyStatementParams converts the Terraform plan's statements list into
// the SDK's PolicyStatementParam slice. Returns diagnostics for any decode errors.
func buildPolicyStatementParams(ctx context.Context, in []policyStatementModel) ([]param.PolicyStatementParam, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := make([]param.PolicyStatementParam, 0, len(in))
	for i, s := range in {
		stmt := param.PolicyStatementParam{
			Name:   s.Name.ValueString(),
			Effect: s.Effect.ValueString(),
		}
		if !s.Principals.IsNull() && !s.Principals.IsUnknown() {
			var principals []string
			d := s.Principals.ElementsAs(ctx, &principals, false)
			diags.Append(d...)
			if d.HasError() {
				diags.AddError(
					fmt.Sprintf("Invalid statements[%d].principals", i),
					"Failed to decode principals list",
				)
				return nil, diags
			}
			stmt.Principals = principals
		}
		if !s.Actions.IsNull() && !s.Actions.IsUnknown() {
			var actions []string
			d := s.Actions.ElementsAs(ctx, &actions, false)
			diags.Append(d...)
			if d.HasError() {
				return nil, diags
			}
			stmt.Actions = actions
		}
		if !s.Resources.IsNull() && !s.Resources.IsUnknown() {
			var resources []string
			d := s.Resources.ElementsAs(ctx, &resources, false)
			diags.Append(d...)
			if d.HasError() {
				return nil, diags
			}
			stmt.Resources = resources
		}
		out = append(out, stmt)
	}
	return out, diags
}

func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := findResourceByQuery(r.client.QueryPolicy, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Policy",
			"Could not read policy UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(policy.Name)
	state.AccountUuid = types.StringValue(policy.AccountUuid)

	// BUG-054: map statements from API view back to state.
	stmts := make([]policyStatementModel, 0, len(policy.Statements))
	for _, s := range policy.Statements {
		principals, pdiags := types.ListValueFrom(ctx, types.StringType, s.Principals)
		resp.Diagnostics.Append(pdiags...)
		actions, adiags := types.ListValueFrom(ctx, types.StringType, s.Actions)
		resp.Diagnostics.Append(adiags...)
		resources, rdiags := types.ListValueFrom(ctx, types.StringType, s.Resources)
		resp.Diagnostics.Append(rdiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		stmts = append(stmts, policyStatementModel{
			Name:       stringValueOrNull(s.Name),
			Effect:     types.StringValue(s.Effect),
			Principals: principals,
			Actions:    actions,
			Resources:  resources,
		})
	}
	state.Statements = stmts

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update operation available for Policy resource
	resp.Diagnostics.AddError(
		"Update not supported",
		"Policy resource does not support updates. Please recreate the resource instead.",
	)
}

func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicy(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Policy",
			"Could not delete policy, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted a policy")
}

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
