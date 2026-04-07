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
	_ resource.Resource                = &snsEmailEndpointResource{}
	_ resource.ResourceWithConfigure   = &snsEmailEndpointResource{}
	_ resource.ResourceWithImportState = &snsEmailEndpointResource{}
)

type snsEmailEndpointResource struct {
	client *client.ZSClient
}

type snsEmailEndpointModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Email        types.String `tfsdk:"email"`
	Description  types.String `tfsdk:"description"`
	PlatformUuid types.String `tfsdk:"platform_uuid"`
	Type         types.String `tfsdk:"type"`
	State        types.String `tfsdk:"state"`
}

func SNSEmailEndpointResource() resource.Resource {
	return &snsEmailEndpointResource{}
}

func (r *snsEmailEndpointResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *snsEmailEndpointResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sns_email_endpoint"
}

func (r *snsEmailEndpointResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage SNS email endpoints in ZStack. " +
			"An SNS email endpoint receives notifications by email.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the SNS email endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SNS email endpoint.",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "The email address for the SNS email endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the SNS email endpoint.",
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
				Description: "The current state of the SNS email endpoint.",
			},
		},
	}
}

func (r *snsEmailEndpointResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan snsEmailEndpointModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateSNSEmailEndpointParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSNSEmailEndpointParamDetail{
			Name:         plan.Name.ValueString(),
			Email:        stringPtr(plan.Email.ValueString()),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			PlatformUuid: stringPtrOrNil(plan.PlatformUuid.ValueString()),
		},
	}

	result, err := r.client.CreateSNSEmailEndpoint(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create SNS email endpoint",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Email = types.StringValue(plan.Email.ValueString())
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

func (r *snsEmailEndpointResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state snsEmailEndpointModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QuerySNSEmailEndpoint(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query SNS email endpoints. It may have been deleted.: "+err.Error())
		state = snsEmailEndpointModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, item := range items {
		if item.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(item.UUID)
			state.Name = types.StringValue(item.Name)
			state.Email = types.StringValue(item.Email)
			state.Description = stringValueOrNull(item.Description)
			state.PlatformUuid = stringValueOrNull(item.PlatformUuid)
			state.Type = types.StringValue(item.Type)
			state.State = types.StringValue(item.State)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "SNS email endpoint not found. It might have been deleted outside of Terraform.")
		state = snsEmailEndpointModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snsEmailEndpointResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {

}

func (r *snsEmailEndpointResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	tflog.Warn(ctx, "SNS email endpoint does not support deletion via API. Removing from Terraform state only.")
}

func (r *snsEmailEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
