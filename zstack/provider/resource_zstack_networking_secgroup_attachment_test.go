// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

func importStateIdSecgroupAttachment(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}
		sgUUID := rs.Primary.Attributes["secgroup_uuid"]
		nicUUID := rs.Primary.Attributes["nic_uuid"]
		if sgUUID == "" || nicUUID == "" {
			return "", fmt.Errorf("secgroup_uuid or nic_uuid attribute is empty for %s", resourceName)
		}
		return sgUUID + ":" + nicUUID, nil
	}
}

func TestAccZStackSecurityGroupAttachment(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC is set")
	}

	nicUUID := createSecurityGroupTestNic(t)
	name := testAccName("sg-attach")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecGroupAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_networking_secgroup" "test_sg" {
  name        = %q
  description = "acceptance security group attachment"
  ip_version  = 4
}

resource "zstack_networking_secgroup_attachment" "test_attach" {
  secgroup_uuid = zstack_networking_secgroup.test_sg.uuid
  nic_uuid      = %q
}
`, name+"-sg", nicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_networking_secgroup.test_sg", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_attachment.test_attach", tfjsonpath.New("nic_uuid"), knownvalue.StringExact(nicUUID)),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_attachment.test_attach", tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:                         "zstack_networking_secgroup_attachment.test_attach",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdSecgroupAttachment("zstack_networking_secgroup_attachment.test_attach"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "id",
			},
		},
	})
}

func createSecurityGroupTestNic(t *testing.T) string {
	t.Helper()

	cli := testAccClientLoggedIn()
	q := param.NewQueryParam()
	q.AddQ("networkServiceType=SecurityGroup")
	refs, err := cli.QueryNetworkServiceL3NetworkRef(&q)
	if err != nil {
		t.Fatalf("query SecurityGroup-enabled L3 refs: %v", err)
	}
	if len(refs) == 0 {
		t.Skip("no L3 network with SecurityGroup service enabled")
	}

	vmQuery := param.NewQueryParam()
	vmQuery.AddQ("state=Running")
	vmQuery.AddQ("type=UserVm")
	vms, err := cli.QueryVmInstance(&vmQuery)
	if err != nil {
		t.Fatalf("query running user VMs: %v", err)
	}
	if len(vms) == 0 {
		t.Skip("no running UserVm available for temporary SecurityGroup NIC")
	}

	for _, vm := range vms {
		for _, ref := range refs {
			hasL3 := false
			before := make(map[string]struct{}, len(vm.VmNics))
			for _, nic := range vm.VmNics {
				before[nic.UUID] = struct{}{}
				if nic.L3NetworkUuid == ref.L3NetworkUuid {
					hasL3 = true
				}
			}
			if hasL3 {
				continue
			}

			updated, err := cli.AttachL3NetworkToVm(vm.UUID, ref.L3NetworkUuid, param.AttachL3NetworkToVmParam{
				BaseParam: param.BaseParam{},
				Params:    param.AttachL3NetworkToVmParamDetail{},
			})
			if err != nil {
				t.Logf("attach SecurityGroup L3 %s to VM %s failed: %v", ref.L3NetworkUuid, vm.UUID, err)
				continue
			}

			for _, nic := range updated.VmNics {
				if _, existed := before[nic.UUID]; existed || nic.L3NetworkUuid != ref.L3NetworkUuid {
					continue
				}
				t.Cleanup(func() {
					detachQuery := param.NewQueryParam()
					detachQuery.AddQ("vmNicUuid=" + nic.UUID)
					sgRefs, err := cli.QueryVmNicInSecurityGroup(&detachQuery)
					if err == nil {
						for _, sgRef := range sgRefs {
							_ = cli.DeleteWithSpec("v1/security-groups", sgRef.SecurityGroupUuid, "vm-instances/nics", "vmNicUuids="+nic.UUID, nil)
						}
					}
					if err := cli.DeleteVmNic(nic.UUID, param.DeleteModePermissive); err != nil {
						t.Logf("cleanup temporary NIC %s failed: %v", nic.UUID, err)
					}
				})
				t.Logf("created temporary SecurityGroup-capable NIC %s on L3 %s for VM %s", nic.UUID, ref.L3NetworkUuid, vm.UUID)
				return nic.UUID
			}
		}
	}

	t.Skip("could not create a temporary NIC on a SecurityGroup-enabled L3")
	return ""
}
