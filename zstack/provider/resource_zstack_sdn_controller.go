// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &sdnControllerResource{}
	_ resource.ResourceWithConfigure   = &sdnControllerResource{}
	_ resource.ResourceWithImportState = &sdnControllerResource{}
)

type sdnControllerResource struct {
	client *client.ZSClient
}

type sdnControllerResourceModel struct {
	Uuid          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	VendorType    types.String `tfsdk:"vendor_type"`
	VendorVersion types.String `tfsdk:"vendor_version"`
	Ip            types.String `tfsdk:"ip"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Status        types.String `tfsdk:"status"`
}

func SdnControllerResource() resource.Resource {
	return &sdnControllerResource{}
}

func (r *sdnControllerResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *sdnControllerResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_sdn_controller"
}

func (r *sdnControllerResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage SDN (Software-Defined Network) controllers in ZStack. " +
			"An SDN controller is used to manage network virtualization and provides centralized network management capabilities.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the SDN controller.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the SDN controller.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the SDN controller.",
			},
			"vendor_type": schema.StringAttribute{
				Required:    true,
				Description: "The vendor type of the SDN controller (e.g., OVN, ODL).",
			},
			"vendor_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The version of the SDN controller vendor software.",
			},
			"ip": schema.StringAttribute{
				Required:    true,
				Description: "The IP address of the SDN controller.",
			},
			"username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The username for authentication with the SDN controller.",
			},
			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The password for authentication with the SDN controller.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the SDN controller.",
			},
		},
	}
}

func (r *sdnControllerResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan sdnControllerResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddSdnControllerParam{
		BaseParam: param.BaseParam{},
		Params: param.AddSdnControllerParamDetail{
			Name:          plan.Name.ValueString(),
			Description:   stringPtrOrNil(plan.Description.ValueString()),
			VendorType:    plan.VendorType.ValueString(),
			VendorVersion: stringPtrOrNil(plan.VendorVersion.ValueString()),
			Ip:            plan.Ip.ValueString(),
			UserName:      stringPtrOrNil(plan.Username.ValueString()),
			Password:      stringPtrOrNil(plan.Password.ValueString()),
		},
	}

	result, err := r.client.AddSdnController(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create SDN controller",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.VendorType = types.StringValue(result.VendorType)
	plan.VendorVersion = types.StringValue(result.VendorVersion)
	plan.Ip = types.StringValue(result.Ip)
	plan.Status = types.StringValue(result.Status)
	// Keep the username/password from the plan since API may return empty
	if plan.Username.IsNull() {
		plan.Username = types.StringValue("")
	}
	if plan.Password.IsNull() {
		plan.Password = types.StringValue("")
	}

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *sdnControllerResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state sdnControllerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	items, err := r.client.QuerySdnController(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query SDN controllers. It may have been deleted.: "+err.Error())
		state = sdnControllerResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, item := range items {
		if item.UUID == state.Uuid.ValueString() {
			// Update state with the matched SDN controller details
			state.Uuid = types.StringValue(item.UUID)
			state.Name = types.StringValue(item.Name)
			state.Description = types.StringValue(item.Description)
			state.VendorType = types.StringValue(item.VendorType)
			state.VendorVersion = types.StringValue(item.VendorVersion)
			state.Ip = types.StringValue(item.Ip)
			state.Status = types.StringValue(item.Status)
			// Don't overwrite username/password from API since they're sensitive/write-only
			found = true
			break
		}
	}
	if !found {
		// If the SDN controller is not found, mark it as unmanaged
		tflog.Warn(ctx, "SDN controller not found. It might have been deleted outside of Terraform.")
		state = sdnControllerResourceModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *sdnControllerResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan sdnControllerResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var state sdnControllerResourceModel
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateSdnControllerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSdnControllerParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	result, err := r.client.UpdateSdnController(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update SDN controller",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(result.UUID)
	plan.Name = types.StringValue(result.Name)
	plan.Description = types.StringValue(result.Description)
	plan.VendorType = types.StringValue(result.VendorType)
	plan.VendorVersion = types.StringValue(result.VendorVersion)
	plan.Ip = types.StringValue(result.Ip)
	plan.Status = types.StringValue(result.Status)
	// Keep the username/password from the plan since API may return empty
	if plan.Username.IsNull() {
		plan.Username = types.StringValue("")
	}
	if plan.Password.IsNull() {
		plan.Password = types.StringValue("")
	}

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *sdnControllerResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state sdnControllerResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "SDN controller UUID is empty, skipping delete.")
		return
	}

	err := r.client.RemoveSdnController(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete SDN controller", ""+err.Error())
		return
	}

}

func (r *sdnControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
