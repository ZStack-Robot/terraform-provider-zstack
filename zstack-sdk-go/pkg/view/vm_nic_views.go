// Copyright (c) ZStack.io, Inc.

package view

import "time"

type VmNicInventoryView struct {
	UUID           string   `json:"uuid"`           // Resource UUID, uniquely identifies the resource
	VMInstanceUUID string   `json:"vmInstanceUuid"` // Cloud host UUID
	L3NetworkUUID  string   `json:"l3NetworkUuid"`  // Layer 3 network UUID
	IP             string   `json:"ip"`             // IP address
	Mac            string   `json:"mac"`            // MAC address
	HypervisorType string   `json:"hypervisorType"` // Hypervisor type
	Netmask        string   `json:"netmask"`        // Subnet mask
	Gateway        string   `json:"gateway"`        // Gateway
	MetaData       string   `json:"metaData"`       // Reserved field for internal use, metadata
	IpVersion      int      `json:"ipVersion"`      // IP address version
	DeviceID       int      `json:"deviceId"`       // Device ID, an integer that identifies the order of the NIC in the guest operating system's Ethernet devices. For example, 0 usually represents eth0, 1 usually represents eth1.
	DriverType     string   `json:"driverType"`     // NIC model
	Type           string   `json:"type"`           // NIC type
	CreateDate     string   `json:"createDate"`     // Creation time
	LastOpDate     string   `json:"lastOpDate"`     // Last modification time
	InternalName   string   `json:"internalName"`
	UsedIps        []UsedIp `json:"usedIps"`
}

type UsedIp struct {
	Uuid          string    `json:"uuid"`          // Resource UUID, uniquely identifies the resource
	IpRangeUuid   string    `json:"ipRangeUuid"`   // IP range UUID
	L3NetworkUuid string    `json:"l3NetworkUuid"` // Layer 3 network UUID
	IpVersion     int       `json:"ipVersion"`     // IP protocol number
	Ip            string    `json:"ip"`            // IP address
	Netmask       string    `json:"netmask"`       // Network mask
	Gateway       string    `json:"gateway"`       // Gateway address
	UsedFor       string    `json:"usedFor"`       //
	IpInLong      int64     `json:"ipInLong"`      //
	VmNicUuid     string    `json:"vmNicUuid"`     // Cloud host NIC UUID
	CreateDate    time.Time `json:"createDate"`    // Creation time
	LastOpDate    time.Time `json:"lastOpDate"`    // Last modification time
}

func GetIpFromUsedIps(usedIps []UsedIp) (ip string, ip6 string) {
	for _, usedIp := range usedIps {
		if usedIp.IpVersion == 4 {
			ip = usedIp.Ip
		}
		if usedIp.IpVersion == 6 {
			ip6 = usedIp.Ip
		}
	}
	return
}

type NicSimpleView struct {
	Ip        string    `json:"ip"`
	IpVersion string    `json:"ipVersion"`
	Uuid      string    `json:"uuid"`
	VmNicUuid string    `json:"vmNicUuid"`
	VmNic     VmNicView `json:"vmNic"`
}

type VmNicView struct {
	InternalName string `json:"internalName"`
	Mac          string `json:"mac"`
	Uuid         string `json:"uuid"`
}
