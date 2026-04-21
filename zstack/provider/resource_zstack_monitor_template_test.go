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

func TestMonitorTemplateResource_Schema(t *testing.T) {
	var r monitorTemplateResource
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

func TestMonitorTemplateResource_Metadata(t *testing.T) {
	var r monitorTemplateResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_monitor_template" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccMonitorTemplateResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMonitorTemplateDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_monitor_template" "test" {
  name = "acc-test-monitor-template"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_monitor_template.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_monitor_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-monitor-template")),
				},
			},
			// Step 2: Update name (note: RequiresReplace pending story-07, triggers destroy+recreate)
			{
				Config: providerConfig() + `
resource "zstack_monitor_template" "test" {
  name = "acc-test-monitor-template-updated"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_monitor_template.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-monitor-template-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_monitor_template.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_monitor_template.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
