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

func TestAccZStackAccountsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Accounts) == 0 {
		t.Skip("no accounts in env data")
	}
	acct := env.Accounts[0]
	name := envStr(acct, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_accounts" "test" {
	name = %q
}`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_accounts.test", tfjsonpath.New("accounts").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_accounts.test", tfjsonpath.New("accounts").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(acct, "uuid"))),
				},
			},
		},
	})
}
