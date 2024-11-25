// Copyright (c) ZStack.io, Inc.

package view

type VpcRouterVmInventoryView struct {
	BaseInfoView
	BaseTimeView

	PublicNetworkUuid         string               `json:"publicNetworkUuid"`
	VirtualRouterVips         []string             `json:"virtualRouterVips"`
	ApplianceVmType           string               `json:"applianceVmType"`
	ManagementNetworkUuid     string               `json:"managementNetworkUuid"`
	DefaultRouteL3NetworkUuid string               `json:"defaultRouteL3NetworkUuid"`
	Status                    string               `json:"status"`
	AgentPort                 int                  `json:"agentPort"`
	ZoneUuid                  string               `json:"zoneUuid,omitempty"`    // Zone UUID. If specified, the cloud host will be created in the specified zone.
	ClusterUUID               string               `json:"clusterUuid,omitempty"` // Cluster UUID. If specified, the cloud host will be created in the specified cluster. This field has higher priority than zoneUuid.
	ImageUUID                 string               `json:"imageUuid"`             // Image UUID. The root cloud disk of the cloud host will be created from the image specified by this field.
	HostUuid                  string               `json:"hostUuid"`              // Physical machine UUID
	LastHostUUID              string               `json:"lastHostUuid"`          // UUID of the physical machine where the cloud host last ran
	InstanceOfferingUUID      string               `json:"instanceOfferingUuid"`  // Compute specification UUID. Specifies the CPU, memory, and other parameters of the cloud host.
	RootVolumeUuid            string               `json:"rootVolumeUuid"`
	Platform                  string               `json:"platform"`
	DefaultL3NetworkUuid      string               `json:"defaultL3NetworkUuid"`
	Type                      string               `json:"type"`
	HypervisorType            string               `json:"hypervisorType"` // Hypervisor type, e.g., KVM, Simulator
	MemorySize                int64                `json:"memorySize"`     // Memory size
	CPUNum                    int                  `json:"cpuNum"`         // Number of CPUs
	CPUSpeed                  int64                `json:"cpuSpeed"`       // CPU frequency
	AllocatorStrategy         string               `json:"allocatorStrategy,omitempty"`
	VMNics                    []VmNicInventoryView `json:"vmNics"`     // All NIC information
	AllVolumes                []VolumeView         `json:"allVolumes"` // All volumes
	Dns                       []Dns                `json:"dns"`
}

type Dns struct {
	VpcRouterUuid string `json:"vpcRouterUuid"`
	Dns           string `json:"dns"`
	CreateDate    string `json:"createDate"` // Creation time
	LastOpDate    string `json:"lastOpDate"` // Last modification time
}
