// Copyright (c) HashiCorp, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// ParseOvf 解析OVF模板信息
func (cli *ZSClient) ParseOvf(params param.ParseOvfParam) (*view.OvfInfo, error) {
	resp := view.OvfInfo{}
	return &resp, cli.PostWithRespKey("v1/ovf/parse", "ovfInfo", params, &resp)
}

// CreateVmInstanceFromOvf 从OVF模板导入云主机
func (cli *ZSClient) CreateVmInstanceFromOvf(params param.CreateVmInstanceFromOvfParam) (*view.VmInstanceInventoryView, error) {
	resp := view.VmInstanceInventoryView{}
	return &resp, cli.Post("v1/ovf/create-vm-instance", params, &resp)
}
