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

type vmResource struct {
	client *client.ZSClient
}

var (
	_ resource.Resource              = &vmResource{}
	_ resource.ResourceWithConfigure = &vmResource{}
)

/*ype vmDataSourceModel struct {
	VmInstance []vmModel `tfsdk:"vminstance"`
}
*/

type vmDataSourceModel struct {
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	ImageUuid      types.String `tfsdk:"imageuuid"`
	L3NetworkUuids types.String `tfsdk:"l3networkuuids"`
	//Type                 types.String `tfsdk:"type"`
	RootDiskOfferingUuid types.String `tfsdk:"rootdiskofferinguuid"`
	RootDiskSize         types.Int64  `tfsdk:"rootdisksize"`
	// AllDataDiskSizes                []dataDiskSizes `tfsdk:"alldatadisksize"`
	//ZoneUuid types.String `tfsdk:"zoneuuid"`
	ClusterUuid types.String `tfsdk:"clusteruuid"`
	//HostUuid                        types.String `tfsdk:"hostuuid"`
	//PrimaryStorageUuidForRootVolume types.String `tfsdk:"pristorageuuidforrootvolume"`
	Description types.String `tfsdk:"description"`
	//DefaultL3NetworkUuid            types.String `tfsdk:"defaultl3networkuuid"`
	//ResourceUuid                    types.String `tfsdk:"resourceuuid"`
	//TagUuid                         types.String `tfsdk:"taguuid"`
	//Strategy                        types.String `tfsdk:"strategy"`
	MemorySize types.Int64  `tfsdk:"memorysize"`
	CPUNum     types.Int64  `tfsdk:"cupnum"`
	IP         types.String `tfsdk:"ip"`
}

/*
type dataDiskSizes struct {
	//dataDiskSizes types.Int64 `tfsdk:"datadisksize"`
}
*/

func ZStackvmResource() resource.Resource {
	return &vmResource{}
}

// Configure implements resource.ResourceWithConfigure.
func (r *vmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements resource.Resource.
func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var createvmplan vmDataSourceModel

	diags := req.Plan.Get(ctx, &createvmplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vmInstanceParam := param.CreateVmInstanceParam{
		BaseParam: param.BaseParam{
			SystemTags: []string{"resourceConfig::vm::vm.clock.track::guest", "cdroms::Empty::None::None"},
			UserTags:   nil,
			RequestIp:  "",
		},
		Params: param.CreateVmInstanceDetailParam{
			Name: createvmplan.Name.ValueString(),
			//InstanceOfferingUUID:            "",
			ImageUUID:      createvmplan.ImageUuid.ValueString(), // "968e87334a12422fbe78c8b72bcfab68",
			L3NetworkUuids: []string{createvmplan.L3NetworkUuids.ValueString()},
			//Type:                 "",
			//RootDiskOfferingUuid: "",
			RootDiskSize: createvmplan.RootDiskSize.ValueInt64Pointer(), // &rootDiskSize,
			//	DataDiskOfferingUuids:           []string{"04229f19712d41cb990ab4b9252d9f93"},
			DataDiskSizes:                   []int64{10240},
			ZoneUuid:                        "",
			ClusterUUID:                     "",
			HostUuid:                        "",
			PrimaryStorageUuidForRootVolume: nil,                                       // createvmplan.PrimaryStorageUuidForRootVolume.ValueStringPointer(), //nil
			Description:                     createvmplan.Description.ValueString(),    //"Description",
			DefaultL3NetworkUuid:            createvmplan.L3NetworkUuids.ValueString(), // network[0].UUID,
			//ResourceUuid:                    createvmplan.ResourceUuid.ValueString(), // "56644230e0384ef6b84764530ef306cd",
			TagUuids:   nil, // createvmplan.TagUuid,                    //nil
			Strategy:   "",
			MemorySize: createvmplan.MemorySize.ValueInt64(), //,
			CpuNum:     createvmplan.CPUNum.ValueInt64(),
		},
	}

	instance, err := r.client.CreateVmInstance(vmInstanceParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not Create vm instance ", "Error "+err.Error(),
		)
		return
	}

	createvmplan.Uuid = types.StringValue(instance.UUID)
	createvmplan.Name = types.StringValue(instance.Name)
	createvmplan.Description = types.StringValue(instance.Description)
	createvmplan.MemorySize = types.Int64Value(instance.MemorySize)
	createvmplan.IP = types.StringValue(instance.VMNics[0].IP)

	//createvmplan.AllDataDiskSizes =types.Int64Value(instance.BaseTimeView.)
	diags = resp.State.Set(ctx, &createvmplan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete implements resource.Resource.
func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "wo kao"+state.Uuid.String())

	//Delete existing vm instance
	err := r.client.DestroyVmInstance(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not destroy vm instance", "Error: "+err.Error(),
		)
		return
	}
}

// Metadata implements resource.Resource.
func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

// Read implements resource.Resource.
func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmDataSourceModel
	req.State.Schema.GetAttributes()

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.client.GetVmInstance(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("cannot read vm uuid", "can not read vm uuid"+state.Uuid.ValueString()+err.Error())
		return
	}
	state.Uuid = types.StringValue(vm.UUID)
	state.Name = types.StringValue(vm.Name)
	state.Description = types.StringValue(vm.Description)
	state.ImageUuid = types.StringValue(vm.ImageUUID)
	state.MemorySize = types.Int64Value(vm.MemorySize)
	state.CPUNum = types.Int64Value(int64(vm.CPUNum))
	state.IP = types.StringValue(vm.VMNics[0].IP)
	state.L3NetworkUuids = types.StringValue(vm.DefaultL3NetworkUUID)

	diags = resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Schema implements resource.Resource.
func (r *vmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"ip": schema.StringAttribute{
				Computed: true,
			},
			"imageuuid": schema.StringAttribute{
				Optional: true,
			},
			"l3networkuuids": schema.StringAttribute{
				Optional: true,
			},
			"rootdiskofferinguuid": schema.StringAttribute{
				Optional: true,
			},
			"rootdisksize": schema.Int64Attribute{
				Optional: true,
			},

			"description": schema.StringAttribute{
				Optional: true,
			},
			"memorysize": schema.Int64Attribute{
				Optional: true,
			},
			"cupnum": schema.Int64Attribute{
				Optional: true,
			},
			"clusteruuid": schema.StringAttribute{
				Optional: true,
			},
		},
	}

}

// Update implements resource.Resource. //error 待修复
func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state vmDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "wo kao"+state.Uuid.String())
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name
	uuid := state.Uuid.ValueString()

	tflog.Info(ctx, "uuid"+uuid)
	//tflog.Info(ctx, (memorysize))
	//Generate vms API request body from plan
	vminstanceparam := param.UpdateVmInstanceParam{
		UpdateVmInstance: param.UpdateVmInstanceDetailParam{
			Name: name.ValueString(),
		},
	}

	_, err := r.client.UpdateVmInstance(uuid, vminstanceparam)
	if err != nil {
		resp.Diagnostics.AddError("", ""+err.Error())
		return
	}

	vm, err := r.client.GetVmInstance(uuid)
	if err != nil {
		resp.Diagnostics.AddError("", "")
		return
	}

	state.Uuid = types.StringValue(vm.UUID)
	state.Name = types.StringValue(vm.Name)
	//plan.MemorySize = types.Int64Value(vm.MemorySize)
	state.IP = types.StringValue(vm.VMNics[0].IP)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
