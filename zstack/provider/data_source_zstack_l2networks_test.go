// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackL2NetworkDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L2Networks) == 0 {
		t.Skip("no l2 networks in env data")
	}
	l2 := env.L2Networks[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_l2networks" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.#", fmt.Sprintf("%d", len(env.L2Networks))),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", envStr(l2, "name")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", envStr(l2, "type")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", envStr(l2, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", envStr(l2, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", envStr(l2, "physical_interface")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l2networks" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", envStr(l2, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", envStr(l2, "type")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", envStr(l2, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", envStr(l2, "physical_interface")),
				),
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_l2networks" "test" { name_pattern = %q }`, pattern),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.uuid", envStr(l2, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.type", envStr(l2, "type")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.zone_uuid", envStr(l2, "zone_uuid")),
					resource.TestCheckResourceAttr("data.zstack_l2networks.test", "l2networks.0.physical_interface", envStr(l2, "physical_interface")),
				),
			},
		},
	})
}
