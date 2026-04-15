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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &containerManagementEndpointResource{}
	_ resource.ResourceWithConfigure   = &containerManagementEndpointResource{}
	_ resource.ResourceWithImportState = &containerManagementEndpointResource{}
)

type containerManagementEndpointResource struct {
	client *client.ZSClient
}

type containerManagementEndpointModel struct {
	Uuid                     types.String `tfsdk:"uuid"`
	Name                     types.String `tfsdk:"name"`
	Description              types.String `tfsdk:"description"`
	ManagementIp             types.String `tfsdk:"management_ip"`
	ManagementPort           types.Int64  `tfsdk:"management_port"`
	Vendor                   types.String `tfsdk:"vendor"`
	ContainerAccessKeyId     types.String `tfsdk:"container_access_key_id"`
	ContainerAccessKeySecret types.String `tfsdk:"container_access_key_secret"`
	AccessKeyId              types.String `tfsdk:"access_key_id"`
}

func ContainerManagementEndpointResource() resource.Resource {
	return &containerManagementEndpointResource{}
}

func (r *containerManagementEndpointResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *containerManagementEndpointResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_container_management_endpoint"
}

func (r *containerManagementEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage container management endpoints in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the container management endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the container management endpoint.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the container management endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"management_ip": schema.StringAttribute{
				Required:    true,
				Description: "The management IP address of the container management endpoint.",
			},
			"management_port": schema.Int64Attribute{
				Required:    true,
				Description: "The management port of the container management endpoint.",
			},
			"vendor": schema.StringAttribute{
				Required:    true,
				Description: "The vendor of the container management endpoint.",
			},
			"container_access_key_id": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The container access key ID (write-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"container_access_key_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The container access key secret (write-only).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_key_id": schema.StringAttribute{
				Computed:    true,
				Description: "The access key ID (read-only, from API response).",
			},
		},
	}
}

func (r *containerManagementEndpointResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan containerManagementEndpointModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.AddContainerManagementEndpointParam{
		BaseParam: param.BaseParam{},
		Params: param.AddContainerManagementEndpointParamDetail{
			Name:                     plan.Name.ValueString(),
			Description:              stringPtrOrNil(plan.Description.ValueString()),
			ManagementIp:             plan.ManagementIp.ValueString(),
			ManagementPort:           int(plan.ManagementPort.ValueInt64()),
			Vendor:                   plan.Vendor.ValueString(),
			ContainerAccessKeyId:     plan.ContainerAccessKeyId.ValueString(),
			ContainerAccessKeySecret: plan.ContainerAccessKeySecret.ValueString(),
		},
	}

	endpoint, err := r.client.AddContainerManagementEndpoint(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Container Management Endpoint",
			"Could not create container management endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(endpoint.UUID)
	plan.Name = types.StringValue(endpoint.Name)
	plan.Description = stringValueOrNull(endpoint.Description)
	plan.ManagementIp = types.StringValue(endpoint.ManagementIp)
	plan.ManagementPort = types.Int64Value(int64(endpoint.ManagementPort))
	plan.Vendor = types.StringValue(endpoint.Vendor)
	plan.AccessKeyId = types.StringValue(endpoint.AccessKeyId)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *containerManagementEndpointResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state containerManagementEndpointModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	endpoint, err := findResourceByQuery(r.client.QueryContainerManagementEndpoint, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Container Management Endpoint",
			"Could not read container management endpoint UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(endpoint.UUID)
	state.Name = types.StringValue(endpoint.Name)
	state.Description = stringValueOrNull(endpoint.Description)
	state.ManagementIp = types.StringValue(endpoint.ManagementIp)
	state.ManagementPort = types.Int64Value(int64(endpoint.ManagementPort))
	state.Vendor = types.StringValue(endpoint.Vendor)
	state.AccessKeyId = types.StringValue(endpoint.AccessKeyId)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *containerManagementEndpointResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan containerManagementEndpointModel
	var state containerManagementEndpointModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateContainerManagementEndpointParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateContainerManagementEndpointParamDetail{
			ManagementIp:   stringPtrOrNil(plan.ManagementIp.ValueString()),
			ManagementPort: intPtr(int(plan.ManagementPort.ValueInt64())),
			Vendor:         stringPtrOrNil(plan.Vendor.ValueString()),
		},
	}

	endpoint, err := r.client.UpdateContainerManagementEndpoint(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Container Management Endpoint",
			"Could not update container management endpoint, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(endpoint.UUID)
	plan.Name = types.StringValue(endpoint.Name)
	plan.Description = stringValueOrNull(endpoint.Description)
	plan.ManagementIp = types.StringValue(endpoint.ManagementIp)
	plan.ManagementPort = types.Int64Value(int64(endpoint.ManagementPort))
	plan.Vendor = types.StringValue(endpoint.Vendor)
	plan.AccessKeyId = types.StringValue(endpoint.AccessKeyId)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *containerManagementEndpointResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state containerManagementEndpointModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DeleteContainerManagementEndpoint(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Container Management Endpoint",
			"Could not delete container management endpoint, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *containerManagementEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
