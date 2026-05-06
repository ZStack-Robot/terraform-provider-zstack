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

func TestAccessKeyResource_Schema(t *testing.T) {
	var r accessKeyResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"account_uuid", "user_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "access_key_id", "access_key_secret", "state"}
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

func TestAccessKeyResource_Metadata(t *testing.T) {
	var r accessKeyResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_access_key" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAccessKeyResource(t *testing.T) {
	envData := loadEnvData(t)

	if len(envData.Accounts) == 0 {
		t.Skip("no accounts in env.json")
	}
	if len(envData.Users) == 0 {
		t.Skip("no users in env.json")
	}

	accountUUID := envData.Accounts[0]["uuid"].(string)
	userUUID := envData.Users[0]["uuid"].(string)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAccessKeyDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_access_key" "test" {
  account_uuid = %q
  user_uuid    = %q
}
`, accountUUID, userUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_access_key.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_access_key.test", tfjsonpath.New("access_key_id"), knownvalue.NotNull()),
				},
			},
			// Step 2: Import
			{
				ResourceName:                         "zstack_access_key.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_access_key.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"access_key_secret"},
			},
		},
	})
}
