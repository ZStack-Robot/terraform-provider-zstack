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

var _ = fmt.Sprintf

func TestAccZStackmnNodeDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.MnNodes) == 0 {
		t.Skip("no mn nodes in env data")
	}
	mn := env.MnNodes[0]

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_mnnodes" "test" {}`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_mnnodes.test", tfjsonpath.New("mn_nodes"), knownvalue.ListSizeExact(len(env.MnNodes))),
					statecheck.ExpectKnownValue("data.zstack_mnnodes.test", tfjsonpath.New("mn_nodes").AtSliceIndex(0).AtMapKey("host_name"), knownvalue.StringExact(envStr(mn, "host_name"))),
					statecheck.ExpectKnownValue("data.zstack_mnnodes.test", tfjsonpath.New("mn_nodes").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(mn, "uuid"))),
				},
			},
		},
	})
}
