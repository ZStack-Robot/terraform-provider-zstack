// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"terraform-provider-zstack/zstack/utils"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
)

type gpuDeviceTyp string

type vmResource struct {
	client *client.ZSClient
}

var (
	_ resource.Resource              = &vmResource{}
	_ resource.ResourceWithConfigure = &vmResource{}
)

var networkModelAttrTypes = map[string]attr.Type{
	"uuid":    types.StringType,
	"ip":      types.StringType,
	"netmask": types.StringType,
	"gateway": types.StringType,
}

const (
	mdevDevice gpuDeviceTyp = "mdevDevice"
	pciDevice  gpuDeviceTyp = "pciDevice"
)

type diskModel struct {
	Size               types.Int64  `tfsdk:"size"`
	OfferingUuid       types.String `tfsdk:"offering_uuid"`
	VirtioSCSI         types.Bool   `tfsdk:"virtio_scsi"`
	PrimaryStorageUuid types.String `tfsdk:"primary_storage_uuid"`
	CephPoolName       types.String `tfsdk:"ceph_pool_name"`
}

type networkModel struct {
	Uuid types.String `tfsdk:"uuid"`
}

type gpuModel struct {
	Uuid types.String `tfsdk:"uuid"`
	Type types.String `tfsdk:"type"`
}

type gpuSpecModel struct {
	Uuid   types.String `tfsdk:"uuid"`
	Type   types.String `tfsdk:"type"`
	Number types.Int64  `tfsdk:"number"`
}

type vmInstanceDataSourceModel struct {
	Uuid                 types.String `tfsdk:"uuid"`
	Name                 types.String `tfsdk:"name"`
	ImageUuid            types.String `tfsdk:"image_uuid"`
	L3NetworkUuids       types.List   `tfsdk:"l3_network_uuids"`
	RootDisk             types.Object `tfsdk:"root_disk"`
	DataDisks            types.List   `tfsdk:"data_disks"`
	Networks             types.List   `tfsdk:"networks"`
	ZoneUuid             types.String `tfsdk:"zone_uuid"`
	ClusterUuid          types.String `tfsdk:"cluster_uuid"`
	HostUuid             types.String `tfsdk:"host_uuid"`
	Description          types.String `tfsdk:"description"`
	InstanceOfferingUuid types.String `tfsdk:"instance_offering_uuid"`
	Strategy             types.String `tfsdk:"strategy"`
	MemorySize           types.Int64  `tfsdk:"memory_size"`
	CPUNum               types.Int64  `tfsdk:"cpu_num"`
	NeverStop            types.Bool   `tfsdk:"never_stop"`
	Marketplace          types.Bool   `tfsdk:"marketplace"`
	GPUDevices           types.List   `tfsdk:"gpu_devices"`
	GPUSpecs             types.Object `tfsdk:"gpu_device_specs"`
	UserData             types.String `tfsdk:"user_data"`
	VMNics               types.List   `tfsdk:"vm_nics"`
	Expunge              types.Bool   `tfsdk:"expunge"`
}

type NicsModel struct {
	Uuid    types.String `tfsdk:"uuid"`
	Ip      types.String `tfsdk:"ip"`
	Netmask types.String `tfsdk:"netmask"`
	Gateway types.String `tfsdk:"gateway"`
}

