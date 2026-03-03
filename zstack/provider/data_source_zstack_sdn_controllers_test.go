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

func TestAccZStackSdnControllersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "zstack_networking_sdn_controllers" "test" {
  name_pattern = "172%"
  filter {
    name   = "status"
    values = ["Disconnected"]
  }
}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.#", "2"),

					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.status", "Disconnected"),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.name", "172.30.3.155"),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.ip", "172.30.3.155"),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.uuid", "65589889039944b5a2efeb2ed4d67594"),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.vendor_type", "Ovn"),
				),
			},
		},
	})
}
