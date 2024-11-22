// Copyright (c) ZStack.io, Inc.

package param

type UpdateGlobalConfigParam struct {
	BaseParam
	UpdateGlobalConfig UpdateConfigDetailParam `json:"updateGlobalConfig"`
}

type UpdateConfigDetailParam struct {
	Value string `json:"value"`
}

type UpdateResourceConfigParam struct {
	BaseParam
	UpdateResourceConfig UpdateConfigDetailParam `json:"updateResourceConfig"`
}
