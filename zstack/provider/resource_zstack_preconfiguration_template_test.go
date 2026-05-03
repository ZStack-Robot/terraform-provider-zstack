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

func TestPreconfigurationTemplateResource_Schema(t *testing.T) {
	var r preconfigurationTemplateResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "distribution", "type", "content"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "md5sum", "is_predefined", "state"}
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

func TestPreconfigurationTemplateResource_Metadata(t *testing.T) {
	var r preconfigurationTemplateResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_preconfiguration_template" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

const testPreconfigurationTemplateContent = `# Preconfiguration template smoke test
# Base ZStack system variables required by the API:
# REPO_URL
# USERNAME
# PASSWORD
# NETWORK_CFGS
# FORCE_INSTALL
# PRE_SCRIPTS
# POST_SCRIPTS
install
text
`

func TestAccPreconfigurationTemplateResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPreconfigurationTemplateDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_preconfiguration_template" "test" {
  name         = "acc-test-preconfig-template"
  distribution = "CentOS7"
  type         = "kickstart"
  content      = %q
}
`, testPreconfigurationTemplateContent),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_preconfiguration_template.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_preconfiguration_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-preconfig-template")),
				},
			},
			// Step 2: Update name (note: RequiresReplace pending story-07, triggers destroy+recreate)
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_preconfiguration_template" "test" {
  name         = "acc-test-preconfig-template-updated"
  distribution = "CentOS7"
  type         = "kickstart"
  content      = %q
}
`, testPreconfigurationTemplateContent),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_preconfiguration_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-preconfig-template-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_preconfiguration_template.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_preconfiguration_template.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
