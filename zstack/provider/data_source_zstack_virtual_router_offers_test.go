// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVirtualRouterOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VirtualRouterOfferings) == 0 {
		t.Skip("no virtual router offerings in env data")
	}
	vro := env.VirtualRouterOfferings[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_router_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.#", fmt.Sprintf("%d", len(env.VirtualRouterOfferings))),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.name", envStr(vro, "name")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.uuid", envStr(vro, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.allocator_strategy", envStr(vro, "allocator_strategy")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.cpu_num", envStr(vro, "cpu_num")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.cpu_speed", "0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.image_uuid", envStr(vro, "image_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.is_default", envStr(vro, "is_default")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.management_network_uuid", envStr(vro, "management_network_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.memory_size", envStr(vro, "memory_size")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.public_network_uuid", envStr(vro, "public_network_uuid")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.reserved_memory_size", envStr(vro, "reserved_memory_size")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.sort_key", envStr(vro, "sort_key")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.state", envStr(vro, "state")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.type", envStr(vro, "type")),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.zone_uuid", envStr(vro, "zone_uuid")),
				),
			},
		},
	})
}
