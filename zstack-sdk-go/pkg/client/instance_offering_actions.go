// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateInstanceOffering creates a cloud host specification
func (cli *ZSClient) CreateInstanceOffering(params *param.CreateInstanceOfferingParam) (*view.InstanceOfferingInventoryView, error) {
	var resp view.InstanceOfferingInventoryView
	if err := cli.Post("v1/instance-offerings", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteInstanceOffering deletes a cloud host specification
func (cli *ZSClient) DeleteInstanceOffering(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/instance-offerings", uuid, string(deleteMode))
}

// GetInstanceOffering Get Instance Offering
func (cli *ZSClient) GetInstanceOffering(uuid string) (*view.InstanceOfferingInventoryView, error) {
	offering := view.InstanceOfferingInventoryView{}
	return &offering, cli.Get("v1/instance-offerings", uuid, nil, &offering)
}

// QueryInstanceOffering Query Instance Offering
func (cli *ZSClient) QueryInstaceOffering(params param.QueryParam) ([]view.InstanceOfferingInventoryView, error) {
	var offering []view.InstanceOfferingInventoryView
	return offering, cli.List("v1/instance-offerings", &params, &offering)
}
