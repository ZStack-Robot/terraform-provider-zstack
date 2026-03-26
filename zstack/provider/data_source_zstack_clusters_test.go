// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackClusterDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Clusters) == 0 {
		t.Skip("no clusters in env data")
	}
	c := env.Clusters[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_clusters" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", fmt.Sprintf("%d", len(env.Clusters))),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisor_type", envStr(c, "hypervisor_type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", envStr(c, "name")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", envStr(c, "state")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", envStr(c, "type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", envStr(c, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", envStr(c, "uuid")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_clusters" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", envStr(c, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisor_type", envStr(c, "hypervisor_type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", envStr(c, "state")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", envStr(c, "type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", envStr(c, "zone_uuid")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_clusters" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.uuid", envStr(c, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.hypervisor_type", envStr(c, "hypervisor_type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.state", envStr(c, "state")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.type", envStr(c, "type")),
					resource.TestCheckResourceAttr("data.zstack_clusters.test", "clusters.0.zone_uuid", envStr(c, "zone_uuid")),
				),
			},
		},
	})
}
