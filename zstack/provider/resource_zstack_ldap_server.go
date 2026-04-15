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
)

var (
	_ resource.Resource                = &ldapServerResource{}
	_ resource.ResourceWithConfigure   = &ldapServerResource{}
	_ resource.ResourceWithImportState = &ldapServerResource{}
)

type ldapServerResource struct {
	client *client.ZSClient
}

type ldapServerModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Url         types.String `tfsdk:"url"`
	Base        types.String `tfsdk:"base"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	Encryption  types.String `tfsdk:"encryption"`
	Scope       types.String `tfsdk:"scope"`
}

func LdapServerResource() resource.Resource {
	return &ldapServerResource{}
}

func (r *ldapServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ldapServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_server"
}

func (r *ldapServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage LDAP server in ZStack.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the LDAP server.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the LDAP server.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the LDAP server.",
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "LDAP URL.",
			},
			"base": schema.StringAttribute{
				Required:    true,
				Description: "LDAP base DN.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "LDAP bind username.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "LDAP bind password.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption": schema.StringAttribute{
				Required:    true,
				Description: "LDAP encryption method.",
			},
			"scope": schema.StringAttribute{
				Required:    true,
				Description: "LDAP search scope.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *ldapServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ldapServerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddWarning("Client Not Configured", "The client was not properly configured.")
		return
	}

	p := param.AddLdapServerParam{
		BaseParam: param.BaseParam{},
		Params: param.AddLdapServerParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Url:         plan.Url.ValueString(),
			Base:        plan.Base.ValueString(),
			Username:    plan.Username.ValueString(),
			Password:    plan.Password.ValueString(),
			Encryption:  plan.Encryption.ValueString(),
			Scope:       plan.Scope.ValueString(),
		},
	}

	item, err := r.client.AddLdapServer(p)
	if err != nil {
		resp.Diagnostics.AddError("Error creating LDAP Server", "Could not create LDAP server, unexpected error: "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Url = stringValueOrNull(item.Url)
	plan.Base = stringValueOrNull(item.Base)
	plan.Username = stringValueOrNull(item.Username)
	plan.Encryption = stringValueOrNull(item.Encryption)
	plan.Scope = stringValueOrNull(item.Scope)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ldapServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ldapServerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, err := findResourceByQuery(r.client.QueryLdapServer, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading LDAP Server",
			"Could not read LDAP Server, unexpected error: "+err.Error(),
		)
		return
	}

	state.Uuid = types.StringValue(item.UUID)
	state.Name = types.StringValue(item.Name)
	state.Description = stringValueOrNull(item.Description)
	state.Url = stringValueOrNull(item.Url)
	state.Base = stringValueOrNull(item.Base)
	state.Username = stringValueOrNull(item.Username)
	state.Encryption = stringValueOrNull(item.Encryption)
	state.Scope = stringValueOrNull(item.Scope)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *ldapServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ldapServerModel
	var state ldapServerModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := param.UpdateLdapServerParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateLdapServerParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
			Url:         stringPtr(plan.Url.ValueString()),
			Base:        stringPtr(plan.Base.ValueString()),
			Username:    stringPtr(plan.Username.ValueString()),
			Password:    stringPtr(plan.Password.ValueString()),
			Encryption:  stringPtr(plan.Encryption.ValueString()),
		},
	}

	item, err := r.client.UpdateLdapServer(state.Uuid.ValueString(), p)
	if err != nil {
		resp.Diagnostics.AddError("Error updating LDAP Server", "Could not update LDAP server, unexpected error: "+err.Error())
		return
	}

	plan.Uuid = types.StringValue(item.UUID)
	plan.Name = types.StringValue(item.Name)
	plan.Description = stringValueOrNull(item.Description)
	plan.Url = stringValueOrNull(item.Url)
	plan.Base = stringValueOrNull(item.Base)
	plan.Username = stringValueOrNull(item.Username)
	plan.Encryption = stringValueOrNull(item.Encryption)
	plan.Scope = stringValueOrNull(item.Scope)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *ldapServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ldapServerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


	if err := r.client.DeleteLdapServer(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError("Error deleting LDAP Server", "Could not delete LDAP server, unexpected error: "+err.Error())
		return
	}
}

func (r *ldapServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
