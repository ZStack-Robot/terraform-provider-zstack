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
	_ resource.Resource                = &iam2VirtualIDResource{}
	_ resource.ResourceWithConfigure   = &iam2VirtualIDResource{}
	_ resource.ResourceWithImportState = &iam2VirtualIDResource{}
)

type iam2VirtualIDResource struct {
	client *client.ZSClient
}

type iam2VirtualIDModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Password    types.String `tfsdk:"password"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
}

func IAM2VirtualIDResource() resource.Resource {
	return &iam2VirtualIDResource{}
}

func (r *iam2VirtualIDResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T.", request.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *iam2VirtualIDResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_iam2_virtual_id"
}

func (r *iam2VirtualIDResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages an IAM2 Virtual ID in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the IAM2 virtual ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the IAM2 virtual ID",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for the IAM2 virtual ID",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the IAM2 virtual ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the IAM2 virtual ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the IAM2 virtual ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *iam2VirtualIDResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan iam2VirtualIDModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	createParam := param.CreateIAM2VirtualIDParam{}
	createParam.Params.Name = plan.Name.ValueString()
	createParam.Params.Password = plan.Password.ValueString()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		createParam.Params.Description = stringPtrOrNil(plan.Description.ValueString())
	}

	result, err := r.client.CreateIAM2VirtualID(createParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating IAM2 Virtual ID",
			"Could not create IAM2 virtual ID, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "created an IAM2 virtual ID")

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *iam2VirtualIDResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state iam2VirtualIDModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	virtualID, err := findResourceByQuery(r.client.QueryIAM2VirtualID, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			tflog.Warn(ctx, "IAM2 virtual ID not found, removing from state", map[string]interface{}{
				"uuid": state.Uuid.ValueString(),
			})
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading IAM2 Virtual ID",
			"Could not read IAM2 virtual ID UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(virtualID.Name)
	state.Description = stringValueOrNull(virtualID.Description)
	state.Type = types.StringValue(virtualID.Type)
	state.State = types.StringValue(virtualID.State)
	// Password is not returned by query, preserve from state

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *iam2VirtualIDResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan iam2VirtualIDModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state iam2VirtualIDModel
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	updateParam := param.UpdateIAM2VirtualIDParam{}
	updateParam.Params.Name = plan.Name.ValueString()

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		updateParam.Params.Description = stringPtrOrNil(plan.Description.ValueString())
	}

	if !plan.Password.Equal(state.Password) {
		updateParam.Params.Password = stringPtr(plan.Password.ValueString())
	}

	result, err := r.client.UpdateIAM2VirtualID(state.Uuid.ValueString(), updateParam)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating IAM2 Virtual ID",
			"Could not update IAM2 virtual ID, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = stringValueOrNull(result.Description)
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)
	// Password is write-only, preserve from plan

	tflog.Trace(ctx, "updated an IAM2 virtual ID")

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *iam2VirtualIDResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state iam2VirtualIDModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIAM2VirtualID(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting IAM2 Virtual ID",
			"Could not delete IAM2 virtual ID, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted an IAM2 virtual ID")
}

func (r *iam2VirtualIDResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
