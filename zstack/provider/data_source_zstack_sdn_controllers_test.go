// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackSdnControllersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SdnControllers) == 0 {
		t.Skip("no sdn controllers in env data")
	}
	sdn := env.SdnControllers[0]
	name := envStr(sdn, "name")
	pattern := string([]rune(name)[:3]) + "%"
	status := envStr(sdn, "status")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_networking_sdn_controllers" "test" {
  name_pattern = %q
  filter {
    name   = "status"
    values = [%q]
  }
}`, pattern, status),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.status", status),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.ip", envStr(sdn, "ip")),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.uuid", envStr(sdn, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_networking_sdn_controllers.test", "sdn_controllers.0.vendor_type", envStr(sdn, "vendor_type")),
				),
			},
		},
	})
}
