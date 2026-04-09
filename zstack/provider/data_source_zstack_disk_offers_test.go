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

func TestAccZStackDiskOffersDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_disk_offers" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers"), knownvalue.ListSizeExact(len(env.DiskOfferings))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(do, "name"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(do, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("allocator_strategy"), knownvalue.StringExact(envStr(do, "allocator_strategy"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact(envStr(do, "description"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("disk_size"), knownvalue.StringExact(envStr(do, "disk_size"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("sort_key"), knownvalue.StringExact(envStr(do, "sort_key"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(do, "type"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(do, "state"))),
				},
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]
	name := envStr(do, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_disk_offers" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(do, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("allocator_strategy"), knownvalue.StringExact(envStr(do, "allocator_strategy"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact(envStr(do, "description"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("disk_size"), knownvalue.StringExact(envStr(do, "disk_size"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("sort_key"), knownvalue.StringExact(envStr(do, "sort_key"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(do, "type"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(do, "state"))),
				},
			},
		},
	})
}

func TestAccZStackDiskOffersDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.DiskOfferings) == 0 {
		t.Skip("no disk offerings in env data")
	}
	do := env.DiskOfferings[0]
	name := envStr(do, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_disk_offers" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(do, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_disk_offers.test", tfjsonpath.New("disk_offers").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(do, "state"))),
				},
			},
		},
	})
}
