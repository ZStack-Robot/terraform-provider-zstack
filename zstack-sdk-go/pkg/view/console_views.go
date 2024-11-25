// Copyright (c) ZStack.io, Inc.

package view

type ConsoleInventoryView struct {
	Scheme   string `json:"scheme" bson:"scheme"`     // Access protocol type
	Hostname string `json:"hostname" bson:"hostname"` // Hostname of the host
	Port     int    `json:"port" bson:"port"`         // Port number
	Token    string `json:"token" bson:"token"`       // Token
}
