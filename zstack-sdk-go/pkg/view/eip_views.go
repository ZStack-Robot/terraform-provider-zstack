// Copyright (c) ZStack.io, Inc.

package view

type EipInventoryView struct {
	BaseInfoView
	BaseTimeView

	VmNicUuid string `json:"vmNicUuid"` // UUID of the VM NIC
	VipUuid   string `json:"vipUuid"`
	State     string `json:"state"`
	VipIp     string `json:"vipIp"`
	GuestIp   string `json:"guestIp"`
}
