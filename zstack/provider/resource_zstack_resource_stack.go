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
	_ resource.Resource                = &resourceStackResource{}
	_ resource.ResourceWithConfigure   = &resourceStackResource{}
	_ resource.ResourceWithImportState = &resourceStackResource{}
)

type resourceStackResource struct {
	client *client.ZSClient
}

type resourceStackModel struct {
	Uuid            types.String `tfsdk:"uuid"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Type            types.String `tfsdk:"type"`
	Rollback        types.Bool   `tfsdk:"rollback"`
	TemplateContent types.String `tfsdk:"template_content"`
	TemplateUuid    types.String `tfsdk:"template_uuid"`
	Parameters      types.String `tfsdk:"parameters"`
	Version         types.String `tfsdk:"version"`
	Status          types.String `tfsdk:"status"`
	Reason          types.String `tfsdk:"reason"`
	Outputs         types.String `tfsdk:"outputs"`
	ParamContent    types.String `tfsdk:"param_content"`
	EnableRollback  types.Bool   `tfsdk:"enable_rollback"`
}

func ResourceStackResource() resource.Resource {
	return &resourceStackResource{}
}

func (r *resourceStackResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *resourceStackResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_stack"
}

func (r *resourceStackResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a resource stack in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the resource stack.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the resource stack.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the resource stack.",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of the resource stack.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rollback": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether rollback is enabled.",
			},
			"template_content": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The template content.",
			},
			"template_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The template UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parameters": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The input parameters.",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the resource stack.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the resource stack.",
			},
			"reason": schema.StringAttribute{
				Computed:    true,
				Description: "The reason for the current status.",
			},
			"outputs": schema.StringAttribute{
				Computed:    true,
				Description: "The outputs of the resource stack.",
			},
			"param_content": schema.StringAttribute{
				Computed:    true,
				Description: "The rendered parameter content.",
			},
			"enable_rollback": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether rollback is enabled on the stack.",
			},
		},
	}
}

func (r *resourceStackResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan resourceStackModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateResourceStackParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateResourceStackParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			Type:            stringPtrOrNil(plan.Type.ValueString()),
			Rollback:        boolPtr(plan.Rollback.ValueBool()),
			TemplateContent: stringPtrOrNil(plan.TemplateContent.ValueString()),
			TemplateUuid:    stringPtrOrNil(plan.TemplateUuid.ValueString()),
			Parameters:      stringPtrOrNil(plan.Parameters.ValueString()),
		},
	}

	resourceStack, err := r.client.CreateResourceStack(p)
	if err != nil {
		response.Diagnostics.AddError("Fail to create resource stack", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(resourceStack.UUID)
	plan.Name = types.StringValue(resourceStack.Name)
	plan.Description = stringValueOrNull(resourceStack.Description)
	plan.Version = stringValueOrNull(resourceStack.Version)
	plan.Type = stringValueOrNull(resourceStack.Type)
	plan.TemplateContent = stringValueOrNull(resourceStack.TemplateContent)
	plan.ParamContent = stringValueOrNull(resourceStack.ParamContent)
	plan.Status = stringValueOrNull(resourceStack.Status)
	plan.Reason = stringValueOrNull(resourceStack.Reason)
	plan.Outputs = stringValueOrNull(resourceStack.Outputs)
	plan.EnableRollback = types.BoolValue(resourceStack.EnableRollback)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *resourceStackResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state resourceStackModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	resourceStacks, err := r.client.QueryResourceStack(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query resource stacks. It may have been deleted.: "+err.Error())
		state = resourceStackModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, resourceStack := range resourceStacks {
		if resourceStack.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(resourceStack.UUID)
			state.Name = types.StringValue(resourceStack.Name)
			state.Description = stringValueOrNull(resourceStack.Description)
			state.Version = stringValueOrNull(resourceStack.Version)
			state.Type = stringValueOrNull(resourceStack.Type)
			state.TemplateContent = stringValueOrNull(resourceStack.TemplateContent)
			state.ParamContent = stringValueOrNull(resourceStack.ParamContent)
			state.Status = stringValueOrNull(resourceStack.Status)
			state.Reason = stringValueOrNull(resourceStack.Reason)
			state.Outputs = stringValueOrNull(resourceStack.Outputs)
			state.EnableRollback = types.BoolValue(resourceStack.EnableRollback)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Resource stack not found. It might have been deleted outside of Terraform.")
		state = resourceStackModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *resourceStackResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan resourceStackModel
	var state resourceStackModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateResourceStackParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateResourceStackParamDetail{
			Name:            plan.Name.ValueString(),
			Description:     stringPtrOrNil(plan.Description.ValueString()),
			Rollback:        boolPtr(plan.Rollback.ValueBool()),
			TemplateContent: stringPtrOrNil(plan.TemplateContent.ValueString()),
			Parameters:      stringPtrOrNil(plan.Parameters.ValueString()),
		},
	}

	resourceStack, err := r.client.UpdateResourceStack(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Fail to update resource stack", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(resourceStack.UUID)
	plan.Name = types.StringValue(resourceStack.Name)
	plan.Description = stringValueOrNull(resourceStack.Description)
	plan.Version = stringValueOrNull(resourceStack.Version)
	plan.Type = stringValueOrNull(resourceStack.Type)
	plan.TemplateContent = stringValueOrNull(resourceStack.TemplateContent)
	plan.ParamContent = stringValueOrNull(resourceStack.ParamContent)
	plan.Status = stringValueOrNull(resourceStack.Status)
	plan.Reason = stringValueOrNull(resourceStack.Reason)
	plan.Outputs = stringValueOrNull(resourceStack.Outputs)
	plan.EnableRollback = types.BoolValue(resourceStack.EnableRollback)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *resourceStackResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state resourceStackModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Resource stack UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteResourceStack(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete resource stack", ""+err.Error())
		return
	}
}

func (r *resourceStackResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
