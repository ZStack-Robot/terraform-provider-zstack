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

func TestAccZStackNetworkingSecGroupsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "zstack_networking_secgroups" "test" {
	name_pattern = "p%"
	filter {
		name = "state"
		values = ["Enabled"]
	}
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.name", "p1"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.state", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_networking_secgroups.test", "networking_secgroups.0.uuid", "f450b20497c34397977091bc1c8f87f9"),
				),
			},
		},
	})
}
