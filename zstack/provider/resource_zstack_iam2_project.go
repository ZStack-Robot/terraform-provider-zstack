// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &iam2ProjectResource{}
	_ resource.ResourceWithConfigure   = &iam2ProjectResource{}
	_ resource.ResourceWithImportState = &iam2ProjectResource{}
)

type iam2ProjectResource struct {
	client *client.ZSClient
}

type iam2ProjectResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.String `tfsdk:"state"`
}

func IAM2ProjectResource() resource.Resource {
	return &iam2ProjectResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *iam2ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *iam2ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam2_project"
}

// Schema implements resource.Resource.
func (r *iam2ProjectResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage IAM2 projects in ZStack. " +
			"An IAM2 project provides multi-tenant isolation and resource management capabilities.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (UUID) of the IAM2 project.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the IAM2 project.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the IAM2 project.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the IAM2 project (Enabled, Disabled).",
			},
		},
	}
}

// Create implements resource.Resource.
func (r *iam2ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan iam2ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating IAM2 project")

	createParam := param.CreateIAM2ProjectParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateIAM2ProjectParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	project, err := r.client.CreateIAM2Project(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create IAM2 project in ZStack", "Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(project.UUID)
	plan.Name = types.StringValue(project.Name)
	plan.Description = types.StringValue(project.Description)
	plan.State = types.StringValue(project.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *iam2ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state iam2ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetIAM2Project(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IAM2 project", "Could not read IAM2 project UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(project.UUID)
	state.Name = types.StringValue(project.Name)
	state.Description = types.StringValue(project.Description)
	state.State = types.StringValue(project.State)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *iam2ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan iam2ProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state iam2ProjectResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateIAM2ProjectParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateIAM2ProjectParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	project, err := r.client.UpdateIAM2Project(state.Uuid.ValueString(), updateParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not update IAM2 project", "Error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(project.UUID)
	plan.Name = types.StringValue(project.Name)
	plan.Description = types.StringValue(project.Description)
	plan.State = types.StringValue(project.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *iam2ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state iam2ProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "IAM2 project uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DeleteIAM2Project(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete IAM2 project", err.Error())
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *iam2ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
