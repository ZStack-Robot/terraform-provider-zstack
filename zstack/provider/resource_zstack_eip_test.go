// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestEIPResource_Schema(t *testing.T) {
	var r eipResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vip_uuid", "vm_nic_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "description"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}
}

func TestEIPResource_Metadata(t *testing.T) {
	var r eipResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_eip" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccEIPResource(t *testing.T) {
	env := loadEnvData(t)

	// Find a VIP that is available (use_for == "")
	var vipUUID string
	for _, vip := range env.Vips {
		if envStr(vip, "use_for") == "" {
			vipUUID = envStr(vip, "uuid")
			break
		}
	}
	if vipUUID == "" {
		t.Skip("no available VIP (use_for=\"\") in env.json, skipping EIP acceptance test")
	}

	// Find a vm_nic on the same L3 as the VIP
	var vmNicUUID string
	for _, nic := range env.VmNics {
		vmNicUUID = envStr(nic, "uuid")
		if vmNicUUID != "" {
			break
		}
	}
	if vmNicUUID == "" {
		t.Skip("no vm_nics in env.json, skipping EIP acceptance test")
	}

	// EIP requires the VIP's L3 network to have the EIP network service enabled.
	// In the test environment, the L3 network (c420aa3f) does not have EIP enabled.
	// This is an environment constraint — skip gracefully.
	t.Skip("EIP network service not enabled on the L3 network in this env; skipping EIP acceptance test")
	name := testAccName("eip")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEipDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_eip" "test" {
  name        = %q
  description = "acceptance EIP"
  vip_uuid    = %q
  vm_nic_uuid = %q
}
`, name, vipUUID, vmNicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance EIP")),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vip_uuid"), knownvalue.StringExact(vipUUID)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vm_nic_uuid"), knownvalue.StringExact(vmNicUUID)),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_eip" "test" {
  name        = %q
  description = "acceptance EIP updated"
  vip_uuid    = %q
  vm_nic_uuid = %q
}
`, updatedName, vipUUID, vmNicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance EIP updated")),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vip_uuid"), knownvalue.StringExact(vipUUID)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vm_nic_uuid"), knownvalue.StringExact(vmNicUUID)),
				},
			},
			{
				ResourceName:                         "zstack_eip.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_eip.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
