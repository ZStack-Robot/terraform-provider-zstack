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

func TestAccZStackInstanceGuestToolsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]
	vmUUID := envStr(vm, "uuid")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_guest_tools" "test" { instance_uuid = %q }`, vmUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_guest_tools.test", tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
		},
	})
}
