package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryHost 查询物理机
func (cli *ZSClient) QueryHost(params param.QueryParam) ([]view.HostInventoryView, error) {
	var resp []view.HostInventoryView
	return resp, cli.List("v1/hosts", &params, &resp)
}

// PageHost 物理机分页
func (cli *ZSClient) PageHost(params param.QueryParam) ([]view.HostInventoryView, int, error) {
	var resp []view.HostInventoryView
	total, err := cli.Page("v1/hosts", &params, &resp)
	return resp, total, err
}

// GetHost 物理机详情
func (cli *ZSClient) GetHost(uuid string) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Get("v1/hosts", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateHost 更新物理机
func (cli *ZSClient) UpdateHost(uuid string, params param.UpdateHostParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChangeHostState 更新物理机启用状态
func (cli *ZSClient) ChangeHostState(uuid string, params *param.ChangeHostStateParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ReconnectHost 重连物理机
func (cli *ZSClient) ReconnectHost(uuid string) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Put("v1/hosts", uuid, map[string]struct{}{
		"reconnectHost": {},
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddKVMHost 添加KVM物理机
func (cli *ZSClient) AddKVMHost(params param.AddKVMHostParam) (*view.HostInventoryView, error) {
	resp := view.HostInventoryView{}
	if err := cli.Post("v1/hosts/kvm", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteHost 删除物理机
func (cli *ZSClient) DeleteHost(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/hosts", uuid, string(deleteMode))
}

// QueryHostNetworkBonding 查询物理机Bond信息
func (cli *ZSClient) QueryHostNetworkBonding(params param.QueryParam) ([]view.HostNetworkBondingInventoryView, error) {
	var resp []view.HostNetworkBondingInventoryView
	return resp, cli.List("v1/hosts/bondings", &params, &resp)

}

// QueryHostNetworkInterface 查询物理机网卡信息
func (cli *ZSClient) QueryHostNetworkInterface(params param.QueryParam) ([]view.HostNetworkInterfaceInventoryView, error) {
	var resp []view.HostNetworkInterfaceInventoryView
	return resp, cli.List("v1/hosts/nics", &params, &resp)
}
