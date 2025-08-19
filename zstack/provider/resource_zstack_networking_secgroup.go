// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

var (
	_ resource.Resource              = &securityGroupResource{}
	_ resource.ResourceWithConfigure = &securityGroupResource{}
)

type securityGroupResource struct {
	client *client.ZSClient
}

type securityGroupModel struct {
	Uuid              types.String `tfsdk:"uuid"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	VSwitchType       types.String `tfsdk:"vswitch_type"`
	SdnControllerUuid types.String `tfsdk:"sdn_controller_uuid"`
	IpVersion         types.Int32  `tfsdk:"ip_version"`
}

func SecurityGroupResource() resource.Resource {
	return &securityGroupResource{}
}

func (r *securityGroupResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *securityGroupResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_networking_secgroup"
}

func (r *securityGroupResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource manages a security group in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the security group.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "A name for the security group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the security group.",
			},
			"sdn_controller_uuid": schema.StringAttribute{
				Optional: true,
				Description: "SDN Controller UUID. Required when vswitch_type is 'OvnDpdk'. " +
					"Used to build SystemTags 'SdnControllerUuid::<uuid>' on create.",
			},
			"ip_version": schema.Int32Attribute{
				Required:    true,
				Description: "",
			},
			"vswitch_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The type of virtual switch to use for the security group, e.g., 'LinuxBridge' or 'OvnDpdk'. Defaults to 'LinuxBridge' if not specified.",
				Validators: []validator.String{
					stringvalidator.OneOf("LinuxBridge", "OvnDpdk"),
				},
			},
		},
	}
}

func (r *securityGroupResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var securityGroupPlan securityGroupModel
	diags := request.Plan.Get(ctx, &securityGroupPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	vType := securityGroupPlan.VSwitchType.ValueString()
	if vType == "" {
		vType = "LinuxBridge" // Default value
	}

	securityGroupPlan.VSwitchType = types.StringValue(vType)
	p := param.CreateSecurityGroupParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSecurityGroupDetailParam{
			Name:        securityGroupPlan.Name.ValueString(),
			Description: securityGroupPlan.Description.ValueString(),
			IpVersion:   int(securityGroupPlan.IpVersion.ValueInt32()),
			VSwitchType: vType,
		},
	}

	//OvnDpdk requires system tag
	if vType == "OvnDpdk" {
		u := securityGroupPlan.SdnControllerUuid.ValueString()
		if u == "" {
			// Double-check; ValidateConfig should have caught this already.
			response.Diagnostics.AddAttributeError(
				path.Root("sdn_controller_uuid"),
				"Missing SDN Controller UUID",
				"When 'vswitch_type' is 'OvnDpdk', 'sdn_controller_uuid' must be provided.",
			)
			return
		}
		p.BaseParam.SystemTags = []string{"SdnControllerUuid::" + u}
	}

	secGroup, err := r.client.CreateSecurityGroup(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to create security group",
			"Error "+err.Error(),
		)
		return
	}

	securityGroupPlan.Uuid = types.StringValue(secGroup.UUID)
	securityGroupPlan.Name = types.StringValue(secGroup.Name)
	securityGroupPlan.Description = types.StringValue(secGroup.Description)
	if secGroup.VSwitchType == "" {
		// Use the value we sent to the API
		tflog.Debug(ctx, fmt.Sprintf("API returned empty VSwitchType, using value from plan: %s", vType))
	} else {
		vType = secGroup.VSwitchType
	}
	securityGroupPlan.VSwitchType = types.StringValue(vType)

	diags = response.State.Set(ctx, securityGroupPlan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

}

func (r *securityGroupResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state securityGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	secGroups, err := r.client.GetSecurityGroup(state.Uuid.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Failed to read security group", err.Error())
		return
	}

	if secGroups == nil {
		response.State.RemoveResource(ctx)
		return
	}

	sg := secGroups

	vsType := sg.VSwitchType
	if vsType == "" {
		if !state.VSwitchType.IsNull() && state.VSwitchType.ValueString() != "" {
			vsType = state.VSwitchType.ValueString()
			tflog.Debug(ctx, fmt.Sprintf("API returned empty VSwitchType, using value from state: %s", vsType))
		} else {
			vsType = "LinuxBridge"
		}

	}

	var descriptionValue types.String
	if sg.Description == "" && state.Description.IsNull() {
		descriptionValue = types.StringNull()
	} else {
		descriptionValue = types.StringValue(sg.Description)
	}

	state.Uuid = types.StringValue(secGroups.UUID)
	state.Description = descriptionValue
	state.VSwitchType = types.StringValue(vsType)
	state.Name = types.StringValue(secGroups.Name)

	// Update State
	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *securityGroupResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	response.Diagnostics.AddError(
		"Update Not Supported",
		"This resource does not support updates. Any changes will require replacement.",
	)
}

func (r *securityGroupResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state securityGroupModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "Security group UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteSecurityGroup(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to delete security group",
			"Error "+err.Error(),
		)
		return
	}

}
