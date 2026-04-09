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

func TestAccZStackL3NetworksDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L3Networks) == 0 {
		t.Skip("no l3 networks in env data")
	}
	l3 := env.L3Networks[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_l3networks" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks"), knownvalue.ListSizeExact(len(env.L3Networks))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(l3, "name"))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l3, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("category"), knownvalue.StringExact(envStr(l3, "category"))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("dns").AtSliceIndex(0).AtMapKey("dns_model"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("ip_range").AtSliceIndex(0).AtMapKey("cidr"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("ip_range").AtSliceIndex(0).AtMapKey("start_ip"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("ip_range").AtSliceIndex(0).AtMapKey("end_ip"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("ip_range").AtSliceIndex(0).AtMapKey("netmask"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("ip_range").AtSliceIndex(0).AtMapKey("gateway"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccZStackL3NetworksDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L3Networks) == 0 {
		t.Skip("no l3 networks in env data")
	}
	l3 := env.L3Networks[0]
	name := envStr(l3, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l3networks" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l3, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("category"), knownvalue.StringExact(envStr(l3, "category"))),
				},
			},
		},
	})
}

func TestAccZStackL3NetworksDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L3Networks) == 0 {
		t.Skip("no l3 networks in env data")
	}
	l3 := env.L3Networks[0]
	name := envStr(l3, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l3networks" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(l3, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_l3networks.test", tfjsonpath.New("l3networks").AtSliceIndex(0).AtMapKey("category"), knownvalue.StringExact(envStr(l3, "category"))),
				},
			},
		},
	})
}
