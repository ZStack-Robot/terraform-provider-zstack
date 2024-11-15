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

func TestAccZStackmnNodeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_mnnodes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.#", "2"),

					// Verify the first image to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.host_name", "192.168.251.81"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.uuid", "94210bf4590538f3bcd8b6012a1280d5"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.1.host_name", "192.168.251.103"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.1.uuid", "d8d7697f72083481a6a45976b075dcdf"),
				),
			},
		},
	})
}
