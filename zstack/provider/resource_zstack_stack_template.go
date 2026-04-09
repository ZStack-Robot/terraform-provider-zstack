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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &stackTemplateResource{}
	_ resource.ResourceWithConfigure   = &stackTemplateResource{}
	_ resource.ResourceWithImportState = &stackTemplateResource{}
)

type stackTemplateResource struct {
	client *client.ZSClient
}

type stackTemplateModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Type            types.String `tfsdk:"type"`
	TemplateContent types.String `tfsdk:"template_content"`
	Version         types.String `tfsdk:"version"`
	State           types.Bool   `tfsdk:"state"`
	Md5sum          types.String `tfsdk:"md5sum"`
}

func StackTemplateResource() resource.Resource {
	return &stackTemplateResource{}
}

func (r *stackTemplateResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *stackTemplateResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_stack_template"
}

func (r *stackTemplateResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a stack template in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the stack template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the stack template.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the stack template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the stack template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template_content": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The template content.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the stack template.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.BoolAttribute{
				Computed:    true,
				Description: "The state of the stack template.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"md5sum": schema.StringAttribute{
				Computed:    true,
				Description: "The MD5 checksum of the template content.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *stackTemplateResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan stackTemplateModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddStackTemplateParam{
		BaseParam: param.BaseParam{},
		Params: param.AddStackTemplateParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			Type:            stringPtrOrNil(plan.Type.ValueString()),
			TemplateContent: stringPtrOrNil(plan.TemplateContent.ValueString()),
		},
	}

	stackTemplate, err := r.client.AddStackTemplate(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Stack Template",
			"Could not create stack template, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(stackTemplate.UUID)
	plan.Name = types.StringValue(stackTemplate.Name)
	plan.Description = stringValueOrNull(stackTemplate.Description)
	plan.Type = stringValueOrNull(stackTemplate.Type)
	plan.TemplateContent = stringValueOrNull(stackTemplate.Content)
	plan.Version = stringValueOrNull(stackTemplate.Version)
	plan.State = types.BoolValue(stackTemplate.State)
	plan.Md5sum = stringValueOrNull(stackTemplate.Md5sum)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *stackTemplateResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state stackTemplateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	stackTemplate, err := findResourceByQuery(r.client.QueryStackTemplate, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		tflog.Warn(ctx, "Unable to query stack templates. It may have been deleted.: "+err.Error())
		state = stackTemplateModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(stackTemplate.UUID)
	state.Name = types.StringValue(stackTemplate.Name)
	state.Description = stringValueOrNull(stackTemplate.Description)
	state.Type = stringValueOrNull(stackTemplate.Type)
	state.TemplateContent = stringValueOrNull(stackTemplate.Content)
	state.Version = stringValueOrNull(stackTemplate.Version)
	state.State = types.BoolValue(stackTemplate.State)
	state.Md5sum = stringValueOrNull(stackTemplate.Md5sum)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *stackTemplateResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan stackTemplateModel
	var state stackTemplateModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateStackTemplateParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateStackTemplateParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			State:           boolPtr(plan.State.ValueBool()),
			TemplateContent: stringPtrOrNil(plan.TemplateContent.ValueString()),
		},
	}

	stackTemplate, err := r.client.UpdateStackTemplate(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Stack Template",
			"Could not update stack template, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(stackTemplate.UUID)
	plan.Name = types.StringValue(stackTemplate.Name)
	plan.Description = stringValueOrNull(stackTemplate.Description)
	plan.Type = stringValueOrNull(stackTemplate.Type)
	plan.TemplateContent = stringValueOrNull(stackTemplate.Content)
	plan.Version = stringValueOrNull(stackTemplate.Version)
	plan.State = types.BoolValue(stackTemplate.State)
	plan.Md5sum = stringValueOrNull(stackTemplate.Md5sum)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *stackTemplateResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state stackTemplateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Stack template UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteStackTemplate(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Stack Template",
			"Could not delete stack template, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *stackTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
