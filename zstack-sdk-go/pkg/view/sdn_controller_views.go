// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

type SdnControllerInventoryView struct {
	BaseInfoView
	BaseTimeView

	Ip         string       `json:"ip"`
	Status     string       `json:"status"`
	UserName   string       `json:"username"`
	Passwordd  string       `json:"password"`
	VendorType string       `json:"verdorType"`
	HostRefs   []HostRefs   `json:"hostRefs"`
	VniRanges  []VniRanges  `json:"vniRanges"`
	VxlanPools []VxlanPools `json:"vxlanPools"`
}

type HostRefs struct {
	SdnControllerUuid string `json:"sdnControllerUuid"`
	HostUuid          string `json:"hostUuid"`
	VSwitchType       string `json:"vSwitchType"`
	VtepIp            string `json:"vtepIp"`
	NicPciAddresses   string `json:"nicPciAddresses"`
	NicDrivers        string `json:"nicDrivers"`
	Netmask           string `json:"netmask"`
	BondMode          string `json:"bondMode"`
	LacpMode          string `json:"lacpMode"`
}

type VniRanges struct {
	StartVni string `json:"startVni"`
	EndVni   string `json:"endVni"`
}

type VxlanPools struct{}
