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

func TestMonitorGroupResource_Schema(t *testing.T) {
	var r monitorGroupResource
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

func TestMonitorGroupResource_Metadata(t *testing.T) {
	var r monitorGroupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_monitor_group" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccMonitorGroupResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckMonitorGroupDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_monitor_group" "test" {
  name = "acc-test-monitor-group"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_monitor_group.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_monitor_group.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-monitor-group")),
				},
			},
			// Step 2: Update name (RequiresReplace — triggers destroy+recreate)
			{
				Config: providerConfig() + `
resource "zstack_monitor_group" "test" {
  name = "acc-test-monitor-group-updated"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_monitor_group.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-monitor-group-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_monitor_group.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_monitor_group.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
