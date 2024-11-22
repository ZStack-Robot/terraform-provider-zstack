// Copyright (c) ZStack.io, Inc.

package param

type ChangeHostStateParam struct {
	BaseParam
	ChangeHostState ChangeHostStateDetailParam `json:"changeHostState"`
}
type ChangeHostStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"` //  enable disable maintain
}

type AddKVMHostParam struct {
	BaseParam
	Params AddKVMHostDetailParam `json:"params"`
}

type AddKVMHostDetailParam struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SshPort      int    `json:"sshPort"`
	Name         string `json:"name"`
	ManagementIp string `json:"managementIp"`
	ClusterUuid  string `json:"clusterUuid"`
}

type UpdateHostParam struct {
	BaseParam
	UpdateHost UpdateHostDetailParam `json:"updateHost"`
}

type UpdateHostDetailParam struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	ManagementIp *string `json:"managementIp"`
}
