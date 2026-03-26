// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackHostsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Hosts) == 0 {
		t.Skip("no hosts in env data")
	}
	h := env.Hosts[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_hosts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", fmt.Sprintf("%d", len(env.Hosts))),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", envStr(h, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", envStr(h, "cluster_uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", envStr(h, "state")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", envStr(h, "name")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", envStr(h, "status")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", envStr(h, "type")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", envStr(h, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", envStr(h, "uuid")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_hosts" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", envStr(h, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", envStr(h, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", envStr(h, "cluster_uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", envStr(h, "state")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", envStr(h, "status")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", envStr(h, "type")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", envStr(h, "zone_uuid")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_hosts" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.uuid", envStr(h, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.architecture", envStr(h, "architecture")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.cluster_uuid", envStr(h, "cluster_uuid")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.state", envStr(h, "state")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.status", envStr(h, "status")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.type", envStr(h, "type")),
					resource.TestCheckResourceAttr("data.zstack_hosts.test", "hosts.0.zone_uuid", envStr(h, "zone_uuid")),
				),
			},
		},
	})
}
