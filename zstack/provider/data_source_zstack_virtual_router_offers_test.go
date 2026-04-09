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

var _ = fmt.Sprintf

func TestAccZStackVirtualRouterOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VirtualRouterOfferings) == 0 {
		t.Skip("no virtual router offerings in env data")
	}
	vro := env.VirtualRouterOfferings[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_virtual_router_offers" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers"), knownvalue.ListSizeExact(len(env.VirtualRouterOfferings))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(vro, "name"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(vro, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("allocator_strategy"), knownvalue.StringExact(envStr(vro, "allocator_strategy"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("cpu_num"), knownvalue.StringExact(envStr(vro, "cpu_num"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("cpu_speed"), knownvalue.StringExact("0")),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("image_uuid"), knownvalue.StringExact(envStr(vro, "image_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("is_default"), knownvalue.StringExact(envStr(vro, "is_default"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("management_network_uuid"), knownvalue.StringExact(envStr(vro, "management_network_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("memory_size"), knownvalue.StringExact(envStr(vro, "memory_size"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("public_network_uuid"), knownvalue.StringExact(envStr(vro, "public_network_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("reserved_memory_size"), knownvalue.StringExact(envStr(vro, "reserved_memory_size"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("sort_key"), knownvalue.StringExact(envStr(vro, "sort_key"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(vro, "state"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(vro, "type"))),
					statecheck.ExpectKnownValue("data.zstack_virtual_router_offers.test", tfjsonpath.New("virtual_router_offers").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(vro, "zone_uuid"))),
				},
			},
		},
	})
}
