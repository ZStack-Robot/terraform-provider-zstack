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

type SecurityGroupInventoryView struct {
	BaseInfoView
	BaseTimeView
	State                  string                  `json:"state"`     // Enabled, Disabled
	IpVersion              string                  `json:"ipVersion"` // IPv4, IPv6
	AttachedL3NetworkUuids []string                `json:"attachedL3NetworkUuids"`
	Rules                  []SecurityGroupRuleView `json:"rules"` // Security group rules
}

type SecurityGroupRuleView struct {
	BaseInfoView
	BaseTimeView
	Type                    string `json:"type"` // Ingress, Egress
	IpVersion               string `json:"ipVersion"`
	StartPort               int64  `json:"startPort"`
	EndPort                 int64  `json:"endPort"`
	Protocol                string `json:"protocol"`                // TCP, UDP, ICMP, ICMPv
	State                   string `json:"state"`                   // Enabled, Disabled
	AllowedCidr             string `json:"allowedCidr"`             // CIDR format, e.g
	RemoteSecurityGroupUuid string `json:"remoteSecurityGroupUuid"` // UUID of the remote security group, if applicable
}
