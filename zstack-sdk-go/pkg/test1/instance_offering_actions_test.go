// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestGetInstanceOffering(t *testing.T) {
	t.Log("TestGetInstanceOffering")
	offering, err := accountLoginCli.GetInstanceOffering("4fb8a154b03d418ea771ec74d3273da3")
	if err != nil {
		t.Errorf("TestGetInstanceOffering error %v", err)
		return
	}
	golog.Println(offering)

}

func TestQueryInstanceOffering(t *testing.T) {
	t.Log("TestQueryInstanceOffering")
	params := param.NewQueryParam()
	params.AddQ("name=InstanceOffering-1")
	offering, err := accountLoginCli.QueryInstaceOffering(params)
	if err != nil {
		t.Errorf("TestQueryInstanceOffering error %v", err)
		return
	}
	golog.Println(offering)

}
