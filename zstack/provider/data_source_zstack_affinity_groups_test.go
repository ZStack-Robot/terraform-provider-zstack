// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackAffinityGroupsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.AffinityGroups) == 0 {
		t.Skip("no affinity groups in env data")
	}
	ag := env.AffinityGroups[0]
	name := envStr(ag, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_affinity_groups" "test" {
	name = %q
}`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_affinity_groups.test", "affinity_groups.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_affinity_groups.test", "affinity_groups.0.uuid", envStr(ag, "uuid")),
				),
			},
		},
	})
}
