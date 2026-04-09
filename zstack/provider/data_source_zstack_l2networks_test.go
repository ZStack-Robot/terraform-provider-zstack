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

func TestAccZStackL2NetworkDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L2Networks) == 0 {
		t.Skip("no l2 networks in env data")
	}
	l2 := env.L2Networks[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_l2networks" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks"), knownvalue.ListSizeExact(len(env.L2Networks))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(l2, "name"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(l2, "type"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l2, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(l2, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("physical_interface"), knownvalue.StringExact(envStr(l2, "physical_interface"))),
				},
			},
		},
	})
}

func TestAccZStackL2NetworkDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L2Networks) == 0 {
		t.Skip("no l2 networks in env data")
	}
	l2 := env.L2Networks[0]
	name := envStr(l2, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l2networks" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l2, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(l2, "type"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(l2, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("physical_interface"), knownvalue.StringExact(envStr(l2, "physical_interface"))),
				},
			},
		},
	})
}

func TestAccZStackL2NetworkDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L2Networks) == 0 {
		t.Skip("no l2 networks in env data")
	}
	l2 := env.L2Networks[0]
	name := envStr(l2, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l2networks" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l2, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(l2, "type"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(l2, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l2networks.test", tfjsonpath.New("l2networks").AtSliceIndex(0).AtMapKey("physical_interface"), knownvalue.StringExact(envStr(l2, "physical_interface"))),
				},
			},
		},
	})
}
