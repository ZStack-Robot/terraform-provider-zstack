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

	// http://172.32.1.214:5000 admin/zstack@2022
	accountLoginHostname        = "172.25.16.104" //基础版-高可用-4.4.24
	accountLoginAccountName     = "admin"         //基础版-高可用-4.4.24
	accountLoginAccountPassword = "password"      //基础版-高可用-4.4.24
	accountLoginMasterHostname  = "172.25.16.104" //基础版-高可用-4.4.24
	accountLoginSlaveHostname   = "172.25.16.104" //基础版-高可用-4.4.24

	// http://172.30.220.215:5000 admin/password
	// accountLoginHostname        = "172.30.220.215" //社区版-4.4.0
	// accountLoginAccountName     = "admin"          //社区版-4.4.0
	// accountLoginAccountPassword = "password"       //社区版-4.4.0

	// http://172.20.20.132:5000 admin/password
	accessKeyAuthHostname        = "172.25.16.104"                            //基础版-4.3.28
	accessKeyAuthAccessKeyId     = "mO6W9gzCQxsfK6OsE7dg"                     //基础版-4.3.28
	accessKeyAuthAccessKeySecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms" //基础版-4.3.28

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
