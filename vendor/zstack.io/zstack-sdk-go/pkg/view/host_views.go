// Copyright (c) HashiCorp, Inc.

package view

type HostInventoryView struct {
	BaseInfoView
	BaseTimeView
	Architecture            string `json:"architecture"` //物理机架构
	ZoneUuid                string `json:"zoneUuid"`     //区域UUID
	ClusterUuid             string `json:"clusterUuid"`  //集群UUID
	ManagementIp            string `json:"managementIp"` //
	HypervisorType          string `json:"hypervisorType"`
	State                   string `json:"state"`  //物理机状态，包括Enabled,Disabled,PreMaintenance,Maintenance
	Status                  string `json:"status"` //Connecting,Connected，Disconnected
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
	HostUuid    string `json:"hostUuid"`    //物理机UUID
	BondingName string `json:"bondingName"` //Bond名称

	Mode            string                              `json:"mode"`            //Bond模式
	XmitHashPolicy  string                              `json:"xmitHashPolicy"`  //哈希策略
	MiiStatus       string                              `json:"miiStatus"`       //MII状态
	Mac             string                              `json:"mac"`             //MAC地址
	IpAddresses     []string                            `json:"ipAddresses"`     //IP地址
	Miimon          int64                               `json:"miimon"`          //MII监控间隔
	AllSlavesActive bool                                `json:"allSlavesActive"` //是否所有从机都处于活动状态
	Slaves          []HostNetworkInterfaceInventoryView `json:"slaves"`          //从机列表
}

type HostNetworkInterfaceInventoryView struct {
	BaseTimeView
	UUID             string   `json:"uuid"`             //网络接口UUID
	HostUuid         string   `json:"hostUuid"`         //物理机UUID
	BondingUuid      string   `json:"bondingUuid"`      //Bond UUID
	InterfaceName    string   `json:"interfaceName"`    //网卡名称
	InterfaceType    string   `json:"interfaceType"`    //网卡应用状态，有nomaster、bridgeSlave、bondSlave
	Speed            int64    `json:"speed"`            //网卡速率
	SlaveActive      bool     `json:"slaveActive"`      //Bond链路状态
	CarrierActive    bool     `json:"carrierActive"`    //物理链路状态
	IpAddresses      []string `json:"ipAddresses"`      //IP地址
	Mac              string   `json:"mac"`              //MAC地址
	PciDeviceAddress string   `json:"pciDeviceAddress"` //PCI地址

}
