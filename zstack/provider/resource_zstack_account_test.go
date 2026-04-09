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

func TestAccountResource_Schema(t *testing.T) {
	var r accountResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "password"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "type"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}

	// Verify password is sensitive
	if pwAttr, ok := resp.Schema.Attributes["password"]; ok {
		if !pwAttr.IsSensitive() {
			t.Error("attribute \"password\" should be sensitive")
		}
	}
}

func TestAccountResource_Metadata(t *testing.T) {
	var r accountResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_account" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAccountResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccountDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_account" "test" {
  name     = "acc-test-account"
  password = "Test@12345"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_account.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_account.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-account")),
				},
			},
			{
				ResourceName:            "zstack_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAccountResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccountDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_account" "test_disappears" {
  name     = "acc-test-acct-disappears"
  password = "Test@12345"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_account.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckAccountDisappears("zstack_account.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
