// Copyright (c) HashiCorp, Inc.

package view

type InstanceOfferingInventoryView struct {
	BaseInfoView
	BaseTimeView

	CpuNum            int    `json:"cpuNum" `            //CPU数量
	CpuSpeed          int    `json:"cpuSpeed" `          //CPU速度
	MemorySize        int64  `json:"memorySize" `        //内存大小
	Type              string `json:"type" `              //类型
	AllocatorStrategy string `json:"allocatorStrategy" ` //分配策略
	SortKey           int    `json:"sortKey" `
	State             string `json:"state" ` //状态（启用，禁用）
}