func InstanceResource() resource.Resource {
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
			fmt.Sprintf("Expected *client.ZSClient, got: %T. Please report this issue to the Provider developer. jiajian.chi@zstack.io", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata implements resource.Resource.
func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

// Schema implements resource.Resource.
func (r *vmResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to manage virtual machine (VM) instances in ZStack. " +
			"A VM instance represents a virtualized compute resource that can be created, updated, and deleted. " +
			"You can define the VM's properties, such as its name, image, network configuration, disks, and GPU devices.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the VM instance.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VM instance.",
			},
			"vm_nics": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Required:    true,
							Description: "The UUID of the network.",
						},
						"ip": schema.StringAttribute{
							Computed:    true,
							Description: "The IP address assigned to the network.",
						},
						"netmask": schema.StringAttribute{
							Computed:    true,
							Description: "The netmask of the network.",
						},
						"gateway": schema.StringAttribute{
							Computed:    true,
							Description: "The gateway of the network.",
						},
					},
				},
				Computed:    true,
				Description: "The IP address assigned to the VM instance.",
			},
			"instance_offering_uuid": schema.StringAttribute{
				Optional: true,
				Description: "The UUID of the instance offering used by the VM. Required if using instance offering uuid to create instances. " +
					"  Mutually exclusive with `cpu_num` and `memory_size`.",
			},
			"image_uuid": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the image used to create the VM instance.",
			},
			"l3_network_uuids": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "A list of UUIDs for the L3 networks associated with the VM instance.",
			},
			"networks": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Required:    true,
							Description: "The UUID of the network.",
						},
					},
				},
				Optional:    true,
				Description: "The network configurations associated with the VM instance.",
			},
			"root_disk": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"offering_uuid": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the disk offering for the root disk.",
					},
					"size": schema.Int64Attribute{
						Optional:    true,
						Description: "The size of the root disk in gigabytes (GB).",
					},
					"primary_storage_uuid": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the primary storage for the root disk.",
					},
					"ceph_pool_name": schema.StringAttribute{
						Optional:    true,
						Description: "The Ceph pool name for the root disk.",
					},
					"virtio_scsi": schema.BoolAttribute{
						Optional:    true,
						Description: "Whether the root disk uses Virtio-SCSI.",
					},
				},
				Optional:    true,
				Description: "The configuration for the root disk of the VM instance.",
			},
			"data_disks": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"offering_uuid": schema.StringAttribute{
							Optional:    true,
							Description: "The UUID of the disk offering for the data disk.",
						},
						"size": schema.Int64Attribute{
							Optional:    true,
							Description: "The size of the data disk in gigabytes (GB).",
						},

						"primary_storage_uuid": schema.StringAttribute{
							Computed:    true,
							Description: "The UUID of the primary storage for the data disk.",
						},

						"ceph_pool_name": schema.StringAttribute{
							Optional:    true,
							Description: "The Ceph pool name for the data disk.",
						},
						"virtio_scsi": schema.BoolAttribute{
							Optional:    true,
							Description: "Whether the data disk uses Virtio-SCSI.",
						},
					},
				},
				Optional:    true,
				Description: "The configuration for additional data disks.",
			},
			"gpu_device_specs": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Optional:    true,
						Description: "The UUID of the GPU device.",
					},
					"type": schema.StringAttribute{
						Optional:    true,
						Description: "The type of the GPU device. Must be one of: `mdevDevice` or `pciDevice`.",
						Validators: []validator.String{
							stringvalidator.OneOf("mdevDevice", "pciDevice"),
						},
					},
					"number": schema.Int64Attribute{
						Optional:    true,
						Description: "The number of GPU devices assigned.",
					},
				},
				Optional:    true,
				Description: "The GPU specifications for the VM instance.",
			},
			"gpu_devices": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							Optional:    true,
							Description: "The UUID of the GPU device.",
						},
						"type": schema.StringAttribute{
							Optional:    true,
							Description: "The type of the GPU device.  Must be one of: `mdevDevice` or `pciDevice`.",
							Validators: []validator.String{
								stringvalidator.OneOf("mdevDevice", "pciDevice"),
							},
						},
					},
				},
				Optional:    true,
				Description: "A list of GPU devices assigned to the VM instance.",
			},
			"zone_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the zone where the VM instance is deployed.",
			},
			"cluster_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the cluster where the VM instance is deployed.",
			},
			"host_uuid": schema.StringAttribute{
				Optional:    true,
				Description: "The UUID of the host where the VM instance is running.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the VM instance.",
			},
			"memory_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The memory size allocated to the VM instance in megabytes (MB). When used together with `cpu_num`, the `instance_offering_uuid` is not required.",
			},
			"cpu_num": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of CPUs allocated to the VM instance.  When used together with `memory_size`, the `instance_offering_uuid` is not required.",
			},
			"strategy": schema.StringAttribute{
				Optional:    true,
				Description: "The deployment strategy for the VM instance.",
			},
			"user_data": schema.StringAttribute{
				Optional:    true,
				Description: "User data injected into the VM instance at boot time.",
			},
			"never_stop": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the VM instance should never stop automatically.",
			},
			"expunge": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates if the instance should be expunged after deletion.",
			},
			"marketplace": schema.BoolAttribute{
				Optional:    true,
				Description: "Indicates whether the VM instance is a marketplace instance.",
			},
		},
	}

}

