// Copyright (c) HashiCorp, Inc.

package client

import (
	"errors"
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// GetLicenseInfo 获取许可证信息
// 此接口仅能获取到主管理节点的申请码及许可证信息
func (cli *ZSClient) GetLicenseInfo(params param.QueryParam) (view.LicenseInventoryView, error) {
	var resp view.LicenseInventoryView
	return resp, cli.GetWithSpec("v1/licenses", "", "", responseKeyInventory, nil, &resp)
}

// GetLicenseRecords 获取许可证历史授权信息
// 获取到的结果仅为主管理节点的许可证历史授权信息
func (cli *ZSClient) GetLicenseRecords(params param.QueryParam) ([]view.LicenseInventoryView, error) {
	var resp []view.LicenseInventoryView
	return resp, cli.List("v1/licenses/records", &params, &resp)
}

// GetLicenseCapabilities 获取许可证容量
func (cli *ZSClient) GetLicenseCapabilities(params param.QueryParam) (map[string]interface{}, error) {
	var resp map[string]interface{}
	return resp, cli.GetWithSpec("v1/licenses/capabilities", "", "", "capabilities", nil, &resp)
}

// GetLicenseAddOns 获取附加功能许可证信息
func (cli *ZSClient) GetLicenseAddOns(params param.QueryParam) ([]view.LicenseAddOnInventoryView, error) {
	var resp []view.LicenseAddOnInventoryView
	return resp, cli.ListWithRespKey("v1/licenses/addons", "addons", &params, &resp)
}

// DeleteLicense 删除许可证文件
// 只能删除主管理节点许可证，无法删除从管理节点许可证
func (cli *ZSClient) DeleteLicense(managementNodeUuid, uuid, module string) error {
	if managementNodeUuid == "" || (uuid == "" && module == "") {
		return errors.New("params error")
	}

	paramsStr := ""
	if uuid != "" {
		paramsStr = fmt.Sprintf("uuid=%s", uuid)
	} else if module != "" {
		paramsStr = fmt.Sprintf("module=%s", module)
	}

	return cli.DeleteWithSpec("v1/licenses/mn", managementNodeUuid, "actions", paramsStr, nil)
}

// ReloadLicense 重新加载许可证
// 重新加载指定managementNodeUuids的管理节点申请码及许可证信息(指定1mn则仅刷该mn，指定多mn则刷多mn，指定所有mn则刷所有mn)
// 返回的结果仅为主管理节点的许可证历史授权信息
func (cli *ZSClient) ReloadLicense(params param.ReloadLicenseParam) (view.LicenseInventoryView, error) {
	var resp view.LicenseInventoryView
	return resp, cli.Put("v1/licenses/actions", "", &params, &resp)
}

// UpdateLicense 更新许可证信息
// 传递集群所有管理节点许可证文件，则可全部管理节点都更新
func (cli *ZSClient) UpdateLicense(managementNodeUuid string, params param.UpdateLicenseParam) (*view.LicenseInventoryView, error) {
	if managementNodeUuid == "" {
		return nil, errors.New("params error")
	}

	var resp view.LicenseInventoryView
	if err := cli.PutWithSpec("v1/licenses/mn", managementNodeUuid, "actions", responseKeyInventory, &params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
