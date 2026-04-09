// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ = fmt.Sprintf

func TestAccZStackZonesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}
	z := env.Zones[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_zones" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_zones.test", tfjsonpath.New("zones"), knownvalue.ListSizeExact(len(env.Zones))),
					statecheck.ExpectKnownValue("data.zstack_zones.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(z, "name"))),
					statecheck.ExpectKnownValue("data.zstack_zones.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(z, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_zones.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(z, "type"))),
					statecheck.ExpectKnownValue("data.zstack_zones.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(z, "state"))),
				},
			},
		},
	})
}