// Create implements resource.Resource.

func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmInstanceDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rootDiskPlan diskModel
	var dataDisksPlan []diskModel
	var l3NetworkUuids []string
	var systemTags []string
	var primaryStorageUuidForRootVolume *string
	hostUuid := ""
	clusterUuid := ""
	zoneUuid := ""
	var rootDiskSystemTags []string
	var dataDiskSizes []int64
	var dataDiskOfferingUuids []string
	var dataVolumeSystemTagsOnIndex []string
	var dataDiskSystemTags []string

	// SET ROOT DISK
	if !plan.RootDisk.IsNull() {
		diags = plan.RootDisk.As(ctx, &rootDiskPlan, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if rootDiskPlan.OfferingUuid.IsNull() && rootDiskPlan.Size.IsNull() {
			resp.Diagnostics.AddError(
				"Params Error",
				"rootDiskPlan offering_uuid and size cannot be null at the same time",
			)
			return
		}
		err := isDiskParamValid(r, rootDiskPlan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Params Error",
				fmt.Sprintf("invalid rootDiskPlan param, err: %v", err),
			)
			return
		}
		if !rootDiskPlan.PrimaryStorageUuid.IsNull() && rootDiskPlan.PrimaryStorageUuid.ValueString() != "" {
			primaryStorageUuidForRootVolume = rootDiskPlan.PrimaryStorageUuid.ValueStringPointer()
		}

		if !rootDiskPlan.CephPoolName.IsNull() && rootDiskPlan.CephPoolName.ValueString() != "" {
			rootDiskSystemTags = append(rootDiskSystemTags, fmt.Sprintf("ceph::rootPoolName::%s", rootDiskPlan.CephPoolName.ValueString()))
		}

		if !rootDiskPlan.Size.IsNull() {
			rootDiskPlan.Size = types.Int64Value(utils.GBToBytes(rootDiskPlan.Size.ValueInt64()))
		}
	}

	// SET DATA DISK
	if !plan.DataDisks.IsNull() {
		plan.DataDisks.ElementsAs(ctx, &dataDisksPlan, false)

		for _, disk := range dataDisksPlan {
			if !disk.OfferingUuid.IsNull() {
				dataDiskOfferingUuids = append(dataDiskOfferingUuids, disk.OfferingUuid.ValueString())
			} else if !disk.Size.IsNull() {
				//	dataDiskSizes = append(dataDiskSizes, disk.Size.ValueInt64())
				dataDiskSizes = append(dataDiskSizes, utils.GBToBytes(disk.Size.ValueInt64()))
				if disk.VirtioSCSI.ValueBool() {
					dataVolumeSystemTagsOnIndex = append(dataVolumeSystemTagsOnIndex, "capability::virtio-scsi")
				}
			} else {
				resp.Diagnostics.AddError(
					"Params Error",
					"dataDisk offering_uuid and size cannot be null at the same time",
				)
				return
			}
		}

		//only support one type data disk now

		if len(dataDisksPlan) > 0 {
			err := isDiskParamValid(r, dataDisksPlan[0])
			if err != nil {
				resp.Diagnostics.AddError(
					"Params Error",
					fmt.Sprintf("invalid dataDisk param, err: %v", err),
				)
				return
			}
			/*
				if !dataDisksPlan[0].PrimaryStorageUuid.IsNull() && dataDisksPlan[0].PrimaryStorageUuid.ValueString() != "" {
					systemTags = append(systemTags, fmt.Sprintf("primaryStorageUuidForDataVolume::%s", dataDisksPlan[0].PrimaryStorageUuid.ValueString()))
				}
			*/
			if !dataDisksPlan[0].CephPoolName.IsNull() && dataDisksPlan[0].CephPoolName.ValueString() != "" {
				dataDiskSystemTags = append(dataDiskSystemTags, fmt.Sprintf("ceph::pool::%s", dataDisksPlan[0].CephPoolName.ValueString()))
			}
			dataDiskSystemTags = append(dataDiskSystemTags, dataVolumeSystemTagsOnIndex...)
		}
	}

	// SET NETWORK

	if !plan.L3NetworkUuids.IsNull() && len(plan.L3NetworkUuids.Elements()) > 0 {
		plan.L3NetworkUuids.ElementsAs(ctx, &l3NetworkUuids, false)
	} else if !plan.Networks.IsNull() && len(plan.Networks.Elements()) > 0 {
		var networks []networkModel
		plan.Networks.ElementsAs(ctx, &networks, false)
		for _, network := range networks {
			l3NetworkUuids = append(l3NetworkUuids, network.Uuid.ValueString())
		}
	} else {
		resp.Diagnostics.AddError(
			"Params Error",
			"l3NetworkUuids or networks cannot be null at the same time",
		)
		return
	}

	// SET IMAGE
	image, err := r.client.GetImage(plan.ImageUuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Params Error",
			fmt.Sprintf("failed to find image %s, err: %v", plan.ImageUuid.ValueString(), err),
		)
		return
	}

	if image.Status != "Ready" {
		resp.Diagnostics.AddError(
			"Params Error",
			fmt.Sprintf("image %s Status is %s, not Ready", plan.ImageUuid.ValueString(), image.State),
		)
		return
	}

	if image.State != "Enabled" {
		resp.Diagnostics.AddError(
			"Params Error",
			fmt.Sprintf("image %s State is %s, not Enabled", plan.ImageUuid.ValueString(), image.State),
		)
		return
	}

	// SET HOST UUID
	if !plan.HostUuid.IsNull() && plan.HostUuid.ValueString() != "" {
		hostUuid = plan.HostUuid.ValueString()
	}

	// SET CLUSTER UUID
	if !plan.ClusterUuid.IsNull() && plan.HostUuid.ValueString() != "" {
		clusterUuid = plan.ClusterUuid.ValueString()
	}

	// SET CLUSTER UUID
	if !plan.ZoneUuid.IsNull() && plan.ZoneUuid.ValueString() != "" {
		zoneUuid = plan.ZoneUuid.ValueString()
	}

	// SET SYSTEM TAG
	//systemTags := []string{"resourceConfig::vm::vm.clock.track::guest", "cdroms::Empty::None::None"}
	if plan.Marketplace.ValueBool() {
		systemTags = append(systemTags, "marketplace::true")
	}
	if !plan.NeverStop.IsNull() && plan.NeverStop.ValueBool() {
		systemTags = append(systemTags, "ha::NeverStop")
	}

	if !plan.UserData.IsNull() && plan.UserData.ValueString() != "" {
		systemTags = append(systemTags, fmt.Sprintf("userdata::%s", plan.UserData.ValueString()))
	}

	// SET GPU
	// when gpu device is set, gpu spec is not work
	if !plan.GPUDevices.IsNull() {
		var gpuDevicesPlan []gpuModel
		plan.GPUDevices.ElementsAs(ctx, &gpuDevicesPlan, false)

		for _, gpuDevice := range gpuDevicesPlan {
			if gpuDevice.Type.ValueString() == string(mdevDevice) {
				systemTags = append(systemTags, fmt.Sprintf("mdevDevice::%s", gpuDevice.Uuid.ValueString()))
			} else if gpuDevice.Type.ValueString() == string(pciDevice) {
				systemTags = append(systemTags, fmt.Sprintf("pciDevice::%s", gpuDevice.Uuid.ValueString()))
			} else {
				resp.Diagnostics.AddError(
					"Params Error",
					fmt.Sprintf("gpu type %s is invalid", gpuDevice.Type.ValueString()),
				)
				return
			}
		}
	} else if !plan.GPUSpecs.IsNull() {
		var gpuSpecPlan gpuSpecModel
		diags = plan.GPUSpecs.As(ctx, &gpuSpecPlan, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		number := 1
		if !gpuSpecPlan.Number.IsNull() {
			number = int(gpuSpecPlan.Number.ValueInt64())
		}

		if gpuSpecPlan.Type.ValueString() == string(mdevDevice) {
			systemTags = append(systemTags, fmt.Sprintf("mdevDeviceSpec::%s::%d", gpuSpecPlan.Uuid.ValueString(), number))
		} else if gpuSpecPlan.Type.ValueString() == string(pciDevice) {
			systemTags = append(systemTags, fmt.Sprintf("pciDeviceSpec::%s::%d", gpuSpecPlan.Uuid.ValueString(), number))
		} else {
			resp.Diagnostics.AddError(
				"Params Error",
				fmt.Sprintf("gpu type %s is invalid", gpuSpecPlan.Type.ValueString()),
			)
			return
		}
	}

	//SET OTHER PARAM
	if !plan.Strategy.IsNull() {
		strategyValue := plan.Strategy.ValueString()
		if strategyValue != string(param.InstantStart) && strategyValue != string(param.CreateStopped) {
			resp.Diagnostics.AddError(
				"Params Error",
				fmt.Sprintf("strategy %s is invalid, valid value is InstantStart or CreateStopped", plan.Strategy.ValueString()),
			)
			return
		}
	}

	// Check if instance_offering_uuid is provided
	var memorySize int64
	var cpuNum int64
	if !plan.InstanceOfferingUuid.IsNull() && plan.InstanceOfferingUuid.ValueString() != "" {
		instanceOffering, err := r.client.GetInstanceOffering(plan.InstanceOfferingUuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Params Error",
				fmt.Sprintf("failed to get instance offering %s, err: %v", plan.InstanceOfferingUuid.ValueString(), err),
			)
			return
		}
		memorySize = instanceOffering.MemorySize
		cpuNum = int64(instanceOffering.CpuNum)
	} else {
		if plan.MemorySize.IsNull() || plan.CPUNum.IsNull() {
			resp.Diagnostics.AddError(
				"Params Error",
				"memory_size and cpu_num must be provided if instance_offering_uuid is not set",
			)
			return
		}

		memorySize = utils.MBToBytes(plan.MemorySize.ValueInt64())
		cpuNum = plan.CPUNum.ValueInt64()
	}

	createVmInstanceParam := param.CreateVmInstanceParam{
		BaseParam: param.BaseParam{
			SystemTags: systemTags,
			UserTags:   nil,
			RequestIp:  "",
		},
		Params: param.CreateVmInstanceDetailParam{
			Name:                            plan.Name.ValueString(),
			InstanceOfferingUUID:            plan.InstanceOfferingUuid.ValueString(),
			ImageUUID:                       plan.ImageUuid.ValueString(),
			L3NetworkUuids:                  l3NetworkUuids,
			Type:                            param.UserVm,
			RootDiskOfferingUuid:            rootDiskPlan.OfferingUuid.ValueString(),
			RootDiskSize:                    rootDiskPlan.Size.ValueInt64Pointer(),
			PrimaryStorageUuidForRootVolume: primaryStorageUuidForRootVolume,
			DataDiskSizes:                   dataDiskSizes,
			DataDiskOfferingUuids:           dataDiskOfferingUuids,
			ZoneUuid:                        zoneUuid,
			ClusterUUID:                     clusterUuid,
			HostUuid:                        hostUuid,
			Description:                     plan.Description.ValueString(),
			DefaultL3NetworkUuid:            l3NetworkUuids[0],
			TagUuids:                        nil,
			Strategy:                        param.InstanceStrategy(plan.Strategy.ValueString()),
			MemorySize:                      memorySize,
			CpuNum:                          cpuNum,
			RootVolumeSystemTags:            rootDiskSystemTags,
			DataVolumeSystemTags:            dataDiskSystemTags,
		},
	}

	instance, err := r.client.CreateVmInstance(createVmInstanceParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create VmInstance Error",
			fmt.Sprintf("failed to create vminstance, err: %v", err),
		)
		return
	}

	plan.Uuid = types.StringValue(instance.UUID)
	plan.Name = types.StringValue(instance.Name)
	plan.Description = types.StringValue(instance.Description)
	plan.MemorySize = types.Int64Value(utils.BytesToMB(instance.MemorySize))
	plan.CPUNum = types.Int64Value(int64(instance.CPUNum))
	//plan.IP = types.StringValue(instance.VMNics[0].IP)

	var diskModelAttrTypes = map[string]attr.Type{
		"offering_uuid":        types.StringType,
		"size":                 types.Int64Type,
		"primary_storage_uuid": types.StringType,
		"ceph_pool_name":       types.StringType,
		"virtio_scsi":          types.BoolType,
	}

	if !plan.DataDisks.IsNull() {
		var dataDisksPlan []diskModel
		plan.DataDisks.ElementsAs(ctx, &dataDisksPlan, false)

		for i, disk := range instance.AllVolumes {
			if i < len(dataDisksPlan) {
				dataDisksPlan[i].PrimaryStorageUuid = types.StringValue(disk.PrimaryStorageUUID)
			}
		}

		dataDisksList, diags := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: diskModelAttrTypes, // 使用 diskModelAttrTypes
		}, dataDisksPlan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.DataDisks = dataDisksList

	}

	// 处理所有网卡信息
	var vmNics []NicsModel
	for _, nic := range instance.VMNics {
		vmNics = append(vmNics, NicsModel{
			Uuid:    types.StringValue(nic.UUID),
			Ip:      types.StringValue(nic.IP),
			Netmask: types.StringValue(nic.Netmask),
			Gateway: types.StringValue(nic.Gateway),
		})
	}

	plan.VMNics, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Read implements resource.Resource.
