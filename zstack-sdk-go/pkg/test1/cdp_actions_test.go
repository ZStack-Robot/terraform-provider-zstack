// # Copyright (c) ZStack.io, Inc.

package test

import (
	"fmt"
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryCdpPolicy(t *testing.T) {
	t.Log("TestQueryCdpPolicy")
	p := param.NewQueryParam()
	policy, err := accountLoginCli.QueryCdpPolicy(p)
	if err != nil {
		t.Errorf("QueryCdpPolicy error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(policy))
}

func TestPageCdpPolicy(t *testing.T) {
	t.Log("TestPageCdpPolicy")
	p := param.NewQueryParam()
	policy, total, err := accountLoginCli.PageCdpPolicy(p)
	if err != nil {
		t.Errorf("PageCdpPolicy error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(policy))
	t.Log(total)

}

func TestGetCdpPolicy(t *testing.T) {
	t.Log("TestGetCdpPolicy")
	policy, err := accountLoginCli.GetCdpPolicy("63c46a65a1c140f0ad4220c49648d9fe")
	if err != nil {
		t.Errorf("GetCdpPolicy error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(policy))

}

func TestCreateCdpPolicy(t *testing.T) {
	t.Log("TestCreateCdpPolicy")
	p := param.CreateCdpPolicyParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateCdpPolicyDetailParam{
			Name:                    "cdp-xjy-test1",
			Description:             "123123",
			RecoveryPointPerSecond:  3,
			HourlyRpSinceDay:        3,
			DailyRpSinceDay:         7,
			ExpireTimeInDay:         0,
			FullBackupIntervalInDay: 7,
		},
	}
	policy, err := accountLoginCli.CreateCdpPolicy(&p)
	if err != nil {
		t.Errorf("CreateCdpPolicy error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(policy))
}

func TestDeleteCdpPolicy(t *testing.T) {
	t.Log("TestDeleteCdpPolicy")
	err := accountLoginCli.DeleteCdpPolicy("01c52615a94b4f118c9e079b5b992c8c", param.DeleteModePermissive)
	if err != nil {
		t.Errorf("DeleteCdpPolicy error %v", err)
		return
	}
}

func TestUpdateCdpPolicy(t *testing.T) {
	t.Log("TestUpdateCdpPolicy")
	name := "cdp-xjy-test00111"
	p := param.UpdateCdpPolicyParam{
		BaseParam: param.BaseParam{},
		UpdateCdpPolicy: param.UpdateCdpPolicyDetailParam{
			Name:        &name,
			Description: &name,
		},
	}
	policy, err := accountLoginCli.UpdateCdpPolicy("63c46a65a1c140f0ad4220c49648d9fe", &p)
	if err != nil {
		t.Errorf("UpdateCdpPolicy error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(policy))

}

func TestQueryCdpTask(t *testing.T) {
	t.Log("TestQueryCdpTask")
	p := param.NewQueryParam()
	p.AddQ("resourceRefs.resourceUuid=d2f936c5c69c4be28bfa1c85d1e4652d")
	task, err := accountLoginCli.QueryCdpTask(p)
	if err != nil {
		t.Errorf("QueryCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
}

func TestPageCdpTask(t *testing.T) {
	t.Log("TestPageCdpTask")
	p := param.NewQueryParam()
	task, total, err := accountLoginCli.PageCdpTask(p)
	if err != nil {
		t.Errorf("PageCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
	t.Log(total)
}

func TestGetCdpTask(t *testing.T) {
	t.Log("TestGetCdpTask")
	task, err := accountLoginCli.GetCdpTask("73257a869a064310905d8663cca79783")
	if err != nil {
		t.Errorf("GetCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
}

func TestCreateCdpTask(t *testing.T) {
	t.Log("TestCreateCdpTask")
	p := param.CreateCdpTaskParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateCdpTaskDetailParam{
			Name:              "cdp-xjy-test1",
			Description:       "",
			TaskType:          param.CdpTaskTypeVM,
			PolicyUuid:        "63c46a65a1c140f0ad4220c49648d9fe",
			BackupStorageUuid: "b98ac147fe0d449cb426b97898fb87de",
			ResourceUuids:     []string{"dcf4a75b0c864993b4932235da0647ec"},
			BackupBandwidth:   65536000,
			MaxCapacity:       92341796864,
			MaxLatency:        600000,
		},
	}
	task, err := accountLoginCli.CreateCdpTask(&p)
	if err != nil {
		t.Errorf("CreateCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
}

func TestDisableCdpTask(t *testing.T) {
	t.Log("TestDisableCdpTask")
	task, err := accountLoginCli.DisableCdpTask("18374e18558d4b7cb4c5724ef02ed69a")
	if err != nil {
		t.Errorf("DisableCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
}

func TestEnableCdpTask(t *testing.T) {
	t.Log("TestEnableCdpTask")
	task, err := accountLoginCli.EnableCdpTask("18374e18558d4b7cb4c5724ef02ed69a")
	if err != nil {
		t.Errorf("EnableCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))
}

func TestUpdateCdpTask(t *testing.T) {
	t.Log("TestUpdateCdpTask")
	name := "cdp-xjy-test0011啊啊啊1"
	p := param.UpdateCdpTaskParam{
		BaseParam: param.BaseParam{},
		UpdateCdpTask: param.UpdateCdpTaskDetailParam{
			Name:        &name,
			Description: &name,
		},
	}
	task, err := accountLoginCli.UpdateCdpTask("18374e18558d4b7cb4c5724ef02ed69a", &p)
	if err != nil {
		t.Errorf("UpdateCdpTask error %v", err)
		return
	}
	t.Log(jsonutils.Marshal(task))

}

func TestDeleteCdpTask(t *testing.T) {
	t.Log("TestDeleteCdpTask")
	err := accountLoginCli.DeleteCdpTask("18374e18558d4b7cb4c5724ef02ed69a", param.DeleteModePermissive)
	if err != nil {
		t.Errorf("DeleteCdpTask error %v", err)
		return
	}
}

func TestMountVmInstanceRecoveryPoint(t *testing.T) {
	t.Log("TestMountVmInstanceRecoveryPoint")
	pa := param.MountVmInstanceRecoveryPointParam{
		Params: param.MountVmInstanceRecoveryPointDetailParam{
			VmUuid:  "d2f936c5c69c4be28bfa1c85d1e4652d",
			GroupId: 1117,
			Https:   false,
		},
	}
	result, err := accountLoginCli.MountVmInstanceRecoveryPoint(pa)
	if err != nil {
		t.Fatalf("MountVmInstanceRecoveryPoint error %v", err)
	}
	fmt.Println(jsonutils.Marshal(result))
}

func TestUnmountVmInstanceRecoveryPoint(t *testing.T) {
	t.Log("TestUnmountVmInstanceRecoveryPoint")
	pa := param.UnmountVmInstanceRecoveryPointParam{
		Params: param.MountVmInstanceRecoveryPointDetailParam{
			VmUuid:  "d2f936c5c69c4be28bfa1c85d1e4652d",
			GroupId: 1117,
		},
	}
	err := accountLoginCli.UnmountVmInstanceRecoveryPoint(pa)
	if err != nil {
		t.Fatalf("UnmountVmInstanceRecoveryPoint error %v", err)
	}
}

func TestDeleteCdpTaskData(t *testing.T) {
	t.Log("TestDeleteCdpTaskData")
	err := accountLoginCli.DeleteCdpTaskData("7671d3990f1e4ef8bbb9aed1455ec6a3")
	if err != nil {
		t.Fatalf("TestDeleteCdpTaskData error %v", err)
	}
}
