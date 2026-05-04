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

func TestVmNicResource_Schema(t *testing.T) {
	var r vmNicResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"l3_network_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}
	// Check computed attributes
	computed := []string{"uuid", "mac", "vm_instance_uuid", "netmask", "gateway", "state"}
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

func TestVmNicResource_Metadata(t *testing.T) {
	var r vmNicResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vm_nic" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVmNicResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env data")
	}

	l3UUID := envStr(env.L3Networks[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckVmNicDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_vm_nic" "test" {
  l3_network_uuid = %q
}
`, l3UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_vm_nic.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_vm_nic.test", tfjsonpath.New("l3_network_uuid"), knownvalue.StringExact(l3UUID)),
				},
			},
			// Step 2: Import (no Update supported)
			{
				ResourceName:                         "zstack_vm_nic.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_vm_nic.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
