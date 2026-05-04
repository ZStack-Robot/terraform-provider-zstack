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

func TestGlobalConfigResource_Schema(t *testing.T) {
	var r globalConfigResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "category", "value"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"default_value", "description"}
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

func TestGlobalConfigResource_Metadata(t *testing.T) {
	var r globalConfigResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_global_config" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccGlobalConfigResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckGlobalConfigDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_global_config" "test" {
  category = "vm"
  name     = "deletionPolicy"
  value    = "Delay"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("name"), knownvalue.StringExact("deletionPolicy")),
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("category"), knownvalue.StringExact("vm")),
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("value"), knownvalue.StringExact("Delay")),
				},
			},
			// Step 2: Update value in place
			{
				Config: providerConfig() + `
resource "zstack_global_config" "test" {
  category = "vm"
  name     = "deletionPolicy"
  value    = "Direct"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("name"), knownvalue.StringExact("deletionPolicy")),
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("category"), knownvalue.StringExact("vm")),
					statecheck.ExpectKnownValue("zstack_global_config.test", tfjsonpath.New("value"), knownvalue.StringExact("Direct")),
				},
			},
		},
	})
}
