// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLicenseAuthorizedNodesDataSource_Schema(t *testing.T) {
	var d licenseAuthorizedNodeDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	if _, ok := resp.Schema.Attributes["uuid"]; !ok {
		t.Fatal("schema should include uuid lookup")
	}
	if _, ok := resp.Schema.Attributes["nodes"]; !ok {
		t.Fatal("schema should include computed nodes")
	}
	if _, ok := resp.Schema.Attributes["name"]; ok {
		t.Fatal("schema should not expose unsupported name filter")
	}
	if _, ok := resp.Schema.Attributes["name_pattern"]; ok {
		t.Fatal("schema should not expose unsupported name_pattern filter")
	}
}

func TestAccZStackLicenseAuthorizedNodesDataSource(t *testing.T) {
	loadEnvData(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `data "zstack_license_authorized_nodes" "test" {}`,
			},
		},
	})
}
