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

func TestAccZStackZonesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_zones" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of image returned
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.#", "1"),

					// Verify the first image to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.name", "zone1"),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.uuid", "1df948fedd3b45dd89e9549348280e17"),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.type", "zstack"),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.state", "Enabled"),
				),
			},
		},
	})
}
