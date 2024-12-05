// Copyright (c) ZStack.io, Inc.

package view

type L3NetworkInventoryView struct {
	BaseInfoView
	BaseTimeView

	Type            string             `json:"type"`          // Layer 3 network type
	ZoneUuid        string             `json:"zoneUuid"`      // Zone UUID. If specified, the cloud host will be created in the specified zone.
	L2NetworkUuid   string             `json:"l2NetworkUuid"` // Layer 2 network UUID
	State           string             `json:"state"`         // Layer 3 network availability state
	DnsDomain       string             `json:"dnsDomain"`     // DNS domain
	System          bool               `json:"system"`        // Whether it is used for system cloud hosts
	Category        string             `json:"category"`      // Network type, needs to be used with the system tag, can be set to Public or Private when system is true
	IpVersion       int                `json:"ipVersion"`     // IP protocol number: 4, 6
	Dns             []string           `json:"dns"`
	EnableIPAM      bool               `json:"enableIPAM"`
	IpRanges        []IpRangeInventory `json:"ipRanges"`
	NetworkServices []NetworkServices  `json:"networkServices"`
	HostRoute       []HostRoute        `json:"hostRoute"`
}

type IpRangeInventory struct {
	UUID          string `json:"uuid"`          // Resource UUID, uniquely identifies the resource
	L3NetworkUuid string `json:"l3NetworkUuid"` // Layer 3 network UUID
	Name          string `json:"name"`          // Resource name
	Description   string `json:"description"`   // Detailed description of the resource
	StartIp       string `json:"StartIp"`
	EndIp         string `json:"EndIp"`
	Netmask       string `json:"netmask"`   // Network mask
	PrefixLen     string `json:"prefixLen"` // Mask length
	Gateway       string `json:"gateway"`   // Gateway address
	NetworkCidr   string `json:"networkCidr"`
	IpVersion     string `json:"ipVersion"`   // IP protocol number: 4, 6
	AddressMode   string `json:"addressMode"` // IPv6 address allocation mode
	CreateDate    string `json:"createDate"`  // Creation time
	LastOpDate    string `json:"lastOpDate"`  // Last modification time
	IpRangeType   string `json:"ipRangeType"`
}

type NetworkServices struct {
	L3NetworkUuid              string `json:"l3NetworkUuid"`              // Layer 3 network UUID
	NetworkServiceProviderUuid string `json:"networkServiceProviderUuid"` // Network service provider module UUID
	NetworkServiceType         string `json:"networkServiceType"`
}

type HostRoute struct {
	Id            string `json:"id"`
	L3NetworkUuid string `json:"l3NetworkUuid"` // Layer 3 network UUID
	Prefix        string `json:"prefix"`
	Nexthop       string `json:"nexthop"`
	CreateDate    string `json:"createDate"` // Creation time
	LastOpDate    string `json:"lastOpDate"` // Last modification time
}

type FreeIpInventoryView struct {
	IpRangeUuid string `json:"ipRangeUuid"` // IP range UUID
	Ip          string `json:"ip"`          // IP
	Netmask     string `json:"netmask"`
	Gateway     string `json:"gateway"`
}

type CheckIpView struct {
	Available bool `json:"available"`
}

type IpAddressCapacityView struct {
	TotalCapacity           int64            `json:"totalCapacity"`           // Total IP address capacity
	AvailableCapacity       int64            `json:"availableCapacity"`       // Available IP address capacity
	UsedIpAddressNumber     int64            `json:"usedIpAddressNumber"`     // Number of used IP addresses
	Ipv4TotalCapacity       int64            `json:"ipv4TotalCapacity"`       // Total IPv4 address capacity
	Ipv4AvailableCapacity   int64            `json:"ipv4AvailableCapacity"`   // Available IPv4 address capacity
	Ipv4UsedIpAddressNumber int64            `json:"ipv4UsedIpAddressNumber"` // Number of used IPv4 addresses
	Ipv6TotalCapacity       int64            `json:"ipv6TotalCapacity"`       // Total IPv6 address capacity
	Ipv6AvailableCapacity   int64            `json:"ipv6AvailableCapacity"`   // Available IPv6 address capacity
	Ipv6UsedIpAddressNumber int64            `json:"ipv6UsedIpAddressNumber"` // Number of used IPv6 addresses
	ResourceType            string           `json:"resourceType"`            // Type of the queried resource (IP range, Layer 3 network, zone)
	Success                 bool             `json:"success"`                 // Success
	CapacityData            []IpCapacityData `json:"capacityData"`
}

type IpCapacityData struct {
	ResourceUuid            string `json:"resourceUuid,omitempty"`  // Resource UUID. If specified, the image will use this value as the UUID.
	TotalCapacity           int64  `json:"totalCapacity"`           // Total IP address capacity
	AvailableCapacity       int64  `json:"availableCapacity"`       // Available IP address capacity
	UsedIpAddressNumber     int64  `json:"usedIpAddressNumber"`     // Number of used IP addresses
	Ipv4TotalCapacity       int64  `json:"ipv4TotalCapacity"`       // Total IPv4 address capacity
	Ipv4AvailableCapacity   int64  `json:"ipv4AvailableCapacity"`   // Available IPv4 address capacity
	Ipv4UsedIpAddressNumber int64  `json:"ipv4UsedIpAddressNumber"` // Number of used IPv4 addresses
	Ipv6TotalCapacity       int64  `json:"ipv6TotalCapacity"`       // Total IPv6 address capacity
	Ipv6AvailableCapacity   int64  `json:"ipv6AvailableCapacity"`   // Available IPv6 address capacity
	Ipv6UsedIpAddressNumber int64  `json:"ipv6UsedIpAddressNumber"` // Number of used IPv6 addresses
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
	Name          string `json:"name"`      // Maximum length of 255 characters
	StartIp       string `json:"startIp"`   // Start IP address, in IPv4
	EndIp         string `json:"endIp"`     // End IP address, in IPv4
	IpVersion     int    `json:"ipVersion"` // IP version (e.g., 4 for IPv4)
}
