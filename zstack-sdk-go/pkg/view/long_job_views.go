// Copyright (c) HashiCorp, Inc.

package view

import (
	"zstack.io/zstack-sdk-go/pkg/param"
)

type LongJobInventoryView struct {
	BaseInfoView
	BaseTimeView

	ApiId              string             `json:"apiId" `              //用于关联TaskProgress的APIID
	JobName            string             `json:"jobName" `            //任务名称
	JobData            string             `json:"jobData" `            //任务数据
	JobResult          string             `json:"jobResult" `          //任务结果
	TargetResourceUuid string             `json:"targetResourceUuid" ` //目标资源UUID
	ManagementNodeUuid string             `json:"managementNodeUuid" ` //管理节点UUID
	State              param.LongJobState `json:"state" `
	ExecuteTime        int64              `json:"executeTime" `
}

type TaskProgressInventoryView struct {
	TaskUuid   string                      `json:"taskUuid" `
	TaskName   string                      `json:"taskName" `
	ParentUuid string                      `json:"parentUuid" `
	Type       string                      `json:"type" `
	Content    string                      `json:"content" `
	Opaque     interface{}                 `json:"opaque" `
	Time       int64                       `json:"time" `
	SubTasks   []TaskProgressInventoryView `json:"subTasks" `
}
