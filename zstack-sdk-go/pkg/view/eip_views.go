// Copyright (c) HashiCorp, Inc.

package view

type EipInventoryView struct {
	BaseInfoView
	BaseTimeView

	VmNicUuid string `json:"vmNicUuid"` //云主机网卡UUID
	VipUuid   string `json:"vipUuid"`
	State     string `json:"state"`
	VipIp     string `json:"vipIp"`
	GuestIp   string `json:"guestIp"`
}
