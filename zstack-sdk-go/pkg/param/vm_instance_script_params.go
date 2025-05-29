// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package param

type RenderParams struct {
	Key         string `json:"key"`         // 渲染参数的键
	Value       string `json:"value"`       // 渲染参数的值
	Description string `json:"description"` // 渲染参数的描述
}

type CreateVmInstanceScriptParam struct {
	BaseParam
	Params CreateVmInstanceScriptDetailParam `json:"params"`
}
type CreateVmInstanceScriptDetailParam struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`            // 描述
	ScriptContent string   `json:"scriptContent"`          // 脚本内容
	EncodingType  string   `json:"encodingType"`           // 编码类型 Supports: UTF8, GBK, GB2312, Big5, ISO8859-1
	Platform      string   `json:"platform"`               // 平台类型
	ScriptType    string   `json:"scriptType"`             // script Type Supports: Shell, Python, Perl, Bat, Powershell
	ScriptTimeout int      `json:"scriptTimeout"`          // 超时时间（秒）
	RenderParams  string   `json:"renderParams,omitempty"` // Optional
	TagUuids      []string `json:"tagUuids"`
}

type UpdateVmInstanceScriptParam struct {
	BaseParam
	Params UpdateVmInstanceScriptDetailParam `json:"updateGuestVmScript"`
}

type UpdateVmInstanceScriptDetailParam struct {
	Name          string `json:"name,omitempty"`          // 资源名称
	Description   string `json:"description,omitempty"`   // 资源描述
	ScriptContent string `json:"scriptContent,omitempty"` // 脚本内容
	EncodingType  string `json:"encodingType,omitempty"`  // 编码类型
	Platform      string `json:"platform,omitempty"`      // 平台类型
	ScriptType    string `json:"scriptType,omitempty"`    // 脚本类型
	ScriptTimeout int    `json:"scriptTimeout,omitempty"` // 超时时间（秒）
	RenderParams  string `json:"renderParams,omitempty"`  // Optional
}

type ExecuteVmInstanceScriptParam struct {
	BaseParam
	Params ExecuteVmInstanceScriptDetailParam `json:"executeGuestVmScript"`
}

type ExecuteVmInstanceScriptDetailParam struct {
	VmInstanceUuids []string `json:"vmInstanceUuids"` // 虚拟机实例 UUID
	ScriptTimeout   int      `json:"scriptTimeout"`   // 脚本超时时间（秒）
}
