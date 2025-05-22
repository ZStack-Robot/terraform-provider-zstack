// Copyright (c) ZStack.io, Inc.

package param

const (
	ResourceTypeVmInstanceVO          = "VmInstanceVO"
	ResourceTypeImageVO               = "ImageVO"
	ResourceTypeVolumeVo              = "VolumeVO"
	ResourceTypeVolumeSnapshotVO      = "VolumeSnapshotVO"
	ResourceTypeVolumeSnapshotGroupVO = "VolumeSnapshotGroupVO"
	ResourceTypeL3NetworkVO           = "L3NetworkVO"
)

type CreateTagParam struct {
	BaseParam

	Params CreateTagDetailParam `json:"params"`
}

type CreateTagDetailParam struct {
	ResourceType string `json:"resourceType"`
	ResourceUuid string `json:"resourceUuid"`
	Tag          string `json:"tag"`
}

type UpdateSystemTagParam struct {
	BaseParam

	UpdateSystemTag UpdateTagDetailParam `json:"updateSystemTag"`
}

type UpdateTagDetailParam struct {
	Tag string `json:"tag"`
}

type CreateResourceTagParam struct {
	BaseParam
	Params CreateResourceTagDetailParam `json:"params"`
}
type CreateResourceTagDetailParam struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Type        string `json:"type"` // type 为 simple 不允许更新其 validValues = {"simple", "withToken"})
}

type UpdateResourceTagParam struct {
	BaseParam
	UpdateResourceTag UpdateResourceTagDetailParam `json:"updateTag"`
}

type UpdateResourceTagDetailParam struct {
	Name        string `json:"name"`
	Value       string `json:"value"` //不允许更改 simple Pattern 的 value，仅允许更改 withToken Pattern 的 key 值，如果是withToken 那么 name::{key1}::{key2} … ::{keyN}
	Description string `json:"description"`
	Color       string `json:"color"`
}

type AttachTagToResourceParam struct {
	BaseParam
	Params AttachTagToResourceDetailParam `json:"params"`
}

type AttachTagToResourceDetailParam struct {
	ResourceUuids []string               `json:"resourceUuids"`
	Tokens        map[string]interface{} `json:"tokens,omitempty"`
}
