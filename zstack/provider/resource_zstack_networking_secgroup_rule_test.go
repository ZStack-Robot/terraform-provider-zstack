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

func TestSecurityGroupRuleResource_Schema(t *testing.T) {
	var r securityGroupRuleResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "security_group_uuid", "priority", "direction", "action", "state", "ip_version", "protocol", "ip_ranges"}
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

	optional := []string{"description", "destination_port_ranges", "remote_security_group_uuid"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestSecurityGroupRuleResource_Metadata(t *testing.T) {
	var r securityGroupRuleResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_networking_secgroup_rule" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSecurityGroupRuleResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupRuleDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_networking_secgroup" "test_for_rule" {
  name       = "acc-test-secgroup-for-rule"
  ip_version = 4
}

resource "zstack_networking_secgroup_rule" "test" {
  name                    = "acc-test-secgroup-rule"
  security_group_uuid     = zstack_networking_secgroup.test_for_rule.uuid
  priority                = 1
  direction               = "Ingress"
  action                  = "ACCEPT"
  state                   = "Enabled"
  ip_version              = 4
  protocol                = "TCP"
  ip_ranges               = "10.0.0.0/8"
  destination_port_ranges = "80"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckSecurityGroupRuleDisappears("zstack_networking_secgroup_rule.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSecurityGroupRuleResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSecurityGroupRuleDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_networking_secgroup" "test_for_rule" {
  name       = "acc-test-secgroup-for-rule"
  ip_version = 4
}

resource "zstack_networking_secgroup_rule" "test" {
  name                    = "acc-test-secgroup-rule"
  security_group_uuid     = zstack_networking_secgroup.test_for_rule.uuid
  priority                = 1
  direction               = "Ingress"
  action                  = "ACCEPT"
  state                   = "Enabled"
  ip_version              = 4
  protocol                = "TCP"
  ip_ranges               = "10.0.0.0/8"
  destination_port_ranges = "80"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_networking_secgroup_rule.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_rule.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-secgroup-rule")),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_rule.test", tfjsonpath.New("direction"), knownvalue.StringExact("Ingress")),
				},
			},
			// Step 2: Update — modify description and priority (in-place updatable)
			// Note: name, security_group_uuid, direction, ip_version have RequiresReplace
			{
				Config: providerConfig() + `
resource "zstack_networking_secgroup" "test_for_rule" {
  name       = "acc-test-secgroup-for-rule"
  ip_version = 4
}

resource "zstack_networking_secgroup_rule" "test" {
  name                    = "acc-test-secgroup-rule"
  description             = "Updated acceptance test secgroup rule"
  security_group_uuid     = zstack_networking_secgroup.test_for_rule.uuid
  priority                = 2
  direction               = "Ingress"
  action                  = "ACCEPT"
  state                   = "Enabled"
  ip_version              = 4
  protocol                = "TCP"
  ip_ranges               = "10.0.0.0/8"
  destination_port_ranges = "80"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_networking_secgroup_rule.test",
						tfjsonpath.New("description"), knownvalue.StringExact("Updated acceptance test secgroup rule")),
					statecheck.ExpectKnownValue("zstack_networking_secgroup_rule.test",
						tfjsonpath.New("priority"), knownvalue.Int32Exact(2)),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_networking_secgroup_rule.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_networking_secgroup_rule.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"name"},
			},
		},
	})
}
