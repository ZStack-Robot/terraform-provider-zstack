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
	_ resource.Resource                = &emailMediaResource{}
	_ resource.ResourceWithConfigure   = &emailMediaResource{}
	_ resource.ResourceWithImportState = &emailMediaResource{}
)

type emailMediaResource struct {
	client *client.ZSClient
}

type emailMediaModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	SmtpServer  types.String `tfsdk:"smtp_server"`
	SmtpPort    types.Int64  `tfsdk:"smtp_port"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Type        types.String `tfsdk:"type"`
	State       types.String `tfsdk:"state"`
}

func EmailMediaResource() resource.Resource {
	return &emailMediaResource{}
}

func (r *emailMediaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = cli
}

func (r *emailMediaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_media"
}

func (r *emailMediaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage email media in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the email media.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the email media.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the email media.",
			},
			"smtp_server": schema.StringAttribute{
				Required:    true,
				Description: "The SMTP server of the email media.",
			},
			"smtp_port": schema.Int64Attribute{
				Required:    true,
				Description: "The SMTP port of the email media.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The username for SMTP authentication.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The password for SMTP authentication.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The media type.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the media.",
			},
		},
	}
}

func (r *emailMediaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan emailMediaModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateEmailMediaParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateEmailMediaParamDetail{
			SmtpServer: plan.SmtpServer.ValueString(),
			SmtpPort:   int(plan.SmtpPort.ValueInt64()),
			Name:       plan.Name.ValueString(),
			Description: stringPtrOrNil(
				plan.Description.ValueString(),
			),
		},
	}

	if !plan.Username.IsNull() && !plan.Username.IsUnknown() {
		p.Params.Username = stringPtrOrNil(plan.Username.ValueString())
	}

	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		p.Params.Password = stringPtrOrNil(plan.Password.ValueString())
	}

	item, err := r.client.CreateEmailMedia(p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to create email media", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Type = stringValueOrNull(item.Type)
	plan.State = stringValueOrNull(item.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *emailMediaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state emailMediaModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QueryEmailMedia(&queryParam)
	if err != nil {
		tflog.Warn(ctx, "Unable to query email media. It may have been deleted.: "+err.Error())
		state = emailMediaModel{Uuid: types.StringValue("")}
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	found := false
	for _, item := range items {
		if item.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(item.UUID)
			state.Name = types.StringValue(item.Name)
			state.Description = stringValueOrNull(item.Description)
			state.Type = stringValueOrNull(item.Type)
			state.State = stringValueOrNull(item.State)
			found = true
			break
		}
	}

	if !found {
		tflog.Warn(ctx, "Email media not found. It might have been deleted outside of Terraform.")
		state = emailMediaModel{Uuid: types.StringValue("")}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *emailMediaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan emailMediaModel
	var state emailMediaModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateEmailMediaParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateEmailMediaParamDetail{
			Name: plan.Name.ValueString(),
		},
	}

	if !plan.Description.IsUnknown() {
		p.Params.Description = stringPtrOrNil(plan.Description.ValueString())
	}

	if !plan.SmtpServer.IsNull() && !plan.SmtpServer.IsUnknown() {
		p.Params.SmtpServer = stringPtrOrNil(plan.SmtpServer.ValueString())
	}

	if !plan.SmtpPort.IsNull() && !plan.SmtpPort.IsUnknown() {
		p.Params.SmtpPort = intPtr(int(plan.SmtpPort.ValueInt64()))
	}

	if !plan.Username.IsNull() && !plan.Username.IsUnknown() {
		p.Params.Username = stringPtrOrNil(plan.Username.ValueString())
	}

	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		p.Params.Password = stringPtrOrNil(plan.Password.ValueString())
	}

	item, err := r.client.UpdateEmailMedia(state.Uuid.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError("Fail to update email media", "Error "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.SmtpServer = stringValueOrNull(item.SmtpServer)
	plan.SmtpPort = types.Int64Value(int64(item.SmtpPort))
	plan.Username = stringValueOrNull(item.Username)
	plan.Type = stringValueOrNull(item.Type)
	plan.State = stringValueOrNull(item.State)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *emailMediaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state emailMediaModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Email media UUID is empty, skipping delete.")
		return
	}

	if err := r.client.DeleteMedia(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Fail to delete email media", err.Error())
		return
	}
}

func (r *emailMediaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
