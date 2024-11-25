// Copyright (c) ZStack.io, Inc.

package view

const (
	VirtualRouter = "VirtualRouter"
	Vrouter       = "vrouter"
	SecurityGroup = "SecurityGroup"
	Flat          = "Flat"
)

type NetworkServiceProviderInventoryView struct {
	AttachedL2NetworkUuids []string `json:"attachedL2NetworkUuids"`
	CreateDate             string   `json:"createDate"`
	Description            string   `json:"description"`
	LastOpDate             string   `json:"lastOpDate"`
	Name                   string   `json:"name"`
	NetworkServiceTypes    []string `json:"networkServiceTypes"`
	Type                   string   `json:"type"` // VirtualRouter  vrouter  SecurityGroup  Flat
	Uuid                   string   `json:"uuid"`
}
