// Copyright (c) ZStack.io, Inc.

package client

import (
	"errors"
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// GetLicenseInfo Retrieves license information
// This interface only retrieves the application code and license information for the primary management node
func (cli *ZSClient) GetLicenseInfo(params param.QueryParam) (view.LicenseInventoryView, error) {
	var resp view.LicenseInventoryView
	return resp, cli.GetWithSpec("v1/licenses", "", "", responseKeyInventory, nil, &resp)
}

// GetLicenseRecords Retrieves historical license authorization information
// The results only include historical license authorization information for the primary management node
func (cli *ZSClient) GetLicenseRecords(params param.QueryParam) ([]view.LicenseInventoryView, error) {
	var resp []view.LicenseInventoryView
	return resp, cli.List("v1/licenses/records", &params, &resp)
}

// GetLicenseCapabilities Retrieves license capacity
func (cli *ZSClient) GetLicenseCapabilities(params param.QueryParam) (map[string]interface{}, error) {
	var resp map[string]interface{}
	return resp, cli.GetWithSpec("v1/licenses/capabilities", "", "", "capabilities", nil, &resp)
}

// GetLicenseAddOns Retrieves add-on license information
func (cli *ZSClient) GetLicenseAddOns(params param.QueryParam) ([]view.LicenseAddOnInventoryView, error) {
	var resp []view.LicenseAddOnInventoryView
	return resp, cli.ListWithRespKey("v1/licenses/addons", "addons", &params, &resp)
}

// DeleteLicense Deletes a license file
// Can only delete the license for the primary management node, cannot delete licenses for secondary management nodes
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

// ReloadLicense Reloads the license
// Reloads the application code and license information for the specified management node UUIDs (reload a single mn, multiple mns, or all mns)
// The result will only include historical license authorization information for the primary management node
func (cli *ZSClient) ReloadLicense(params param.ReloadLicenseParam) (view.LicenseInventoryView, error) {
	var resp view.LicenseInventoryView
	return resp, cli.Put("v1/licenses/actions", "", &params, &resp)
}

// UpdateLicense Updates license information
// When the license file for all management nodes in the cluster is provided, it will update all management nodes
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
