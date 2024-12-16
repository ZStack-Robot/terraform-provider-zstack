// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryRecoveryPoint(t *testing.T) {
	t.Log("TestQueryRecoveryPoint")

	p := param.QueryRecoveryPointParam{
		//StartTime: "2023-06-06T00:00:00+08:00",
		//EndTime:   "2023-06-06T23:59:59+08:00",
		Scale: param.RecoveryPointScaleHour,
	}

	result, err := accountLoginCli.QueryRecoveryPoint("e6f6c2f3dcd746edaf26f3c036ac7f76", p)
	if err != nil {
		t.Fatalf("QueryRecoveryPoint error %v", err)
	}
	t.Log(jsonutils.Marshal(result))
}

func TestQueryProtectRecoveryPoint(t *testing.T) {
	t.Log("TestQueryProtectRecoveryPoint")
	p := param.QueryProtectRecoveryPointParam{}
	result, err := accountLoginCli.QueryProtectRecoveryPoint("e6f6c2f3dcd746edaf26f3c036ac7f76", p)
	if err != nil {
		t.Fatalf("TestQueryProtectRecoveryPoint error %v", err)
	}
	t.Log(jsonutils.Marshal(result))
}

func TestProtectVmInstanceRecoveryPoint(t *testing.T) {
	t.Log("TestProtectVmInstanceRecoveryPoint")
	p := param.ProtectVmInstanceRecoveryPointParam{
		ProtectVmInstanceRecoveryPoint: param.ProtectVmInstanceRecoveryDetailPointParam{
			GroupId:     351592,
			Description: "protect from api test",
		},
	}
	err := accountLoginCli.ProtectVmInstanceRecoveryPoint("9931279886da463db7aa38c6b34b99f5", p)
	if err != nil {
		t.Fatalf("TestProtectVmInstanceRecoveryPoint error %v", err)
	}
}

func TestUnprotectVmInstanceRecoveryPoint(t *testing.T) {
	t.Log("TestUnprotectVmInstanceRecoveryPoint")
	p := param.UnprotectVmInstanceRecoveryPointParam{
		UnprotectVmInstanceRecoveryPoint: param.UnprotectVmInstanceRecoveryPointDetailParam{
			GroupId: 351592,
		},
	}
	err := accountLoginCli.UnprotectVmInstanceRecoveryPoint("9931279886da463db7aa38c6b34b99f5", p)
	if err != nil {
		t.Fatalf("TestUnprotectVmInstanceRecoveryPoint error %v", err)
	}
}
