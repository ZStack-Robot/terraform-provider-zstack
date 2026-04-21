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

func TestLogServerResource_Schema(t *testing.T) {
	var r logServerResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "category", "type", "configuration"}
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

func TestLogServerResource_Metadata(t *testing.T) {
	var r logServerResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_log_server" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccLogServerResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLogServerDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + `
resource "zstack_log_server" "test" {
  name          = "acc-test-log-server"
  category      = "ManagementNodeLog"
  type          = "Log4j2"
  configuration = "{\"host\":\"127.0.0.1\",\"port\":9200}"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_log_server.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_log_server.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-log-server")),
				},
			},
			// Step 2: Update name (category/type/configuration are RequiresReplace)
			{
				Config: providerConfig() + `
resource "zstack_log_server" "test" {
  name          = "acc-test-log-server-updated"
  category      = "ManagementNodeLog"
  type          = "Log4j2"
  configuration = "{\"host\":\"127.0.0.1\",\"port\":9200}"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_log_server.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-log-server-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_log_server.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_log_server.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
