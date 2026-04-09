// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

var (
	_ resource.Resource                = &zoneResource{}
	_ resource.ResourceWithConfigure   = &zoneResource{}
	_ resource.ResourceWithImportState = &zoneResource{}
)

type zoneResource struct {
	client *client.ZSClient
}

type zoneResourceModel struct {
	Uuid        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	State       types.String `tfsdk:"state"`
	Type        types.String `tfsdk:"type"`
}

func ZoneResource() resource.Resource {
	return &zoneResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *zoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ZSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. ", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *zoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

// Schema implements resource.Resource.
func (r *zoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manage ZStack zones. A zone is a logical grouping of resources such as clusters, hosts, and primary storage.",
		MarkdownDescription: "Manage ZStack zones. A zone is a logical grouping of resources such as clusters, hosts, and primary storage.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The UUID of the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the zone.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The state of the zone (Enabled, Disabled).",
				Validators: []validator.String{
					stringvalidator.OneOf("Enabled", "Disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan zoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating zone", map[string]any{"name": plan.Name.ValueString()})

	createParam := param.CreateZoneParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateZoneParamDetail{
			Name: plan.Name.ValueString(),
		},
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		createParam.Params.Description = stringPtr(plan.Description.ValueString())
	}

	zone, err := r.client.CreateZone(createParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Zone",
			"Could not create zone, unexpected error: "+err.Error(),
		)
		return
	}

	// If the desired state is not the default "Enabled", change it
	if !plan.State.IsNull() && plan.State.ValueString() != "" && plan.State.ValueString() != "Enabled" {
		stateEvent := deriveZoneStateEvent(plan.State.ValueString())
		_, err := r.client.ChangeZoneState(zone.UUID, param.ChangeZoneStateParam{
			Params: param.ChangeZoneStateParamDetail{
				StateEvent: stateEvent,
			},
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error changing Zone state",
				"Could not change zone state, unexpected error: "+err.Error(),
			)
			return
		}

		// Re-read the zone to get the updated state
		zone, err = r.client.GetZone(zone.UUID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading Zone",
				"Could not read zone after state change: "+err.Error(),
			)
			return
		}
	}

	state := zoneModelFromView(zone)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Read implements resource.Resource.
func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state zoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zone, err := findResourceByGet(r.client.GetZone, state.Uuid.ValueString())
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading Zone",
			"Could not read zone UUID "+state.Uuid.ValueString()+": "+err.Error(),
		)
		return
	}

	refreshedState := zoneModelFromView(zone)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Update implements resource.Resource.
func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan zoneResourceModel
	var state zoneResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := state.Uuid.ValueString()

	updateParam := param.UpdateZoneParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateZoneParamDetail{
			Name:        plan.Name.ValueString(),
			Description: stringPtrOrNil(plan.Description.ValueString()),
		},
	}

	if _, err := r.client.UpdateZone(uuid, updateParam); err != nil {
		resp.Diagnostics.AddError(
			"Error updating Zone",
			"Could not update zone, unexpected error: "+err.Error(),
		)
		return
	}

	// If the state has changed, update it
	if plan.State.ValueString() != state.State.ValueString() {
		stateEvent := deriveZoneStateEvent(plan.State.ValueString())
		if _, err := r.client.ChangeZoneState(uuid, param.ChangeZoneStateParam{
			Params: param.ChangeZoneStateParamDetail{
				StateEvent: stateEvent,
			},
		}); err != nil {
			resp.Diagnostics.AddError(
				"Error changing Zone state",
				"Could not change zone state, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Read back the updated resource
	zone, err := r.client.GetZone(uuid)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Zone",
			"Could not read zone UUID "+uuid+" after update: "+err.Error(),
		)
		return
	}

	refreshedState := zoneModelFromView(zone)

	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
}

// Delete implements resource.Resource.
func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state zoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.IsNull() || state.Uuid.ValueString() == "" {
		tflog.Warn(ctx, "zone uuid is empty, skip delete")
		return
	}

	if err := r.client.DeleteZone(state.Uuid.ValueString(), param.DeleteModePermissive); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Zone",
			"Could not delete zone, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState implements resource.ResourceWithImportState.
func (r *zoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func zoneModelFromView(z *view.ZoneInventoryView) zoneResourceModel {
	return zoneResourceModel{
		Uuid:        types.StringValue(z.UUID),
		Name:        types.StringValue(z.Name),
		Description: stringValueOrNull(z.Description),
		State:       stringValueOrNull(z.State),
		Type:        stringValueOrNull(z.Type),
	}
}

// deriveZoneStateEvent converts a desired state (e.g. "Enabled", "Disabled") to a state event (e.g. "enable", "disable").
func deriveZoneStateEvent(desiredState string) string {
	return strings.ToLower(strings.TrimSuffix(desiredState, "d"))
}
