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

func TestAccZStackAffinityGroupsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.AffinityGroups) == 0 {
		t.Skip("no affinity groups in env data")
	}
	ag := env.AffinityGroups[0]
	name := envStr(ag, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_affinity_groups" "test" {
	name = %q
}`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_affinity_groups.test", tfjsonpath.New("affinity_groups").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_affinity_groups.test", tfjsonpath.New("affinity_groups").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(ag, "uuid"))),
				},
			},
		},
	})
}
