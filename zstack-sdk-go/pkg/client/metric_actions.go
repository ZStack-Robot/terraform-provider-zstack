// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// GetAllMetricMetadata Retrieve all metric metadata
func (cli *ZSClient) GetAllMetricMetadata(params param.GetAllMetricMetadataParam) ([]view.MetricMetadataView, error) {
	queryParams := params.ToQueryParam()
	var resp []view.MetricMetadataView
	return resp, cli.ListWithRespKey("v1/zwatch/metrics/meta-data", "metrics", &queryParams, &resp)
}

// GetMetricLabelValue Retrieve the label values for a metric
func (cli *ZSClient) GetMetricLabelValue(params param.GetMetricLabelValueParam) ([]map[string]interface{}, error) {
	queryParams := params.ToQueryParam()
	var resp []map[string]interface{}
	return resp, cli.ListWithRespKey("v1/zwatch/metrics/label-values", "labels", &queryParams, &resp)
}

// GetMetricData Retrieve metric data
func (cli *ZSClient) GetMetricData(params param.GetMetricDataParam) ([]view.MetricDataView, error) {
	queryParams := params.ToQueryParam()
	var resp []view.MetricDataView
	return resp, cli.ListWithRespKey("v1/zwatch/metrics", "data", &queryParams, &resp)
}
