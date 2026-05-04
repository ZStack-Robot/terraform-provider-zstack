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

func TestAlarmResource_Schema(t *testing.T) {
	var r alarmResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "comparison_operator", "namespace", "metric_name", "threshold"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "status", "state"}
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

func TestAlarmResource_Metadata(t *testing.T) {
	var r alarmResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_alarm" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAlarmResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAlarmDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_alarm" "test" {
  name                = "acc-test-alarm-disappears"
  comparison_operator = "GreaterThanOrEqualTo"
  namespace           = "ZStack/VM"
  metric_name         = "CPUAverageUsedUtilization"
  threshold           = 90
  period              = 60
  repeat_interval     = 600
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_alarm.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckAlarmDisappears("zstack_alarm.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAlarmResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAlarmDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_alarm" "test" {
  name                = "acc-test-alarm"
  comparison_operator = "GreaterThanOrEqualTo"
  namespace           = "ZStack/VM"
  metric_name         = "CPUAverageUsedUtilization"
  threshold           = 90
  period              = 60
  repeat_interval     = 600
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_alarm.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_alarm.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-alarm")),
				},
			},
			// Step 2: Update — modify name and description
			{
				Config: providerConfig() + `
resource "zstack_alarm" "test" {
  name                = "acc-test-alarm-updated"
  description         = "Updated acceptance test alarm"
  comparison_operator = "GreaterThanOrEqualTo"
  namespace           = "ZStack/VM"
  metric_name         = "CPUAverageUsedUtilization"
  threshold           = 90
  period              = 60
  repeat_interval     = 600
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_alarm.test",
						tfjsonpath.New("name"), knownvalue.StringExact("acc-test-alarm-updated")),
					statecheck.ExpectKnownValue("zstack_alarm.test",
						tfjsonpath.New("description"), knownvalue.StringExact("Updated acceptance test alarm")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_alarm.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_alarm.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
