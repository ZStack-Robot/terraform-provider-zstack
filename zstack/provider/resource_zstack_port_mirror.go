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
	_ resource.Resource                = &portMirrorResource{}
	_ resource.ResourceWithConfigure   = &portMirrorResource{}
	_ resource.ResourceWithImportState = &portMirrorResource{}
)

type portMirrorResource struct {
	client *client.ZSClient
}

type portMirrorModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	MirrorNetworkUuid types.String `tfsdk:"mirror_network_uuid"`
	State             types.String `tfsdk:"state"`
}

func PortMirrorResource() resource.Resource {
	return &portMirrorResource{}
}

func (r *portMirrorResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *portMirrorResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_port_mirror"
}

func (r *portMirrorResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a port mirror in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the port mirror.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The name of the port mirror.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the port mirror.",
			},
			"mirror_network_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the mirror network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the port mirror.",
			},
		},
	}
}

func (r *portMirrorResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan portMirrorModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreatePortMirrorParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePortMirrorParamDetail{
			MirrorNetworkUuid: plan.MirrorNetworkUuid.ValueString(),
			Name:              plan.Name.ValueString(),
			Description:       stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	portMirror, err := r.client.CreatePortMirror(p)
	if err != nil {
		response.Diagnostics.AddError("Fail to create port mirror", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(portMirror.UUID)
	plan.Name = stringValueOrNull(portMirror.Name)
	plan.Description = stringValueOrNull(portMirror.Description)
	plan.MirrorNetworkUuid = types.StringValue(portMirror.MirrorNetworkUuid)
	plan.State = stringValueOrNull(portMirror.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *portMirrorResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state portMirrorModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	portMirrors, err := r.client.QueryPortMirror(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query port mirrors. It may have been deleted.: "+err.Error())
		state = portMirrorModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, portMirror := range portMirrors {
		if portMirror.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(portMirror.UUID)
			state.Name = stringValueOrNull(portMirror.Name)
			state.Description = stringValueOrNull(portMirror.Description)
			state.MirrorNetworkUuid = types.StringValue(portMirror.MirrorNetworkUuid)
			state.State = stringValueOrNull(portMirror.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Port mirror not found. It might have been deleted outside of Terraform.")
		state = portMirrorModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *portMirrorResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan portMirrorModel
	var state portMirrorModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdatePortMirrorParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdatePortMirrorParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	portMirror, err := r.client.UpdatePortMirror(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError("Fail to update port mirror", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(portMirror.UUID)
	plan.Name = stringValueOrNull(portMirror.Name)
	plan.Description = stringValueOrNull(portMirror.Description)
	plan.MirrorNetworkUuid = types.StringValue(portMirror.MirrorNetworkUuid)
	plan.State = stringValueOrNull(portMirror.State)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *portMirrorResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state portMirrorModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Port mirror UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeletePortMirror(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete port mirror", ""+err.Error())
		return
	}
}

func (r *portMirrorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
