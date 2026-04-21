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

func TestL2vxlanNetworkResource_Schema(t *testing.T) {
	var r l2vxlanNetworkResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "pool_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "physical_interface", "type"}
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

func TestL2vxlanNetworkResource_Metadata(t *testing.T) {
	var r l2vxlanNetworkResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_l2vxlan_network" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccL2VxlanNetworkResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.SdnControllers) == 0 {
		t.Skip("no SDN controllers / VXLAN pool available in env, skipping L2 VXLAN network acceptance test")
	}

	// SdnControllers is present — extract pool_uuid from the first controller
	poolUUID := envStr(env.SdnControllers[0], "pool_uuid")
	if poolUUID == "" {
		t.Skip("no pool_uuid found in SDN controller, skipping L2 VXLAN network acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckL2VxlanNetworkDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create (Update not implemented)
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l2vxlan_network" "test" {
  name      = "acc-test-l2vxlan"
  pool_uuid = %q
}
`, poolUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_l2vxlan_network.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_l2vxlan_network.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-l2vxlan")),
					statecheck.ExpectKnownValue("zstack_l2vxlan_network.test", tfjsonpath.New("pool_uuid"), knownvalue.StringExact(poolUUID)),
				},
			},
			// Step 2: Import
			{
				ResourceName:                        "zstack_l2vxlan_network.test",
				ImportState:                         true,
				ImportStateIdFunc:                   importStateIdFromUUID("zstack_l2vxlan_network.test"),
				ImportStateVerify:                   true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
