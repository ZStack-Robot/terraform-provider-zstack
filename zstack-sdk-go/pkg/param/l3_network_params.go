// Copyright (c) ZStack.io, Inc.

package param

type L3Category string

const (
	Public  L3Category = "Public"
	Private L3Category = "Private"
	System  L3Category = "System"
)

// QueryL3NetworkRequest queries a layer 3 network
type QueryL3NetworkRequest struct {
	UUID string `json:"uuid"` // Resource UUID, uniquely identifies the resource
}

type UpdateL3NetworkParam struct {
	BaseParam
	UpdateL3Network UpdateL3NetworkDetailParam `json:"updateL3Network"`
}

type UpdateL3NetworkDetailParam struct {
	BaseParam
	Name        string      `json:"name"`        // Layer 3 network name
	Description *string     `json:"description"` // Layer 3 network description
	System      *bool       `json:"system"`      // Whether it is for system cloud hosts
	DnsDomain   *string     `json:"dnsDomain"`   // DNS domain of the layer 3 network
	Category    *L3Category `json:"category"`    // Category of the layer 3 network
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
	Name        string `json:"name"`
	StartIp     string `json:"startIp"`
	EndIp       string `json:"endIp"`
	Gateway     string `json:"gateway"`
	PrefixLen   int    `json:"prefixLen"`
	AddressMode string `json:"addressMode"` // SLAAC, Stateful-DHCP, Stateless-DHCP
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
	AddressMode string `json:"addressMode"` // SLAAC, Stateful-DHCP, Stateless-DHCP
}

type CreateL3NetworkParam struct {
	BaseParam
	Params CreateL3NetworkDetailParam `json:"params"`
}

type CreateL3NetworkDetailParam struct {
	Name          string `json:"name"`
	Description   string `json:"description"` // Layer 3 network description
	Type          string `json:"type"`
	L2NetworkUuid string `json:"l2NetworkUuid"`
	Category      string `json:"category"`
	System        bool   `json:"system"`
	EnableIPAM    bool   `json:"enableIPAM"`
}

type AddReservedIpRangeParam struct {
	BaseParam
	Params AddReservedIpRangeDetailParam `json:"params"`
}

type AddReservedIpRangeDetailParam struct {
	StartIp string `json:"startIp"`
	EndIp   string `json:"endIp"`
}
