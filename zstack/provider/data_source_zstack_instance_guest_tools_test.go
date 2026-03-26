// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackInstanceGuestToolsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm instances in env data")
	}
	vm := env.VmInstances[0]
	vmUUID := envStr(vm, "uuid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`data "zstack_guest_tools" "test" { instance_uuid = %q }`, vmUUID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zstack_guest_tools.test", "status"),
				),
			},
		},
	})
}
