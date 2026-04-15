// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestCertificateResource_Schema(t *testing.T) {
	var r certificateResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "certificate"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid"}
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

func TestCertificateResource_Metadata(t *testing.T) {
	var r certificateResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_certificate" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccCertificateResource(t *testing.T) {
	_ = loadEnvData(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCertificateDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + `
resource "zstack_certificate" "test" {
  name        = "acc-test-certificate"
  certificate = "-----BEGIN CERTIFICATE-----\nMIIBkTCB+wIUEh8m0yZ4X1GbWOKQoSqxVh7gkfUwDQYJKoZIhvcNAQELBQAwEjEQ\nMA4GA1UEAwwHdGVzdC1jYTAeFw0yNDA0MTMwMDAwMDBaFw0zNDA0MTMwMDAwMDBa\nMBIxEDAOBgNVBAMMB3Rlc3QtY2EwXDANBgkqhkiG9w0BAQEFAANLADBIAkEA0Z3V\nS3MwRXfHOVMnz0pHRvPqNsLffO9DeXSGPnkHMWVFnkFPnAGI+ZhouBnfZMwBY0Mj\nJmpRGsXSlYMsqDMNFwIDAQABoyMwITAfBgNVHREEGDAWhwR/AAABhwTAqAEBhwQK\nAAEBMA0GCSqGSIb4DQEBCQUAA0EAkcPGWFv43IkajKMmp/CjTOrLEMFSiHBr7hIK\nBGVEkVDqMe/dIYMe+bpFZ23LalFa4p27iE1uchIMjjkE2adkJQ==\n-----END CERTIFICATE-----"
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_certificate.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_certificate.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-certificate")),
				},
			},
			{
				ResourceName:                        "zstack_certificate.test",
				ImportState:                         true,
				ImportStateIdFunc:                   importStateIdFromUUID("zstack_certificate.test"),
				ImportStateVerify:                   true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"certificate"},
			},
		},
	})
}
