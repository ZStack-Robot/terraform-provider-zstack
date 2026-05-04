// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVipDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Vips) == 0 {
		t.Skip("no vips in env data")
	}
	item := env.Vips[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_vips" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.#", fmt.Sprintf("%d", len(env.Vips))),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.name", envStr(item, "name")),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}

func TestAccZStackVipDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Vips) == 0 {
		t.Skip("no vips in env data")
	}
	item := env.Vips[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_vips" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_vips.test", "vips.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}
