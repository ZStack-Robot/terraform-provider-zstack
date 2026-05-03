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

func TestSchedulerJobResource_Schema(t *testing.T) {
	var r schedulerJobResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "target_resource_uuid", "type"}
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
	computed := []string{"uuid", "state", "job_data", "job_class_name"}
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

func TestSchedulerJobResource_Metadata(t *testing.T) {
	var r schedulerJobResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_scheduler_job" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSchedulerJobResource_disappears(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm_instances in env data, required for scheduler job")
	}
	vmUUID := envStr(env.VmInstances[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchedulerJobDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_scheduler_job" "test" {
  name                 = "acc-test-scheduler-job"
  target_resource_uuid = %q
  type                 = "stopVm"
}
`, vmUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckSchedulerJobDisappears("zstack_scheduler_job.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSchedulerJobResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.VmInstances) == 0 {
		t.Skip("no vm_instances in env data, required for scheduler job")
	}
	vmUUID := envStr(env.VmInstances[0], "uuid")
	name := testAccName("scheduler-job")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSchedulerJobDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_scheduler_job" "test" {
  name                 = %q
  description          = "acceptance scheduler job"
  target_resource_uuid = %q
  type                 = "stopVm"
}
`, name, vmUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_scheduler_job.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_scheduler_job.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_scheduler_job.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance scheduler job")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_scheduler_job" "test" {
  name                 = %q
  description          = "acceptance scheduler job updated"
  target_resource_uuid = %q
  type                 = "stopVm"
}
`, updatedName, vmUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_scheduler_job.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_scheduler_job.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance scheduler job updated")),
				},
			},
		},
	})
}
