// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackInstanceGuestToolsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_guest_tools" "test" { instance_uuid = "1e3a89a3dbfe434cbc8571af0e27e711" }`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_guest_tools.test", "version", ""),
					resource.TestCheckResourceAttr("data.zstack_guest_tools.test", "status", "Not Connected"),
				),
			},
		},
	})
}
