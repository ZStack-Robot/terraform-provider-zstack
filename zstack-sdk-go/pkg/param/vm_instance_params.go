// Copyright (c) ZStack.io, Inc.

package param

type InstanceType string
type InstanceStrategy string
type InstanceStopType string
type HA string
type ClockTrack string

const (
	UserVm      InstanceType = "UserVm"
	ApplianceVm InstanceType = "ApplianceVm"

	InstantStart  InstanceStrategy = "InstantStart"
	CreateStopped InstanceStrategy = "CreateStopped"

	Grace InstanceStopType = "grace" // Graceful shutdown, requires ACPI driver installed in the VM
	Cold  InstanceStopType = "cold"  // Cold shutdown, equivalent to directly cutting power

	NeverStop HA = "NeverStop" // Enable high availability
	None      HA = "None"      // High availability not enabled

	Host  ClockTrack = "host"
	Guest ClockTrack = "guest"
)

type CreateVmInstanceParam struct {
	BaseParam
	Params CreateVmInstanceDetailParam `json:"params" `
}

type CreateVmInstanceDetailParam struct {
	Name                            string       `json:"name" `                            // VM instance name
	InstanceOfferingUUID            string       `json:"instanceOfferingUuid" `            // Instance offering UUID, specifies CPU, memory, etc.
	CpuNum                          int64        `json:"cpuNum"`                           // CPU number
	MemorySize                      int64        `json:"memorySize"`                       // Memory size
	ImageUUID                       string       `json:"imageUuid" `                       // Image UUID, the root volume of the VM will be created from this image.
	L3NetworkUuids                  []string     `json:"l3NetworkUuids" `                  // List of layer 3 network UUIDs, one NIC will be created for each network.
	Type                            InstanceType `json:"type" `                            // VM instance type, reserved field, no need to specify. UserVm/ApplianceVm
	RootDiskOfferingUuid            string       `json:"rootDiskOfferingUuid" `            // Root disk offering UUID, must be specified if the image type is ISO.
	RootDiskSize                    *int64       `json:"rootDiskSize"`                     // Root disk size
	DataDiskOfferingUuids           []string     `json:"dataDiskOfferingUuids" `           // List of data disk offering UUIDs, one or more data disks will be created.
	DataDiskSizes                   []int64      `json:"dataDiskSizes"`                    // Data disk sizes
	ZoneUuid                        string       `json:"zoneUuid" `                        // Zone UUID, if specified, the VM will be created in the specified zone.
	ClusterUUID                     string       `json:"clusterUuid" `                     // Cluster UUID, if specified, the VM will be created in the specified cluster, higher priority than zoneUuid.
	HostUuid                        string       `json:"hostUuid" `                        // Host UUID, if specified, the VM will be created on the specified host, higher priority than zoneUuid and clusterUuid.
	PrimaryStorageUuidForRootVolume *string      `json:"primaryStorageUuidForRootVolume" ` // Primary storage UUID, if specified, the root volume will be created on the specified primary storage.
	Description                     string       `json:"description" `                     // Detailed description of the VM instance
	DefaultL3NetworkUuid            string       `json:"defaultL3NetworkUuid" `            // Default layer 3 network UUID, specifies the default network for routing when multiple networks are specified.
	ResourceUuid                    string       `json:"resourceUuid" `                    // Resource UUID, if specified, the VM will use this value as its UUID.

	TagUuids             []string         `json:"tagUuids" ` // List of tag UUIDs
	Strategy             InstanceStrategy `json:"strategy" ` // VM creation strategy, InstantStart for immediate start, CreateStopped for stopped after creation.
	RootVolumeSystemTags []string         `json:"rootVolumeSystemTags"`
	DataVolumeSystemTags []string         `json:"dataVolumeSystemTags"`
}

type CreateVmFromVolumeParam struct {
	BaseParam
	Params CreateVmFromVolumeDetailParams `json:"params"`
}

type CreateVmFromVolumeDetailParams struct {
	Name                 string   `json:"name"`                 // VM instance name
	Description          string   `json:"description"`          // Detailed description of the resource
	InstanceOfferingUuid string   `json:"instanceOfferingUuid"` // Instance offering UUID, note: this parameter is mutually exclusive with CPU number and memory size.
	CpuNum               int      `json:"cpuNum"`               // CPU number/memory size, note: this parameter is mutually exclusive with instanceOfferingUuid.
	MemorySize           int64    `json:"memorySize"`           // CPU number/memory size, note: this parameter is mutually exclusive with instanceOfferingUuid.
	L3NetworkUuids       []string `json:"l3NetworkUuids"`       // List of layer 3 network UUIDs, one NIC will be created for each network.
	Type                 string   `json:"type"`                 // VM instance type, reserved field, no need to specify.
	VolumeUuid           string   `json:"volumeUuid"`           // Volume UUID
	Platform             string   `json:"platform"`             // Volume system platform
	ZoneUuid             string   `json:"zoneUuid"`             // Zone UUID, if specified, the VM will be created in the specified zone.
	ClusterUuid          string   `json:"clusterUuid"`          // Cluster UUID, if specified, the VM will be created in the specified cluster, higher priority than zoneUuid.
	HostUuid             string   `json:"hostUuid"`             // Host UUID, if specified, the VM will be created on the specified host, higher priority than zoneUuid and clusterUuid.
	PrimaryStorageUuid   string   `json:"primaryStorageUuid"`   // Primary storage UUID, if specified, the root volume will be created on the specified primary storage.
	DefaultL3NetworkUuid string   `json:"defaultL3NetworkUuid"` // Default layer 3 network UUID, specifies the default network for routing when multiple networks are specified.
	Strategy             string   `json:"strategy"`             // VM creation strategy, 1. Start immediately after creation, 2. Do not start after creation.
	ResourceUuid         string   `json:"resourceUuid"`         // Resource UUID, if specified, the VM will use this value as its UUID.
	TagUuids             []string `json:"tagUuids"`             // List of tag UUIDs
}

