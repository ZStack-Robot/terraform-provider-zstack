// Copyright (c) HashiCorp, Inc.

package view

type L3NetworkInventoryView struct {
	BaseInfoView
	BaseTimeView

	Type            string             `json:"type"`          //三层网络类型
	ZoneUuid        string             `json:"zoneUuid"`      //区域UUID 若指定，云主机会在指定区域创建。
	L2NetworkUuid   string             `json:"l2NetworkUuid"` //二层网络UUID
	State           string             `json:"state"`         //三层网络的可用状态
	DnsDomain       string             `json:"dnsDomain"`     //DNS域
	System          bool               `json:"system"`        //是否用于系统云主机
	Category        string             `json:"category"`      //网络类型，需要与system标签搭配使用，system为true时可设置为Public、Private
	IpVersion       int                `json:"ipVersion"`     //ip协议号:4,6
	Dns             []string           `json:"dns"`
	IpRanges        []IpRangeInventory `json:"ipRanges"`
	NetworkServices []NetworkServices  `json:"networkServices"`
	HostRoute       []HostRoute        `json:"hostRoute"`
}

type IpRangeInventory struct {
	UUID          string `json:"uuid"`          //资源的UUID，唯一标示该资源
	L3NetworkUuid string `json:"l3NetworkUuid"` //三层网络UUID
	Name          string `json:"name"`          //资源名称
	Description   string `json:"description"`   //资源的详细描述
	StartIp       string `json:"StartIp"`
	EndIp         string `json:"EndIp"`
	Netmask       string `json:"netmask"`   //网络掩码
	PrefixLen     string `json:"prefixLen"` //掩码长度
	Gateway       string `json:"gateway"`   //网关地址
	NetworkCidr   string `json:"networkCidr"`
	IpVersion     string `json:"ipVersion"`   //ip协议号:4,6
	AddressMode   string `json:"addressMode"` //IPv6地址分配模式
	CreateDate    string `json:"createDate"`  //创建时间
	LastOpDate    string `json:"lastOpDate"`  //最后一次修改时间
	IpRangeType   string `json:"ipRangeType"`
}

type NetworkServices struct {
	L3NetworkUuid              string `json:"l3NetworkUuid"`              //三层网络UUID
	NetworkServiceProviderUuid string `json:"networkServiceProviderUuid"` //网络服务提供模块UUID
	NetworkServiceType         string `json:"networkServiceType"`
}

type HostRoute struct {
	Id            string `json:"id"`
	L3NetworkUuid string `json:"l3NetworkUuid"` //三层网络UUID
	Prefix        string `json:"prefix"`
	Nexthop       string `json:"nexthop"`
	CreateDate    string `json:"createDate"` //创建时间
	LastOpDate    string `json:"lastOpDate"` //最后一次修改时间
}

type FreeIpInventoryView struct {
	IpRangeUuid string `json:"ipRangeUuid"` //IP段UUID
	Ip          string `json:"ip"`          //ip
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
}

type CheckIpView struct {
	Available bool `json:"available"`
}

type IpAddressCapacityView struct {
	TotalCapacity           int64            `json:"totalCapacity" `           //IP地址容量
	AvailableCapacity       int64            `json:"availableCapacity" `       //可用IP地址容量
	UsedIpAddressNumber     int64            `json:"usedIpAddressNumber" `     //已使用IP数量
	Ipv4TotalCapacity       int64            `json:"ipv4TotalCapacity" `       //IPv4地址容量
	Ipv4AvailableCapacity   int64            `json:"ipv4AvailableCapacity" `   //可用IPv4地址容量
	Ipv4UsedIpAddressNumber int64            `json:"ipv4UsedIpAddressNumber" ` //已使用IPv4数量
	Ipv6TotalCapacity       int64            `json:"ipv6TotalCapacity" `       //IPv6地址容量
	Ipv6AvailableCapacity   int64            `json:"ipv6AvailableCapacity" `   //可用IPv6地址容量
	Ipv6UsedIpAddressNumber int64            `json:"ipv6UsedIpAddressNumber" ` //已使用IPv6数量
	ResourceType            string           `json:"resourceType" `            //所查询资源的类型（地址范围、三层网络、区域）
	Success                 bool             `json:"success" `                 //成功
	CapacityData            []IpCapacityData `json:"capacityData" `
}

