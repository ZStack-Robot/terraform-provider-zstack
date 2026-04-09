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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &baremetalChassisResource{}
	_ resource.ResourceWithConfigure   = &baremetalChassisResource{}
	_ resource.ResourceWithImportState = &baremetalChassisResource{}
)

type baremetalChassisResource struct {
	client *client.ZSClient
}

type baremetalChassisModel struct {
	Uuid         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	ClusterUuid  types.String `tfsdk:"cluster_uuid"`
	IpmiAddress  types.String `tfsdk:"ipmi_address"`
	IpmiPort     types.Int64  `tfsdk:"ipmi_port"`
	IpmiUsername types.String `tfsdk:"ipmi_username"`
	IpmiPassword types.String `tfsdk:"ipmi_password"`
	ZoneUuid     types.String `tfsdk:"zone_uuid"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
}

func BaremetalChassisResource() resource.Resource {
	return &baremetalChassisResource{}
}

func (r *baremetalChassisResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *baremetalChassisResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_baremetal_chassis"
}

func (r *baremetalChassisResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage baremetal chassis in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the baremetal chassis.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the baremetal chassis.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the baremetal chassis.",
			},
			"cluster_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The cluster UUID for the chassis.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ipmi_address": schema.StringAttribute{
				Required:    true,
				Description: "The IPMI address.",
			},
			"ipmi_port": schema.Int64Attribute{
				Optional:    true,
				Description: "The IPMI port.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ipmi_username": schema.StringAttribute{
				Required:    true,
				Description: "The IPMI username.",
			},
			"ipmi_password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The IPMI password.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The zone UUID.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the chassis.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the chassis.",
			},
		},
	}
}

func (r *baremetalChassisResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan baremetalChassisModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateBaremetalChassisParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateBaremetalChassisParamDetail{
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			ClusterUuid:  plan.ClusterUuid.ValueString(),
			IpmiAddress:  plan.IpmiAddress.ValueString(),
			IpmiUsername: plan.IpmiUsername.ValueString(),
			IpmiPassword: plan.IpmiPassword.ValueString(),
		},
	}

	if !plan.IpmiPort.IsNull() && !plan.IpmiPort.IsUnknown() {
		p.Params.IpmiPort = intPtr(int(plan.IpmiPort.ValueInt64()))
	}

	chassis, err := r.client.CreateBaremetalChassis(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Baremetal Chassis",
			"Could not create baremetal chassis, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(chassis.UUID)
	plan.Name = types.StringValue(chassis.Name)
	plan.Description = stringValueOrNull(chassis.Description)
	plan.ClusterUuid = types.StringValue(chassis.ClusterUuid)
	plan.IpmiAddress = types.StringValue(chassis.IpmiAddress)
	plan.IpmiPort = types.Int64Value(int64(chassis.IpmiPort))
	plan.IpmiUsername = types.StringValue(chassis.IpmiUsername)
	plan.ZoneUuid = stringValueOrNull(chassis.ZoneUuid)
	plan.State = stringValueOrNull(chassis.State)
	plan.Status = stringValueOrNull(chassis.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalChassisResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state baremetalChassisModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	chassis, err := findResourceByQuery(r.client.QueryBaremetalChassis, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Baremetal Chassis",
			"Could not read baremetal chassis UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(chassis.UUID)
	state.Name = types.StringValue(chassis.Name)
	state.Description = stringValueOrNull(chassis.Description)
	state.ClusterUuid = types.StringValue(chassis.ClusterUuid)
	state.IpmiAddress = types.StringValue(chassis.IpmiAddress)
	state.IpmiPort = types.Int64Value(int64(chassis.IpmiPort))
	state.IpmiUsername = types.StringValue(chassis.IpmiUsername)
	state.ZoneUuid = stringValueOrNull(chassis.ZoneUuid)
	state.State = stringValueOrNull(chassis.State)
	state.Status = stringValueOrNull(chassis.Status)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalChassisResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan baremetalChassisModel
	var state baremetalChassisModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateBaremetalChassisParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateBaremetalChassisParamDetail{
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			IpmiAddress:  stringPtrOrNil(plan.IpmiAddress.ValueString()),
			IpmiUsername: stringPtrOrNil(plan.IpmiUsername.ValueString()),
			IpmiPassword: stringPtrOrNil(plan.IpmiPassword.ValueString()),
		},
	}

	if !plan.IpmiPort.IsNull() && !plan.IpmiPort.IsUnknown() {
		p.Params.IpmiPort = intPtr(int(plan.IpmiPort.ValueInt64()))
	}

	chassis, err := r.client.UpdateBaremetalChassis(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Baremetal Chassis",
			"Could not update baremetal chassis, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(chassis.UUID)
	plan.Name = types.StringValue(chassis.Name)
	plan.Description = stringValueOrNull(chassis.Description)
	plan.ClusterUuid = types.StringValue(chassis.ClusterUuid)
	plan.IpmiAddress = types.StringValue(chassis.IpmiAddress)
	plan.IpmiPort = types.Int64Value(int64(chassis.IpmiPort))
	plan.IpmiUsername = types.StringValue(chassis.IpmiUsername)
	plan.ZoneUuid = stringValueOrNull(chassis.ZoneUuid)
	plan.State = stringValueOrNull(chassis.State)
	plan.Status = stringValueOrNull(chassis.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalChassisResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state baremetalChassisModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "Baremetal chassis UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteBaremetalChassis(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting Baremetal Chassis", "Could not delete baremetal chassis, unexpected error: "+err.Error())
		return
	}
}

func (r *baremetalChassisResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
