// Copyright (c) ZStack.io, Inc.

package view

import "time"

type PointResourceInfoView struct {
	VirtualSizes map[string]interface{} `json:"virtualSizes" `
	RealSizes    map[string]interface{} `json:"realSizes" `
	Timestamp    time.Time              `json:"timestamp" `
}

type RecoveryPointRespView map[string][]RecoveryPointView

type RecoveryPointView struct {
	Id    int64  `json:"id,omitempty"`    // id
	GrpId int64  `json:"grpId,omitempty"` // group id
	Size  int64  `json:"size,omitempty"`  // size
	VolId string `json:"volId,omitempty"` // volume id
	Ts    string `json:"ts,omitempty"`    // timestamp
	Desc  string `json:"desc,omitempty"`  // description
}

type ProtectRecoveryPointRespView RecoveryPointRespView
