// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
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
	_ resource.Resource                = &snmpAgentResource{}
	_ resource.ResourceWithConfigure   = &snmpAgentResource{}
	_ resource.ResourceWithImportState = &snmpAgentResource{}
)

type snmpAgentResource struct {
	client *client.ZSClient
}

type snmpAgentModel struct {
	Uuid             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Version          types.String `tfsdk:"version"`
	ReadCommunity    types.String `tfsdk:"read_community"`
	UserName         types.String `tfsdk:"user_name"`
	AuthAlgorithm    types.String `tfsdk:"auth_algorithm"`
	AuthPassword     types.String `tfsdk:"auth_password"`
	PrivacyAlgorithm types.String `tfsdk:"privacy_algorithm"`
	PrivacyPassword  types.String `tfsdk:"privacy_password"`
	Port             types.Int64  `tfsdk:"port"`
	Status           types.String `tfsdk:"status"`
	SecurityLevel    types.String `tfsdk:"security_level"`
}

func SnmpAgentResource() resource.Resource {
	return &snmpAgentResource{}
}

func (r *snmpAgentResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *snmpAgentResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_snmp_agent"
}

func (r *snmpAgentResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage SNMP agents in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the SNMP agent.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the SNMP agent.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "SNMP version.",
			},
			"read_community": schema.StringAttribute{
				Optional:    true,
				Description: "SNMP read community.",
			},
			"user_name": schema.StringAttribute{
				Optional:    true,
				Description: "SNMP username.",
			},
			"auth_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "SNMP authentication algorithm.",
			},
			"auth_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "SNMP authentication password.",
			},
			"privacy_algorithm": schema.StringAttribute{
				Optional:    true,
				Description: "SNMP privacy algorithm.",
			},
			"privacy_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "SNMP privacy password.",
			},
			"port": schema.Int64Attribute{
				Required:    true,
				Description: "SNMP agent port.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "SNMP agent status.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"security_level": schema.StringAttribute{
				Computed:    true,
				Description: "SNMP security level.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *snmpAgentResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan snmpAgentModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateSnmpAgentParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateSnmpAgentParamDetail{
			Version:          plan.Version.ValueString(),
			ReadCommunity:    stringPtrOrNil(plan.ReadCommunity.ValueString()),
			UserName:         stringPtrOrNil(plan.UserName.ValueString()),
			AuthAlgorithm:    stringPtrOrNil(plan.AuthAlgorithm.ValueString()),
			AuthPassword:     stringPtrOrNil(plan.AuthPassword.ValueString()),
			PrivacyAlgorithm: stringPtrOrNil(plan.PrivacyAlgorithm.ValueString()),
			PrivacyPassword:  stringPtrOrNil(plan.PrivacyPassword.ValueString()),
			Port:             int(plan.Port.ValueInt64()),
		},
	}

	snmpAgent, err := r.client.CreateSnmpAgent(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating SNMP Agent",
			"Could not create snmp agent, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(snmpAgent.UUID)
	plan.Name = types.StringValue(snmpAgent.Name)
	plan.Version = types.StringValue(snmpAgent.Version)
	plan.ReadCommunity = stringValueOrNull(snmpAgent.ReadCommunity)
	plan.UserName = stringValueOrNull(snmpAgent.UserName)
	plan.AuthAlgorithm = stringValueOrNull(snmpAgent.AuthAlgorithm)
	plan.AuthPassword = stringValueOrNull(snmpAgent.AuthPassword)
	plan.PrivacyAlgorithm = stringValueOrNull(snmpAgent.PrivacyAlgorithm)
	plan.PrivacyPassword = stringValueOrNull(snmpAgent.PrivacyPassword)
	plan.Port = types.Int64Value(int64(snmpAgent.Port))
	plan.Status = stringValueOrNull(snmpAgent.Status)
	plan.SecurityLevel = stringValueOrNull(snmpAgent.SecurityLevel)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snmpAgentResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state snmpAgentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	snmpAgent, err := findResourceByQuery(r.client.QuerySnmpAgent, state.Uuid.ValueString())

	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError(
			"Error reading SNMP Agent",
			"Could not read SNMP Agent, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(snmpAgent.UUID)
	state.Name = types.StringValue(snmpAgent.Name)
	state.Version = types.StringValue(snmpAgent.Version)
	state.ReadCommunity = stringValueOrNull(snmpAgent.ReadCommunity)
	state.UserName = stringValueOrNull(snmpAgent.UserName)
	state.AuthAlgorithm = stringValueOrNull(snmpAgent.AuthAlgorithm)
	state.AuthPassword = stringValueOrNull(snmpAgent.AuthPassword)
	state.PrivacyAlgorithm = stringValueOrNull(snmpAgent.PrivacyAlgorithm)
	state.PrivacyPassword = stringValueOrNull(snmpAgent.PrivacyPassword)
	state.Port = types.Int64Value(int64(snmpAgent.Port))
	state.Status = stringValueOrNull(snmpAgent.Status)
	state.SecurityLevel = stringValueOrNull(snmpAgent.SecurityLevel)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snmpAgentResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan snmpAgentModel
	var state snmpAgentModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	p := param.UpdateSnmpAgentParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateSnmpAgentParamDetail{
			Uuid:             state.Uuid.ValueString(),
			Version:          plan.Version.ValueString(),
			ReadCommunity:    stringPtrOrNil(plan.ReadCommunity.ValueString()),
			UserName:         stringPtrOrNil(plan.UserName.ValueString()),
			AuthAlgorithm:    stringPtrOrNil(plan.AuthAlgorithm.ValueString()),
			AuthPassword:     stringPtrOrNil(plan.AuthPassword.ValueString()),
			PrivacyAlgorithm: stringPtrOrNil(plan.PrivacyAlgorithm.ValueString()),
			PrivacyPassword:  stringPtrOrNil(plan.PrivacyPassword.ValueString()),
			Port:             int(plan.Port.ValueInt64()),
		},
	}

	snmpAgent, err := r.client.UpdateSnmpAgent(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Error updating SNMP Agent",
			"Could not update snmp agent, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(snmpAgent.UUID)
	plan.Name = types.StringValue(snmpAgent.Name)
	plan.Version = types.StringValue(snmpAgent.Version)
	plan.ReadCommunity = stringValueOrNull(snmpAgent.ReadCommunity)
	plan.UserName = stringValueOrNull(snmpAgent.UserName)
	plan.AuthAlgorithm = stringValueOrNull(snmpAgent.AuthAlgorithm)
	plan.AuthPassword = stringValueOrNull(snmpAgent.AuthPassword)
	plan.PrivacyAlgorithm = stringValueOrNull(snmpAgent.PrivacyAlgorithm)
	plan.PrivacyPassword = stringValueOrNull(snmpAgent.PrivacyPassword)
	plan.Port = types.Int64Value(int64(snmpAgent.Port))
	plan.Status = stringValueOrNull(snmpAgent.Status)
	plan.SecurityLevel = stringValueOrNull(snmpAgent.SecurityLevel)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *snmpAgentResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state snmpAgentModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, "ZStack API does not support deleting SNMP agents. The resource will be removed from Terraform state but will remain in ZStack.")
}

func (r *snmpAgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
