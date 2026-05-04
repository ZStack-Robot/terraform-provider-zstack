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

func TestStackTemplateResource_Schema(t *testing.T) {
	var r stackTemplateResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
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
	// Check computed attributes
	computed := []string{"uuid", "version", "state", "md5sum"}
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

func TestStackTemplateResource_Metadata(t *testing.T) {
	var r stackTemplateResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_stack_template" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestHasStackTemplateFormatVersionMarker(t *testing.T) {
	if !hasStackTemplateFormatVersionMarker(testStackTemplateContent) {
		t.Fatal("expected test template content to contain ZStackTemplateFormatVersion")
	}
	if hasStackTemplateFormatVersionMarker(`{"Resources":{}}`) {
		t.Fatal("expected template content without marker to be rejected")
	}
}

const testStackTemplateContent = `{
  "ZStackTemplateFormatVersion": "2018-06-18",
  "Resources": {}
}`

func TestAccStackTemplateResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckStackTemplateDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_stack_template" "test" {
  name             = "acc-test-stack-template"
  template_content = <<EOT
{
  "ZStackTemplateFormatVersion": "2018-06-18",
  "Resources": {}
}
EOT
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_stack_template.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_stack_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-stack-template")),
				},
			},
			// Step 2: Update name (type is RequiresReplace, name is not)
			{
				Config: providerConfig() + `
resource "zstack_stack_template" "test" {
  name             = "acc-test-stack-template-updated"
  template_content = <<EOT
{
  "ZStackTemplateFormatVersion": "2018-06-18",
  "Resources": {}
}
EOT
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_stack_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-stack-template-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_stack_template.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_stack_template.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
