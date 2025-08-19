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
	State                  string                           `json:"state"`     // Enabled, Disabled
	IpVersion              int32                            `json:"ipVersion"` // IPv4, IPv6
	AttachedL3NetworkUuids []string                         `json:"attachedL3NetworkUuids"`
	VSwitchType            string                           `json:"vswitchType"` // LinuxBridge, OvnDpdk
	Rules                  []SecurityGroupRuleInventoryView `json:"rules"`       // Security group rules
}

type SecurityGroupRuleInventoryView struct {
	BaseTimeView
	UUID                    string `json:"uuid"`
	Action                  string `json:"action"` // Allow, Deny
	Description             string `json:"description,omitempty"`
	StartPort               int64  `json:"startPort"`
	EndPort                 int64  `json:"endPort"`
	IpVersion               int    `json:"ipVersion"`              // 4 or 6
	Priority                int    `json:"priority"`               // Optional, default is 0
	Protocol                string `json:"protocol"`               // TCP, UDP, ICMP, ALL
	AllowedCidr             string `json:"allowedCidr"`            // CIDR format, e.g., "
	DstPortRange            string `json:"dstPortRange,omitempty"` // e.g., "21, 80-443" for TCP/UDP
	SecurityGroupUuid       string `json:"securityGroupUuid"`      // UUID of the security group this rule belongs to
	SrcIpRange              string `json:"srcIpRange,omitempty"`   // Ingress Only. CIDR format, e.g., 192.168.1.1-192.168.1.10,192.168.1.11"
	DstIpRange              string `json:"dstIpRange,omitempty"`   // Egress Only. CIDR format, e.g.,  192.168.1.1-192.168.1.10,192.168.1.11""
	State                   string `json:"state"`                  // Enabled, Disabled
	Type                    string `json:"type"`                   // Ingress, Egress
	RemoteSecurityGroupUuid string `json:"remoteSecurityGroupUuid"`
}
