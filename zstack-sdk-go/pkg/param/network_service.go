// Copyright (c) ZStack.io, Inc.

package param

type AttachNetworkServiceToL3NetworkParam struct {
	BaseParam
	Params AttachNetworkServiceToL3NetworkDetailParam `json:"params"`
}

type AttachNetworkServiceToL3NetworkDetailParam struct {
	NetworkServices map[string][]string `json:"networkServices"`
}
