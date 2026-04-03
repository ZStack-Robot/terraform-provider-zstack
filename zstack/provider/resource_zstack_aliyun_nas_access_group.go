// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource              = &aliyunNasAccessGroupResource{}
	_ resource.ResourceWithConfigure = &aliyunNasAccessGroupResource{}
)

func AliyunNasAccessGroupResource() resource.Resource {
	return &aliyunNasAccessGroupResource{}
}

type aliyunNasAccessGroupResource struct {
	client *client.ZSClient
}

type aliyunNasAccessGroupModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	DataCenterUuid types.String `tfsdk:"data_center_uuid"`
	NetworkType    types.String `tfsdk:"network_type"`
	Type           types.String `tfsdk:"type"`
}

func (r *aliyunNasAccessGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aliyun_nas_access_group"
}

func (r *aliyunNasAccessGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Aliyun NAS access group",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the Aliyun NAS access group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the access group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the access group",
			},
			"data_center_uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID of the data center",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Network type of the access group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Type of the access group",
			},
		},
	}
}

func (r *aliyunNasAccessGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *aliyunNasAccessGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aliyunNasAccessGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.CreateAliyunNasAccessGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateAliyunNasAccessGroupParamDetail{
			Name:           plan.Name.ValueString(),
			DataCenterUuid: plan.DataCenterUuid.ValueString(),
			Description:    stringPtrOrNil(plan.Description.ValueString()),
			NetworkType:    stringPtrOrNil(plan.NetworkType.ValueString()),
		},
	}

	item, err := r.client.CreateAliyunNasAccessGroup(p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Aliyun NAS access group",
			err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Type = types.StringValue(item.Type)
	if item.Description != "" {
		plan.Description = types.StringValue(item.Description)
	}

	tflog.Trace(ctx, "created Aliyun NAS access group", map[string]any{"uuid": plan.Uuid.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *aliyunNasAccessGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aliyunNasAccessGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := r.client.GetAliyunNasAccessGroup(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Aliyun NAS access group",
			err.Error(),
		)
		return
	}

	if item == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = types.StringValue(item.Description)
	state.DataCenterUuid = types.StringValue(item.DataCenterUuid)
	state.Type = types.StringValue(item.Type)

	tflog.Trace(ctx, "read Aliyun NAS access group", map[string]any{"uuid": state.Uuid.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *aliyunNasAccessGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aliyunNasAccessGroupModel
	var state aliyunNasAccessGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateAliyunNasAccessGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateAliyunNasAccessGroupParamDetail{
			Uuid:        state.Uuid.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	item, err := r.client.UpdateAliyunNasAccessGroup(p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Aliyun NAS access group",
			err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Type = types.StringValue(item.Type)
	if item.Description != "" {
		plan.Description = types.StringValue(item.Description)
	}

	tflog.Trace(ctx, "updated Aliyun NAS access group", map[string]any{"uuid": plan.Uuid.ValueString()})

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *aliyunNasAccessGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aliyunNasAccessGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAliyunNasAccessGroup(state.Uuid.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Aliyun NAS access group",
			err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "deleted Aliyun NAS access group", map[string]any{"uuid": state.Uuid.ValueString()})
}
