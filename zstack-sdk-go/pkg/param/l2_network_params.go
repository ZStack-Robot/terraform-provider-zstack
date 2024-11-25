// Copyright (c) ZStack.io, Inc.

package param

type UpdateL2NetworkParam struct {
	BaseParam
	UpdateL2Network UpdateL2NetworkDetailParam `json:"updateL2Network"`
}

type UpdateL2NetworkDetailParam struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type CreateL2NoVlanNetworkParam struct {
	BaseParam
	Params CreateL2NoVlanNetworkDetailParam `json:"params"`
}

type CreateL2NoVlanNetworkDetailParam struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	ZoneUuid          string `json:"zoneUuid"`          // Zone UUID
	PhysicalInterface string `json:"physicalInterface"` // Physical network interface
	Type              string `json:"type"`              // Layer 2 network type
	ResourceUuid      string `json:"resourceUuid"`      // Resource UUID
}

type CreateL2VlanNetworkParam struct {
	BaseParam
	Params CreateL2VlanNetworkDetailParam `json:"params"`
}

type CreateL2VlanNetworkDetailParam struct {
	Vlan              int    `json:"vlan"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	ZoneUuid          string `json:"zoneUuid"`
	PhysicalInterface string `json:"physicalInterface"` // Physical network interface
	Type              string `json:"type"`              // Layer 2 network type
	ResourceUuid      string `json:"resourceUuid"`      // Resource UUID
}

type AttachL2NetworkToClusterParam struct {
	BaseParam
}
