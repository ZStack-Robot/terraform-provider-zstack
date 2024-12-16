// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryManagementNode(t *testing.T) {
	data, err := accountLoginCli.QueryManagementNode(param.NewQueryParam())
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))
}

func TestGetVersion(t *testing.T) {
	data, err := accountLoginCli.GetVersion()
	if err != nil {
		t.Errorf("error %v ", err)
	}
	golog.Infof("%v", data)
}
