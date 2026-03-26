// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackDiskOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_disk_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.#", fmt.Sprintf("%d", len(env.DiskOfferings))),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", envStr(do, "name")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", envStr(do, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.allocator_strategy", envStr(do, "allocator_strategy")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.description", envStr(do, "description")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.disk_size", envStr(do, "disk_size")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.sort_key", envStr(do, "sort_key")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.type", envStr(do, "type")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", envStr(do, "state")),
				),
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]
	name := envStr(do, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_disk_offers" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", envStr(do, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.allocator_strategy", envStr(do, "allocator_strategy")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.description", envStr(do, "description")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.disk_size", envStr(do, "disk_size")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.sort_key", envStr(do, "sort_key")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.type", envStr(do, "type")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", envStr(do, "state")),
				),
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]
	name := envStr(do, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_disk_offers" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.uuid", envStr(do, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_disk_offers.test", "disk_offers.0.state", envStr(do, "state")),
				),
			},
		},
	})
}
