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

func TestPolicyResource_Schema(t *testing.T) {
	var r policyResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "account_uuid"}
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

func TestPolicyResource_Metadata(t *testing.T) {
	var r policyResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_policy" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccPolicyResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_policy" "test" {
  name = "acc-test-policy"

  statements = [
    {
      effect    = "Allow"
      actions   = ["vm:StartVmInstance"]
      resources = ["*"]
    }
  ]
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckPolicyDisappears("zstack_policy.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPolicyResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_policy" "test" {
  name = "acc-test-policy"

  statements = [
    {
      effect    = "Allow"
      actions   = ["vm:StartVmInstance"]
      resources = ["*"]
    }
  ]
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_policy.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_policy.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-policy")),
				},
			},
			{
				ResourceName:                         "zstack_policy.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_policy.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
