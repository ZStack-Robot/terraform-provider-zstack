// Copyright (c) HashiCorp, Inc.

package param

type RequestConsoleAccessParam struct {
	BaseParam
	Params RequestConsoleAccessDetailParam `json:"params" bson:"params"`
}

type RequestConsoleAccessDetailParam struct {
	VMInstanceUUID string `json:"vmInstanceUuid" bson:"vmInstanceUuid"` //云主机UUID
}
