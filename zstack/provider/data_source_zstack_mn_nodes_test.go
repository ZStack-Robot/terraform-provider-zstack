// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackmnNodeDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.MnNodes) == 0 {
		t.Skip("no mn nodes in env data")
	}
	mn := env.MnNodes[0]

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_mnnodes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.#", fmt.Sprintf("%d", len(env.MnNodes))),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.host_name", envStr(mn, "host_name")),
					resource.TestCheckResourceAttr("data.zstack_mnnodes.test", "mn_nodes.0.uuid", envStr(mn, "uuid")),
				),
			},
		},
	})
}
