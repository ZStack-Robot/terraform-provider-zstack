// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackInstanceOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_instance_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.#", fmt.Sprintf("%d", len(env.InstanceOfferings))),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", envStr(io, "name")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", envStr(io, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.allocator_strategy", envStr(io, "allocator_strategy")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_num", envStr(io, "cpu_num")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.memory_size", envStr(io, "memory_size")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.sort_key", envStr(io, "sort_key")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", envStr(io, "state")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.type", envStr(io, "type")),
				),
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]
	name := envStr(io, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instance_offers" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", envStr(io, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.allocator_strategy", envStr(io, "allocator_strategy")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.cpu_num", envStr(io, "cpu_num")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.memory_size", envStr(io, "memory_size")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.sort_key", envStr(io, "sort_key")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", envStr(io, "state")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.type", envStr(io, "type")),
				),
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]
	name := envStr(io, "name")
	pattern := name[:3] + "%"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instance_offers" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.uuid", envStr(io, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_instance_offers.test", "instance_offers.0.state", envStr(io, "state")),
				),
			},
		},
	})
}
