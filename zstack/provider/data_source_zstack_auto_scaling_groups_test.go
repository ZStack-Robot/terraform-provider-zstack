// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAutoScalingGroupDataSource_Schema(t *testing.T) {
	var d autoScalingGroupDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	expectedAttrs := []string{"name", "name_pattern", "auto_scaling_groups"}
	for _, attr := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing attribute %q", attr)
		}
	}
}

func TestAutoScalingGroupDataSource_Metadata(t *testing.T) {
	var d autoScalingGroupDataSource
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_auto_scaling_groups" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccAutoScalingGroupDataSource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_auto_scaling_groups" "test" {
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_auto_scaling_groups.test", tfjsonpath.New("auto_scaling_groups"), knownvalue.NotNull()),
				},
			},
		},
	})
}
