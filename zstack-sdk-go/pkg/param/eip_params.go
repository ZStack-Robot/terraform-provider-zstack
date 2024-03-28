// Copyright (c) HashiCorp, Inc.

package param

type CreateEipParam struct {
	BaseParam

	Params CreateEipDetailParam `json:"params"`
}

type CreateEipDetailParam struct {
	Name         string `json:"name"`                  //资源名称
	Description  string `json:"description,omitempty"` //详细描述
	VipUuid      string `json:"vipUuid"`
	VmNicUuid    string `json:"vmNicUuid,omitempty"`
	UsedIpUuid   *int   `json:"usedIpUuid,omitempty"`   //亲和组策略
	ResourceUuid string `json:"resourceUuid,omitempty"` //资源UUID。若指定，镜像会使用该字段值作为UUID。
}

type UpdateEipParam struct {
	BaseParam

	UUID      string               `json:"uuid"` //资源的UUID，唯一标示该资源
	UpdateEip UpdateEipDetailParam `json:"updateEip"`
}

type UpdateEipDetailParam struct {
	Name        string `json:"name,omitempty"`        //资源名称
	Description string `json:"description,omitempty"` //详细描述
}

type ChangeEipStateParam struct {
	BaseParam

	UUID           string                    `json:"uuid"` //资源的UUID，唯一标示该资源
	ChangeEipState ChangeEipStateDetailParam `json:"changeEipState"`
}

type ChangeEipStateDetailParam struct {
	StateEvent StateEvent `json:"stateEvent"`
}

type GetEipAttachableVmNicsParam struct {
	BaseParam

	EipUuid string `json:"eipUuid,omitempty"` //弹性IP UUID
	VipUuid string `json:"vipUuid,omitempty"` //VIP UUID
}

type GetVmNicAttachableEipsParam struct {
	BaseParam

	VmNicUuid string `json:"vmNicUuid"`
	IpVersion int    `json:"ipVersion,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Start     int    `json:"start,omitempty"`
}
