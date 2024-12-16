// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
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
