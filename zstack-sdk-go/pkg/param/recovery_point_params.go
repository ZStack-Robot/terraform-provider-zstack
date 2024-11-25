// Copyright (c) ZStack.io, Inc.

package param

const (
	RecoveryPointScaleDay    = "day"
	RecoveryPointScaleHour   = "hour"
	RecoveryPointScaleMinute = "minute"
)

type GetRecoveryPointParam struct {
	BaseParam `json:",inline,omitempty"`
	GroupId   int64 `json:"groupId" validate:"required"`
	// Start     interface{} `json:"start" `
	// Limit     interface{} `json:"limit" `
}

type QueryRecoveryPointParam struct {
	BaseParam `json:",inline,omitempty"`
	PageVar   `json:",inline,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Scale     string `json:"scale,omitempty"`
}

type QueryProtectRecoveryPointParam struct {
	BaseParam `json:",inline,omitempty"`
	PageVar   `json:",inline,omitempty"`
}

type ProtectVmInstanceRecoveryPointParam struct {
	BaseParam                      `json:",inline,omitempty"`
	ProtectVmInstanceRecoveryPoint ProtectVmInstanceRecoveryDetailPointParam `json:"protectVmInstanceRecoveryPoint"`
}

type ProtectVmInstanceRecoveryDetailPointParam struct {
	GroupId     int64  `json:"groupId,omitempty"`
	Description string `json:"description,omitempty"`
}

type UnprotectVmInstanceRecoveryPointParam struct {
	BaseParam                        `json:",inline,omitempty"`
	UnprotectVmInstanceRecoveryPoint UnprotectVmInstanceRecoveryPointDetailParam `json:"unprotectVmInstanceRecoveryPoint"`
}

type UnprotectVmInstanceRecoveryPointDetailParam struct {
	GroupId int64 `json:"groupId,omitempty"`
}
