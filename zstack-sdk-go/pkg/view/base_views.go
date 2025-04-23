// Copyright (c) ZStack.io, Inc.

package view

import "time"

type BaseInfoView struct {
	UUID        string `json:"uuid"`        // Resource UUID, unique identifier for the resource
	Name        string `json:"name"`        // Resource name
	Description string `json:"description"` // Detailed description of the resource
}

type BaseTimeView struct {
	CreateDate time.Time `json:"createDate"` // Creation time
	LastOpDate time.Time `json:"lastOpDate"` // Last modification time
}

type BaseExecTimeView struct {
	StartTime time.Time `json:"startTime"` // Start time
	EndTime   time.Time `json:"endTime"`   // End time
}

type ErrorCodeView struct {
	Code        string                 `json:"code"`        // Error code number, global unique identifier for the error, e.g., SYS.1000, HOST.1001
	Description string                 `json:"description"` // Brief description of the error
	Details     string                 `json:"details"`     // Detailed information about the error
	Elaboration string                 `json:"elaboration"` // Reserved field, defaults to null
	Location    string                 `json:"location"`
	Cost        string                 `json:"cost"`
	Opaque      map[string]interface{} `json:"opaque"` // Reserved field, defaults to null
	Cause       *ErrorCodeView         `json:"cause"`  // Root error, the source error that caused the current error, null if there is no source error
}
