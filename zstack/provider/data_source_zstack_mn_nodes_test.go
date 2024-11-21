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
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.host_name", "172.26.107.217"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.uuid", "5f46f797e5ed3ab6950b5e55e786c7af"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.1.host_name", "172.26.100.198"),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.1.uuid", "932ba6b5d8fb3506b7c23506e7317942"),
				),
			},
		},
	})
}
