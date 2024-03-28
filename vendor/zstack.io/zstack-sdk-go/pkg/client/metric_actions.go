// Copyright (c) HashiCorp, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// GetAllMetricMetadata 获取所有metric元数据
func (cli *ZSClient) GetAllMetricMetadata(params param.GetAllMetricMetadataParam) ([]view.MetricMetadataView, error) {
	queryParams := params.ToQueryParam()
	var resp []view.MetricMetadataView
	return resp, cli.ListWithRespKey("v1/zwatch/metrics/meta-data", "metrics", &queryParams, &resp)
}

// GetMetricLabelValue 获取metric的标签值
func (cli *ZSClient) GetMetricLabelValue(params param.GetMetricLabelValueParam) ([]map[string]interface{}, error) {
	queryParams := params.ToQueryParam()
	var resp []map[string]interface{}
	return resp, cli.ListWithRespKey("v1/zwatch/metrics/label-values", "labels", &queryParams, &resp)
}

// GetMetricData 获取metric数据
func (cli *ZSClient) GetMetricData(params param.GetMetricDataParam) ([]view.MetricDataView, error) {
	queryParams := params.ToQueryParam()
	var resp []view.MetricDataView
	return resp, cli.ListWithRespKey("v1/zwatch/metrics", "data", &queryParams, &resp)
}
