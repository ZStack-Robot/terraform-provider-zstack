// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestNetworkingSecgroupResource_Schema(t *testing.T) {
	var r securityGroupResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "ip_version"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid"}
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

func TestNetworkingSecgroupResource_Metadata(t *testing.T) {
	var r securityGroupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_networking_secgroup" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSecurityGroupResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_networking_secgroup" "test" {
  name       = "acc-test-secgroup"
  ip_version = 4
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckSecurityGroupDisappears("zstack_networking_secgroup.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSecurityGroupResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_networking_secgroup" "test" {
  name       = "acc-test-secgroup"
  ip_version = 4
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_networking_secgroup.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_networking_secgroup.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-secgroup")),
				},
			},
			// Step 2: Update — N/A
			// networking_secgroup name/description have RequiresReplace (Category C in REQUIRES-REPLACE-AUDIT.md:
			// SDK supports UpdateSecurityGroup but Provider has not implemented Update).
			// No in-place updatable attributes exist. Update Step skipped until Provider Update is implemented.
			// Step 3: Import
			{
				ResourceName:                         "zstack_networking_secgroup.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_networking_secgroup.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"description", "ip_version"},
			},
		},
	})
}
