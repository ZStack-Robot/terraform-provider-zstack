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

func TestAccZStackHostsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Hosts) == 0 {
		t.Skip("no hosts in env data")
	}
	h := env.Hosts[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_hosts" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts"), knownvalue.ListSizeExact(len(env.Hosts))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(h, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("cluster_uuid"), knownvalue.StringExact(envStr(h, "cluster_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(h, "state"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(h, "name"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(h, "status"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(h, "type"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(h, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(h, "uuid"))),
				},
			},
		},
	})
}

func TestAccZStackHostDataSourceFilterByname(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Hosts) == 0 {
		t.Skip("no hosts in env data")
	}
	h := env.Hosts[0]
	name := envStr(h, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_hosts" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(h, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(h, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("cluster_uuid"), knownvalue.StringExact(envStr(h, "cluster_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(h, "state"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(h, "status"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(h, "type"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(h, "zone_uuid"))),
				},
			},
		},
	})
}

func TestAccZStackHostDataSourceFilterBynamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Hosts) == 0 {
		t.Skip("no hosts in env data")
	}
	h := env.Hosts[0]
	name := envStr(h, "name")
	pattern := name[:3] + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_hosts" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(h, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("architecture"), knownvalue.StringExact(envStr(h, "architecture"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("cluster_uuid"), knownvalue.StringExact(envStr(h, "cluster_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(h, "state"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("status"), knownvalue.StringExact(envStr(h, "status"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(h, "type"))),
					statecheck.ExpectKnownValue("data.zstack_hosts.test", tfjsonpath.New("hosts").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(h, "zone_uuid"))),
				},
			},
		},
	})
}
