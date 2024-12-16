// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryHost(t *testing.T) {
	host, err := accountLoginCli.QueryHost(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryHost %v", err)
	}

	golog.Info(jsonutils.Marshal(host))
}

func TestQueryHostNetworkBonding(t *testing.T) {
	host, err := accountLoginCli.QueryHostNetworkBonding(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryHostNetworkBonding %v", err)
	}

	golog.Info(jsonutils.Marshal(host))
}

func TestQueryHostNetworkInterface(t *testing.T) {
	host, err := accountLoginCli.QueryHostNetworkInterface(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryHostNetworkInterface %v", err)
	}

	golog.Info(jsonutils.Marshal(host))
}

func TestPageHost(t *testing.T) {
	host, num, err := accountLoginCli.PageHost(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestPageHost %v", err)
	}

	golog.Info(jsonutils.Marshal(host), num)
}

func TestGetHost(t *testing.T) {
	host, err := accountLoginCli.GetHost("4f8a562c919041e6980dc9f523ee0e8e")
	if err != nil {
		golog.Errorf("TestGetHost %v", err)
	}

	golog.Info(jsonutils.Marshal(host))

}

func TestChangeHostState(t *testing.T) {
	host, err := accountLoginCli.ChangeHostState("5d08acb763a448b9ab5bfe64753f2f7f", &param.ChangeHostStateParam{
		BaseParam:       param.BaseParam{},
		ChangeHostState: param.ChangeHostStateDetailParam{StateEvent: param.StateEventEnable},
	})
	if err != nil {
		golog.Errorf("TestChangeHostState %v", err)
	}

	golog.Info(jsonutils.Marshal(host))
}

func TestReconnectHost(t *testing.T) {
	host, err := accountLoginCli.ReconnectHost("5d08acb763a448b9ab5bfe64753f2f7f")
	if err != nil {
		golog.Errorf("TestReconnectHost %v", err)
	}

	golog.Info(jsonutils.Marshal(host))
}
