// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryNetworkServiceProvider Query network service module
func (cli *ZSClient) QueryNetworkServiceProvider(params param.QueryParam) ([]view.NetworkServiceProviderInventoryView, error) {
	var resp []view.NetworkServiceProviderInventoryView
	return resp, cli.List("v1/network-services/providers", &params, &resp)
}

// AttachNetworkServiceToL3Network Attach network service to L3 network
func (cli *ZSClient) AttachNetworkServiceToL3Network(l3NetworkUuid string, p param.AttachNetworkServiceToL3NetworkParam) error {
	return cli.Post("v1/l3-networks/"+l3NetworkUuid+"/network-services", p, nil)
}
