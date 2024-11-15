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
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "8"),

					// Verify the first image to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", "CentOS7.9"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guestostype", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", "Disabled"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", "18e898cf4afa472cadec9d767bbace22"),
				),
			},
		},
	})
}

func TestAccZStackImageDataSourceFilterByname_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_images" "test" { name_regex="RDS-3.13.10" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of clusters returned
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.#", "1"),

					// Verify the first Cluster to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.name", "RDS-3.13.10"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.format", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.guestostype", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.status", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_images.test", "images.0.uuid", "85abc3c6cfcd4bf683444fd022fc2c86"),
				),
			},
		},
	})
}
