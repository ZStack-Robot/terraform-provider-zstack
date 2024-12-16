// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestGetAllMetricMetadata(t *testing.T) {
	// data, err := accountLoginCli.GetAllMetricMetadata(param.GetAllMetricMetadataParam{})
	// data, err := accountLoginCli.GetAllMetricMetadata(param.GetAllMetricMetadataParam{Namespace: "ZStack/VM"})
	data, err := accountLoginCli.GetAllMetricMetadata(param.GetAllMetricMetadataParam{Namespace: "ZStack/VM", Name: "CPUUsedUtilization"})
	if err != nil {
		golog.Errorf("TestGetAllMetricMetadata %v", err)
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestGetMetricLabelValue(t *testing.T) {
	// data, err := accountLoginCli.GetMetricLabelValue(param.GetMetricLabelValueParam{
	// 	Namespace:    "ZStack/VM",
	// 	MetricName:   "CPUUsedUtilization",
	// 	LabelNames:   []string{"CPUNum"},
	// 	FilterLabels: []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	// data, err := accountLoginCli.GetMetricLabelValue(param.GetMetricLabelValueParam{
	// 	Namespace:    "ZStack/VM",
	// 	MetricName:   "DiskReadBytes",
	// 	LabelNames:   []string{"DiskDeviceLetter"},
	// 	FilterLabels: []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	data, err := accountLoginCli.GetMetricLabelValue(param.GetMetricLabelValueParam{
		Namespace:    "ZStack/VM",
		MetricName:   "NetworkOutBytes",
		LabelNames:   []string{"NetworkDeviceLetter"},
		FilterLabels: []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	})
	if err != nil {
		golog.Errorf("TestGetMetricLabelValue %v", err)
	}
	golog.Info(jsonutils.Marshal(data))
}

func TestGetMetricData(t *testing.T) {
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "CPUAverageUsedUtilization",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "CPUUsedUtilization",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721", "CPUNum=1"},
	// })
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "CPUUsedUtilization",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "MemoryFreeBytes",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "MemoryUsedBytes",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	// })
	data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
		Namespace:  "ZStack/VM",
		MetricName: "DiskReadBytes",
		StartTime:  1664352933,
		EndTime:    1664353833,
		Period:     256,
		Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721"},
	})
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "DiskReadBytes",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721", "DiskDeviceLetter=vda"},
	// })
	// data, err := accountLoginCli.GetMetricData(param.GetMetricDataParam{
	// 	Namespace:  "ZStack/VM",
	// 	MetricName: "NetworkOutBytes",
	// 	StartTime:  1664352933,
	// 	EndTime:    1664353833,
	// 	Period:     256,
	// 	Labels:     []string{"VMUuid=06018f0e95d147bfbda69cb75ad1e721", "NetworkDeviceLetter=vnic239.0"},
	// })
	if err != nil {
		golog.Errorf("TestGetMetricData %v", err)
	}
	golog.Info(jsonutils.Marshal(data))
}
