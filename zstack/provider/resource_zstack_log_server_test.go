// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	required := []string{"name", "category", "type"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}
	optional := []string{"configuration", "appender_type", "appender_configuration"}
	for _, attr := range optional {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
		if !a.IsOptional() {
			t.Errorf("attribute %q should be optional", attr)
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

func TestBuildLogServerConfigurationFromAppender(t *testing.T) {
	appenderConfig, diags := types.MapValue(types.StringType, map[string]attr.Value{
		"hostname": types.StringValue("127.0.0.1"),
		"port":     types.StringValue("514"),
		"protocol": types.StringValue("UDP"),
		"facility": types.StringValue("LOCAL5"),
	})
	if diags.HasError() {
		t.Fatalf("unexpected map diagnostics: %v", diags)
	}

	var buildDiags diag.Diagnostics
	got, ok := buildLogServerConfiguration(context.Background(), logServerModel{
		AppenderType:   types.StringValue("Syslog"),
		AppenderConfig: appenderConfig,
	}, &buildDiags)
	if !ok || buildDiags.HasError() {
		t.Fatalf("buildLogServerConfiguration failed: %v", buildDiags)
	}

	want := `{"appenderType":"Syslog","configuration":{"facility":"LOCAL5","hostname":"127.0.0.1","port":"514","protocol":"UDP"}}`
	if got != want {
		t.Fatalf("unexpected configuration:\n got: %s\nwant: %s", got, want)
	}
}

func TestBuildLogServerConfigurationRejectsSimpleRawJSON(t *testing.T) {
	var buildDiags diag.Diagnostics
	_, ok := buildLogServerConfiguration(context.Background(), logServerModel{
		Configuration: types.StringValue(`{"host":"127.0.0.1","port":9200}`),
	}, &buildDiags)
	if ok {
		t.Fatal("expected simple raw JSON to be rejected")
	}
	if !buildDiags.HasError() {
		t.Fatal("expected diagnostics for simple raw JSON")
	}
}

func TestBuildLogServerConfigurationAcceptsRawNestedJSON(t *testing.T) {
	var buildDiags diag.Diagnostics
	got, ok := buildLogServerConfiguration(context.Background(), logServerModel{
		Configuration: types.StringValue(`{"appenderType":"Syslog","configuration":{"hostname":"127.0.0.1","port":"514"}}`),
	}, &buildDiags)
	if !ok || buildDiags.HasError() {
		t.Fatalf("expected raw nested JSON to be accepted: %v", buildDiags)
	}
	if got == "" {
		t.Fatal("expected non-empty configuration")
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
  appender_type = "Syslog"
  appender_configuration = {
    hostname = "127.0.0.1"
    port     = "514"
    protocol = "UDP"
    facility = "LOCAL5"
  }
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
  appender_type = "Syslog"
  appender_configuration = {
    hostname = "127.0.0.1"
    port     = "514"
    protocol = "UDP"
    facility = "LOCAL5"
  }
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
