// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestRequestConsoleAccess(t *testing.T) {
	params := param.RequestConsoleAccessParam{
		Params: param.RequestConsoleAccessDetailParam{
			VMInstanceUUID: "10c9a1aa52f5404899586ff6bd484685",
		},
	}
	access, err := accountLoginCli.RequestConsoleAccess(params)
	if err != nil {
		t.Errorf("RequestConsoleAccess error: %v", err)
	}
	golog.Info(access)
}
