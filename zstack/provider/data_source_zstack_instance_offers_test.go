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

func TestAccZStackInstanceOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_instance_offers" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers"), knownvalue.ListSizeExact(len(env.InstanceOfferings))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(io, "name"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(io, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("allocator_strategy"), knownvalue.StringExact(envStr(io, "allocator_strategy"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("cpu_num"), knownvalue.StringExact(envStr(io, "cpu_num"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("memory_size"), knownvalue.StringExact(envStr(io, "memory_size"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("sort_key"), knownvalue.StringExact(envStr(io, "sort_key"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(io, "state"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(io, "type"))),
				},
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]
	name := envStr(io, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instance_offers" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(io, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("allocator_strategy"), knownvalue.StringExact(envStr(io, "allocator_strategy"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("cpu_num"), knownvalue.StringExact(envStr(io, "cpu_num"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("memory_size"), knownvalue.StringExact(envStr(io, "memory_size"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("sort_key"), knownvalue.StringExact(envStr(io, "sort_key"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(io, "state"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(io, "type"))),
				},
			},
		},
	})
}

func TestAccZStackInstanceOffersDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.InstanceOfferings) == 0 {
		t.Skip("no instance offerings in env data")
	}
	io := env.InstanceOfferings[0]
	name := envStr(io, "name")
	pattern := name[:3] + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_instance_offers" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(io, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_instance_offers.test", tfjsonpath.New("instance_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(io, "state"))),
				},
			},
		},
	})
}
