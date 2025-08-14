// Copyright (c) ZStack.io, Inc.

package param

type AttachNetworkServiceToL3NetworkParam struct {
	BaseParam
	Params AttachNetworkServiceToL3NetworkDetailParam `json:"params"`
}

type AttachNetworkServiceToL3NetworkDetailParam struct {
	NetworkServices map[string][]string `json:"networkServices"`
}

type AddVmNicToSecurityGroupParam struct {
	BaseParam
	Params AddVmNicToSecurityGroupDetailParam `json:"params"`
}

type AddVmNicToSecurityGroupDetailParam struct {
	//	SecurityGroupUuid string   `json:"securityGroupUuid"`
	VmNicUuids []string `json:"vmNicUuids"`
}

type CreateSecurityGroupParam struct {
	BaseParam
	Params CreateSecurityGroupDetailParam `json:"params"`
}

type CreateSecurityGroupDetailParam struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IpVersion   int    `json:"ipVersion"`   // 4 or 6 default is 4
	VSwitchType string `json:"vSwitchType"` // "LinuxBridge" or "OvnDpdk" default is "LinuxBridge"
}

type ChangeSecurityGroupStateParam struct {
	BaseParam
	SecurityGroupUuid string                              `json:"securityGroupUuid"` // securityGroup UUID
	ChangeImageState  ChangeSecurityGroupStateDetailParam `json:"changeSecurityGroupState"`
}

type ChangeSecurityGroupStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}

type AddSecurityGroupRuleParam struct {
	BaseParam
	Params AddSecurityGroupRuleDetailParam `json:"params"`
}

type AddSecurityGroupRuleDetailParam struct {
	Rules    []AddSecurityGroupRule `json:"rules"`
	Priority int                    `json:"priority"` // Optional
}

type AddSecurityGroupRule struct {
	RuleType                string `json:"type"`  // "Ingress" or "Egress
	State                   string `json:"state"` // "Enabled" or "Disabled"
	Description             string `json:"description"`
	IpVersion               int    `json:"ipVersion"`               // "4" or 6"
	Protocol                string `json:"protocol"`                // "TCP", "UDP", "ICMP", "ALL"
	SrcIpRange              string `json:"srcIpRange"`              // CIDR format, e.g., "
	DstIpRange              string `json:"dstIpRange"`              // CIDR format, e.g., "
	Action                  string `json:"action"`                  // "Allow" or "Deny"
	DstPortRange            string `json:"dstPortRange"`            // e.g., "21, 80-443" for TCP/UDP
	RemoteSecurityGroupUuid string `json:"remoteSecurityGroupUuid"` // Optional, for cross-security group rules
}

type UpdateSecurityGroupRuleParam struct {
	BaseParam
	ChangeSecurityGroupRule UpdateSecurityGroupRuleDetailParam `json:"changeSecurityGroupRule"`
}

type UpdateSecurityGroupRuleDetailParam struct {
	Description             string `json:"description,omitempty"`             // Optional, update rule description
	RemoteSecurityGroupUuid string `json:"remoteSecurityGroupUuid,omitempty"` // Optional, update remote security group UUID
	Action                  string `json:"action,omitempty"`                  // Optional, update action (Allow or Deny)
	State                   string `json:"state,omitempty"`                   // Optional, update state (Enabled or Disabled)
	Priority                int    `json:"priority"`                          // Require, update priority
	Protocol                string `json:"protocol,omitempty"`                // Optional, update protocol (TCP, UDP, ICMP, ALL)
	SrcIpRange              string `json:"srcIpRange,omitempty"`              // Optional, update source IP
	DstIpRange              string `json:"dstIpRange,omitempty"`              // Optional, update destination IP
	DstPortRange            string `json:"dstPortRange,omitempty"`            // Optional, update destination port
}
