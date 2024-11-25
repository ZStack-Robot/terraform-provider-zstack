// Copyright (c) ZStack.io, Inc.

package view

type VmInstanceInventoryView struct {
	BaseInfoView
	BaseTimeView

	ZoneUUID             string               `json:"zoneUuid"`             // Zone UUID
	ClusterUUID          string               `json:"clusterUuid"`          // Cluster UUID
	ImageUUID            string               `json:"imageUuid"`            // Image UUID
	HostUUID             string               `json:"hostUuid"`             // Physical machine UUID
	LastHostUUID         string               `json:"lastHostUuid"`         // Physical machine UUID where the cloud host last ran
	InstanceOfferingUUID string               `json:"instanceOfferingUuid"` // Compute specification UUID
	RootVolumeUUID       string               `json:"rootVolumeUuid"`       // Root cloud disk UUID
	Platform             string               `json:"platform"`             // Cloud host running platform
	Architecture         string               `json:"architecture"`         // Architecture type
	GuestOsType          string               `json:"guestOsType" `         // Guest OS type corresponding to the image
	DefaultL3NetworkUUID string               `json:"defaultL3NetworkUuid"` // Default layer 3 network UUID
	Type                 string               `json:"type"`                 // Cloud host type
	HypervisorType       string               `json:"hypervisorType"`       // Hypervisor type of the cloud host
	MemorySize           int64                `json:"memorySize"`           // Memory size
	CPUNum               int                  `json:"cpuNum"`               // Number of CPUs
	CPUSpeed             int64                `json:"cpuSpeed"`             // CPU frequency
	AllocatorStrategy    string               `json:"allocatorStrategy"`    // Allocation strategy
	State                string               `json:"state"`                // Availability status of the cloud host
	VMNics               []VmNicInventoryView `json:"vmNics"`               // All NIC information
	AllVolumes           []VolumeView         `json:"allVolumes"`           // All volumes
	VmCdRoms             []VmCdRom            `json:"vmCdRoms"`             // CD-ROMs
}

type CloneVmInstanceResult struct {
	NumberOfClonedVm int                        `json:"numberOfClonedVm"`
	Inventories      []CloneVmInstanceInventory `json:"inventories"`
}

type CloneVmInstanceInventory struct {
	Error     *ErrorCodeView          `json:"error"`
	Inventory VmInstanceInventoryView `json:"inventory"`
}

type VmCdRom struct {
	BaseInfoView
	BaseTimeView

	DeviceId       int    `json:"deviceId"`       // Device ID
	VmInstanceUuid string `json:"vmInstanceUuid"` // VM instance UUID
}

type VMConsoleAddressView struct {
	HostIp      string      `json:"hostIp" bson:"hostIp"`           // IP of the physical machine running the cloud host
	Port        string      `json:"port" bson:"port"`               // Console port of the cloud host
	Protocol    string      `json:"protocol" bson:"protocol"`       // Console protocol of the cloud host, e.g., vnc or spice or vncAndSpice
	Success     bool        `json:"success" bson:"success"`         // Whether the operation was successful
	VdiPortInfo VdiPortInfo `json:"vdiPortInfo" bson:"vdiPortInfo"` // Port group
}

type VdiPortInfo struct {
	VncPort      int `json:"vncPort" bson:"vncPort"`           // VNC port number
	SpicePort    int `json:"spicePort" bson:"spicePort"`       // SPICE port number
	SpiceTlsPort int `json:"spiceTlsPort" bson:"spiceTlsPort"` // SPICE TLS port number, used when SPICE is encrypted with TLS
}

type GetVmConsolePasswordView struct {
	ConsolePassword string `json:"consolePassword" bson:"consolePassword"` // Password
}

type VmGuestToolsInfoView struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

type LatestGuestToolsView struct {
	BaseInfoView
	BaseTimeView

	ManagementNodeUuid string      `json:"managementNodeUuid" `
	AgentType          interface{} `json:"agentType" `
	HypervisorType     string      `json:"hypervisorType" ` // Hypervisor type
	Version            interface{} `json:"version" `        // Version
	Architecture       string      `json:"architecture" `   // Architecture
}

type VMQgaView struct {
	UUID   string `json:"uuid" `
	Enable bool   `json:"enable" `
}

type VMSshKeyView struct {
	SshKey string `json:"sshKey" `
}

type VMCDRomView struct {
	BaseInfoView
	BaseTimeView
	VmInstanceUuid string  `json:"vmInstanceUuid"`
	DeviceId       float64 `json:"deviceId"`
	IsoUuid        string  `json:"isoUuid"`
	IsoInstallPath string  `json:"isoInstallPath"`
}
