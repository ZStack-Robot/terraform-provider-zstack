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
	_ resource.Resource                = &preconfigurationTemplateResource{}
	_ resource.ResourceWithConfigure   = &preconfigurationTemplateResource{}
	_ resource.ResourceWithImportState = &preconfigurationTemplateResource{}
)

type preconfigurationTemplateResource struct {
	client *client.ZSClient
}

type preconfigurationTemplateModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Distribution types.String `tfsdk:"distribution"`
	Type         types.String `tfsdk:"type"`
	Content      types.String `tfsdk:"content"`
	Md5sum       types.String `tfsdk:"md5sum"`
	IsPredefined types.Bool   `tfsdk:"is_predefined"`
	State        types.String `tfsdk:"state"`
}

func PreconfigurationTemplateResource() resource.Resource {
	return &preconfigurationTemplateResource{}
}

func (r *preconfigurationTemplateResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *preconfigurationTemplateResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_preconfiguration_template"
}

func (r *preconfigurationTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage preconfiguration templates in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the preconfiguration template.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the preconfiguration template.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the preconfiguration template.",
			},
			"distribution": schema.StringAttribute{
				Required:    true,
				Description: "Distribution of the preconfiguration template.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of the preconfiguration template.",
			},
			"content": schema.StringAttribute{
				Required:    true,
				Description: "Template content.",
			},
			"md5sum": schema.StringAttribute{
				Computed:    true,
				Description: "MD5 checksum of the template content.",
			},
			"is_predefined": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this template is predefined by the system.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Current state of the preconfiguration template.",
			},
		},
	}
}

func (r *preconfigurationTemplateResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan preconfigurationTemplateModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddPreconfigurationTemplateParam{
		BaseParam: param.BaseParam{},
		Params: param.AddPreconfigurationTemplateParamDetail{
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			Distribution: plan.Distribution.ValueString(),
			Type:         plan.Type.ValueString(),
			Content:      plan.Content.ValueString(),
		},
	}

	preconfigurationTemplate, err := r.client.AddPreconfigurationTemplate(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create preconfiguration template",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(preconfigurationTemplate.UUID)
	plan.Name = types.StringValue(preconfigurationTemplate.Name)
	plan.Description = stringValueOrNull(preconfigurationTemplate.Description)
	plan.Distribution = types.StringValue(preconfigurationTemplate.Distribution)
	plan.Type = types.StringValue(preconfigurationTemplate.Type)
	plan.Content = types.StringValue(preconfigurationTemplate.Content)
	plan.Md5sum = stringValueOrNull(preconfigurationTemplate.Md5sum)
	plan.IsPredefined = types.BoolValue(preconfigurationTemplate.IsPredefined)
	plan.State = stringValueOrNull(preconfigurationTemplate.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *preconfigurationTemplateResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state preconfigurationTemplateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	preconfigurationTemplates, err := r.client.QueryPreconfigurationTemplate(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query preconfiguration templates. It may have been deleted.: "+err.Error())
		state = preconfigurationTemplateModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, preconfigurationTemplate := range preconfigurationTemplates {
		if preconfigurationTemplate.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(preconfigurationTemplate.UUID)
			state.Name = types.StringValue(preconfigurationTemplate.Name)
			state.Description = stringValueOrNull(preconfigurationTemplate.Description)
			state.Distribution = types.StringValue(preconfigurationTemplate.Distribution)
			state.Type = types.StringValue(preconfigurationTemplate.Type)
			state.Content = types.StringValue(preconfigurationTemplate.Content)
			state.Md5sum = stringValueOrNull(preconfigurationTemplate.Md5sum)
			state.IsPredefined = types.BoolValue(preconfigurationTemplate.IsPredefined)
			state.State = stringValueOrNull(preconfigurationTemplate.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Preconfiguration template not found. It might have been deleted outside of Terraform.")
		state = preconfigurationTemplateModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *preconfigurationTemplateResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan preconfigurationTemplateModel
	var state preconfigurationTemplateModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdatePreconfigurationTemplateParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdatePreconfigurationTemplateParamDetail{
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			Distribution: stringPtrOrNil(plan.Distribution.ValueString()),
			Type:         stringPtrOrNil(plan.Type.ValueString()),
			Content:      stringPtrOrNil(plan.Content.ValueString()),
		},
	}

	preconfigurationTemplate, err := r.client.UpdatePreconfigurationTemplate(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update preconfiguration template",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(preconfigurationTemplate.UUID)
	plan.Name = types.StringValue(preconfigurationTemplate.Name)
	plan.Description = stringValueOrNull(preconfigurationTemplate.Description)
	plan.Distribution = types.StringValue(preconfigurationTemplate.Distribution)
	plan.Type = types.StringValue(preconfigurationTemplate.Type)
	plan.Content = types.StringValue(preconfigurationTemplate.Content)
	plan.Md5sum = stringValueOrNull(preconfigurationTemplate.Md5sum)
	plan.IsPredefined = types.BoolValue(preconfigurationTemplate.IsPredefined)
	plan.State = stringValueOrNull(preconfigurationTemplate.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *preconfigurationTemplateResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state preconfigurationTemplateModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Preconfiguration template UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeletePreconfigurationTemplate(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Fail to delete preconfiguration template", "Error "+err.Error())
		return
	}
}

func (r *preconfigurationTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
