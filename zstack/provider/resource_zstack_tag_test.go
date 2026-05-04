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

func TestTagResource_Schema(t *testing.T) {
	var r tagResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "value"}
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

func TestTagResource_Metadata(t *testing.T) {
	var r tagResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_tag" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccTagResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_tag" "test" {
  name  = "acc-test-tag"
  value = "test-value"
  type  = "simple"
  color = "#FF0000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckTagDisappears("zstack_tag.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTagResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_tag" "test" {
  name  = "acc-test-tag"
  value = "test-value"
  type  = "simple"
  color = "#FF0000"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_tag.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_tag.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-tag")),
					statecheck.ExpectKnownValue("zstack_tag.test", tfjsonpath.New("value"), knownvalue.StringExact("test-value")),
				},
			},
			// Step 2: Update — N/A
			// ZStack API rejects UpdateTag for type=simple tags entirely
			// (error: "cannot update simple tag pattern format", code ORG_ZSTACK_TAG2_10002).
			// type=withToken requires complex pattern format with {tokenName} placeholders.
			// Update Step skipped until tag update semantics are clarified with ZStack team.
			// Step 3: Import
			{
				ResourceName:                        "zstack_tag.test",
				ImportState:                         true,
				ImportStateIdFunc:                   importStateIdFromUUID("zstack_tag.test"),
				ImportStateVerify:                   true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
