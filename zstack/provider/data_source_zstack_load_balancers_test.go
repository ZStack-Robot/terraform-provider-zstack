// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLoadBalancerDataSource_Schema(t *testing.T) {
	var d loadBalancerDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	expectedAttrs := []string{"name", "name_pattern", "load_balancers"}
	for _, attr := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing attribute %q", attr)
		}
	}
}

func TestLoadBalancerDataSource_Metadata(t *testing.T) {
	var d loadBalancerDataSource
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_load_balancers" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccLoadBalancerDataSource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping load balancer data source acceptance test")
	}

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_load_balancers" "test" {
}
`,
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("data.zstack_load_balancers.test", "load_balancers.#"),
				),
			},
		},
	})
}
