// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func TestAccZStackNetworkingSecGroupRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "zstack_networking_secgroup_rules" "test" {
	filter {
		name   = "action"
		values = ["DROP"]
	}
	filter {
        name   = "protocol"
        values = ["TCP"]
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.#", "3"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.action", "DROP"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.uuid", "3d4c2f323ebf4066bf0ab5d2039f3548"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.dst_port_range", "80,443"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.type", "Ingress"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroup_rules.test", "rules.0.src_ip_range", "13.13.15.13"),
				),
			},
		},
	})
}
