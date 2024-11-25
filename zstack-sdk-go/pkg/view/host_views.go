// Copyright (c) ZStack.io, Inc.

package view

type HostInventoryView struct {
	BaseInfoView
	BaseTimeView
	Architecture            string `json:"architecture"` // Physical machine architecture
	ZoneUuid                string `json:"zoneUuid"`     // Zone UUID
	ClusterUuid             string `json:"clusterUuid"`  // Cluster UUID
	ManagementIp            string `json:"managementIp"` // Management IP
	HypervisorType          string `json:"hypervisorType"`
	State                   string `json:"state"`  // Physical machine state, including Enabled, Disabled, PreMaintenance, Maintenance
	Status                  string `json:"status"` // Connecting, Connected, Disconnected
	TotalCpuCapacity        int64  `json:"totalCpuCapacity"`
	AvailableCpuCapacity    int64  `json:"availableCpuCapacity"`
	CpuSockets              int    `json:"cpuSockets"`
	TotalMemoryCapacity     int64  `json:"totalMemoryCapacity"`
	AvailableMemoryCapacity int64  `json:"availableMemoryCapacity"`
	CpuNum                  int    `json:"cpuNum"`
	Username                string `json:"username"`
	SshPort                 int    `json:"sshPort"`
}

type HostNetworkBondingInventoryView struct {
	BaseInfoView
	BaseTimeView
	HostUuid    string `json:"hostUuid"`    // Physical machine UUID
	BondingName string `json:"bondingName"` // Bond name

	Mode            string                              `json:"mode"`            // Bond mode
	XmitHashPolicy  string                              `json:"xmitHashPolicy"`  // Hash policy
	MiiStatus       string                              `json:"miiStatus"`       // MII status
	Mac             string                              `json:"mac"`             // MAC address
	IpAddresses     []string                            `json:"ipAddresses"`     // IP addresses
	Miimon          int64                               `json:"miimon"`          // MII monitoring interval
	AllSlavesActive bool                                `json:"allSlavesActive"` // Whether all slaves are active
	Slaves          []HostNetworkInterfaceInventoryView `json:"slaves"`          // Slaves list
}

type HostNetworkInterfaceInventoryView struct {
	BaseTimeView
	UUID             string   `json:"uuid"`             // Network interface UUID
	HostUuid         string   `json:"hostUuid"`         // Physical machine UUID
	BondingUuid      string   `json:"bondingUuid"`      // Bond UUID
	InterfaceName    string   `json:"interfaceName"`    // Network interface name
	InterfaceType    string   `json:"interfaceType"`    // Network interface application status, including nomaster, bridgeSlave, bondSlave
	Speed            int64    `json:"speed"`            // Network interface speed
	SlaveActive      bool     `json:"slaveActive"`      // Bond link status
	CarrierActive    bool     `json:"carrierActive"`    // Physical link status
	IpAddresses      []string `json:"ipAddresses"`      // IP addresses
	Mac              string   `json:"mac"`              // MAC address
	PciDeviceAddress string   `json:"pciDeviceAddress"` // PCI address
}
