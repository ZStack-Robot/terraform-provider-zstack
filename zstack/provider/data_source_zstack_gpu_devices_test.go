// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestGpuDeviceDataSource_Schema(t *testing.T) {
	var d gpuDeviceDataSource
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	expectedAttrs := []string{"name", "name_pattern", "gpu_devices"}
	for _, attr := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing attribute %q", attr)
		}
	}
}

func TestGpuDeviceDataSource_Metadata(t *testing.T) {
	var d gpuDeviceDataSource
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_gpu_devices" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccGpuDeviceDataSource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
data "zstack_gpu_devices" "test" {
}
`,
				Check: tfresource.ComposeAggregateTestCheckFunc(
					tfresource.TestCheckResourceAttrSet("data.zstack_gpu_devices.test", "gpu_devices.#"),
				),
			},
		},
	})
}
