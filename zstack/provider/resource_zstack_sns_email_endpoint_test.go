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

func TestSnsEmailEndpointResource_Schema(t *testing.T) {
	var r snsEmailEndpointResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "email", "platform_uuid"}
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

func TestSnsEmailEndpointResource_Metadata(t *testing.T) {
	var r snsEmailEndpointResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_sns_email_endpoint" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSNSEmailEndpointResource(t *testing.T) {
	envData := loadEnvData(t)
	if len(envData.SNSEmailPlatforms) == 0 {
		t.Skip("no SNS email platforms in env.json")
	}
	platformUUID := envData.SNSEmailPlatforms[0]["uuid"].(string)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSNSEmailEndpointDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_sns_email_endpoint" "test" {
  name          = "acc-test-email-endpoint"
  email         = "test@example.com"
  platform_uuid = %q
}
`, platformUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_email_endpoint.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_sns_email_endpoint.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-email-endpoint")),
				},
			},
			{
				ResourceName:                         "zstack_sns_email_endpoint.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_sns_email_endpoint.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
