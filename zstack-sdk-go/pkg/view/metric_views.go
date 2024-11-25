// Copyright (c) ZStack.io, Inc.

package view

type MetricMetadataView struct {
	Namespace   string   `json:"namespace"`   // Namespace
	Name        string   `json:"name"`        // Resource name
	Description string   `json:"description"` // Resource description
	LabelNames  []string `json:"labelNames"`  // Label names
	Driver      string   `json:"driver"`
}

type MetricDataView struct {
	Value  float64                `json:"value"`  // Monitoring value
	Time   int64                  `json:"time"`   // Record generation time, timestamp in seconds
	Labels map[string]interface{} `json:"labels"` // Labels
}
