// Copyright (c) ZStack.io, Inc.

package view

import "time"

type ManagementNodeInventoryView struct {
	UUID      string    `json:"uuid"`      // Resource UUID, uniquely identifies the resource
	HostName  string    `json:"hostName"`  // Host name
	JoinDate  time.Time `json:"joinDate"`  // Join date
	HeartBeat time.Time `json:"heartBeat"` // Heartbeat time
}
