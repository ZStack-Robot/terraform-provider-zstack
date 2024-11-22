// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateCdpPolicy 创建CDP策略
func (cli *ZSClient) CreateCdpPolicy(params *param.CreateCdpPolicyParam) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Post("v1/cdp-backup-storage/policy", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCdpPolicy 删除CDP策略
func (cli *ZSClient) DeleteCdpPolicy(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/cdp-backup-storage/policy", uuid, string(deleteMode))
}

// QueryCdpPolicy 查询CDP策略
func (cli *ZSClient) QueryCdpPolicy(params param.QueryParam) ([]view.CdpPolicyInventoryView, error) {
	resp := make([]view.CdpPolicyInventoryView, 0)
	return resp, cli.List("v1/cdp-backup-storage/policy", &params, &resp)
}

// PageCdpPolicy 分页查询CDP策略
func (cli *ZSClient) PageCdpPolicy(params param.QueryParam) ([]view.CdpPolicyInventoryView, int, error) {
	var resp []view.CdpPolicyInventoryView
	total, err := cli.Page("v1/cdp-backup-storage/policy", &params, &resp)
	return resp, total, err
}

// GetCdpPolicy 查询CDP策略
func (cli *ZSClient) GetCdpPolicy(uuid string) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Get("v1/cdp-backup-storage/policy", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCdpPolicy 更新CDP策略
func (cli *ZSClient) UpdateCdpPolicy(uuid string, params *param.UpdateCdpPolicyParam) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Put("v1/cdp-backup-storage/policy", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCdpTask 创建CDP任务
func (cli *ZSClient) CreateCdpTask(params *param.CreateCdpTaskParam) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-backup-storage/task", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCdpTask 删除CDP任务
func (cli *ZSClient) DeleteCdpTask(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/cdp-task", uuid, string(deleteMode))
}

// QueryCdpTask 查询CDP任务
func (cli *ZSClient) QueryCdpTask(params param.QueryParam) ([]view.CdpTaskInventoryView, error) {
	resp := make([]view.CdpTaskInventoryView, 0)
	return resp, cli.List("v1/cdp-task", &params, &resp)
}

// PageCdpTask 分页查询CDP任务
func (cli *ZSClient) PageCdpTask(params param.QueryParam) ([]view.CdpTaskInventoryView, int, error) {
	var resp []view.CdpTaskInventoryView
	total, err := cli.Page("v1/cdp-task", &params, &resp)
	return resp, total, err
}

// GetCdpTask 查询CDP任务
func (cli *ZSClient) GetCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Get("v1/cdp-task", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCdpTask 更新CDP任务
func (cli *ZSClient) UpdateCdpTask(uuid string, params *param.UpdateCdpTaskParam) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Put("v1/cdp-backup-storage/task", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// EnableCdpTask 启用CDP任务
func (cli *ZSClient) EnableCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-task/enable/"+uuid, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DisableCdpTask 禁用CDP任务
func (cli *ZSClient) DisableCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-task/disable/"+uuid, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MountVmInstanceRecoveryPoint 挂载CDP恢复点
func (cli *ZSClient) MountVmInstanceRecoveryPoint(params param.MountVmInstanceRecoveryPointParam) (*view.MountVmInstanceRecoveryPointView, error) {
	var resp view.MountVmInstanceRecoveryPointView
	if err := cli.PostWithRespKey("v1/cdp-backup-storage/mount-recovery-point", "", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnmountVmInstanceRecoveryPoint 卸载CDP恢复点
func (cli *ZSClient) UnmountVmInstanceRecoveryPoint(params param.UnmountVmInstanceRecoveryPointParam) error {
	return cli.PostWithRespKey("v1/cdp-backup-storage/unmount-recovery-point", "", params, nil)
}

// DeleteCdpTaskData 删除CDP任务的数据
func (cli *ZSClient) DeleteCdpTaskData(uuid string) error {
	reqUri := fmt.Sprintf("v1/cdp-task/%s/data", uuid)
	return cli.PostWithRespKey(reqUri, "", map[string]interface{}{}, nil)
}

// CreateVmFromCdpBackup 从CDP恢复点创建虚拟机
func (cli *ZSClient) CreateVmFromCdpBackup(params *param.CreateVmFromCdpBackupParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/cdp-backups/actions", "", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevertVmFromCdpBackup 从CDP恢复点创建虚拟机
func (cli *ZSClient) RevertVmFromCdpBackup(uuid string, params *param.RevertVmFromCdpBackupParam) error {
	return cli.PutWithSpec("v1/cdp-backups", uuid, "actions", "", params, nil)
}
