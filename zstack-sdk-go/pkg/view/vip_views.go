// Copyright (c) ZStack.io, Inc.

package view

type VipInventoryView struct {
	BaseInfoView
	BaseTimeView

	L3NetworkUUID      string         `json:"l3NetworkUuid"` // Layer 3 network UUID
	Ip                 string         `json:"ip"`
	State              string         `json:"state"`
	Gateway            string         `json:"gateway"`            // Gateway
	Netmask            string         `json:"netmask"`            // Subnet mask
	PrefixLen          string         `json:"prefixLen"`          // Mask length
	ServiceProvider    string         `json:"serviceProvider"`    // Service provider offering VIP service
	PeerL3NetworkUuids string         `json:"peerL3NetworkUuids"` // L3 network UUIDs offering VIP service
	UseFor             string         `json:"useFor"`             // Usage, e.g., port forwarding
	System             bool           `json:"system"`             // Whether created by the system
	ServicesRefs       []ServicesRefs `json:"servicesRefs"`
}

type ServicesRefs struct {
	UUID        string `json:"uuid"`        // Resource UUID, uniquely identifies the resource
	ServiceType string `json:"serviceType"` // Service type
	VipUuid     string `json:"vipUuid"`
	CreateDate  string `json:"createDate"` // Creation time
	LastOpDate  string `json:"lastOpDate"` // Last modification time
}
