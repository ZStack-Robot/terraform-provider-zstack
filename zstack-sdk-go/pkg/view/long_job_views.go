// Copyright (c) ZStack.io, Inc.

package view

import (
	"zstack.io/zstack-sdk-go/pkg/param"
)

type LongJobInventoryView struct {
	BaseInfoView
	BaseTimeView

	ApiId              string             `json:"apiId"`              // API ID used to associate with TaskProgress
	JobName            string             `json:"jobName"`            // Job name
	JobData            string             `json:"jobData"`            // Job data
	JobResult          string             `json:"jobResult"`          // Job result
	TargetResourceUuid string             `json:"targetResourceUuid"` // Target resource UUID
	ManagementNodeUuid string             `json:"managementNodeUuid"` // Management node UUID
	State              param.LongJobState `json:"state"`
	ExecuteTime        int64              `json:"executeTime"`
}

type TaskProgressInventoryView struct {
	TaskUuid   string                      `json:"taskUuid"`
	TaskName   string                      `json:"taskName"`
	ParentUuid string                      `json:"parentUuid"`
	Type       string                      `json:"type"`
	Content    string                      `json:"content"`
	Opaque     interface{}                 `json:"opaque"`
	Time       int64                       `json:"time"`
	SubTasks   []TaskProgressInventoryView `json:"subTasks"`
}
