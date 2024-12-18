// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestGetDiskOffering(t *testing.T) {
	t.Log("TestGetDiskOffering")
	offering, err := accountLoginCli.GetDiskOffering("be0d064c6cc34fc9a70b4ae72c1c89ef")
	if err != nil {
		t.Errorf("TestGetDiskOffering error %v", err)
		return
	}
	golog.Println(offering)

}

func TestQueryDiskOffering(t *testing.T) {
	t.Log("TestQueryDiskOffering")
	params := param.NewQueryParam()
	params.AddQ("name=mediumDiskOffering")
	offering, err := accountLoginCli.QueryDiskOffering(params)
	if err != nil {
		t.Errorf("TestQueryInstanceOffering error %v", err)
		return
	}
	golog.Println(offering)

}

func TestCreateDiskOffering(t *testing.T) {
	t.Log("TestCreateDiskOffering")
	params := param.NewQueryParam()
	params.AddQ("name=InstanceOffering-1")
	offering, err := accountLoginCli.QueryDiskOffering(params)
	if err != nil {
		t.Errorf("TestQueryInstanceOffering error %v", err)
		return
	}
	golog.Println(offering)

}

func TestDeleteDiskOffering(t *testing.T) {
	t.Log("TestQueryDiskOffering")
	params := param.NewQueryParam()
	params.AddQ("name=InstanceOffering-1")
	offering, err := accountLoginCli.QueryDiskOffering(params)
	if err != nil {
		t.Errorf("TestQueryInstanceOffering error %v", err)
		return
	}
	golog.Println(offering)

}
