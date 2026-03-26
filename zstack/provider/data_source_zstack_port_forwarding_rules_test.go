// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPortForwardingRuleDataSource_Schema(t *testing.T) {
	var d portForwardingRuleDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	optional := []string{"name", "name_pattern"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}

	if _, ok := resp.Schema.Attributes["port_forwarding_rules"]; !ok {
		t.Fatal("schema missing computed attribute port_forwarding_rules")
	}
}

func TestPortForwardingRuleDataSource_Metadata(t *testing.T) {
	var d portForwardingRuleDataSource
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_port_forwarding_rules" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccPortForwardingRuleDataSource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_port_forwarding_rules" "test" {
}
`,
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("data.zstack_port_forwarding_rules.test", "port_forwarding_rules.#"),
				),
			},
		},
	})
}
