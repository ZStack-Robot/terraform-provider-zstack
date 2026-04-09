// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccZStackBackupStorageDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bs := env.BackupStorages[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_backupstorages" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages"), knownvalue.ListSizeExact(len(env.BackupStorages))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(bs, "name"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(bs, "status"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(bs, "state"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(bs, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("total_capacity"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("available_capacity"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccZStackBackupStorageDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bs := env.BackupStorages[0]
	name := envStr(bs, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_backupstorages" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(bs, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(bs, "status"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(bs, "state"))),
				},
			},
		},
	})
}

func TestAccZStackBackupStorageDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bs := env.BackupStorages[0]
	name := envStr(bs, "name")
	pattern := name[:4] + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_backupstorages" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(bs, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(bs, "status"))),
					statecheck.ExpectKnownValue("data.zstack_backupstorages.test", tfjsonpath.New("backup_storages").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(bs, "state"))),
				},
			},
		},
	})
}
