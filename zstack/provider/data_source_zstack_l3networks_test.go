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
					data "zstack_l3networks" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the number of l3networks returned
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.#", "1"),

					// Verify attributes of the first l3network
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", "public network"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", "a5e77b2972e64316878993af7a695400"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", "Public"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.dns.0.dns_model", "223.5.5.5"),

					//Verify attribute free ips
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip", "172.26.111.241"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip_range_uuid", "8caf96a048f0423f9178e54d29b36b86"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.netmask", "255.255.0.0"),

					//Verify attribute ip range
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.cidr", "172.26.0.0/16"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.start_ip", "172.26.111.240"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.end_ip", "172.26.111.254"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.ip_range_name", "host IP range"),
				),
			},
		},
	})
}

func TestAccZStackL3NetworksDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "zstack_l3networks" "test" {
						name = "public network"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the number of l3networks returned
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.#", "1"),

					// Verify attributes of the first l3network
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", "public network"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", "a5e77b2972e64316878993af7a695400"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", "Public"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.dns.0.dns_model", "223.5.5.5"),

					//Verify attribute free ips
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip", "172.26.111.241"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip_range_uuid", "8caf96a048f0423f9178e54d29b36b86"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.netmask", "255.255.0.0"),

					//Verify attribute ip range
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.cidr", "172.26.0.0/16"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.start_ip", "172.26.111.240"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.end_ip", "172.26.111.254"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.ip_range_name", "host IP range"),
				),
			},
		},
	})
}

func TestAccZStackL3NetworksDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
					data "zstack_l3networks" "test" {
						name_pattern = "p%"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the number of l3networks returned
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.#", "1"),

					// Verify attributes of the first l3network
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", "public network"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", "a5e77b2972e64316878993af7a695400"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", "Public"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.dns.0.dns_model", "223.5.5.5"),

					//Verify attribute free ips
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip", "172.26.111.241"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.ip_range_uuid", "8caf96a048f0423f9178e54d29b36b86"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.free_ips.0.netmask", "255.255.0.0"),

					//Verify attribute ip range
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.cidr", "172.26.0.0/16"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.start_ip", "172.26.111.240"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.end_ip", "172.26.111.254"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.gateway", "172.26.0.1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.ip_range.0.ip_range_name", "host IP range"),
				),
			},
		},
	})
}
