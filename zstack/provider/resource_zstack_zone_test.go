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

func TestZoneResource_Schema(t *testing.T) {
	var r zoneResource
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

	computed := []string{"uuid", "type"}
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

func TestZoneResource_Metadata(t *testing.T) {
	var r zoneResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_zone" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccZoneResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckZoneDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_zone" "test" {
  name        = "acc-test-zone"
  description = "Acceptance test zone"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_zone.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_zone.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-zone")),
					statecheck.ExpectKnownValue("zstack_zone.test", tfjsonpath.New("state"), knownvalue.StringExact("Enabled")),
				},
			},
			{
				ResourceName:      "zstack_zone.test",
				ImportState:       true,
				ImportStateIdFunc:       importStateUUID("zstack_zone.test"),
				ImportStateVerify: true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccZoneResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckZoneDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_zone" "test" {
  name        = "acc-test-zone-disappears"
  description = "Disappears test zone"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_zone.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckZoneDisappears("zstack_zone.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
