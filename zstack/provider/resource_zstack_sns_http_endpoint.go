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
	_ resource.Resource                = &snsHttpEndpointResource{}
	_ resource.ResourceWithConfigure   = &snsHttpEndpointResource{}
	_ resource.ResourceWithImportState = &snsHttpEndpointResource{}
)

type snsHttpEndpointResource struct {
	client *client.ZSClient
}

type snsHttpEndpointModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Url          types.String `tfsdk:"url"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Description  types.String `tfsdk:"description"`
	PlatformUuid types.String `tfsdk:"platform_uuid"`
	Type         types.String `tfsdk:"type"`
	State        types.String `tfsdk:"state"`
}

func SNSHttpEndpointResource() resource.Resource {
	return &snsHttpEndpointResource{}
}

func (r *snsHttpEndpointResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *snsHttpEndpointResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sns_http_endpoint"
}

func (r *snsHttpEndpointResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage SNS HTTP endpoints in ZStack. " +
			"An SNS HTTP endpoint receives notifications by HTTP callback.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the SNS HTTP endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SNS HTTP endpoint.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The callback URL of the SNS HTTP endpoint.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The username for the SNS HTTP endpoint.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The password for the SNS HTTP endpoint. This is a write-only sensitive field.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the SNS HTTP endpoint.",
			},
			"platform_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the SNS platform.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The endpoint type.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The current state of the SNS HTTP endpoint.",
			},
		},
	}
}

func (r *snsHttpEndpointResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan snsHttpEndpointModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateSNSHttpEndpointParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSNSHttpEndpointParamDetail{
			Url:          plan.Url.ValueString(),
			Username:     stringPtrOrNil(plan.Username.ValueString()),
			Password:     stringPtrOrNil(plan.Password.ValueString()),
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			PlatformUuid: stringPtrOrNil(plan.PlatformUuid.ValueString()),
		},
	}

	result, err := r.client.CreateSNSHttpEndpoint(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create SNS HTTP endpoint",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Url = types.StringValue(result.Url)
	plan.Username = stringValueOrNull(result.Username)
	plan.Description = stringValueOrNull(result.Description)
	plan.PlatformUuid = stringValueOrNull(result.PlatformUuid)
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snsHttpEndpointResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state snsHttpEndpointModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QuerySNSHttpEndpoint(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query SNS HTTP endpoints. It may have been deleted.: "+err.Error())
		state = snsHttpEndpointModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, item := range items {
		if item.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(item.UUID)
			state.Name = types.StringValue(item.Name)
			state.Url = types.StringValue(item.Url)
			state.Username = stringValueOrNull(item.Username)
			state.Description = stringValueOrNull(item.Description)
			state.PlatformUuid = stringValueOrNull(item.PlatformUuid)
			state.Type = types.StringValue(item.Type)
			state.State = types.StringValue(item.State)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "SNS HTTP endpoint not found. It might have been deleted outside of Terraform.")
		state = snsHttpEndpointModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snsHttpEndpointResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan snsHttpEndpointModel
	var state snsHttpEndpointModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateSNSHttpEndpointParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSNSHttpEndpointParamDetail{
			Url:         stringPtr(plan.Url.ValueString()),
			Username:    stringPtrOrNil(plan.Username.ValueString()),
			Password:    stringPtrOrNil(plan.Password.ValueString()),
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdateSNSHttpEndpoint(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update SNS HTTP endpoint",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Url = plan.Url
	plan.Username = plan.Username
	plan.Password = plan.Password
	plan.Description = stringValueOrNull(result.Description)
	plan.PlatformUuid = stringValueOrNull(result.PlatformUuid)
	plan.Type = types.StringValue(result.Type)
	plan.State = types.StringValue(result.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snsHttpEndpointResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Warn(ctx, "SNS HTTP endpoint does not support deletion via API. Removing from Terraform state only.")
}

func (r *snsHttpEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
