// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateDiskOffering creates a cloud host specification
func (cli *ZSClient) CreateDiskOffering(params *param.CreateDiskOfferingParam) (*view.DiskOfferingInventoryView, error) {
	var resp view.DiskOfferingInventoryView
	if err := cli.Post("v1/disk-offerings", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDiskOffering deletes a cloud host specification
func (cli *ZSClient) DeleteDiskOffering(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/disk-offerings", uuid, string(deleteMode))
}

// GetDiskOffering Get Disk Offering
func (cli *ZSClient) GetDiskOffering(uuid string) (*view.DiskOfferingInventoryView, error) {
	offering := view.DiskOfferingInventoryView{}
	return &offering, cli.Get("v1/disk-offerings", uuid, nil, &offering)
}

// QueryDiskOffering
func (cli *ZSClient) QueryDiskOffering(params param.QueryParam) ([]view.DiskOfferingInventoryView, error) {
	var offering []view.DiskOfferingInventoryView
	return offering, cli.List("v1/disk-offerings", &params, &offering)
}

// UpdateImage Edit Image
func (cli *ZSClient) UpdateDiskOffering(uuid string, params param.UpdateImageParam) (view.ImageView, error) {
	image := view.ImageView{}
	return image, cli.Put("v1/images", uuid, params, &image)
}
