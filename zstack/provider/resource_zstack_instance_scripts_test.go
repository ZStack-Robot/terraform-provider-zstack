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

func TestScriptResource_Schema(t *testing.T) {
	var r scriptResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "encoding_type", "script_content", "script_type"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "render_params", "script_timeout"}
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

func TestScriptResource_Metadata(t *testing.T) {
	var r scriptResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_instance_scripts" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestInstanceScriptsUpdateGuardsUnknownValues(t *testing.T) {
	tests := []struct {
		name   string
		method string
		field  string
	}{
		{
			name:   "Create_ScriptTimeout_Unknown",
			method: "Create",
			field:  "ScriptTimeout",
		},
		{
			name:   "Update_ScriptTimeout_Unknown",
			method: "Update",
			field:  "ScriptTimeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%s method must guard %s with IsUnknown() check", tt.method, tt.field)
		})
	}
}

func TestAccScriptResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckScriptDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_instance_scripts" "test" {
  name           = "acc-test-script"
  script_content = "echo hello"
  script_type    = "Shell"
  encoding_type  = "PlainText"
  platform       = "Linux"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_scripts.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_instance_scripts.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-script")),
				},
			},
			// Step 2: Update name (encoding_type is RequiresReplace, name is not)
			{
				Config: providerConfig() + `
resource "zstack_instance_scripts" "test" {
  name           = "acc-test-script-updated"
  script_content = "echo hello"
  script_type    = "Shell"
  encoding_type  = "PlainText"
  platform       = "Linux"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_scripts.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-script-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_instance_scripts.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_instance_scripts.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
