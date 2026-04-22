// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackPrimaryStorageDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.PrimaryStorages) == 0 {
		t.Skip("no primary_storages in env data")
	}
	item := env.PrimaryStorages[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_primary_storages" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.#", fmt.Sprintf("%d", len(env.PrimaryStorages))),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.name", envStr(item, "name")),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.uuid", envStr(item, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.state", envStr(item, "state")),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.status", envStr(item, "status")),
				),
			},
		},
	})
}

func TestAccZStackPrimaryStorageDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.PrimaryStorages) == 0 {
		t.Skip("no primary_storages in env data")
	}
	item := env.PrimaryStorages[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_primary_storages" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.uuid", envStr(item, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.state", envStr(item, "state")),
					resource.TestCheckResourceAttr("data.zstack_primary_storages.test", "primary_storages.0.status", envStr(item, "status")),
				),
			},
		},
	})
}