type IpCapacityData struct {
	ResourceUuid            string `json:"resourceUuid,omitempty"`   //资源UUID。若指定，镜像会使用该字段值作为UUID。
	TotalCapacity           int64  `json:"totalCapacity" `           //IP地址总容量
	AvailableCapacity       int64  `json:"availableCapacity" `       //可用IP地址容量
	UsedIpAddressNumber     int64  `json:"usedIpAddressNumber" `     //已用IP地址容量
	Ipv4TotalCapacity       int64  `json:"ipv4TotalCapacity"`        //IPv4地址总容量
	Ipv4AvailableCapacity   int64  `json:"ipv4AvailableCapacity" `   //可用IPv4地址容量
	Ipv4UsedIpAddressNumber int64  `json:"ipv4UsedIpAddressNumber" ` //已用IPv4地址容量
	Ipv6TotalCapacity       int64  `json:"ipv6TotalCapacity" `       //IPv6地址总容量
	Ipv6AvailableCapacity   int64  `json:"ipv6AvailableCapacity" `   //	可用IPv6地址容量
	Ipv6UsedIpAddressNumber int64  `json:"ipv6UsedIpAddressNumber" ` //已用IPv6地址容量
}

type DnsInventoryView struct {
	Name          string   `json:"name"`
	L2NetworkUuid string   `json:"l2NetworkUuid"`
	Dns           []string `json:"dns"`
}

type IpRangeInventoryView struct {
	CreateDate    string `json:"createDate"`
	EndIp         string `json:"endIp"`
	Gateway       string `json:"gateway"`
	IpVersion     int    `json:"ipVersion"`
	L3NetworkUuid string `json:"l3NetworkUuid"`
	LastOpDate    string `json:"lastOpDate"`
	Name          string `json:"name"`
	Netmask       string `json:"netmask"`
	NetworkCidr   string `json:"networkCidr"`
	PrefixLen     int    `json:"prefixLen"`
	StartIp       string `json:"startIp"`
	Uuid          string `json:"uuid"`
	AddressMode   string `json:"addressMode"`
}

type DhcpIpAddressView struct {
	Ip  string `json:"ip"`
	Ip6 string `json:"ip6"`
}

type MtuView struct {
	Mtu int64 `json:"mtu"`
}

type IpStatisticView struct {
	Ip             string   `json:"ip"`
	VmInstanceName string   `json:"vmInstanceName"`
	VmInstanceType string   `json:"vmInstanceType"`
	VmInstanceUuid string   `json:"vmInstanceUuid"`
	ResourceTypes  []string `json:"resourceTypes"`
}

type IpAddressInventoryView struct {
	Uuid          string  `json:"uuid"`
	IpRangeUuid   string  `json:"ipRangeUuid"`
	L3NetworkUuid string  `json:"l3NetworkUuid"`
	IpVersion     float64 `json:"ipVersion"`
	Ip            string  `json:"ip"`
	Netmask       string  `json:"netmask"`
	Gateway       string  `json:"gateway"`
	IpInLong      float64 `json:"ipInLong"`
	VmNicUuid     string  `json:"vmNicUuid"`
}

type ReservedIpRangeInventoryView struct {
	Uuid          string `json:"uuid"`
	L3NetworkUuid string `json:"l3NetworkUuid"`
	Name          string `json:"name"`      // max length of 255 characters
	StartIp       string `json:"startIp"`   // start IP address, in IPv4
	EndIp         string `json:"endIp"`     // end IP address, in IPv4
	IpVersion     int    `json:"ipVersion"` // IP version (e.g., 4 for IPv4)
}
