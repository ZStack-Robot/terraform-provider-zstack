// Copyright (c) HashiCorp, Inc.

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
	ZoneUuid                  string               `json:"zoneUuid,omitempty"`    //区域UUID 若指定，云主机会在指定区域创建。
	ClusterUUID               string               `json:"clusterUuid,omitempty"` //集群UUID 若指定，云主机会在指定集群创建，该字段优先级高于zoneUuid。
	ImageUUID                 string               `json:"imageUuid"`             //镜像UUID 云主机的根云盘会从该字段指定的镜像创建。
	HostUuid                  string               `json:"hostUuid"`              //物理机UUID
	LastHostUUID              string               `json:"lastHostUuid"`          //上一次运行云主机的物理机UUID
	InstanceOfferingUUID      string               `json:"instanceOfferingUuid"`  //计算规格UUID 指定云主机的CPU、内存等参数。
	RootVolumeUuid            string               `json:"rootVolumeUuid"`
	Platform                  string               `json:"platform"`
	DefaultL3NetworkUuid      string               `json:"defaultL3NetworkUuid"`
	Type                      string               `json:"type"`
	HypervisorType            string               `json:"hypervisorType"` //虚拟机管理程序类型,KVM Simulator
	MemorySize                int64                `json:"memorySize"`     //内存大小
	CPUNum                    int                  `json:"cpuNum"`         //cpu数量
	CPUSpeed                  int64                `json:"cpuSpeed"`       //cpu主频
	AllocatorStrategy         string               `json:"allocatorStrategy,omitempty"`
	VMNics                    []VmNicInventoryView `json:"vmNics"`     //所有网卡信息
	AllVolumes                []VolumeView         `json:"allVolumes"` //所有卷
	Dns                       []Dns                `json:"dns"`
}

type Dns struct {
	VpcRouterUuid string `json:"vpcRouterUuid"`
	Dns           string `json:"dns"`
	CreateDate    string `json:"createDate"` //创建时间
	LastOpDate    string `json:"lastOpDate"` //最后一次修改时间
}
