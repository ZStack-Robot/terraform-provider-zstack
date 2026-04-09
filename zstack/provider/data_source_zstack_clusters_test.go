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

func TestAccZStackClusterDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Clusters) == 0 {
		t.Skip("no clusters in env data")
	}
	c := env.Clusters[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_clusters" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters"), knownvalue.ListSizeExact(len(env.Clusters))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("hypervisor_type"), knownvalue.StringExact(envStr(c, "hypervisor_type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(envStr(c, "name"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(c, "state"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(c, "type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(c, "zone_uuid"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(c, "uuid"))),
				},
			},
		},
	})
}

func TestAccZStackClusterDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Clusters) == 0 {
		t.Skip("no clusters in env data")
	}
	c := env.Clusters[0]
	name := envStr(c, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_clusters" "test" { name = %q }`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(c, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("hypervisor_type"), knownvalue.StringExact(envStr(c, "hypervisor_type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(c, "state"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(c, "type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(c, "zone_uuid"))),
				},
			},
		},
	})
}

func TestAccZStackClusterDataSourceFilterByNamePattern(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Clusters) == 0 {
		t.Skip("no clusters in env data")
	}
	c := env.Clusters[0]
	name := envStr(c, "name")
	pattern := name[:3] + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_clusters" "test" { name_pattern = %q }`, pattern),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(c, "uuid"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("hypervisor_type"), knownvalue.StringExact(envStr(c, "hypervisor_type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(c, "state"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("type"), knownvalue.StringExact(envStr(c, "type"))),
					statecheck.ExpectKnownValue("data.zstack_clusters.test", tfjsonpath.New("clusters").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.StringExact(envStr(c, "zone_uuid"))),
				},
			},
		},
	})
}
