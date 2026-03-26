// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackL3NetworksDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L3Networks) == 0 {
		t.Skip("no l3 networks in env data")
	}
	l3 := env.L3Networks[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_l3networks" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.#", fmt.Sprintf("%d", len(env.L3Networks))),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", envStr(l3, "name")),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", envStr(l3, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", envStr(l3, "category")),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.dns.0.dns_model"),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.ip_range.0.cidr"),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.ip_range.0.start_ip"),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.ip_range.0.end_ip"),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.ip_range.0.netmask"),
					resource.TestCheckResourceAttrSet("data.zstack_l3networks.test", "l3networks.0.ip_range.0.gateway"),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l3networks" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", envStr(l3, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", envStr(l3, "category")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l3networks" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.uuid", envStr(l3, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l3networks.test", "l3networks.0.category", envStr(l3, "category")),
				),
			},
		},
	})
}
