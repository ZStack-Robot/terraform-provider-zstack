// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackAccountsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Accounts) == 0 {
		t.Skip("no accounts in env data")
	}
	acct := env.Accounts[0]
	name := envStr(acct, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_accounts" "test" {
	name = %q
}`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_accounts.test", "accounts.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_accounts.test", "accounts.0.uuid", envStr(acct, "uuid")),
				),
			},
		},
	})
}
