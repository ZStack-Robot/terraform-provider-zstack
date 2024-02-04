package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryPrimaryStorage 查询主存储
func (cli *ZSClient) QueryPrimaryStorage(params param.QueryParam) ([]view.PrimaryStorageInventoryView, error) {
	var views []view.PrimaryStorageInventoryView
	return views, cli.List("v1/primary-storage", &params, &views)
}

// PagePrimaryStorage 分页查询主存储
func (cli *ZSClient) PagePrimaryStorage(params param.QueryParam) ([]view.PrimaryStorageInventoryView, int, error) {
	var views []view.PrimaryStorageInventoryView
	total, err := cli.Page("v1/primary-storage", &params, &views)
	return views, total, err
}

// QueryCephPrimaryStoragePool 查询Ceph主存储池
func (cli *ZSClient) QueryCephPrimaryStoragePool(params param.QueryParam) ([]view.CephPrimaryStoragePoolInventoryView, error) {
	var views []view.CephPrimaryStoragePoolInventoryView
	return views, cli.List("v1/primary-storage/ceph/pools", &params, &views)
}

// PageCephPrimaryStoragePool 分页查询Ceph主存储池
func (cli *ZSClient) PageCephPrimaryStoragePool(params param.QueryParam) ([]view.CephPrimaryStoragePoolInventoryView, int, error) {
	var views []view.CephPrimaryStoragePoolInventoryView
	total, err := cli.Page("v1/primary-storage/ceph/pools", &params, &views)
	return views, total, err
}
