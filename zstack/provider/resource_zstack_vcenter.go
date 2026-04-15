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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vcenterResource{}
	_ resource.ResourceWithConfigure   = &vcenterResource{}
	_ resource.ResourceWithImportState = &vcenterResource{}
)

type vcenterResource struct {
	client *client.ZSClient
}

type vcenterModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	ZoneUuid    types.String `tfsdk:"zone_uuid"`
	DomainName  types.String `tfsdk:"domain_name"`
	Https       types.Bool   `tfsdk:"https"`
	Port        types.Int64  `tfsdk:"port"`
	Version     types.String `tfsdk:"version"`
	State       types.String `tfsdk:"state"`
	Status      types.String `tfsdk:"status"`
}

func VCenterResource() resource.Resource {
	return &vcenterResource{}
}

func (r *vcenterResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *vcenterResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vcenter"
}

func (r *vcenterResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage vCenters in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the vCenter.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The vCenter name.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The vCenter username.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The vCenter password.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The zone UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The vCenter domain name.",
			},
			"https": schema.BoolAttribute{
				Optional:    true,
				Description: "Use HTTPS.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.Int64Attribute{
				Optional:    true,
				Description: "The service port.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The vCenter version.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status.",
			},
		},
	}
}

func (r *vcenterResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vcenterModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.AddVCenterParam{
		BaseParam: param.BaseParam{},
		Params: param.AddVCenterParamDetail{
			Username:    plan.Username.ValueString(),
			Password:    plan.Password.ValueString(),
			ZoneUuid:    plan.ZoneUuid.ValueString(),
			Name:        plan.Name.ValueString(),
			DomainName:  plan.DomainName.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if !plan.Https.IsNull() && !plan.Https.IsUnknown() {
		p.Params.Https = boolPtr(plan.Https.ValueBool())
	}
	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = intPtr(int(plan.Port.ValueInt64()))
	}

	vcenter, err := r.client.AddVCenter(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating vCenter",
			"Could not create vcenter, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vcenter.UUID)
	plan.Name = types.StringValue(vcenter.Name)
	plan.Description = stringValueOrNull(vcenter.Description)
	plan.Username = stringValueOrNull(vcenter.UserName)
	plan.ZoneUuid = types.StringValue(vcenter.ZoneUuid)
	plan.DomainName = stringValueOrNull(vcenter.DomainName)
	plan.Port = types.Int64Value(int64(vcenter.Port))
	plan.Version = stringValueOrNull(vcenter.Version)
	plan.Https = types.BoolValue(vcenter.Https)
	plan.State = stringValueOrNull(vcenter.State)
	plan.Status = stringValueOrNull(vcenter.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *vcenterResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vcenterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	vcenter, err := findResourceByQuery(r.client.QueryVCenter, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading vCenter",
			"Could not read vCenter, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(vcenter.UUID)
	state.Name = types.StringValue(vcenter.Name)
	state.Description = stringValueOrNull(vcenter.Description)
	state.Username = stringValueOrNull(vcenter.UserName)
	state.ZoneUuid = types.StringValue(vcenter.ZoneUuid)
	state.DomainName = stringValueOrNull(vcenter.DomainName)
	state.Port = types.Int64Value(int64(vcenter.Port))
	state.Version = stringValueOrNull(vcenter.Version)
	state.Https = types.BoolValue(vcenter.Https)
	state.State = stringValueOrNull(vcenter.State)
	state.Status = stringValueOrNull(vcenter.Status)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *vcenterResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan vcenterModel
	var state vcenterModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateVCenterParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVCenterParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Username:    stringPtrOrNil(plan.Username.ValueString()),
			Password:    stringPtrOrNil(plan.Password.ValueString()),
			DomainName:  stringPtrOrNil(plan.DomainName.ValueString()),
			State:       stringPtrOrNil(plan.State.ValueString()),
		},
	}

	if !plan.Port.IsNull() && !plan.Port.IsUnknown() {
		p.Params.Port = intPtr(int(plan.Port.ValueInt64()))
	}

	vcenter, err := r.client.UpdateVCenter(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating vCenter",
			"Could not update vcenter, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vcenter.UUID)
	plan.Name = types.StringValue(vcenter.Name)
	plan.Description = stringValueOrNull(vcenter.Description)
	plan.Username = stringValueOrNull(vcenter.UserName)
	plan.ZoneUuid = types.StringValue(vcenter.ZoneUuid)
	plan.DomainName = stringValueOrNull(vcenter.DomainName)
	plan.Port = types.Int64Value(int64(vcenter.Port))
	plan.Version = stringValueOrNull(vcenter.Version)
	plan.Https = types.BoolValue(vcenter.Https)
	plan.State = stringValueOrNull(vcenter.State)
	plan.Status = stringValueOrNull(vcenter.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *vcenterResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vcenterModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVCenter(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting vCenter",
			"Could not delete vcenter, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *vcenterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
