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

func TestSnsTopicResource_Schema(t *testing.T) {
	var r snsTopicResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
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
	// Check computed attributes
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

func TestSnsTopicResource_Metadata(t *testing.T) {
	var r snsTopicResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_sns_topic" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSNSTopicResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSNSTopicDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_sns_topic" "test" {
  name = "acc-test-sns-topic"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckSNSTopicDisappears("zstack_sns_topic.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSNSTopicResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("sns-topic")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSNSTopicDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_sns_topic" "test" {
  name        = %q
  description = "acceptance SNS topic"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_topic.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_sns_topic.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_sns_topic.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SNS topic")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_sns_topic" "test" {
  name        = %q
  description = "acceptance SNS topic updated"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_sns_topic.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_sns_topic.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance SNS topic updated")),
				},
			},
			{
				ResourceName:                         "zstack_sns_topic.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_sns_topic.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
