// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackvmInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_vminstances" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of vm instances returned
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.#", "5"),

					// Verify the first vm instances to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.architecture", "x86_64"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.clusteruuid", "b286789480254a208e6327136bb3dcb3"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.cupnum", "2"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.hostuuid", "818c9469d7ae4cab9bb7516753543d3f"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.hypervisortype", "KVM"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.imageuuid", "ffd420a9259645c7a3ef8594d2d1f56e"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.memorysize", "34359738368"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.name", "Zaku"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.state", "Running"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.type", "UserVm"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.uuid", "07d7f69688f24099b0e569dfdbb15121"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.zoneuuid", "4981061bd27c42c7bc063c4c4529518a"),

					// Verify the first volume of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumeactualsize", "28502401024"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumedescription", "Root volume for VM[uuid:07d7f69688f24099b0e569dfdbb15121]"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumeformat", "qcow2"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumesize", "515396075520"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumestate", "Enabled"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumestatus", "Ready"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumetype", "Root"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.allvolumes.0.volumeuuid", "498398ac73d14f14a0bfca1be6d32c90"),

					// Verify the first nic of vm instance to ensure all attributes are set
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.vmnics.0.ip", "172.25.126.218"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.vmnics.0.mac", "fa:6d:c8:20:fe:00"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.vmnics.0.netmask", "255.255.0.0"),
					resource.TestCheckResourceAttr("data.zstack_vminstances.test", "vminstances.0.vmnics.0.gateway", "172.25.0.1"),
				),
			},
		},
	})
}
