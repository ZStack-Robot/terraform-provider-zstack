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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &flkSecSecurityMachineResource{}
	_ resource.ResourceWithConfigure   = &flkSecSecurityMachineResource{}
	_ resource.ResourceWithImportState = &flkSecSecurityMachineResource{}
)

type flkSecSecurityMachineResource struct {
	client *client.ZSClient
}

type flkSecSecurityMachineModel struct {
	Uuid                   types.String `tfsdk:"uuid"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	ManagementIp           types.String `tfsdk:"management_ip"`
	Model                  types.String `tfsdk:"model"`
	Type                   types.String `tfsdk:"type"`
	ZoneUuid               types.String `tfsdk:"zone_uuid"`
	SecretResourcePoolUuid types.String `tfsdk:"secret_resource_pool_uuid"`
	Port                   types.Int64  `tfsdk:"port"`
	State                  types.String `tfsdk:"state"`
	Status                 types.String `tfsdk:"status"`
}

func FlkSecSecurityMachineResource() resource.Resource {
	return &flkSecSecurityMachineResource{}
}

func (r *flkSecSecurityMachineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *flkSecSecurityMachineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flk_sec_security_machine"
}

func (r *flkSecSecurityMachineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage FLK security machines in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the security machine.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the security machine.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the security machine.",
			},
			"management_ip": schema.StringAttribute{
				Required:    true,
				Description: "The management IP of the security machine.",
			},
			"model": schema.StringAttribute{
				Required:    true,
				Description: "The model of the security machine.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the security machine.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The zone UUID of the security machine.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_resource_pool_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The secret resource pool UUID of the security machine.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Required:    true,
				Description: "The management port of the security machine.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the security machine.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the security machine.",
			},
		},
	}
}

func (r *flkSecSecurityMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan flkSecSecurityMachineModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddFlkSecSecurityMachineParam{
		BaseParam: param.BaseParam{},
		Params: param.AddFlkSecSecurityMachineParamDetail{
			Port:                   int(plan.Port.ValueInt64()),
			Name:                   plan.Name.ValueString(),
			Description:            stringPtrOrNil(plan.Description.ValueString()),
			ManagementIp:           plan.ManagementIp.ValueString(),
			Model:                  plan.Model.ValueString(),
			Type:                   plan.Type.ValueString(),
			ZoneUuid:               plan.ZoneUuid.ValueString(),
			SecretResourcePoolUuid: plan.SecretResourcePoolUuid.ValueString(),
		},
	}

	item, err := r.client.AddFlkSecSecurityMachine(p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating FLK Security Machine",
			"Could not create FLK security machine, unexpected error: "+err.Error(),
		)
		return
	}

	state := flkSecSecurityMachineModelFromView(item, plan)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *flkSecSecurityMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state flkSecSecurityMachineModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := findResourceByQuery(r.client.QuerySecurityMachine, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Security Machine",
			"Could not read Security Machine, unexpected error: "+err.Error(),
		)
		return
	}

	state = flkSecSecurityMachineModelFromView(item, state)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *flkSecSecurityMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan flkSecSecurityMachineModel
	var state flkSecSecurityMachineModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateFlkSecSecurityMachineParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateFlkSecSecurityMachineParamDetail{
			Port:         intPtr(int(plan.Port.ValueInt64())),
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			ManagementIp: stringPtrOrNil(plan.ManagementIp.ValueString()),
			Model:        stringPtrOrNil(plan.Model.ValueString()),
		},
	}

	item, err := r.client.UpdateFlkSecSecurityMachine(state.Uuid.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating FLK Security Machine",
			"Could not update FLK security machine, unexpected error: "+err.Error(),
		)
		return
	}

	state = flkSecSecurityMachineModelFromView(item, plan)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *flkSecSecurityMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state flkSecSecurityMachineModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


	if err := r.client.DeleteSecurityMachine(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting FLK Security Machine",
			"Could not delete FLK security machine, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *flkSecSecurityMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func flkSecSecurityMachineModelFromView(item *view.SecurityMachineInventoryView, prior flkSecSecurityMachineModel) flkSecSecurityMachineModel {
	state := prior

	state.Uuid = types.StringValue(item.UUID)
	if item.Name != "" {
		state.Name = types.StringValue(item.Name)
	}
	state.Description = stringValueOrNull(item.Description)
	if item.ManagementIp != "" {
		state.ManagementIp = types.StringValue(item.ManagementIp)
	}
	if item.Model != "" {
		state.Model = types.StringValue(item.Model)
	}
	if item.Type != "" {
		state.Type = types.StringValue(item.Type)
	}
	if item.ZoneUuid != "" {
		state.ZoneUuid = types.StringValue(item.ZoneUuid)
	}
	if item.SecretResourcePoolUuid != "" {
		state.SecretResourcePoolUuid = types.StringValue(item.SecretResourcePoolUuid)
	}
	state.State = stringValueOrNull(item.State)
	state.Status = stringValueOrNull(item.Status)

	return state
}
