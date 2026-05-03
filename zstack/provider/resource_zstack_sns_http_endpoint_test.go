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

func TestSnsHttpEndpointResource_Schema(t *testing.T) {
	var r snsHttpEndpointResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "url"}
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
	computed := []string{"uuid", "type", "state"}
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

func TestSnsHttpEndpointResource_Metadata(t *testing.T) {
	var r snsHttpEndpointResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_sns_http_endpoint" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSNSHttpEndpointResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("http-endpoint")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSNSHttpEndpointDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_sns_http_endpoint" "test" {
  name        = %q
  description = "acceptance SNS HTTP endpoint"
  url         = "http://example.com/sns"
  username    = "create-user"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SNS HTTP endpoint")),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("url"), knownvalue.StringExact("http://example.com/sns")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_sns_http_endpoint" "test" {
  name        = %q
  description = "acceptance SNS HTTP endpoint updated"
  url         = "http://example.com/sns-updated"
  username    = "update-user"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SNS HTTP endpoint updated")),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("url"), knownvalue.StringExact("http://example.com/sns-updated")),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("username"), knownvalue.StringExact("update-user")),
				},
			},
			{
				ResourceName:                         "zstack_sns_http_endpoint.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_sns_http_endpoint.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"password"},
			},
		},
	})
}