type CloneVmInstanceParam struct {
	BaseParam
	CloneVmInstance CloneVmInstanceDetailParam `json:"cloneVmInstance"`
}

type CloneVmInstanceDetailParam struct {
	Names                           []string         `json:"names"`    // VM instance names
	Strategy                        InstanceStrategy `json:"strategy"` // Strategy, InstantStart for immediate start after cloning, JustCreate for not starting after cloning.
	Full                            *bool            `json:"full"`     // Whether to clone mounted data disks
	PrimaryStorageUuidForRootVolume *string          `json:"primaryStorageUuidForRootVolume" `
	PrimaryStorageUuidForDataVolume *string          `json:"primaryStorageUuidForDataVolume" `
	RootVolumeSystemTags            []string         `json:"rootVolumeSystemTags" `
	DataVolumeSystemTags            []string         `json:"dataVolumeSystemTags" `
}

type StartVmInstanceParam struct {
	BaseParam
	StartVmInstance StartVmInstanceDetailParam `json:"startVmInstance"` // Can pass hostUuid
}

type StartVmInstanceDetailParam struct {
	HostUuid string `json:"hostUuid"` // Host UUID
}

type StopVmInstanceParam struct {
	BaseParam
	StopVmInstance StopVmInstanceDetailParam `json:"stopVmInstance"` // Requires uuid and type
}

type StopVmInstanceDetailParam struct {
	Type   InstanceStopType `json:"type"`   // Default is grace: graceful shutdown; cold: cold shutdown (power off)
	StopHA bool             `json:"stopHa"` // Completely shut down HA VM
}

type UpdateVmInstanceParam struct {
	BaseParam
	UpdateVmInstance UpdateVmInstanceDetailParam `json:"updateVmInstance"`
}

type UpdateVmInstanceDetailParam struct {
	Name                 string  `json:"name"`        // VM instance name
	Description          *string `json:"description"` // Detailed description of the VM instance
	State                string  `json:"state"`
	DefaultL3NetworkUuid string  `json:"defaultL3NetworkUuid"` // Default layer 3 network UUID, specifies the default network for routing when multiple networks are specified.
	Platform             string  `json:"platform"`             // Volume system platform
	CpuNum               *int    `json:"cpuNum"`               // Number of CPUs
	MemorySize           *int64  `json:"memorySize"`           // CPU number/memory size, note: this parameter is mutually exclusive with instanceOfferingUuid.
	GuestOsType          string  `json:"guestOsType"`
}

type UpdateVmInstanceQgaParam struct {
	BaseParam
	SetVmQga SetVmQgaParam `json:"setVmQga"`
	UUID     string        `json:"UUID"`
}

type SetVmQgaParam struct {
	Enable bool `json:"enable"`
}

type SetVmBootModeParam struct {
	BaseParam
	SetVmBootMode SetVmBootModeDetailParam `json:"setVmBootMode"`
}

type SetVmBootModeDetailParam struct {
	BootMode BootMode `json:"bootMode"` // Boot mode: Legacy, UEFI, UEFI_WITH_CSM
}

type UpdateVmInstanceSshKeyParam struct {
	UUID        string         `json:"uuid"`
	SetVmSshKey SetSshKeyParam `json:"setVmSshKey"`
}
type SetSshKeyParam struct {
	SshKey string `json:"SshKey"`
}

type UpdateVmInstanceChangePwdParam struct {
	UUID             string                `json:"uuid"`
	ChangeVmPassword ChangeVmPasswordParam `json:"changeVmPassword"`
}
type ChangeVmPasswordParam struct {
	Password string `json:"password"`
	Account  string `json:"account"`
}

type UpdateVmInstanceClockTrackParam struct {
	BaseParam
	SetVmClockTrack UpdateVmInstanceClockTrackDetailParam `json:"setVmClockTrack"`
}

type UpdateVmInstanceClockTrackDetailParam struct {
	Track             ClockTrack `json:"track"`             // Clock synchronization method, optional values: guest, host
	SyncAfterVMResume bool       `json:"syncAfterVMResume"` // Whether to synchronize the clock when the VM resumes
	IntervalInSeconds float64    `json:"intervalInSeconds"` // Clock synchronization interval, unit: seconds (0, 60, 600, 1800, 3600, 7200, 21600, 43200, 86400)
}

type UpdateVmCdRomParam struct {
	BaseParam
	UpdateVmCdRom UpdateVmCdRomDetailParam `json:"updateVmCdRom"`
}

type UpdateVmCdRomDetailParam struct {
	Name string `json:"name"`
}

type CreateVmCdRomParam struct {
	BaseParam
	Params CreateVmCdRomDetailParam `json:"params"`
}

type CreateVmCdRomDetailParam struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	VmInstanceUuid string `json:"vmInstanceUuid"`
	IsoUuid        string `json:"isoUuid"`
	ResourceUuid   string `json:"resourceUuid"`
}
