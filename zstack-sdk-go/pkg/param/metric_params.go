// Copyright (c) ZStack.io, Inc.

package param

import "strconv"

type GetAllMetricMetadataParam struct {
	BaseParam
	Namespace string `json:"namespace"` // Metric namespace
	Name      string `json:"name"`      // Metric name
}

type GetMetricLabelValueParam struct {
	BaseParam
	Namespace    string   `json:"namespace"`    // Namespace name
	MetricName   string   `json:"metricName"`   // Metric name
	StartTime    int64    `json:"startTime"`    // Start time, timestamp in seconds
	EndTime      int64    `json:"endTime"`      // End time, timestamp in seconds
	LabelNames   []string `json:"labelNames"`   // List of label names to get values for
	FilterLabels []string `json:"filterLabels"` // List of label filters, e.g., HostUuid=e47f7145f4cd4fca8e2856038ecdf3e1 to select a specific physical machine
}

type GetMetricDataParam struct {
	BaseParam
	Namespace                string   `json:"namespace"`                // Namespace
	MetricName               string   `json:"metricName"`               // Metric name
	StartTime                int64    `json:"startTime"`                // Start time, timestamp in seconds
	EndTime                  int64    `json:"endTime"`                  // End time, timestamp in seconds
	Period                   int32    `json:"period"`                   // Data precision
	Labels                   []string `json:"labels"`                   // Filter labels
	ValueConditions          []string `json:"valueConditions"`          // Value conditions (TODO: clarify)
	Functions                []string `json:"functions"`                // Function list
	OffsetAheadOfCurrentTime int64    `json:"offsetAheadOfCurrentTime"` // Offset ahead of current time (TODO: clarify)
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
