// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestGetLicenseInfo(t *testing.T) {
	// 此接口仅能获取到主管理节点的许可证信息
	data, err := accountLoginCli.GetLicenseInfo(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestGetAllManagementNodeLicenseInfo(t *testing.T) {
	accountLoginMasterCli := client.NewZSClient(
		client.DefaultZSConfig(accountLoginMasterHostname).
			LoginAccount(accountLoginAccountName, accountLoginAccountPassword).
			ReadOnly(readOnly).
			Debug(debug),
	)
	_, err := accountLoginMasterCli.Login()
	if err != nil {
		golog.Errorf("TestGetAllManagementNodeLicenseInfo.Master.Login err %v", err)
	}
	defer accountLoginMasterCli.Logout()

	accountLoginSlaveCli := client.NewZSClient(
		client.DefaultZSConfig(accountLoginSlaveHostname).
			LoginAccount(accountLoginAccountName, accountLoginAccountPassword).
			ReadOnly(readOnly).
			Debug(debug),
	)
	_, err = accountLoginSlaveCli.Login()
	if err != nil {
		golog.Errorf("TestGetAllManagementNodeLicenseInfo.Slave.Login err %v", err)
	}
	defer accountLoginSlaveCli.Logout()

	data, err := accountLoginMasterCli.GetLicenseInfo(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("Master ==> %v", jsonutils.Marshal(data))

	data, err = accountLoginSlaveCli.GetLicenseInfo(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("Slave ==> %v", jsonutils.Marshal(data))
}

func TestGetLicenseRecords(t *testing.T) {
	// 获取到的结果仅为主管理节点的许可证历史授权信息
	data, err := accountLoginCli.GetLicenseRecords(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestGetLicenseCapabilities(t *testing.T) {
	data, err := accountLoginCli.GetLicenseCapabilities(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestGetLicenseAddOns(t *testing.T) {
	data, err := accountLoginCli.GetLicenseAddOns(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestDeleteLicense(t *testing.T) {
	// 只能删除主管理节点许可证，无法删除从管理节点许可证
	err := accountLoginCli.DeleteLicense("4fc8ad1e000430d598ab1f89ff83d581", "c3531997757b454c8be84e0601e789a4", "")
	if err != nil {
		golog.Errorf("TestDeleteLicense error %v ", err)
	}
	// 无法删除从管理节点许可证
	err = accountLoginCli.DeleteLicense("da2a1d71761d31d0b100d95d2074ac61", "2a03cea1dbee40faafa0756244afb584", "")
	if err != nil {
		golog.Errorf("TestDeleteLicense error %v ", err)
	}
}

func TestReloadLicense(t *testing.T) {
	// 重新加载指定managementNodeUuids的管理节点申请码及许可证信息(指定1mn则刷该mn，指定多mn则刷多mn)
	// 返回的结果仅为主管理节点的许可证历史授权信息
	params := param.ReloadLicenseParam{
		ReloadLicense: param.ReloadLicenseDetailParam{
			ManagementNodeUuids: []string{"b1c9d0bca6d73d57a69cc27d2686eb90", "77ced08b52c034aaa784f2c516c4d094"},
		},
	}
	data, err := accountLoginCli.ReloadLicense(params)
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestUpdateLicense(t *testing.T) {
	// 传递集群所有管理节点许可证文件(轴心申请通过后的许可证压缩包做base64)，则可全部管理节点都更新
	params := param.UpdateLicenseParam{
		UpdateLicense: param.UpdateLicenseDetailParam{
			License: "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",
		},
	}
	data, err := accountLoginCli.UpdateLicense("4fc8ad1e000430d598ab1f89ff83d581", params)
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}
