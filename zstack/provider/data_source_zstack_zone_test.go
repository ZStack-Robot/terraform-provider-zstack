// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackZonesDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data")
	}
	z := env.Zones[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_zones" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.#", fmt.Sprintf("%d", len(env.Zones))),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.name", envStr(z, "name")),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.uuid", envStr(z, "uuid")),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.type", envStr(z, "type")),
					resource.TestCheckResourceAttr("data.zstack_zones.test", "zones.0.state", envStr(z, "state")),
				),
			},
		},
	})
}
