// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVolumeDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	_ = env

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_volumes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.#"),
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.0.name"),
					resource.TestCheckResourceAttrSet("data.zstack_volumes.test", "volumes.0.uuid"),
				),
			},
		},
	})
}

func TestAccZStackVolumeDataSourceFilterByName(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Volumes) == 0 {
		t.Skip("no volumes in env data")
	}
	item := env.Volumes[0]
	name := envStr(item, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_volumes" "test" { name = %q }`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_volumes.test", "volumes.0.uuid", envStr(item, "uuid")),
				),
			},
		},
	})
}
