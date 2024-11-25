// Copyright (c) ZStack.io, Inc.

package param

type AttachType string

const (
	PassThrough AttachType = "PassThrough"
	Redirect    AttachType = "Redirect"
)

type UpdateUsbDeviceParam struct {
	BaseParam
	UpdateUsbDevice UpdateUsbDeviceDetailParam `json:"updateUsbDevice"`
}

type UpdateUsbDeviceDetailParam struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	State       *string `json:"state"`
}

type AttachUsbDeviceToVmParam struct {
	BaseParam
	Params AttachUsbDeviceToVmDetailParam `json:"params"`
}
type AttachUsbDeviceToVmDetailParam struct {
	VmInstanceUuid string `json:"vmInstanceUuid"`
	AttachType     string `json:"attachType"` // PassThrough Redirect
}

type DetachUsbDeviceFromVmParam struct {
	BaseParam
	Params interface{} `json:"params"`
}
