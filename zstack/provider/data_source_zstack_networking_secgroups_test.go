// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccZStackNetworkingSecGroupsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SecurityGroups) == 0 {
		t.Skip("no security groups in env data")
	}
	sg := env.SecurityGroups[0]
	name := envStr(sg, "name")
	pattern := string([]rune(name)[:1]) + "%"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_networking_secgroups" "test" {
	name_pattern = %q
	filter {
		name = "state"
		values = [%q]
	}
}`, pattern, envStr(sg, "state")),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_networking_secgroups.test", tfjsonpath.New("networking_secgroups").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroups.test", tfjsonpath.New("networking_secgroups").AtSliceIndex(0).AtMapKey("state"), knownvalue.StringExact(envStr(sg, "state"))),
					statecheck.ExpectKnownValue("data.zstack_networking_secgroups.test", tfjsonpath.New("networking_secgroups").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(sg, "uuid"))),
				},
			},
		},
	})
}
