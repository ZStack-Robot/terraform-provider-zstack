// Copyright (c) ZStack.io, Inc.

package client

import (
	"fmt"

	"github.com/pkg/errors"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

const RecoverPotListRespKey = "recoveryPoints"
const ProtectRecoverPotListRespKey = "recoveryPoints"

// DescribeVmInstanceRecoveryPoint Get recovery point information for a VM instance
func (cli *ZSClient) DescribeVmInstanceRecoveryPoint(uuid string, params *param.GetRecoveryPointParam) (*view.PointResourceInfoView, error) {
	if len(uuid) == 0 {
		return nil, fmt.Errorf("uuid should not empty")
	}
	var resp view.PointResourceInfoView
	reqUri := fmt.Sprintf("v1/vm-instances/%s/recovery-point", uuid)
	if err := cli.GetWithRespKey(reqUri, "", "", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryRecoveryPoint Query recovery points for a VM instance
func (cli *ZSClient) QueryRecoveryPoint(uuid string, params param.QueryRecoveryPointParam) (*view.RecoveryPointRespView, error) {
	if len(uuid) == 0 {
		return nil, fmt.Errorf("uuid should not empty")
	}
	resp := view.RecoveryPointRespView{}
	reqUri := fmt.Sprintf("v1/vm-instances/%s/recovery-points", uuid)
	if err := cli.GetWithRespKey(reqUri, "", RecoverPotListRespKey, &params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryProtectRecoveryPoint Query protected recovery points for a VM instance
func (cli *ZSClient) QueryProtectRecoveryPoint(uuid string, params param.QueryProtectRecoveryPointParam) (*view.ProtectRecoveryPointRespView, error) {
	if len(uuid) == 0 {
		return nil, errors.New("uuid should not empty")
	}
	resp := view.ProtectRecoveryPointRespView{}
	reqUri := fmt.Sprintf("v1/vm-instances/%s/protected-recovery-points", uuid)
	if err := cli.GetWithRespKey(reqUri, "", ProtectRecoverPotListRespKey, &params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ProtectVmInstanceRecoveryPoint Set protected recovery point for a VM instance
func (cli *ZSClient) ProtectVmInstanceRecoveryPoint(uuid string, params param.ProtectVmInstanceRecoveryPointParam) error {
	if len(uuid) == 0 {
		return errors.New("uuid should not empty")
	}
	return cli.PutWithSpec("v1/vm-instances", uuid, "protect-recovery-point", "", &params, nil)
}

// UnprotectVmInstanceRecoveryPoint Remove protected recovery point for a VM instance
func (cli *ZSClient) UnprotectVmInstanceRecoveryPoint(uuid string, params param.UnprotectVmInstanceRecoveryPointParam) error {
	if len(uuid) == 0 {
		return errors.New("uuid should not empty")
	}
	return cli.PutWithSpec("v1/vm-instances", uuid, "unprotect-recovery-point", "", &params, nil)
}
