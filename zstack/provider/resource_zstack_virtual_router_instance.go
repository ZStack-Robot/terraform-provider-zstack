// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &virtualRouterInstanceResource{}
	_ resource.ResourceWithConfigure = &virtualRouterInstanceResource{}
)

type virtualRouterInstanceResource struct {
	client *client.ZSClient
}

type virtualRouterInstanceResourceModel struct {
	Uuid                      types.String `tfsdk:"uuid"`
	Name                      types.String `tfsdk:"name"`
	State                     types.String `tfsdk:"state"`
	Status                    types.String `tfsdk:"status"`
	Description               types.String `tfsdk:"description"`
	VirtualRouterOfferingUuid types.String `tfsdk:"virtual_router_offering_uuid"`
	//ResourceUuid                    types.String `tfsdk:"resource_uuid" `                       // Resource UUID, if specified, the VM will use this value as its UUID.
	ZoneUuid                        types.String `tfsdk:"zone_uuid" `                           // Zone UUID, if specified, the VM will be created in the specified zone.
	ClusterUUID                     types.String `tfsdk:"cluster_uuid" `                        // Cluster UUID, if specified, the VM will be created in the specified cluster, higher priority than zoneUuid.
	HostUuid                        types.String `tfsdk:"host_uuid" `                           // Host UUID, if specified, the VM will be created on the specified host, higher priority than zoneUuid and clusterUuid.
	PrimaryStorageUuidForRootVolume types.String `tfsdk:"primary_storage_uuid_for_rootvolume" ` // Primary storage UUID, if specified, the root volume will be created on the specified primary storage.
}

// Configure implements resource.ResourceWithConfigure.
func (r *virtualRouterInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

func VirtualRouterInstanceResource() resource.Resource {
	return &virtualRouterInstanceResource{}
}

// Create implements resource.Resource.
func (r *virtualRouterInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan virtualRouterInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//systemTags = append(systemTags, param.vr)
	tflog.Info(ctx, "Configuring ZStack client")
	virtualRouterInstanceParam := param.CreateVirtualRouterInstanceParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVirtualRouterInstanceDetailParam{
			Name:                      plan.Name.ValueString(),
			Description:               plan.Description.ValueString(),
			VirtualRouterOfferingUuid: plan.VirtualRouterOfferingUuid.ValueString(),
			//PrimaryStorageUuidForRootVolume: plan.PrimaryStorageUuidForRootVolume.ValueStringPointer(),
		},
	}

	if !plan.Description.IsNull() {
		virtualRouterInstanceParam.Params.Description = plan.Description.ValueString()
	}

	if !plan.PrimaryStorageUuidForRootVolume.IsNull() {
		virtualRouterInstanceParam.Params.PrimaryStorageUuidForRootVolume = plan.PrimaryStorageUuidForRootVolume.ValueStringPointer()
	}

	vrInstance, err := r.client.CreateVirtualRouterInstance(virtualRouterInstanceParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Create virtual router instance", "Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(vrInstance.UUID)
	plan.Name = types.StringValue(vrInstance.Name)
	plan.VirtualRouterOfferingUuid = types.StringValue(vrInstance.InstanceOfferingUUID)

	if !plan.Description.IsNull() {
		plan.Description = types.StringValue(vrInstance.Description)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete implements resource.Resource.
func (r *virtualRouterInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state virtualRouterInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "virtual router image uuid is empty, so nothing to delete, skip it")
		return
	}

	err := r.client.DestroyVmInstance(state.Uuid.ValueString(), param.DeleteModeEnforcing)

	if err != nil {
		resp.Diagnostics.AddError("fail to delete virtual router image", ""+err.Error())
		return
	}
}

// Metadata implements resource.Resource.
func (r *virtualRouterInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_router_instance"
}

// Read implements resource.Resource.
func (r *virtualRouterInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state virtualRouterInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	vrInstance, err := r.client.GetVirtualRouterVm(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting virtual router image uuid", "Could not read image uuid "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(vrInstance.UUID)
	state.Name = types.StringValue(vrInstance.Name)

	if vrInstance.Description != "" {
		state.Description = types.StringValue(vrInstance.Description)
	} else {
		state.Description = types.StringNull()
	}
	//state.Description = types.StringValue(vrInstance.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema implements resource.Resource.
func (r *virtualRouterInstanceResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage virtual router instances in ZStack. " +
			"A virtual router instance is a virtual machine that provides network services such as routing, NAT, and VPN. " +
			"You can define the instance's properties, such as its name, description, associated virtual router offering, and deployment location.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the virtual router instance, uniquely identifying this resource in ZStack. Automatically generated upon creation.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual router instance. This is a required field in your environment.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "An optional description of the virtual router instance. Provides additional context or purpose of this instance.",
			},
			"virtual_router_offering_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the virtual router offering associated with this instance. Specifies the configuration and resource settings for the virtual router.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The current state of the virtual router instance. Possible values include 'Enabled', 'Disabled', etc.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The operational status of the virtual router instance. Indicates whether the instance is running, stopped, or in an error Status.",
			},
			"zone_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the zone where the virtual router instance will be deployed. Ensures the instance is placed within a specific zone.",
			},
			"cluster_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the cluster where the virtual router instance will be deployed. Takes precedence over 'zone_uuid' if both are specified.",
			},
			"host_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the host where the virtual router instance will be deployed. Takes precedence over both 'zone_uuid' and 'cluster_uuid' if specified.",
			},
			"primary_storage_uuid_for_rootvolume": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the primary storage where the root volume of the virtual router instance will be created. Ensures the root volume is placed on the specified storage.",
			},
		},
	}
}

func (r *virtualRouterInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}
