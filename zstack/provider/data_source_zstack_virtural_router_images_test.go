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

func TestAccZStackVirtualRouterImagesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_virtual_router_images" "test" {}`,
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

func TestAccZStackVirtualRouterImagesDataSourceFilterByName(t *testing.T) {
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

func TestAccZStackVirtualRouterImagesDataSourceFilterByNamePattern(t *testing.T) {
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
