// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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

	resource.ParallelTest(t, resource.TestCase{
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_networking_sdn_controllers.test", tfjsonpath.New("sdn_controllers").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(status)),
					statecheck.ExpectKnownValue("data.zstack_networking_sdn_controllers.test", tfjsonpath.New("sdn_controllers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_networking_sdn_controllers.test", tfjsonpath.New("sdn_controllers").AtSliceIndex(0).AtMapKey("ip"), knownvalue.StringExact(envStr(sdn, "ip"))),
					statecheck.ExpectKnownValue("data.zstack_networking_sdn_controllers.test", tfjsonpath.New("sdn_controllers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(sdn, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_networking_sdn_controllers.test", tfjsonpath.New("sdn_controllers").AtSliceIndex(0).AtMapKey("vendor_type"), knownvalue.StringExact(envStr(sdn, "vendor_type"))),
				},
			},
		},
	})
}
