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

func TestAccZStackImageDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_images" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "6"),

					// Verify the first image to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", "image_for_sg_test"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", "3eff18db19d74568abb5e1ec0759379a"),
				),
			},
		},
	})
}

func TestAccZStackImageDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_images" "test" { name ="Marketplace_zstack_io_servicemonitor_11.2.0_aarch64" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", "Marketplace_zstack_io_servicemonitor_11.2.0_aarch64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", "4c2b64fe5d1b43ccb15410be32227076"),
				),
			},
		},
	})
}

func TestAccZStackImageDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_images" "test" { name_pattern ="Mar%" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", "Marketplace_zstack_io_servicemonitor_11.2.0_aarch64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", "aarch64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guest_os_type", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", "4c2b64fe5d1b43ccb15410be32227076"),
				),
			},
		},
	})
}
