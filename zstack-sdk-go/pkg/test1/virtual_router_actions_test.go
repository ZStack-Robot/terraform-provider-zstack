// Copyright (c) HashiCorp, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQueryVirtualRouterVm(t *testing.T) {
	vm, err := accountLoginCli.QueryVirtualRouterVm(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryVirtualRouterVm error %v", err)
		return
	}
	golog.Println(vm)
}

func TestGetVirtualRouterVm(t *testing.T) {
	vm, err := accountLoginCli.GetVirtualRouterVm("")
	if err != nil {
		golog.Errorf("TestGetVirtualRouterVm error %v", err)
		return
	}
	golog.Println(vm)
}
