// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackVirtualRouterOffersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_virtual_router_offers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.#", "1"),

					// Verify the first virtual router instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.name", "vr"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.uuid", "9b56381daf204657adba1a4f0ec7ef39"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.allocator_strategy", "LeastVmPreferredHostAllocatorStrategy"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.cpu_num", "1"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.cpu_speed", "0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.image_uuid", "93005c8a2a314a489635eca8c30794d4"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.is_default", "false"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.management_network_uuid", "50e8c0d69681447fbe347c8dae2b1bef"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.memory_size", "1073741824"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.public_network_uuid", "50e8c0d69681447fbe347c8dae2b1bef"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.reserved_memory_size", "0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.sort_key", "0"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.type", "VirtualRouter"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_offers.test", "virtual_router_offers.0.zone_uuid", "d29f4847a99f4dea83bc446c8fe6e64c"),
				),
			},
		},
	})
}

func TestAccZStackVirtualRouterOffersDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_virtual_router_images" "test" { name = "fd"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.#", "1"),

					// Verify the first virtual router instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.name", "fd"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.uuid", "f0fe97d9b1a649f787927b86808b739a"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.guest_os_type", "openEuler 22.03"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.status", "Ready"),
				),
			},
		},
	})
}

func TestAccZStackVirtualRouterOffersDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_virtual_router_images" "test" { name_pattern = "%"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.#", "2"),

					// Verify the first virtual router instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.name", "fd"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.uuid", "f0fe97d9b1a649f787927b86808b739a"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.guest_os_type", "openEuler 22.03"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_virtual_router_images.test", "images.0.status", "Ready"),
				),
			},
		},
	})
}
