// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryPrimaryStorage Query primary storage
func (cli *ZSClient) QueryPrimaryStorage(params param.QueryParam) ([]view.PrimaryStorageInventoryView, error) {
	var views []view.PrimaryStorageInventoryView
	return views, cli.List("v1/primary-storage", &params, &views)
}

// PagePrimaryStorage Paginate query primary storage
func (cli *ZSClient) PagePrimaryStorage(params param.QueryParam) ([]view.PrimaryStorageInventoryView, int, error) {
	var views []view.PrimaryStorageInventoryView
	total, err := cli.Page("v1/primary-storage", &params, &views)
	return views, total, err
}

// QueryCephPrimaryStoragePool Query Ceph primary storage pool
func (cli *ZSClient) QueryCephPrimaryStoragePool(params param.QueryParam) ([]view.CephPrimaryStoragePoolInventoryView, error) {
	var views []view.CephPrimaryStoragePoolInventoryView
	return views, cli.List("v1/primary-storage/ceph/pools", &params, &views)
}

// PageCephPrimaryStoragePool Paginate query Ceph primary storage pool
func (cli *ZSClient) PageCephPrimaryStoragePool(params param.QueryParam) ([]view.CephPrimaryStoragePoolInventoryView, int, error) {
	var views []view.CephPrimaryStoragePoolInventoryView
	total, err := cli.Page("v1/primary-storage/ceph/pools", &params, &views)
	return views, total, err
}
