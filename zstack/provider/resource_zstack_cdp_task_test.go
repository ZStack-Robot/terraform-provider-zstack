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

func TestCdpTaskResource_Schema(t *testing.T) {
	var r cdpTaskResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "task_type", "policy_uuid", "backup_storage_uuid", "resource_uuids"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "status", "state"}
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

func TestCdpTaskResource_Metadata(t *testing.T) {
	var r cdpTaskResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_cdp_task" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccCdpTaskResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.BackupStorages) == 0 {
		t.Fatal("env.json has no backup_storages")
	}
	bsUUID := env.BackupStorages[0]["uuid"].(string)

	// Find a Ready volume for resource_uuids
	var volumeUUID string
	for _, v := range env.Volumes {
		if v["status"] == "Ready" {
			volumeUUID = v["uuid"].(string)
			break
		}
	}
	if volumeUUID == "" {
		t.Fatal("env.json has no Ready volume for cdp_task resource_uuids")
	}

	createConfig := providerConfig() + fmt.Sprintf(`
resource "zstack_cdp_policy" "dep" {
  name                      = "acc-test-cdp-policy-for-task"
  recovery_point_per_second = 1
}

resource "zstack_cdp_task" "test" {
  name               = "acc-test-cdp-task"
  task_type          = "Volume"
  policy_uuid        = zstack_cdp_policy.dep.uuid
  backup_storage_uuid = "%s"
  resource_uuids     = ["%s"]
}
`, bsUUID, volumeUUID)

	updateConfig := providerConfig() + fmt.Sprintf(`
resource "zstack_cdp_policy" "dep" {
  name                      = "acc-test-cdp-policy-for-task"
  recovery_point_per_second = 1
}

resource "zstack_cdp_task" "test" {
  name               = "acc-test-cdp-task-updated"
  task_type          = "Volume"
  policy_uuid        = zstack_cdp_policy.dep.uuid
  backup_storage_uuid = "%s"
  resource_uuids     = ["%s"]
}
`, bsUUID, volumeUUID)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCdpTaskDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create
			{
				Config: createConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cdp_task.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_cdp_task.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cdp-task")),
				},
			},
			// Step 2: Update name (in-place, only name is not ForceNew)
			{
				Config: updateConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_cdp_task.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-cdp-task-updated")),
				},
			},
			// Step 3: Import
			{
				ResourceName:                         "zstack_cdp_task.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_cdp_task.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
