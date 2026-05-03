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

func TestGlobalConfigsDataSource_Schema(t *testing.T) {
	var d globalConfigsDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	expectedAttrs := []string{"category", "name", "name_pattern", "global_configs"}
	for _, attr := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing attribute %q", attr)
		}
	}
}

func TestGlobalConfigsDataSource_Metadata(t *testing.T) {
	var d globalConfigsDataSource
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_global_configs" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccGlobalConfigsDataSource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_global_configs" "test" {
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_global_configs.test", tfjsonpath.New("global_configs"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccGlobalConfigsDataSourceFilterByCategoryAndName(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_global_configs" "test" {
  category = "vm"
  name     = "deletionPolicy"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_global_configs.test", tfjsonpath.New("global_configs"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_global_configs.test", tfjsonpath.New("global_configs").AtSliceIndex(0).AtMapKey("category"), knownvalue.StringExact("vm")),
					statecheck.ExpectKnownValue("data.zstack_global_configs.test", tfjsonpath.New("global_configs").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact("deletionPolicy")),
				},
			},
		},
	})
}
