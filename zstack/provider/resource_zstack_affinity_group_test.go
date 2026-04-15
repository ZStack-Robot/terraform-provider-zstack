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

func TestAffinityGroupResource_Schema(t *testing.T) {
	var r affinityGroupResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "policy"}
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

	optional := []string{"description", "type", "zone_uuid"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestAffinityGroupResource_Metadata(t *testing.T) {
	var r affinityGroupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_affinity_group" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAffinityGroupResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAffinityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_affinity_group" "test" {
  name   = "acc-test-affinity-group"
  policy = "antiSoft"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-affinity-group")),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("policy"), knownvalue.StringExact("antiSoft")),
				},
			},
			{
				ResourceName:      "zstack_affinity_group.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_affinity_group.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccAffinityGroupResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAffinityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_affinity_group" "test_disappears" {
  name   = "acc-test-ag-disappears"
  policy = "antiSoft"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_affinity_group.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckAffinityGroupDisappears("zstack_affinity_group.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
