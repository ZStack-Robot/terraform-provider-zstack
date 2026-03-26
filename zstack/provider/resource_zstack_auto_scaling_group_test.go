// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAutoScalingGroupResource_Schema(t *testing.T) {
	var r autoScalingGroupResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "scaling_resource_type", "default_cooldown", "min_resource_size", "max_resource_size", "removal_policy"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

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

	optional := []string{"description"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestAutoScalingGroupResource_Metadata(t *testing.T) {
	var r autoScalingGroupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_auto_scaling_group" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAutoScalingGroupResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_auto_scaling_group" "test" {
  name                  = "acc-test-scaling-group"
  scaling_resource_type = "VmInstance"
  default_cooldown      = 60
  min_resource_size     = 0
  max_resource_size     = 5
  removal_policy        = "OldestInstance"
}
`,
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("zstack_auto_scaling_group.test", "uuid"),
					tfresource.TestCheckResourceAttr("zstack_auto_scaling_group.test", "name", "acc-test-scaling-group"),
					tfresource.TestCheckResourceAttr("zstack_auto_scaling_group.test", "scaling_resource_type", "VmInstance"),
					tfresource.TestCheckResourceAttr("zstack_auto_scaling_group.test", "min_resource_size", "0"),
					tfresource.TestCheckResourceAttr("zstack_auto_scaling_group.test", "max_resource_size", "5"),
				),
			},
		},
	})
}
