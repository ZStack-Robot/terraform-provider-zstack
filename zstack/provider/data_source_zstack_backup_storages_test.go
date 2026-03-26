// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackBackupStorageDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.BackupStorages) == 0 {
		t.Skip("no backup storages in env data")
	}
	bs := env.BackupStorages[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_backupstorages" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.#", fmt.Sprintf("%d", len(env.BackupStorages))),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", envStr(bs, "name")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", envStr(bs, "status")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", envStr(bs, "state")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", envStr(bs, "uuid")),
					resource.TestCheckResourceAttrSet("data.zstack_backupstorages.test", "backup_storages.0.total_capacity"),
					resource.TestCheckResourceAttrSet("data.zstack_backupstorages.test", "backup_storages.0.available_capacity"),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_backupstorages" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", envStr(bs, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", envStr(bs, "status")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", envStr(bs, "state")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_backupstorages" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.uuid", envStr(bs, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.status", envStr(bs, "status")),
					resource.TestCheckResourceAttr("data.zstack_backupstorages.test", "backup_storages.0.state", envStr(bs, "state")),
				),
			},
		},
	})
}
