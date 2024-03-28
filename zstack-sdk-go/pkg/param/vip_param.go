// Copyright (c) HashiCorp, Inc.

package param

type VipAllocatorStrategy string

const (
	DefaultHostAllocatorStrategy            VipAllocatorStrategy = "DefaultHostAllocatorStrategy"
	LastHostPreferredAllocatorStrategy                           = "LastHostPreferredAllocatorStrategy"
	LeastVmPreferredHostAllocatorStrategy                        = "LeastVmPreferredHostAllocatorStrategy"
	MinimumCPUUsageHostAllocatorStrategy                         = "MinimumCPUUsageHostAllocatorStrategy"
	MinimumMemoryUsageHostAllocatorStrategy                      = "MinimumMemoryUsageHostAllocatorStrategy"
	MaxInstancePerHostHostAllocatorStrategy                      = "MaxInstancePerHostHostAllocatorStrategy"
)

type CreateVipParam struct {
	BaseParam

	Params CreateVipDetailParam `json:"params"`
}

type CreateVipDetailParam struct {
	Name              string               `json:"name"`                        //资源名称
	Description       string               `json:"description,omitempty"`       //详细描述
	L3NetworkUUID     string               `json:"l3NetworkUuid"`               //	三层网络UUID
	IpRangeUUID       string               `json:"ipRangeUuid,omitempty"`       //IP段UUID
	AllocatorStrategy VipAllocatorStrategy `json:"allocatorStrategy,omitempty"` //分配策略
	RequiredIp        string               `json:"requiredIp,omitempty"`        //请求的IP
	ResourceUuid      string               `json:"resourceUuid,omitempty"`      //资源UUID。若指定，镜像会使用该字段值作为UUID。
}

type UpdateVipParam struct {
	BaseParam

	UUID      string               `json:"uuid"` //资源的UUID，唯一标示该资源
	UpdateVip UpdateVipDetailParam `json:"updateVip"`
}

type UpdateVipDetailParam struct {
	Name        string `json:"name"`                  //资源名称
	Description string `json:"description,omitempty"` //详细描述
}
