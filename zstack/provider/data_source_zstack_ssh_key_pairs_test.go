// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackSshKeyPairsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.SshKeyPairs) == 0 {
		t.Skip("no SSH key pairs in env data")
	}
	skp := env.SshKeyPairs[0]
	name := envStr(skp, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_ssh_key_pairs" "test" {
	name = %q
}`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_ssh_key_pairs.test", "ssh_key_pairs.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_ssh_key_pairs.test", "ssh_key_pairs.0.uuid", envStr(skp, "uuid")),
				),
			},
		},
	})
}
