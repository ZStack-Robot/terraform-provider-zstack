// Copyright (c) ZStack.io, Inc.

package param

type VipAllocatorStrategy string

const (
	DefaultHostAllocatorStrategy            VipAllocatorStrategy = "DefaultHostAllocatorStrategy"
	LastHostPreferredAllocatorStrategy      VipAllocatorStrategy = "LastHostPreferredAllocatorStrategy"
	LeastVmPreferredHostAllocatorStrategy   VipAllocatorStrategy = "LeastVmPreferredHostAllocatorStrategy"
	MinimumCPUUsageHostAllocatorStrategy    VipAllocatorStrategy = "MinimumCPUUsageHostAllocatorStrategy"
	MinimumMemoryUsageHostAllocatorStrategy VipAllocatorStrategy = "MinimumMemoryUsageHostAllocatorStrategy"
	MaxInstancePerHostHostAllocatorStrategy VipAllocatorStrategy = "MaxInstancePerHostHostAllocatorStrategy"
)

type CreateVipParam struct {
	BaseParam
	Params CreateVipDetailParam `json:"params"`
}

type CreateVipDetailParam struct {
	Name              string               `json:"name"`                        // Resource name
	Description       string               `json:"description,omitempty"`       // Detailed description
	L3NetworkUUID     string               `json:"l3NetworkUuid"`               // Layer 3 network UUID
	IpRangeUUID       string               `json:"ipRangeUuid,omitempty"`       // IP range UUID
	AllocatorStrategy VipAllocatorStrategy `json:"allocatorStrategy,omitempty"` // Allocation strategy
	RequiredIp        string               `json:"requiredIp,omitempty"`        // Requested IP
	ResourceUuid      string               `json:"resourceUuid,omitempty"`      // Resource UUID. If specified, the image will use this field value as the UUID.
}

type UpdateVipParam struct {
	BaseParam
	UUID      string               `json:"uuid"` // Resource UUID, uniquely identifies the resource
	UpdateVip UpdateVipDetailParam `json:"updateVip"`
}

type UpdateVipDetailParam struct {
	Name        string `json:"name"`                  // Resource name
	Description string `json:"description,omitempty"` // Detailed description
}
