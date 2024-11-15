// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Set environment variable TF_ACC to run acceptance tests.
// In VSCode, edit settings.json and add:
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackL3NetworksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "zstack_l3network" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the number of l3networks returned
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.#", "2"),

					// Verify attributes of the first l3network
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.name", "pvcl3"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.category", "Private"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.dns.0.dnsmodel", "192.166.255.254"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.cidr", "192.166.255.0/24"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.startip", "192.166.255.2"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.endip", "192.166.255.62"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.gateway", "192.166.255.254"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.iprangename", "192.166.255.2-192.166.255.62"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.uuid", "c9c2451b7fe1457699abd98a182df95b"),
				),
			},
		},
	})
}

func TestAccZStackL3NetworksDataSourceFilterByNameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "zstack_l3network" "test" {
						name_regex = "public"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the number of l3networks returned
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.#", "1"),

					// Verify attributes of the filtered l3network
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.name", "public"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.category", "Public"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.dns.0.dnsmodel", "223.5.5.5"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.cidr", "172.25.0.0/16"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.startip", "172.25.126.192"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.endip", "172.25.126.223"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.gateway", "172.25.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.iprange.0.iprangename", "172.25.126.192-172.25.126.223"),
					resource.TestCheckResourceAttr("data.zstack_l3network.test", "l3networks.0.uuid", "de7f26a7304d45aea9e9871a1ba7dbae"),
				),
			},
		},
	})
}
