package view

type VmInstanceInventoryView struct {
	BaseInfoView
	BaseTimeView

	ZoneUUID             string               `json:"zoneUuid"`             //区域UUID
	ClusterUUID          string               `json:"clusterUuid"`          //集群UUID
	ImageUUID            string               `json:"imageUuid"`            //镜像UUID
	HostUUID             string               `json:"hostUuid"`             //物理机UUID
	LastHostUUID         string               `json:"lastHostUuid"`         //上一次运行云主机的物理机UUID
	InstanceOfferingUUID string               `json:"instanceOfferingUuid"` //计算规格UUID
	RootVolumeUUID       string               `json:"rootVolumeUuid"`       //根云盘UUID
	Platform             string               `json:"platform"`             //云主机运行平台
	Architecture         string               `json:"architecture"`         //架构类型
	GuestOsType          string               `json:"guestOsType" `         //镜像对应客户机操作系统的类型
	DefaultL3NetworkUUID string               `json:"defaultL3NetworkUuid"` //默认三层网络UUID
	Type                 string               `json:"type"`                 //云主机类型
	HypervisorType       string               `json:"hypervisorType"`       //云主机的hypervisor类型
	MemorySize           int64                `json:"memorySize"`           //内存大小
	CPUNum               int                  `json:"cpuNum"`               //cpu数量
	CPUSpeed             int64                `json:"cpuSpeed"`             //cpu主频
	AllocatorStrategy    string               `json:"allocatorStrategy"`    //分配策略
	State                string               `json:"state"`                //云主机的可用状态
	VMNics               []VmNicInventoryView `json:"vmNics"`               //所有网卡信息
	AllVolumes           []VolumeView         `json:"allVolumes"`           //所有卷
	VmCdRoms             []VmCdRom            `json:"vmCdRoms"`             //驱动
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

	DeviceId       int    `json:"deviceId"`       //
	VmInstanceUuid string `json:"vmInstanceUuid"` //
}

type VMConsoleAddressView struct {
	HostIp      string      `json:"hostIp" bson:"hostIp"`           //云主机运行物理机IP
	Port        string      `json:"port" bson:"port"`               //云主机控制台端口
	Protocol    string      `json:"protocol" bson:"protocol"`       //云主机控制台协议，vnc或spice或vncAndSpice
	Success     bool        `json:"success" bson:"success"`         //操作是否成功
	VdiPortInfo VdiPortInfo `json:"vdiPortInfo" bson:"vdiPortInfo"` //端口组
}
type VdiPortInfo struct {
	VncPort      int `json:"vncPort" bson:"vncPort"`           //	vnc端口号
	SpicePort    int `json:"spicePort" bson:"spicePort"`       //spice端口号
	SpiceTlsPort int `json:"spiceTlsPort" bson:"spiceTlsPort"` //spice开启Tls加密，会存在spiceTlsPort和spicePort两个端口号，通过spice客户端连接云主机需要使用spiceTlsPort端口号
}

type GetVmConsolePasswordView struct {
	ConsolePassword string `json:"consolePassword" bson:"consolePassword"` //密码
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
	HypervisorType     string      `json:"hypervisorType" ` //虚拟化类型
	Version            interface{} `json:"version" `        //版本
	Architecture       string      `json:"architecture" `   //架构
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
