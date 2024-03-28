// Copyright (c) HashiCorp, Inc.

package param

type L3Category string

const (
	Public  L3Category = "Public"
	Private L3Category = "Private"
	System  L3Category = "System"
)

// QueryL3NetworkRequest 查询三层网络
type QueryL3NetworkRequest struct {
	UUID string `json:"uuid"` //资源的UUID，唯一标示该资源
}

type UpdateL3NetworkParam struct {
	BaseParam
	UpdateL3Network UpdateL3NetworkDetailParam `json:"updateL3Network"`
}

type UpdateL3NetworkDetailParam struct {
	BaseParam
	Name        string      `json:"name"`        //三层网络名称
	Description *string     `json:"description"` //三层网络描述
	System      *bool       `json:"system"`      //是否用于系统云主机
	DnsDomain   *string     `json:"dnsDomain"`   //三层网络的DNS域名
	Category    *L3Category `json:"category"`    //三层网络的分类
}

type AddDnsToL3NetworkParam struct {
	BaseParam
	Params AddDnsToL3NetworkDetailParam `json:"params"`
}
type AddDnsToL3NetworkDetailParam struct {
	Dns string `json:"dns"`
}

type AddIpRangeParam struct {
	BaseParam
	Params AddIpRangeDetailParam `json:"params"`
}

type AddIpRangeDetailParam struct {
	Name        string `json:"name"`
	StartIp     string `json:"startIp"`
	EndIp       string `json:"endIp"`
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
	IpRangeType string `json:"ipRangeType"`
}

type AddIpv6RangeParam struct {
	BaseParam
	Params AddIpv6RangeDetailParam `json:"params"`
}
type AddIpv6RangeDetailParam struct {
	Name    string `json:"name"`
	StartIp string `json:"startIp"`
	EndIp   string `json:"endIp"`

	Gateway   string `json:"gateway"`
	PrefixLen int    `json:"prefixLen"`

	AddressMode string `json:"addressMode"` // SLAAC	Stateful-DHCP	Stateless-DHCP
}

type AddIpRangeByNetworkCidrParam struct {
	BaseParam
	Params AddIpRangeByNetworkCidrDetailParam `json:"params"`
}

type AddIpRangeByNetworkCidrDetailParam struct {
	Name        string `json:"name"`
	NetworkCidr string `json:"networkCidr"`
	Gateway     string `json:"gateway"`
	IpRangeType string `json:"ipRangeType"`
}

type AddIpv6RangeByNetworkCidrParam struct {
	BaseParam
	Params AddIpv6RangeByNetworkCidrDetailParam `json:"params"`
}
type AddIpv6RangeByNetworkCidrDetailParam struct {
	Name        string `json:"name"`
	NetworkCidr string `json:"networkCidr"`
	AddressMode string `json:"addressMode"` // SLAAC	Stateful-DHCP	Stateless-DHCP
}

type CreateL3NetworkParam struct {
	BaseParam
	Params CreateL3NetworkDetailParam `json:"params"`
}
type CreateL3NetworkDetailParam struct {
	Name          string `json:"name"`
	Description   string `json:"description"` //三层网络描述
	Type          string `json:"type"`
	L2NetworkUuid string `json:"l2NetworkUuid"`
	Category      string `json:"category"`
	System        bool   `json:"system"`
}
