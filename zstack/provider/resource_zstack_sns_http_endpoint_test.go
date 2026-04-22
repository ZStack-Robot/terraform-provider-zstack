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

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSNSHttpEndpointDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_sns_http_endpoint" "test" {
  name = "acc-test-http-endpoint"
  url  = "http://example.com/sns"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_sns_http_endpoint.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-http-endpoint")),
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
