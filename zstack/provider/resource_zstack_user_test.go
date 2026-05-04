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

func TestUserResource_Schema(t *testing.T) {
	var r userResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
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
	// Check computed attributes
	computed := []string{"uuid", "account_uuid"}
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

func TestUserResource_Metadata(t *testing.T) {
	var r userResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_user" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccUserResource_disappears(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_user" "test" {
  name     = "acc-test-user"
  password = "Test@12345"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckUserDisappears("zstack_user.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccUserResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("user")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_user" "test" {
  name        = %q
  password    = "Test@12345"
  description = "acceptance user"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_user.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_user.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_user.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance user")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_user" "test" {
  name        = %q
  password    = "Test@12345"
  description = "acceptance user updated"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_user.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_user.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance user updated")),
				},
			},
			{
				ResourceName:                         "zstack_user.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_user.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"password"},
			},
		},
	})
}
