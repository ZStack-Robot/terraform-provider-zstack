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
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &vmCdRomResource{}
	_ resource.ResourceWithConfigure   = &vmCdRomResource{}
	_ resource.ResourceWithImportState = &vmCdRomResource{}
)

type vmCdRomResource struct {
	client *client.ZSClient
}

type vmCdRomResourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	VmInstanceUuid types.String `tfsdk:"vm_instance_uuid"`
	IsoUuid        types.String `tfsdk:"iso_uuid"`
	Description    types.String `tfsdk:"description"`
	DeviceId       types.Int64  `tfsdk:"device_id"`
}

func VmCdRomResource() resource.Resource {
	return &vmCdRomResource{}
}

func (r *vmCdRomResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*client.ZSClient)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", request.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *vmCdRomResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_vm_cdrom"
}

func (r *vmCdRomResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Manages a VM CD-ROM in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the VM CD-ROM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VM CD-ROM.",
			},
			"vm_instance_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VM instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"iso_uuid": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the ISO image.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the VM CD-ROM.",
			},
			"device_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The device ID of the CD-ROM.",
			},
		},
	}
}

func (r *vmCdRomResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan vmCdRomResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.CreateVmCdRomParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVmCdRomParamDetail{
			Name:           plan.Name.ValueString(),
			VmInstanceUuid: plan.VmInstanceUuid.ValueString(),
			IsoUuid:        stringPtrOrNil(plan.IsoUuid.ValueString()),
			Description:    stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	cdrom, err := r.client.CreateVmCdRom(params)
	if err != nil {
		response.Diagnostics.AddError("Error creating VM CD-ROM", err.Error())
		return
	}

	plan.Uuid = types.StringValue(cdrom.UUID)
	plan.Name = types.StringValue(cdrom.Name)
	plan.VmInstanceUuid = types.StringValue(cdrom.VmInstanceUuid)
	plan.IsoUuid = stringValueOrNull(cdrom.IsoUuid)
	plan.Description = stringValueOrNull(cdrom.Description)
	plan.DeviceId = types.Int64Value(int64(cdrom.DeviceId))

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vmCdRomResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state vmCdRomResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	queryParam.AddQ("uuid=" + state.Uuid.ValueString())

	cdroms, err := r.client.QueryVmCdRom(&queryParam)
	if err != nil {
		response.Diagnostics.AddError("Error reading VM CD-ROM", err.Error())
		return
	}

	if len(cdroms) == 0 {
		response.State.RemoveResource(ctx)
		return
	}

	cdrom := cdroms[0]
	state.Name = types.StringValue(cdrom.Name)
	state.VmInstanceUuid = types.StringValue(cdrom.VmInstanceUuid)
	state.IsoUuid = stringValueOrNull(cdrom.IsoUuid)
	state.Description = stringValueOrNull(cdrom.Description)
	state.DeviceId = types.Int64Value(int64(cdrom.DeviceId))

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vmCdRomResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var state vmCdRomResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var plan vmCdRomResourceModel
	diags = request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	params := param.UpdateVmCdRomParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVmCdRomParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	cdrom, err := r.client.UpdateVmCdRom(state.Uuid.ValueString(), params)
	if err != nil {
		response.Diagnostics.AddError("Error updating VM CD-ROM", err.Error())
		return
	}

	plan.Uuid = types.StringValue(cdrom.UUID)
	plan.Name = types.StringValue(cdrom.Name)
	plan.VmInstanceUuid = types.StringValue(cdrom.VmInstanceUuid)
	plan.IsoUuid = stringValueOrNull(cdrom.IsoUuid)
	plan.Description = stringValueOrNull(cdrom.Description)
	plan.DeviceId = types.Int64Value(int64(cdrom.DeviceId))

	diags = response.State.Set(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *vmCdRomResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state vmCdRomResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVmCdRom(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError("Error deleting VM CD-ROM", err.Error())
		return
	}
}

func (r *vmCdRomResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), request, response)
}
