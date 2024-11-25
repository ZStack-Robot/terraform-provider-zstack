// Copyright (c) ZStack.io, Inc.

package view

type ClusterInventoryView struct {
	Architecture   string `json:"architecture"`
	CreateDate     string `json:"createDate"`
	HypervisorType string `json:"hypervisorType"`
	LastOpDate     string `json:"lastOpDate"`
	Name           string `json:"name"`
	State          string `json:"state"`
	Type           string `json:"type"`
	Uuid           string `json:"uuid"`
	ZoneUuid       string `json:"zoneUuid"`
}
