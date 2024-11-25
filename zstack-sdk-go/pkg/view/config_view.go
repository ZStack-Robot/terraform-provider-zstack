// Copyright (c) ZStack.io, Inc.

package view

type GlobalConfigView struct {
	BaseInfoView

	Category     string `json:"category"`
	DefaultValue string `json:"defaultValue"`
	Value        string `json:"value"`
}

type ResourceConfigView struct {
	BaseInfoView
	BaseTimeView

	ResourceUuid string `json:"resourceUuid"` // UUID of the resource corresponding to the configuration
	ResourceType string `json:"resourceType"` // Type of the resource corresponding to the configuration
	Category     string `json:"category"`     // Configuration category
	Value        string `json:"value"`        // Value of the configuration
}
