// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"regexp"
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
				Config: providerConfig() + `data "zstack_zone" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones"), knownvalue.ListSizeExact(len(env.Zones))),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(z, "name"))),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(z, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(z, "type"))),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(z, "state"))),
				},
			},
		},
	})
}

// AI-style lookup: pin the data source by uuid so the result is deterministic.
func TestAccZStackZonesDataSource_byUUID(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}
	z := env.Zones[0]
	uuid := envStr(z, "uuid")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_zone" "test" { uuid = %q }`, uuid),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(uuid)),
					statecheck.ExpectKnownValue("data.zstack_zone.test", tfjsonpath.New("zones").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(z, "name"))),
				},
			},
		},
	})
}

// uuid + name in the same config should fail validation at plan time (ConflictsWith).
func TestAccZStackZonesDataSource_uuidNameMutualExclusion(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}
	z := env.Zones[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(
					`
data "zstack_zone" "test" {
  uuid = %q
  name = %q
}
`,
					envStr(z, "uuid"), envStr(z, "name"),
				),
				ExpectError: regexp.MustCompile(`(?s)Invalid Attribute Combination.*cannot be specified`),
			},
		},
	})
}
