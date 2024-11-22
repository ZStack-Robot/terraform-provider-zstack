// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateEip Create Elastic IP
func (cli *ZSClient) CreateEip(params param.CreateEipParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Post("v1/eips", params, &resp)
}

// DeleteEip Delete Elastic IP
func (cli *ZSClient) DeleteEip(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/eips", uuid, string(deleteMode))
}

// PageEip Pagination
func (cli *ZSClient) PageEip(params param.QueryParam) ([]view.EipInventoryView, int, error) {
	var eips []view.EipInventoryView
	total, err := cli.Page("v1/eips", &params, &eips)
	return eips, total, err
}

// QueryEip Query Elastic IP
func (cli *ZSClient) QueryEip(params param.QueryParam) ([]view.EipInventoryView, error) {
	resp := make([]view.EipInventoryView, 0)
	return resp, cli.List("v1/eips", &params, &resp)
}

// GetEip Query Elastic IP by UUID
func (cli *ZSClient) GetEip(uuid string) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Get("v1/eips", uuid, nil, &resp)
}

// UpdateEip Update Elastic IP
func (cli *ZSClient) UpdateEip(params param.UpdateEipParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Put("v1/eips", params.UUID, params, &resp)
}

// ChangeEipState Change Virtual IP enable state
func (cli *ZSClient) ChangeEipState(params param.ChangeEipStateParam) (view.EipInventoryView, error) {
	var resp view.EipInventoryView
	return resp, cli.Put("v1/eips", params.UUID, params, &resp)
}

// GetEipAttachableVmNics Get VM NICs that can be attached to the specified Elastic IP
func (cli *ZSClient) GetEipAttachableVmNics(params param.GetEipAttachableVmNicsParam) ([]view.VmNicInventoryView, error) {
	resp := make([]view.VmNicInventoryView, 0)
	return resp, cli.GetWithSpec("v1/eips", params.EipUuid, "vm-instances/candidate-nics", responseKeyInventories, nil, &resp)
}

// GetVmNicAttachableEips Get Elastic IPs that can be attached to a VM NIC
func (cli *ZSClient) GetVmNicAttachableEips(params param.GetVmNicAttachableEipsParam) ([]view.EipInventoryView, error) {
	resp := make([]view.EipInventoryView, 0)
	return resp, cli.GetWithSpec("v1/vm-instances/nics", params.VmNicUuid, "candidate-eips", responseKeyInventories, nil, &resp)
}

// AttachEip Attach Elastic IP
func (cli *ZSClient) AttachEip(eipUuid, vmNicUuid string) error {
	return cli.PutWithSpec("v1/eips", eipUuid, fmt.Sprintf("vm-instances/nics/%s", vmNicUuid), "", map[string]string{}, nil)
}

// DetachEip Detach Elastic IP
func (cli *ZSClient) DetachEip(eipUuid string) error {
	return cli.Delete("v1/eips", fmt.Sprintf("%s/vm-instances/nics", eipUuid), "")
}
