// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/client"
)

const (
	// The ZStack Cloud Community Edition only supports login authentication for super admin accounts.
	// The ZStack Cloud Basic Edition supports login authentication for AccessKey, super admin, and sub-accounts.
	// The ZStack Cloud Enterprise Edition supports login authentication for AccessKey, super admin, sub-accounts, and enterprise users.

	accountLoginHostname        = "172.30.3.3" //ZStack Cloud API endpoint IP address
	accountLoginAccountName     = "admin"
	accountLoginAccountPassword = "password"

	accountLoginMasterHostname = "IPOfCloudAPIEndpoint"
	accountLoginSlaveHostname  = "IPOfCloudAPIEndpoint"

	accessKeyAuthHostname        = "IPOfCloudAPIEndpoint"
	accessKeyAuthAccessKeyId     = "AccessKeyId"
	accessKeyAuthAccessKeySecret = "AccessKeySecret"

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

func TestMain(m *testing.M) {
	_, err := accountLoginCli.Login()
	if err != nil {
		golog.Errorf("TestMain err %v", err)
	}
	defer accountLoginCli.Logout()

	m.Run()
}
