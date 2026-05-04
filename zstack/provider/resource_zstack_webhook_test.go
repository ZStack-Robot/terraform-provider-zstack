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

func TestWebhookResource_Schema(t *testing.T) {
	var r webhookResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "url", "type"}
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

func TestWebhookResource_Metadata(t *testing.T) {
	var r webhookResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_webhook" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccWebhookResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_webhook" "test" {
  name = "acc-test-webhook"
  url  = "http://example.com/webhook"
  type = "custom"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckWebhookDisappears("zstack_webhook.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWebhookResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("webhook")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckWebhookDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_webhook" "test" {
  name        = %q
  description = "acceptance webhook"
  url         = "http://example.com/webhook"
  type        = "custom"
  opaque      = "create"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance webhook")),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("url"), knownvalue.StringExact("http://example.com/webhook")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_webhook" "test" {
  name        = %q
  description = "acceptance webhook updated"
  url         = "http://example.com/webhook-updated"
  type        = "custom"
  opaque      = "update"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance webhook updated")),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("url"), knownvalue.StringExact("http://example.com/webhook-updated")),
					statecheck.ExpectKnownValue("zstack_webhook.test", tfjsonpath.New("opaque"), knownvalue.StringExact("update")),
				},
			},
			{
				ResourceName:                         "zstack_webhook.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_webhook.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
