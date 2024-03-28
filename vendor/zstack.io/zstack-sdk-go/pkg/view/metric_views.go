// Copyright (c) HashiCorp, Inc.

package view

type MetricMetadataView struct {
	Namespace   string   `json:"namespace"`   //名字空间
	Name        string   `json:"name"`        //资源名称
	Description string   `json:"description"` //资源名称
	LabelNames  []string `json:"labelNames"`  //标签名
	Driver      string   `json:"driver"`
}

type MetricDataView struct {
	Value  float64                `json:"value"`  //监控值
	Time   int64                  `json:"time"`   //记录生成时间，时间戳，秒
	Labels map[string]interface{} `json:"labels"` //标签
}
