// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQueryVpcRouter(t *testing.T) {
	router, err := accountLoginCli.QueryVpcRouter(param.NewQueryParam())
	if err != nil {
		golog.Errorf("TestQueryVpcRouter error %v", err)
	}
	golog.Println(router)
}

func TestGetVpcRouter(t *testing.T) {
	router, err := accountLoginCli.GetVpcRouter("aaaa")
	if err != nil {
		golog.Errorf("TestGetVpcRouter error %v", err)
	}
	golog.Println(router)
}
