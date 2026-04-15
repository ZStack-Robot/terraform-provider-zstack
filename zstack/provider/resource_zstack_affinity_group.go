// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &affinityGroupResource{}
	_ resource.ResourceWithConfigure   = &affinityGroupResource{}
	_ resource.ResourceWithImportState = &affinityGroupResource{}
)

type affinityGroupResource struct {
	client *client.ZSClient
}

type affinityGroupResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Policy      types.String `tfsdk:"policy"`
	Type        types.String `tfsdk:"type"`
	ZoneUuid    types.String `tfsdk:"zone_uuid"`
	State       types.String `tfsdk:"state"`
}

func AffinityGroupResource() resource.Resource {
	return &affinityGroupResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *affinityGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *affinityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_affinity_group"
}

// Schema implements resource.Resource.
func (r *affinityGroupResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage affinity groups in ZStack. " +
			"An affinity group defines placement policies for VM instances, such as anti-affinity (spreading VMs across hosts) " +
			"or affinity (keeping VMs on the same host).",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the affinity group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the affinity group.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the affinity group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy": schema.StringAttribute{
				Required:    true,
				Description: "The placement policy: antiSoft, antiHard, proSoft, or proHard.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("antiSoft", "antiHard", "proSoft", "proHard"),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of affinity group. Defaults to 'host'.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the zone for this affinity group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the affinity group (Enabled, Disabled).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *affinityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan affinityGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating affinity group")

	createParam := param.CreateAffinityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAffinityGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Policy:      stringPtr(plan.Policy.ValueString()),
		},
	}

	if !plan.Type.IsNull() && plan.Type.ValueString() != "" {
		createParam.Params.Type = stringPtr(plan.Type.ValueString())
	}

	if !plan.ZoneUuid.IsNull() && plan.ZoneUuid.ValueString() != "" {
		createParam.Params.ZoneUuid = stringPtr(plan.ZoneUuid.ValueString())
	}

	affinityGroup, err := r.client.CreateAffinityGroup(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Affinity Group", "Could not create affinity group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(affinityGroup.UUID)
	plan.Name = types.StringValue(affinityGroup.Name)
	plan.Description = types.StringValue(affinityGroup.Description)
	// API may uppercase the policy (e.g. "antiSoft" -> "ANTISOFT"), preserve user input if case-insensitive match
	if !strings.EqualFold(plan.Policy.ValueString(), affinityGroup.Policy) {
		plan.Policy = types.StringValue(affinityGroup.Policy)
	}
	plan.Type = types.StringValue(affinityGroup.Type)
	plan.State = types.StringValue(affinityGroup.State)
	if affinityGroup.ZoneUuid != "" {
		plan.ZoneUuid = types.StringValue(affinityGroup.ZoneUuid)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *affinityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state affinityGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	affinityGroup, err := findResourceByGet(r.client.GetAffinityGroup, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Affinity Group", "Could not read affinity group UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(affinityGroup.UUID)
	state.Name = types.StringValue(affinityGroup.Name)
	state.Description = types.StringValue(affinityGroup.Description)
	// Preserve existing policy value if case-insensitive match (API uppercases)
	if !strings.EqualFold(state.Policy.ValueString(), affinityGroup.Policy) {
		state.Policy = types.StringValue(affinityGroup.Policy)
	}
	state.Type = types.StringValue(affinityGroup.Type)
	state.State = types.StringValue(affinityGroup.State)
	if affinityGroup.ZoneUuid != "" {
		state.ZoneUuid = types.StringValue(affinityGroup.ZoneUuid)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *affinityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan affinityGroupResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state affinityGroupResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateAffinityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAffinityGroupParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	affinityGroup, err := r.client.UpdateAffinityGroup(state.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Affinity Group", "Could not update affinity group, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(affinityGroup.UUID)
	plan.Name = types.StringValue(affinityGroup.Name)
	plan.Description = types.StringValue(affinityGroup.Description)
	if !strings.EqualFold(plan.Policy.ValueString(), affinityGroup.Policy) {
		plan.Policy = types.StringValue(affinityGroup.Policy)
	}
	plan.Type = types.StringValue(affinityGroup.Type)
	plan.State = types.StringValue(affinityGroup.State)
	if affinityGroup.ZoneUuid != "" {
		plan.ZoneUuid = types.StringValue(affinityGroup.ZoneUuid)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *affinityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state affinityGroupResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteAffinityGroup(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Affinity Group", "Could not delete affinity group, unexpected error: "+err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *affinityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
