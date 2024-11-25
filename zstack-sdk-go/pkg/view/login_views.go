// Copyright (c) ZStack.io, Inc.

package view

import "time"

type SessionView struct {
	UUID        string    `json:"uuid"`        // Resource UUID, uniquely identifies the resource
	AccountUuid string    `json:"accountUuid"` // Account UUID
	UserUuid    string    `json:"userUuid"`    // User UUID
	ExpiredDate time.Time `json:"expiredDate"` // Session expiration date
	CreateDate  time.Time `json:"createDate"`  // Creation time
}

type WebUISessionView struct {
	SessionId       string `json:"sessionId"`   // Resource UUID
	AccountUuid     string `json:"accountUuid"` // Account UUID
	UserUuid        string `json:"userUuid"`    // User UUID
	UserName        string `json:"username"`    // Username
	LoginType       string `json:"loginType"`
	CurrentIdentity string `json:"currentIdentity"`
	ZSVersion       string `json:"zsVersion"` // ZStack Cloud detailed version
}
