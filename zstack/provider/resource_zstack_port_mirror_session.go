// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &portMirrorSessionResource{}
	_ resource.ResourceWithConfigure   = &portMirrorSessionResource{}
	_ resource.ResourceWithImportState = &portMirrorSessionResource{}
)

type portMirrorSessionResource struct {
	client *client.ZSClient
}

type portMirrorSessionModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	PortMirrorUuid types.String `tfsdk:"port_mirror_uuid"`
	Type           types.String `tfsdk:"type"`
	SrcEndPoint    types.String `tfsdk:"src_end_point"`
	DstEndPoint    types.String `tfsdk:"dst_end_point"`
	Status         types.String `tfsdk:"status"`
	InternalId     types.Int64  `tfsdk:"internal_id"`
}

func PortMirrorSessionResource() resource.Resource {
	return &portMirrorSessionResource{}
}

func (r *portMirrorSessionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *portMirrorSessionResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_port_mirror_session"
}

func (r *portMirrorSessionResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a port mirror session in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the port mirror session.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the port mirror session.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the port mirror session.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port_mirror_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the port mirror.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the port mirror session.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"src_end_point": schema.StringAttribute{
				Required:    true,
				Description: "The source endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dst_end_point": schema.StringAttribute{
				Required:    true,
				Description: "The destination endpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the port mirror session.",
			},
			"internal_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The internal ID of the port mirror session.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *portMirrorSessionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan portMirrorSessionModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreatePortMirrorSessionParam{
		BaseParam: param.BaseParam{},
		Params: param.CreatePortMirrorSessionParamDetail{
			PortMirrorUuid: plan.PortMirrorUuid.ValueString(),
			Name:           plan.Name.ValueString(),
			Description:    stringPtrOrNil(plan.Description.ValueString()),
			Type:           plan.Type.ValueString(),
			SrcEndPoint:    plan.SrcEndPoint.ValueString(),
			DstEndPoint:    plan.DstEndPoint.ValueString(),
		},
	}

	portMirrorSession, err := r.client.CreatePortMirrorSession(p)
	if err != nil {
		response.Diagnostics.AddError("Fail to create port mirror session", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(portMirrorSession.UUID)
	plan.Name = types.StringValue(portMirrorSession.Name)
	plan.Description = stringValueOrNull(portMirrorSession.Description)
	plan.PortMirrorUuid = types.StringValue(portMirrorSession.PortMirrorUuid)
	plan.Type = types.StringValue(portMirrorSession.Type)
	plan.SrcEndPoint = types.StringValue(portMirrorSession.SrcEndPoint)
	plan.DstEndPoint = types.StringValue(portMirrorSession.DstEndPoint)
	plan.Status = stringValueOrNull(portMirrorSession.Status)
	plan.InternalId = types.Int64Value(portMirrorSession.InternalId)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *portMirrorSessionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state portMirrorSessionModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	portMirrorSessions, err := r.client.QueryPortMirrorSession(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query port mirror sessions. It may have been deleted.: "+err.Error())
		state = portMirrorSessionModel{Uuid: types.StringValue("")}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, portMirrorSession := range portMirrorSessions {
		if portMirrorSession.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(portMirrorSession.UUID)
			state.Name = types.StringValue(portMirrorSession.Name)
			state.Description = stringValueOrNull(portMirrorSession.Description)
			state.PortMirrorUuid = types.StringValue(portMirrorSession.PortMirrorUuid)
			state.Type = types.StringValue(portMirrorSession.Type)
			state.SrcEndPoint = types.StringValue(portMirrorSession.SrcEndPoint)
			state.DstEndPoint = types.StringValue(portMirrorSession.DstEndPoint)
			state.Status = stringValueOrNull(portMirrorSession.Status)
			state.InternalId = types.Int64Value(portMirrorSession.InternalId)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Port mirror session not found. It might have been deleted outside of Terraform.")
		state = portMirrorSessionModel{Uuid: types.StringValue("")}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *portMirrorSessionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
}

func (r *portMirrorSessionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state portMirrorSessionModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Port mirror session UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeletePortMirrorSession(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("fail to delete port mirror session", ""+err.Error())
		return
	}
}

func (r *portMirrorSessionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
