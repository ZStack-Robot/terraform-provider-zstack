// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package test

import (
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestCreateVmInstanceScript(t *testing.T) {

	instanceScripts, err := accountLoginCli.CreateVmInstanceScript(param.CreateVmInstanceScriptParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateVmInstanceScriptDetailParam{
			Name:        "testScript",
			Description: "testScriptDesc",
			ScriptContent: `#!/bin/bash

# 创建实际升级脚本
cat << 'EOF' > /tmp/do-upgrade-qga.sh
#!/bin/bash
set +e
mkdir -p /mnt/cdrom
mount /dev/cdrom /mnt/cdrom || true
if [ -f /mnt/cdrom/zs-tools-install.sh ]; then
    cd /mnt/cdrom/
   echo Y | bash ./zs-tools-install.sh
    echo $? > /tmp/zs-tools-exit-code
    cd ~
else
    echo "zs-tools-install.sh not found" > /tmp/zs-tools-error.log
    echo 1 > /tmp/zs-tools-exit-code
fi
umount /mnt/cdrom || true
EOF

chmod +x /tmp/do-upgrade-qga.sh

# 延迟 60 秒后执行实际脚本，避免当前 QGA session 被重启
nohup bash -c "sleep 60 && bash /tmp/do-upgrade-qga.sh" > /tmp/zs-tools-install.log 2>&1 &

echo "QGA upgrade will be executed in background after 60 seconds"
exit 0
			`,
			ScriptType:    "Shell",
			ScriptTimeout: 60,
			Platform:      "Linux",
		},
	})
	if err != nil {
		t.Errorf("TestCreateVmInstance %v", err)
	}
	golog.Println(instanceScripts)
}

func TestDeleteVmInstanceScript(t *testing.T) {
	scriptUuid := "1b396c2f19774c1fb2b2285edebe71a2"
	err := accountLoginCli.DeleteVmInstanceScrpt(scriptUuid, param.DeleteModePermissive)
	if err != nil {
		t.Errorf("TestDeleteVmInstanceScript %v", err)
	}
	golog.Println("DeleteVmInstanceScript success")
}

func TestGetVmInstanceScript(t *testing.T) {
	scriptUuid := "159051a78fbb49d3b72436a0de6c3c2b"
	instanceScripts, err := accountLoginCli.GetVmInstanceScript(scriptUuid)
	if err != nil {
		t.Errorf("TestGetVmInstanceScript %v", err)
	}
	golog.Println(instanceScripts)
}

func TestQueryVmInstanceScript(t *testing.T) {

	instanceScripts, err := accountLoginCli.QueryVmInstanceScript(param.NewQueryParam())
	if err != nil {
		t.Errorf("TestQueryVmInstanceScript %v", err)
	}
	golog.Println(instanceScripts)
}

func TestUpdateVmInstanceScript(t *testing.T) {
	scriptUuid := "8f9ca7f7e74647708ec806882a24b89e"
	err := accountLoginCli.UpdateVmInstanceScript(scriptUuid, param.UpdateVmInstanceScriptParam{
		BaseParam: param.BaseParam{},
		Params: param.UpdateVmInstanceScriptDetailParam{
			Name:        "testScript",
			Description: "testScriptDesc",
			ScriptContent: `
#!/bin/bash
echo hello world
`,
			ScriptType:    "Shell",
			ScriptTimeout: 60,
			Platform:      "Linux",
		},
	})
	if err != nil {
		t.Errorf("TestUpdateVmInstanceScript %v", err)
	}
	golog.Println("UpdateVmInstanceScript success")
}

func TestExecuteVmInstanceScript(t *testing.T) {
	scriptUuid := "8f9ca7f7e74647708ec806882a24b89e"
	instanceScripts, err := accountLoginCli.ExecuteVmInstanceScript(scriptUuid, param.ExecuteVmInstanceScriptParam{
		BaseParam: param.BaseParam{},
		Params: param.ExecuteVmInstanceScriptDetailParam{
			VmInstanceUuids: []string{"26008c4f1e794c08bcfc3b270fab3e83", "28a0a09073d944d88a4aaaf59ce2ef60"},
			ScriptTimeout:   60,
		},
	})
	if err != nil {
		t.Errorf("TestExecuteVmInstanceScript %v", err)
	}
	golog.Println(instanceScripts)
}

func TestGetVmInstanceScriptExecutedRecord(t *testing.T) {
	scriptUuid := "7d7a02df12a74a04800b448458bff44e"
	instanceScripts, err := accountLoginCli.GetVmInstanceScriptExecutedRecord(scriptUuid)
	if err != nil {
		t.Errorf("TestGetVmInstanceScriptExecutedRecord %v", err)
	}
	golog.Println(instanceScripts)
}

func TestQueryVmInstanceScriptExecutedRecord(t *testing.T) {
	instanceScripts, err := accountLoginCli.QueryVmInstanceScriptExecutedRecord(param.NewQueryParam())
	if err != nil {
		t.Errorf("TestQueryVmInstanceScriptExecutedRecord %v", err)
	}
	golog.Println(instanceScripts)
}

func TestQueryGuestVmScriptExecutedRecordDetail(t *testing.T) {
	instanceScripts, err := accountLoginCli.QueryGuestVmScriptExecutedRecordDetail(param.NewQueryParam())
	if err != nil {
		t.Errorf("TestQueryGuestVmScriptExecutedRecordDetail %v", err)
	}
	golog.Println(instanceScripts)
}
