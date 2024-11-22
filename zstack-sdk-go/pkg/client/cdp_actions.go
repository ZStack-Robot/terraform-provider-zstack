// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateCdpPolicy creates a CDP policy.
func (cli *ZSClient) CreateCdpPolicy(params *param.CreateCdpPolicyParam) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Post("v1/cdp-backup-storage/policy", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCdpPolicy deletes a CDP policy.
func (cli *ZSClient) DeleteCdpPolicy(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/cdp-backup-storage/policy", uuid, string(deleteMode))
}

// QueryCdpPolicy queries CDP policies.
func (cli *ZSClient) QueryCdpPolicy(params param.QueryParam) ([]view.CdpPolicyInventoryView, error) {
	resp := make([]view.CdpPolicyInventoryView, 0)
	return resp, cli.List("v1/cdp-backup-storage/policy", &params, &resp)
}

// PageCdpPolicy performs paginated querying of CDP policies.
func (cli *ZSClient) PageCdpPolicy(params param.QueryParam) ([]view.CdpPolicyInventoryView, int, error) {
	var resp []view.CdpPolicyInventoryView
	total, err := cli.Page("v1/cdp-backup-storage/policy", &params, &resp)
	return resp, total, err
}

// GetCdpPolicy retrieves a specific CDP policy by UUID.
func (cli *ZSClient) GetCdpPolicy(uuid string) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Get("v1/cdp-backup-storage/policy", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCdpPolicy updates a specific CDP policy.
func (cli *ZSClient) UpdateCdpPolicy(uuid string, params *param.UpdateCdpPolicyParam) (*view.CdpPolicyInventoryView, error) {
	var resp view.CdpPolicyInventoryView
	if err := cli.Put("v1/cdp-backup-storage/policy", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateCdpTask creates a CDP task.
func (cli *ZSClient) CreateCdpTask(params *param.CreateCdpTaskParam) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-backup-storage/task", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCdpTask deletes a CDP task.
func (cli *ZSClient) DeleteCdpTask(uuid string, deleteMode param.DeleteMode) error {
	return cli.Delete("v1/cdp-task", uuid, string(deleteMode))
}

// QueryCdpTask queries CDP tasks.
func (cli *ZSClient) QueryCdpTask(params param.QueryParam) ([]view.CdpTaskInventoryView, error) {
	resp := make([]view.CdpTaskInventoryView, 0)
	return resp, cli.List("v1/cdp-task", &params, &resp)
}

// PageCdpTask performs paginated querying of CDP tasks.
func (cli *ZSClient) PageCdpTask(params param.QueryParam) ([]view.CdpTaskInventoryView, int, error) {
	var resp []view.CdpTaskInventoryView
	total, err := cli.Page("v1/cdp-task", &params, &resp)
	return resp, total, err
}

// GetCdpTask retrieves a specific CDP task by UUID.
func (cli *ZSClient) GetCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Get("v1/cdp-task", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateCdpTask updates a specific CDP task.
func (cli *ZSClient) UpdateCdpTask(uuid string, params *param.UpdateCdpTaskParam) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Put("v1/cdp-backup-storage/task", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// EnableCdpTask enables a specific CDP task.
func (cli *ZSClient) EnableCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-task/enable/"+uuid, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DisableCdpTask disables a specific CDP task.
func (cli *ZSClient) DisableCdpTask(uuid string) (*view.CdpTaskInventoryView, error) {
	var resp view.CdpTaskInventoryView
	if err := cli.Post("v1/cdp-task/disable/"+uuid, map[string]interface{}{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MountVmInstanceRecoveryPoint mounts a CDP recovery point.
func (cli *ZSClient) MountVmInstanceRecoveryPoint(params param.MountVmInstanceRecoveryPointParam) (*view.MountVmInstanceRecoveryPointView, error) {
	var resp view.MountVmInstanceRecoveryPointView
	if err := cli.PostWithRespKey("v1/cdp-backup-storage/mount-recovery-point", "", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnmountVmInstanceRecoveryPoint unmounts a CDP recovery point.
func (cli *ZSClient) UnmountVmInstanceRecoveryPoint(params param.UnmountVmInstanceRecoveryPointParam) error {
	return cli.PostWithRespKey("v1/cdp-backup-storage/unmount-recovery-point", "", params, nil)
}

// DeleteCdpTaskData deletes data associated with a specific CDP task.
func (cli *ZSClient) DeleteCdpTaskData(uuid string) error {
	reqUri := fmt.Sprintf("v1/cdp-task/%s/data", uuid)
	return cli.PostWithRespKey(reqUri, "", map[string]interface{}{}, nil)
}

// CreateVmFromCdpBackup creates a VM from a CDP backup.
func (cli *ZSClient) CreateVmFromCdpBackup(params *param.CreateVmFromCdpBackupParam) (*view.VmInstanceInventoryView, error) {
	var resp view.VmInstanceInventoryView
	if err := cli.Put("v1/cdp-backups/actions", "", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevertVmFromCdpBackup reverts a VM to a CDP backup state.
func (cli *ZSClient) RevertVmFromCdpBackup(uuid string, params *param.RevertVmFromCdpBackupParam) error {
	return cli.PutWithSpec("v1/cdp-backups", uuid, "actions", "", params, nil)
}
