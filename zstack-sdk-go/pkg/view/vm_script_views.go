// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

type VmInstanceScriptInventoryView struct {
	BaseInfoView
	BaseTimeView

	ScriptContent string `json:"scriptContent"` // 脚本内容，可能是明文或已解码后的内容
	RenderParams  string `json:"renderParams"`  // JSON 字符串，表示渲染参数数组
	Platform      string `json:"platform"`      // 平台类型，如 Linux
	ScriptType    string `json:"scriptType"`    // 脚本类型，如 Shell
	ScriptTimeout int    `json:"scriptTimeout"` // 超时时间（秒）
	EncodingType  string `json:"encodingType"`  // 编码类型，如 base64, planText
}

type VmInstanceScriptResultInventoryView struct {
	BaseInfoView
	BaseExecTimeView
	ScriptUUID     string `json:"scriptUuid"` // 脚本 UUID
	RecordName     string `json:"recordName"` // 记录名称
	ScriptTimeout  int    `json:"scriptTimeout"`
	Status         string `json:"status"`         // 状态，如 Running, Succeeded, Failed
	Executor       string `json:"executor"`       // 执行者
	ExecutionCount int    `json:"executionCount"` // 执行次数
	Version        int    `json:"version"`        // 脚本版本
	ScriptContent  string `json:"scriptContent"`  // 脚本内容，可能是明文或已解码后的内容
	RenderParams   string `json:"renderParams"`   // JSON 字符串，表示渲染参数数组
}

type VmInstanceScriptResultDetailInventoryView struct {
	BaseExecTimeView
	RecordUuid     string `json:"recordUuid"`     // 记录 UUID
	VmInstanceUuid string `json:"vmInstanceUuid"` // 虚拟机实例 UUID
	VmName         string `json:"vmName"`
	Status         string `json:"status"`
	ExitCode       int    `json:"exitCode"` // 退出码
	Stdout         string `json:"stdout"`   // 标准输出
	Stderr         string `json:"stderr"`   // 标准错误输出
	ErrCause       string `json:"errCause"` // 错误原因
}
