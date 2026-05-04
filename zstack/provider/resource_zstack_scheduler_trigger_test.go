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

func TestSchedulerTriggerResource_Schema(t *testing.T) {
	var r schedulerTriggerResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "scheduler_type"}
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
	computed := []string{"uuid", "stop_time"}
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

func TestSchedulerTriggerResource_Metadata(t *testing.T) {
	var r schedulerTriggerResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_scheduler_trigger" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSchedulerTriggerResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchedulerTriggerDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_scheduler_trigger" "test" {
  name           = "acc-test-trigger"
  scheduler_type = "cron"
  cron           = "0 0 0 * * ?"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckSchedulerTriggerDisappears("zstack_scheduler_trigger.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSchedulerTriggerResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("trigger")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchedulerTriggerDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_scheduler_trigger" "test" {
  name           = %q
  description    = "acceptance scheduler trigger"
  scheduler_type = "cron"
  cron           = "0 0 0 * * ?"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance scheduler trigger")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_scheduler_trigger" "test" {
  name           = %q
  description    = "acceptance scheduler trigger updated"
  scheduler_type = "cron"
  cron           = "0 5 0 * * ?"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance scheduler trigger updated")),
					statecheck.ExpectKnownValue("zstack_scheduler_trigger.test", tfjsonpath.New("cron"), knownvalue.StringExact("0 5 0 * * ?")),
				},
			},
			{
				ResourceName:                         "zstack_scheduler_trigger.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_scheduler_trigger.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
