// Copyright (c) ZStack.io, Inc.

package param

type CreateEipParam struct {
	BaseParam

	Params CreateEipDetailParam `json:"params"`
}

type CreateEipDetailParam struct {
	Name         string `json:"name"`                  // Resource name
	Description  string `json:"description,omitempty"` // Detailed description
	VipUuid      string `json:"vipUuid"`
	VmNicUuid    string `json:"vmNicUuid,omitempty"`
	UsedIpUuid   *int   `json:"usedIpUuid,omitempty"`   // Affinity group policy
	ResourceUuid string `json:"resourceUuid,omitempty"` // Resource UUID. If specified, the image will use this value as the UUID.
}

type UpdateEipParam struct {
	BaseParam

	UUID      string               `json:"uuid"` // Resource UUID, uniquely identifies the resource
	UpdateEip UpdateEipDetailParam `json:"updateEip"`
}

type UpdateEipDetailParam struct {
	Name        string `json:"name,omitempty"`        // Resource name
	Description string `json:"description,omitempty"` // Detailed description
}

type ChangeEipStateParam struct {
	BaseParam

	UUID           string                    `json:"uuid"` // Resource UUID, uniquely identifies the resource
	ChangeEipState ChangeEipStateDetailParam `json:"changeEipState"`
}

type ChangeEipStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}

type GetEipAttachableVmNicsParam struct {
	BaseParam

	EipUuid string `json:"eipUuid,omitempty"` // Elastic IP UUID
	VipUuid string `json:"vipUuid,omitempty"` // VIP UUID
}

type GetVmNicAttachableEipsParam struct {
	BaseParam

	VmNicUuid string `json:"vmNicUuid"`
	IpVersion int    `json:"ipVersion,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Start     int    `json:"start,omitempty"`
}
