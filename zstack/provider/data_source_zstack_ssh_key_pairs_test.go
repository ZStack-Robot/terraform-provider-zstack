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

func TestAccZStackSshKeyPairsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SshKeyPairs) == 0 {
		t.Skip("no SSH key pairs in env data")
	}
	skp := env.SshKeyPairs[0]
	name := envStr(skp, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_ssh_key_pairs" "test" {
	name = %q
}`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_ssh_key_pairs.test", tfjsonpath.New("ssh_key_pairs").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_ssh_key_pairs.test", tfjsonpath.New("ssh_key_pairs").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(skp, "uuid"))),
				},
			},
		},
	})
}
