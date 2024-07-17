// Copyright (c) HashiCorp, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/client"
)

const (
	//ZStack Cloud社区版仅支持超管账户登录认证
	//ZStack Cloud基础版支持AccessKey、超管及子账户登录认证
	//ZStack Cloud企业版支持AccessKey、超管及子账户登录认证、企业用户登录认证

	accountLoginHostname        = "IPOfCloudAPIEndpoint" //基础版-高可用-4.4.24
	accountLoginAccountName     = "admin"                //基础版-高可用-4.4.24
	accountLoginAccountPassword = "password"             //基础版-高可用-4.4.24
	accountLoginMasterHostname  = "IPOfCloudAPIEndpoint" //基础版-高可用-4.4.24
	accountLoginSlaveHostname   = "IPOfCloudAPIEndpoint" //基础版-高可用-4.4.24

	accessKeyAuthHostname        = "IPOfCloudAPIEndpoint" //基础版-4.3.28
	accessKeyAuthAccessKeyId     = "AccessKeyId"          //基础版-4.3.28
	accessKeyAuthAccessKeySecret = "AccessKeySecret"      //基础版-4.3.28

	// userLoginHostname            = "" //企业版
	// userLoginAccountName         = "" //企业版
	// userLoginAccountUserName     = "" //企业版
	// userLoginAccountUserPassword = "" //企业版

	readOnly = false
	debug    = false
)

var accountLoginCli = client.NewZSClient(
	client.DefaultZSConfig(accountLoginHostname).
		LoginAccount(accountLoginAccountName, accountLoginAccountPassword).
		ReadOnly(readOnly).
		Debug(true),
)

var accessKeyAuthCli = client.NewZSClient(
	client.DefaultZSConfig(accessKeyAuthHostname).
		AccessKey(accessKeyAuthAccessKeyId, accessKeyAuthAccessKeySecret).
		ReadOnly(readOnly).
		Debug(debug),
)

// var userLoginCli = client.NewZSClient(
// 	client.DefaultZSConfig(accountLoginHostname).
// 		LoginAccountUser(userLoginAccountName, userLoginAccountUserName, userLoginAccountUserPassword).
// 		ReadOnly(readOnly).
// 		Debug(debug),
// )

func TestMain(m *testing.M) {
	_, err := accountLoginCli.Login()
	if err != nil {
		golog.Errorf("TestMain err %v", err)
	}
	defer accountLoginCli.Logout()

	m.Run()
}