func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmInstanceDataSourceModel
	req.State.Schema.GetAttributes()

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.client.GetVmInstance(state.Uuid.ValueString())
	if err != nil {
		tflog.Warn(ctx, "cannot read vm, maybe it has been deleted, set uuid to 'empty'. vm was no longer managed by terraform. error: "+err.Error())
		state.Uuid = types.StringValue("")
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	state.Uuid = types.StringValue(vm.UUID)
	state.Name = types.StringValue(vm.Name)
	state.Description = types.StringValue(vm.Description)
	state.ImageUuid = types.StringValue(vm.ImageUUID)
	state.MemorySize = types.Int64Value(utils.BytesToMB(vm.MemorySize))
	state.CPUNum = types.Int64Value(int64(vm.CPUNum))

	var vmNics []NicsModel
	for _, nic := range vm.VMNics {
		vmNics = append(vmNics, NicsModel{
			Uuid:    types.StringValue(nic.UUID),
			Ip:      types.StringValue(nic.IP),
			Netmask: types.StringValue(nic.Netmask),
			Gateway: types.StringValue(nic.Gateway),
		})
	}

	state.VMNics, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)

	state.L3NetworkUuids, _ = types.ListValue(types.StringType, []attr.Value{types.StringValue(vm.DefaultL3NetworkUUID)})
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state vmInstanceDataSourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Uuid.ValueString() == "" {
		resp.Diagnostics.AddError("Parameter Error",
			"uuid of vm is empty, cannot upgrade vm.")
		return
	}

	//"uuid of vm is empty, cannot upgrade vm."

	uuid := state.Uuid.ValueString()

	updateVmInstanceParam := param.UpdateVmInstanceParam{}
	updateVm := false

	if plan.Name.ValueString() != state.Name.ValueString() {
		updateVmInstanceParam.UpdateVmInstance.Name = plan.Name.ValueString()
		updateVm = true
	}
	if plan.Description.ValueString() != state.Description.ValueString() {
		updateVmInstanceParam.UpdateVmInstance.Description = plan.Description.ValueStringPointer()
		updateVm = true

	}
	if plan.CPUNum.ValueInt64() != state.CPUNum.ValueInt64() {
		updateVmInstanceParam.UpdateVmInstance.CpuNum = utils.TfInt64ToIntPointer(plan.CPUNum)
		updateVm = true
	}
	if plan.MemorySize.ValueInt64() != state.MemorySize.ValueInt64() {
		memorySizeBytes := utils.MBToBytes(plan.MemorySize.ValueInt64())
		updateVmInstanceParam.UpdateVmInstance.MemorySize = &memorySizeBytes
		updateVm = true
	}

	if updateVm {
		instance, err := r.client.UpdateVmInstance(uuid, updateVmInstanceParam)
		if err != nil {
			resp.Diagnostics.AddError(
				"Update VmInstance Error",
				"failed to update vm instance, err:"+err.Error())
			return
		}

		plan.Uuid = types.StringValue(instance.UUID)
		plan.Name = types.StringValue(instance.Name)
		plan.Description = types.StringValue(instance.Description)
		plan.MemorySize = types.Int64Value(utils.BytesToMB(instance.MemorySize))
		plan.CPUNum = types.Int64Value(int64(instance.CPUNum))
		//plan.IP = types.StringValue(instance.VMNics[0].IP)

		// 更新 vm_nics 信息
		var vmNics []NicsModel
		for _, nic := range instance.VMNics {
			vmNics = append(vmNics, NicsModel{
				Uuid:    types.StringValue(nic.UUID),
				Ip:      types.StringValue(nic.IP),
				Netmask: types.StringValue(nic.Netmask),
				Gateway: types.StringValue(nic.Gateway),
			})
		}

		// 将 vm_nics 信息设置到状态中
		plan.VMNics, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: networkModelAttrTypes}, vmNics)

		diags := resp.State.Set(ctx, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	//return
}

