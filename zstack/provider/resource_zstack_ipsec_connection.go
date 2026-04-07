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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

var (
	_ resource.Resource                = &ipsecConnectionResource{}
	_ resource.ResourceWithConfigure   = &ipsecConnectionResource{}
	_ resource.ResourceWithImportState = &ipsecConnectionResource{}
)

type ipsecConnectionResource struct {
	client *client.ZSClient
}

type ipsecConnectionResourceModel struct {
	Uuid                      types.String `tfsdk:"uuid"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	VipUuid                   types.String `tfsdk:"vip_uuid"`
	PeerAddress               types.String `tfsdk:"peer_address"`
	AuthKey                   types.String `tfsdk:"auth_key"`
	AuthMode                  types.String `tfsdk:"auth_mode"`
	IkeAuthAlgorithm          types.String `tfsdk:"ike_auth_algorithm"`
	IkeEncryptionAlgorithm    types.String `tfsdk:"ike_encryption_algorithm"`
	PolicyAuthAlgorithm       types.String `tfsdk:"policy_auth_algorithm"`
	PolicyEncryptionAlgorithm types.String `tfsdk:"policy_encryption_algorithm"`
	Pfs                       types.String `tfsdk:"pfs"`
	PolicyMode                types.String `tfsdk:"policy_mode"`
	TransformProtocol         types.String `tfsdk:"transform_protocol"`
	State                     types.String `tfsdk:"state"`
	Status                    types.String `tfsdk:"status"`
}

func IPsecConnectionResource() resource.Resource {
	return &ipsecConnectionResource{}
}

func (r *ipsecConnectionResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *ipsecConnectionResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_ipsec_connection"
}

func (r *ipsecConnectionResource) Schema(_ context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "This resource allows you to manage IPsec connections in ZStack. " +
			"An IPsec connection provides secure site-to-site VPN connectivity.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the IPsec connection.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description for the IPsec connection.",
			},
			"vip_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the VIP.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"peer_address": schema.StringAttribute{
				Required:    true,
				Description: "The peer address for the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The authentication key for the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The authentication mode for the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ike_auth_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The IKE authentication algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ike_encryption_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The IKE encryption algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_auth_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The policy authentication algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_encryption_algorithm": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The policy encryption algorithm.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pfs": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The Perfect Forward Secrecy setting.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The policy mode for the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"transform_protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The transform protocol for the IPsec connection.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The state of the IPsec connection.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of the IPsec connection.",
			},
		},
	}
}

func (r *ipsecConnectionResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan ipsecConnectionResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.CreateIPsecConnectionParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateIPsecConnectionParamDetail{
			Name:                      plan.Name.ValueString(),
			Description:               stringPtrOrNil(plan.Description.ValueString()),
			VipUuid:                   plan.VipUuid.ValueString(),
			PeerAddress:               plan.PeerAddress.ValueString(),
			AuthKey:                   plan.AuthKey.ValueString(),
			AuthMode:                  stringPtrOrNil(plan.AuthMode.ValueString()),
			IkeAuthAlgorithm:          stringPtrOrNil(plan.IkeAuthAlgorithm.ValueString()),
			IkeEncryptionAlgorithm:    stringPtrOrNil(plan.IkeEncryptionAlgorithm.ValueString()),
			PolicyAuthAlgorithm:       stringPtrOrNil(plan.PolicyAuthAlgorithm.ValueString()),
			PolicyEncryptionAlgorithm: stringPtrOrNil(plan.PolicyEncryptionAlgorithm.ValueString()),
			Pfs:                       stringPtrOrNil(plan.Pfs.ValueString()),
			PolicyMode:                stringPtrOrNil(plan.PolicyMode.ValueString()),
			TransformProtocol:         stringPtrOrNil(plan.TransformProtocol.ValueString()),
		},
	}

	ipsecConnection, err := r.client.CreateIPsecConnection(p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to create IPsec connection",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(ipsecConnection.UUID)
	plan.Name = types.StringValue(ipsecConnection.Name)
	plan.Description = types.StringValue(ipsecConnection.Description)
	plan.VipUuid = types.StringValue(ipsecConnection.VipUuid)
	plan.PeerAddress = types.StringValue(ipsecConnection.PeerAddress)
	plan.AuthMode = types.StringValue(ipsecConnection.AuthMode)
	plan.IkeAuthAlgorithm = types.StringValue(ipsecConnection.IkeAuthAlgorithm)
	plan.IkeEncryptionAlgorithm = types.StringValue(ipsecConnection.IkeEncryptionAlgorithm)
	plan.PolicyAuthAlgorithm = types.StringValue(ipsecConnection.PolicyAuthAlgorithm)
	plan.PolicyEncryptionAlgorithm = types.StringValue(ipsecConnection.PolicyEncryptionAlgorithm)
	plan.Pfs = types.StringValue(ipsecConnection.Pfs)
	plan.PolicyMode = types.StringValue(ipsecConnection.PolicyMode)
	plan.TransformProtocol = types.StringValue(ipsecConnection.TransformProtocol)
	plan.State = types.StringValue(ipsecConnection.State)
	plan.Status = types.StringValue(ipsecConnection.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *ipsecConnectionResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state ipsecConnectionResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	queryParam := param.NewQueryParam()
	ipsecConnections, err := r.client.QueryIPSecConnection(&queryParam)

	if err != nil {
		tflog.Warn(ctx, "Unable to query IPsec connections. It may have been deleted.: "+err.Error())
		state = ipsecConnectionResourceModel{
			Uuid: types.StringValue(""),
		}
		diags = response.State.Set(ctx, &state)
		response.Diagnostics.Append(diags...)
		return
	}

	found := false

	for _, ipsecConnection := range ipsecConnections {
		if ipsecConnection.UUID == state.Uuid.ValueString() {
			state.Uuid = types.StringValue(ipsecConnection.UUID)
			state.Name = types.StringValue(ipsecConnection.Name)
			state.Description = types.StringValue(ipsecConnection.Description)
			state.VipUuid = types.StringValue(ipsecConnection.VipUuid)
			state.PeerAddress = types.StringValue(ipsecConnection.PeerAddress)
			state.AuthMode = types.StringValue(ipsecConnection.AuthMode)
			state.IkeAuthAlgorithm = types.StringValue(ipsecConnection.IkeAuthAlgorithm)
			state.IkeEncryptionAlgorithm = types.StringValue(ipsecConnection.IkeEncryptionAlgorithm)
			state.PolicyAuthAlgorithm = types.StringValue(ipsecConnection.PolicyAuthAlgorithm)
			state.PolicyEncryptionAlgorithm = types.StringValue(ipsecConnection.PolicyEncryptionAlgorithm)
			state.Pfs = types.StringValue(ipsecConnection.Pfs)
			state.PolicyMode = types.StringValue(ipsecConnection.PolicyMode)
			state.TransformProtocol = types.StringValue(ipsecConnection.TransformProtocol)
			state.State = types.StringValue(ipsecConnection.State)
			state.Status = types.StringValue(ipsecConnection.Status)
			found = true
			break
		}
	}
	if !found {
		tflog.Warn(ctx, "IPsec connection not found. It might have been deleted outside of Terraform.")
		state = ipsecConnectionResourceModel{
			Uuid: types.StringValue(""),
		}
	}

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *ipsecConnectionResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan ipsecConnectionResourceModel
	var state ipsecConnectionResourceModel

	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	diags = request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		response.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.UpdateIPsecConnectionParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateIPsecConnectionParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	ipsecConnection, err := r.client.UpdateIPsecConnection(state.Uuid.ValueString(), p)
	if err != nil {
		response.Diagnostics.AddError(
			"Fail to update IPsec connection",
			"Error "+err.Error(),
		)
		return
	}

	plan.Uuid = types.StringValue(ipsecConnection.UUID)
	plan.Name = types.StringValue(ipsecConnection.Name)
	plan.Description = types.StringValue(ipsecConnection.Description)
	plan.VipUuid = types.StringValue(ipsecConnection.VipUuid)
	plan.PeerAddress = types.StringValue(ipsecConnection.PeerAddress)
	plan.AuthMode = types.StringValue(ipsecConnection.AuthMode)
	plan.IkeAuthAlgorithm = types.StringValue(ipsecConnection.IkeAuthAlgorithm)
	plan.IkeEncryptionAlgorithm = types.StringValue(ipsecConnection.IkeEncryptionAlgorithm)
	plan.PolicyAuthAlgorithm = types.StringValue(ipsecConnection.PolicyAuthAlgorithm)
	plan.PolicyEncryptionAlgorithm = types.StringValue(ipsecConnection.PolicyEncryptionAlgorithm)
	plan.Pfs = types.StringValue(ipsecConnection.Pfs)
	plan.PolicyMode = types.StringValue(ipsecConnection.PolicyMode)
	plan.TransformProtocol = types.StringValue(ipsecConnection.TransformProtocol)
	plan.State = types.StringValue(ipsecConnection.State)
	plan.Status = types.StringValue(ipsecConnection.Status)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *ipsecConnectionResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state ipsecConnectionResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "IPsec connection UUID is empty, skipping delete.")
		return
	}

	err := r.client.DeleteIPsecConnection(state.Uuid.ValueString(), param.DeleteModePermissive)

	if err != nil {
		response.Diagnostics.AddError("fail to delete IPsec connection", ""+err.Error())
		return
	}
}

func (r *ipsecConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
