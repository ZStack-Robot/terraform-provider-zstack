// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateInstanceOffering creates a cloud host specification
func (cli *ZSClient) CreateVirtualRouterOffering(params param.CreateVirtualRouterOfferingParam) (*view.VirtualRouterOfferingInventoryView, error) {
	var resp view.VirtualRouterOfferingInventoryView
	if err := cli.Post("v1/instance-offerings/virtual-routers", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryVirtualRouter Offering Query VPC Virtual Router Offering
func (cli *ZSClient) QueryVirtualRouterOffering(params param.QueryParam) ([]view.VirtualRouterOfferingInventoryView, error) {
	resp := make([]view.VirtualRouterOfferingInventoryView, 0)
	return resp, cli.List("v1/instance-offerings/virtual-routers", &params, &resp)
}

// GetVirtualRouterOffering Query a specific VPC Virtual Router Offering
func (cli *ZSClient) GetVirtualRouterOffering(uuid string) (view.VirtualRouterOfferingInventoryView, error) {
	resp := view.VirtualRouterOfferingInventoryView{}
	return resp, cli.Get("v1/instance-offerings/virtual-routers", uuid, nil, &resp)
}

// DeleteInstanceOffering deletes a cloud host specification
/*
func (cli *ZSClient) DeleteVirtualRouterOffering(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/instance-offerings", uuid, string(deleteMode))
}
*/
