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

func TestAccZStackL2NetworkDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_l2networks" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of l2network returned
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.#", "1"),

					// Verify the first l2network to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", "L2"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", "L2VlanNetwork"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", "6de83607f46544e497c84c7eb085b498"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.vlan", "36"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", "d29f4847a99f4dea83bc446c8fe6e64c"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", "ens29f1"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.attached_cluster_uuids.0", "37c25209578c495ca176f60ad0cd97fa"),
				),
			},
		},
	})
}

func TestAccZStackL2NetworkDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_l2networks" "test" { name ="L2" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of l2network returned
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.#", "1"),

					// Verify the first l2network to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", "L2"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", "L2VlanNetwork"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", "6de83607f46544e497c84c7eb085b498"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.vlan", "36"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", "d29f4847a99f4dea83bc446c8fe6e64c"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", "ens29f1"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.attached_cluster_uuids.0", "37c25209578c495ca176f60ad0cd97fa"),
				),
			},
		},
	})
}

func TestAccZStackL2NetworkDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_l2networks" "test" { name_pattern ="L%" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of l2network returned
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.#", "1"),

					// Verify the first l2network to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", "L2"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", "L2VlanNetwork"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", "6de83607f46544e497c84c7eb085b498"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.vlan", "36"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", "d29f4847a99f4dea83bc446c8fe6e64c"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", "ens29f1"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.attached_cluster_uuids.0", "37c25209578c495ca176f60ad0cd97fa"),
				),
			},
		},
	})
}