// Delete implements resource.Resource.
func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmInstanceDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if state.Uuid == types.StringValue("") {
		tflog.Warn(ctx, "vm uuid is empty, so nothing to delete, skip it")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO: query vm instance again in delete function is not smart. Update vm instance's data disk state in read function is a better way
	vm, err := r.client.GetVmInstance(state.Uuid.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not read vm instance", "Error: "+err.Error(),
		)
		return
	}

	var volumeUuids []string
	for _, volume := range vm.AllVolumes {
		if volume.Type != "Data" {
			continue
		}
		volumeUuids = append(volumeUuids, volume.UUID)
	}

	tflog.Info(ctx, "Deleting vm instance "+state.Uuid.String())

	//Delete existing vm instance
	err = r.client.DestroyVmInstance(state.Uuid.ValueString(), param.DeleteModePermissive)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not destroy vm instance", "Error: "+err.Error(),
		)
		return
	}

	//Delete vm data volume
	for _, uuid := range volumeUuids {
		err = r.client.DeleteDataVolume(uuid, param.DeleteModePermissive)
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not delete data volume", "Error: "+err.Error(),
			)
			return
		}
	}

	expunge := false
	if !state.Expunge.IsNull() && !state.Expunge.IsUnknown() {
		expunge = state.Expunge.ValueBool()
	}

	if expunge {
		tflog.Info(ctx, fmt.Sprintf("expunge instance %s", state.Uuid.ValueString()))
		//Expunge vm instance
		err = r.client.ExpungeVmInstance(state.Uuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not expunge vm instance", "Error: "+err.Error(),
			)
			return
		}

		//Expunge vm data volume
		for _, uuid := range volumeUuids {
			err = r.client.ExpungeDataVolume(uuid)
			if err != nil {
				resp.Diagnostics.AddError(
					"Could not expunge data volume", "Error: "+err.Error(),
				)
				return
			}
		}
	}

}

func isDiskParamValid(r *vmResource, model diskModel) error {
	if model.PrimaryStorageUuid.IsNull() || model.PrimaryStorageUuid.ValueString() == "" {
		return nil
	}

	dataDiskPrimaryStorageUuid := model.PrimaryStorageUuid.ValueString()
	dataDiskCephPoolName := model.CephPoolName.ValueString()

	qparam := param.NewQueryParam()
	qparam.AddQ("uuid=" + dataDiskPrimaryStorageUuid)
	qparam.AddQ("state=Enabled")
	qparam.Limit(1)
	primaryStorages, err := r.client.QueryPrimaryStorage(qparam)
	if err != nil {
		return fmt.Errorf("failed to get primary storage %s, err: %v", dataDiskPrimaryStorageUuid, err)
	}

	if len(primaryStorages) == 0 {
		return fmt.Errorf("unable to find primary storage %s, err: %v", dataDiskPrimaryStorageUuid, err)
	}

	if dataDiskCephPoolName != "" {
		found := false
		for _, pool := range primaryStorages[0].Pools {
			if pool.PoolName == dataDiskCephPoolName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unable to find pool name %s", dataDiskCephPoolName)
		}
	}
	return nil
}
