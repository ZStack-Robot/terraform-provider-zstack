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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &iam2OrganizationResource{}
	_ resource.ResourceWithConfigure   = &iam2OrganizationResource{}
	_ resource.ResourceWithImportState = &iam2OrganizationResource{}
)

type iam2OrganizationResource struct {
	client *client.ZSClient
}

type iam2OrganizationResourceModel struct {
	Uuid                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Type                 types.String `tfsdk:"type"`
	ParentUuid           types.String `tfsdk:"parent_uuid"`
	State                types.String `tfsdk:"state"`
	SrcType              types.String `tfsdk:"src_type"`
	RootOrganizationUuid types.String `tfsdk:"root_organization_uuid"`
}

func IAM2OrganizationResource() resource.Resource {
	return &iam2OrganizationResource{}
}

func (r *iam2OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *iam2OrganizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam2_organization"
}

func (r *iam2OrganizationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an IAM2 Organization in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the IAM2 organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the IAM2 organization",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the IAM2 organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the IAM2 organization (e.g., Department)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parent_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the parent organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the IAM2 organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"src_type": schema.StringAttribute{
				Computed:    true,
				Description: "The source type of the IAM2 organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"root_organization_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the root organization",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *iam2OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam2OrganizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := param.CreateIAM2OrganizationParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateIAM2OrganizationParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Type:        plan.Type.ValueString(),
			ParentUuid:  stringPtrOrNil(plan.ParentUuid.ValueString()),
		},
	}

	result, err := r.client.CreateIAM2Organization(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating IAM2 Organization",
			"Could not create IAM2 organization, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.ParentUuid = stringValueOrNull(result.ParentUuid)
	plan.State = types.StringValue(result.State)
	plan.SrcType = types.StringValue(result.SrcType)
	plan.RootOrganizationUuid = types.StringValue(result.RootOrganizationUuid)

	tflog.Trace(ctx, "created an IAM2 organization")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iam2OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam2OrganizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := findResourceByQuery(r.client.QueryIAM2Organization, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "IAM2 organization not found, removing from state", map[string]interface{}{
				"uuid": state.Uuid.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading IAM2 Organization",
			"Could not read IAM2 organization UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(org.Name)
	state.Description = stringValueOrNull(org.Description)
	state.Type = types.StringValue(org.Type)
	state.ParentUuid = stringValueOrNull(org.ParentUuid)
	state.State = types.StringValue(org.State)
	state.SrcType = types.StringValue(org.SrcType)
	state.RootOrganizationUuid = types.StringValue(org.RootOrganizationUuid)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *iam2OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam2OrganizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state iam2OrganizationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateIAM2OrganizationParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateIAM2OrganizationParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdateIAM2Organization(state.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating IAM2 Organization",
			"Could not update IAM2 organization, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.ParentUuid = stringValueOrNull(result.ParentUuid)
	plan.State = types.StringValue(result.State)
	plan.SrcType = types.StringValue(result.SrcType)
	plan.RootOrganizationUuid = types.StringValue(result.RootOrganizationUuid)

	tflog.Trace(ctx, "updated an IAM2 organization")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *iam2OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam2OrganizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIAM2Organization(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting IAM2 Organization",
			"Could not delete IAM2 organization, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted an IAM2 organization")
}

func (r *iam2OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
