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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &baremetalInstanceResource{}
	_ resource.ResourceWithConfigure   = &baremetalInstanceResource{}
	_ resource.ResourceWithImportState = &baremetalInstanceResource{}
)

type baremetalInstanceResource struct {
	client *client.ZSClient
}

type baremetalInstanceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	ChassisUuid  types.String `tfsdk:"chassis_uuid"`
	ImageUuid    types.String `tfsdk:"image_uuid"`
	Password     types.String `tfsdk:"password"`
	TemplateUuid types.String `tfsdk:"template_uuid"`
	Username     types.String `tfsdk:"username"`
	Strategy     types.String `tfsdk:"strategy"`

	Platform     types.String `tfsdk:"platform"`
	ManagementIp types.String `tfsdk:"management_ip"`
	Port         types.Int64  `tfsdk:"port"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
}

func BaremetalInstanceResource() resource.Resource {
	return &baremetalInstanceResource{}
}

func (r *baremetalInstanceResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *baremetalInstanceResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_baremetal_instance"
}

func (r *baremetalInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage baremetal instances in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the baremetal instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the baremetal instance.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description for the baremetal instance.",
			},
			"chassis_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The chassis UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The image UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The instance login password.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"template_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The template UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Description: "The login username.",
			},
			"strategy": schema.StringAttribute{
				Optional:    true,
				Description: "The provisioning strategy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"platform": schema.StringAttribute{
				Computed:    true,
				Description: "The platform.",
			},
			"management_ip": schema.StringAttribute{
				Computed:    true,
				Description: "The management IP.",
			},
			"port": schema.Int64Attribute{
				Computed:    true,
				Description: "The management port.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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

func (r *baremetalInstanceResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan baremetalInstanceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.CreateBaremetalInstanceParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateBaremetalInstanceParamDetail{
			Name:         plan.Name.ValueString(),
			Description:  stringPtrOrNil(plan.Description.ValueString()),
			ChassisUuid:  plan.ChassisUuid.ValueString(),
			ImageUuid:    plan.ImageUuid.ValueString(),
			Password:     plan.Password.ValueString(),
			TemplateUuid: stringPtrOrNil(plan.TemplateUuid.ValueString()),
			Username:     stringPtrOrNil(plan.Username.ValueString()),
			Strategy:     stringPtrOrNil(plan.Strategy.ValueString()),
		},
	}

	instance, err := r.client.CreateBaremetalInstance(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Baremetal Instance",
			"Could not create baremetal instance, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(instance.UUID)
	plan.Name = types.StringValue(instance.Name)
	plan.Description = stringValueOrNull(instance.Description)
	plan.ChassisUuid = types.StringValue(instance.ChassisUuid)
	plan.ImageUuid = types.StringValue(instance.ImageUuid)
	plan.TemplateUuid = stringValueOrNull(instance.TemplateUuid)
	plan.Platform = stringValueOrNull(instance.Platform)
	plan.ManagementIp = stringValueOrNull(instance.ManagementIp)
	plan.Username = stringValueOrNull(instance.Username)
	plan.Port = types.Int64Value(int64(instance.Port))
	plan.State = stringValueOrNull(instance.State)
	plan.Status = stringValueOrNull(instance.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalInstanceResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state baremetalInstanceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	instance, err := findResourceByQuery(r.client.QueryBaremetalInstance, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading Baremetal Instance",
			"Could not read baremetal instance UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(instance.UUID)
	state.Name = types.StringValue(instance.Name)
	state.Description = stringValueOrNull(instance.Description)
	state.ChassisUuid = types.StringValue(instance.ChassisUuid)
	state.ImageUuid = types.StringValue(instance.ImageUuid)
	state.TemplateUuid = stringValueOrNull(instance.TemplateUuid)
	state.Platform = stringValueOrNull(instance.Platform)
	state.ManagementIp = stringValueOrNull(instance.ManagementIp)
	state.Username = stringValueOrNull(instance.Username)
	state.Port = types.Int64Value(int64(instance.Port))
	state.State = stringValueOrNull(instance.State)
	state.Status = stringValueOrNull(instance.Status)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalInstanceResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan baremetalInstanceModel
	var state baremetalInstanceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateBaremetalInstanceParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateBaremetalInstanceParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Password:    stringPtrOrNil(plan.Password.ValueString()),
			Platform:    stringPtrOrNil(plan.Platform.ValueString()),
		},
	}

	instance, err := r.client.UpdateBaremetalInstance(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating Baremetal Instance",
			"Could not update baremetal instance, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(instance.UUID)
	plan.Name = types.StringValue(instance.Name)
	plan.Description = stringValueOrNull(instance.Description)
	plan.ChassisUuid = types.StringValue(instance.ChassisUuid)
	plan.ImageUuid = types.StringValue(instance.ImageUuid)
	plan.TemplateUuid = stringValueOrNull(instance.TemplateUuid)
	plan.Platform = stringValueOrNull(instance.Platform)
	plan.ManagementIp = stringValueOrNull(instance.ManagementIp)
	plan.Username = stringValueOrNull(instance.Username)
	plan.Port = types.Int64Value(int64(instance.Port))
	plan.State = stringValueOrNull(instance.State)
	plan.Status = stringValueOrNull(instance.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *baremetalInstanceResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state baremetalInstanceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}


	err := r.client.DestroyBaremetalInstance(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting Baremetal Instance", "Could not delete baremetal instance, unexpected error: "+err.Error())
		return
	}
}

func (r *baremetalInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
