// Copyright (c) ZStack.io, Inc.

package param

type AttachL3NetworkToVmParam struct {
	BaseParam
	Params AttachL3NetworkToVmDetailParam `json:"params"`
}

type AttachL3NetworkToVmDetailParam struct {
	StaticIp string `json:"staticIp"` // Specified IP address to be assigned to the VM
}

type UpdateVmNicMacParam struct {
	BaseParam
	UpdateVmNicMac UpdateVmNicMacDetailParam `json:"updateVmNicMac"`
}

type UpdateVmNicMacDetailParam struct {
	Mac string `json:"mac"` // MAC address
}

type SetVmStaticIpParam struct {
	BaseParam
	SetVmStaticIp SetVmStaticIpDetailParam `json:"setVmStaticIp"`
}

type SetVmStaticIpDetailParam struct {
	L3NetworkUuid string `json:"l3NetworkUuid"` // Layer 3 network UUID
	Ip            string `json:"ip"`            // Specified IP address
	Ip6           string `json:"ip6"`           // Specified IPv6 address
}

type DeleteVmStaticIpParam struct {
	BaseParam
	Params DeleteVmStaticIpDetailParam `json:"params"`
}

type DeleteVmStaticIpDetailParam struct {
	L3NetworkUuid string     `json:"l3NetworkUuid"` // Layer 3 network UUID
	DeleteMode    DeleteMode `json:"deleteMode"`
}

type ChangeVmNicNetworkParam struct {
	BaseParam
	Params ChangeVmNicNetworkDetailParam `json:"params"`
}

type ChangeVmNicNetworkDetailParam struct {
	DestL3NetworkUuid string `json:"destL3NetworkUuid"` // Specified layer 3 network UUID
}
