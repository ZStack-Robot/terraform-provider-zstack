// Copyright (c) HashiCorp, Inc.

package param

import "strconv"

type GetAllMetricMetadataParam struct {
	BaseParam

	Namespace string `json:"namespace"` //监控指标命名空间
	Name      string `json:"name"`      //监控指标名称
}

type GetMetricLabelValueParam struct {
	BaseParam

	Namespace    string   `json:"namespace"`    //名字空间名称
	MetricName   string   `json:"metricName"`   //监控指标名称
	StartTime    int64    `json:"startTime"`    //起始时间，时间戳，秒
	EndTime      int64    `json:"endTime"`      //结束时间，时间戳，秒
	LabelNames   []string `json:"labelNames"`   //要获取值得标签名列表
	FilterLabels []string `json:"filterLabels"` //标签过滤器列表，例如可以指定标签HostUuid=e47f7145f4cd4fca8e2856038ecdf3e1来选择特定物理机的，labelNames中指定标签的值
}

type GetMetricDataParam struct {
	BaseParam

	Namespace                string   `json:"namespace"`                //名字空间
	MetricName               string   `json:"metricName"`               //监控项
	StartTime                int64    `json:"startTime"`                //起始时间，时间戳，秒
	EndTime                  int64    `json:"endTime"`                  //结束时间，时间戳，秒
	Period                   int32    `json:"period"`                   //数据精度
	Labels                   []string `json:"labels"`                   //过滤标签
	ValueConditions          []string `json:"valueConditions"`          //未知TODO
	Functions                []string `json:"functions"`                //函数列表
	OffsetAheadOfCurrentTime int64    `json:"offsetAheadOfCurrentTime"` //未知TODO
}

func (p GetAllMetricMetadataParam) ToQueryParam() QueryParam {
	result := NewQueryParam()
	if p.Namespace != "" {
		result.Set("namespace", p.Namespace)
	}
	if p.Name != "" {
		result.Set("name", p.Name)
	}
	if p.RequestIp != "" {
		result.Set("requestIp", p.RequestIp)
	}
	for _, systemTag := range p.SystemTags {
		if result.Get("systemTags") == "" {
			result.Set("systemTags", systemTag)
		} else {
			result.Add("systemTags", systemTag)
		}
	}
	for _, userTag := range p.UserTags {
		if result.Get("userTags") == "" {
			result.Set("userTags", userTag)
		} else {
			result.Add("userTags", userTag)
		}
	}
	return result
}

func (p GetMetricLabelValueParam) ToQueryParam() QueryParam {
	result := NewQueryParam()
	result.Set("namespace", p.Namespace)
	result.Set("metricName", p.MetricName)
	if p.StartTime != 0 {
		result.Set("startTime", strconv.FormatInt(p.StartTime, 10))
	}
	if p.EndTime != 0 {
		result.Set("endTime", strconv.FormatInt(p.EndTime, 10))
	}
	if p.RequestIp != "" {
		result.Set("requestIp", p.RequestIp)
	}
	for _, labelName := range p.LabelNames {
		if result.Get("labelNames") == "" {
			result.Set("labelNames", labelName)
		} else {
			result.Add("labelNames", labelName)
		}
	}
	for _, filterLabel := range p.FilterLabels {
		if result.Get("filterLabels") == "" {
			result.Set("filterLabels", filterLabel)
		} else {
			result.Add("filterLabels", filterLabel)
		}
	}
	for _, systemTag := range p.SystemTags {
		if result.Get("systemTags") == "" {
			result.Set("systemTags", systemTag)
		} else {
			result.Add("systemTags", systemTag)
		}
	}
	for _, userTag := range p.UserTags {
		if result.Get("userTags") == "" {
			result.Set("userTags", userTag)
		} else {
			result.Add("userTags", userTag)
		}
	}
	return result
}

func (p GetMetricDataParam) ToQueryParam() QueryParam {
	result := NewQueryParam()
	result.Set("namespace", p.Namespace)
	result.Set("metricName", p.MetricName)
	if p.StartTime != 0 {
		result.Set("startTime", strconv.FormatInt(p.StartTime, 10))
	}
	if p.EndTime != 0 {
		result.Set("endTime", strconv.FormatInt(p.EndTime, 10))
	}
	if p.Period != 0 {
		result.Set("period", strconv.FormatInt(int64(p.Period), 10))
	}
	if p.OffsetAheadOfCurrentTime != 0 {
		result.Set("offsetAheadOfCurrentTime", strconv.FormatInt(p.OffsetAheadOfCurrentTime, 10))
	}
	if p.RequestIp != "" {
		result.Set("requestIp", p.RequestIp)
	}
	for _, label := range p.Labels {
		if result.Get("labels") == "" {
			result.Set("labels", label)
		} else {
			result.Add("labels", label)
		}
	}
	for _, valueCondition := range p.ValueConditions {
		if result.Get("valueConditions") == "" {
			result.Set("valueConditions", valueCondition)
		} else {
			result.Add("valueConditions", valueCondition)
		}
	}
	for _, function := range p.Functions {
		if result.Get("functions") == "" {
			result.Set("functions", function)
		} else {
			result.Add("functions", function)
		}
	}
	for _, systemTag := range p.SystemTags {
		if result.Get("systemTags") == "" {
			result.Set("systemTags", systemTag)
		} else {
			result.Add("systemTags", systemTag)
		}
	}
	for _, userTag := range p.UserTags {
		if result.Get("userTags") == "" {
			result.Set("userTags", userTag)
		} else {
			result.Add("userTags", userTag)
		}
	}
	return result
}
