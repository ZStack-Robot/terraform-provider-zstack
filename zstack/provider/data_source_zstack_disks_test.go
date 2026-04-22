// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackDiskDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	_ = env

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_disks" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_disks.test", "disks.#"),
					resource.TestCheckResourceAttrSet("data.zstack_disks.test", "disks.0.name"),
					resource.TestCheckResourceAttrSet("data.zstack_disks.test", "disks.0.uuid"),
				),
			},
		},
	})
}

func TestAccZStackDiskDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	item := env.Volumes[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_disks" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_disks.test", "disks.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_disks.test", "disks.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_disks.test", "disks.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}
