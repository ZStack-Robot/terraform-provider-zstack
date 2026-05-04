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

func TestCdpPolicyResource_Schema(t *testing.T) {
	var r cdpPolicyResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "recovery_point_per_second"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state"}
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

func TestCdpPolicyResource_Metadata(t *testing.T) {
	var r cdpPolicyResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_cdp_policy" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccCdpPolicyResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCdpPolicyDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_cdp_policy" "test" {
  name                     = "acc-test-cdp-policy"
  recovery_point_per_second = 1
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cdp_policy.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_cdp_policy.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cdp-policy")),
				},
			},
			// Step 2: Update name (no ForceNew on name — true in-place update)
			{
				Config: providerConfig() + `
resource "zstack_cdp_policy" "test" {
  name                     = "acc-test-cdp-policy-updated"
  recovery_point_per_second = 1
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cdp_policy.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cdp-policy-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_cdp_policy.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_cdp_policy.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
