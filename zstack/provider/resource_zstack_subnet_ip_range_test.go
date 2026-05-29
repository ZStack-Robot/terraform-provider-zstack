// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestSubnetIpRangeResource_Schema(t *testing.T) {
	var r subnetResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"l3_network_uuid", "name", "start_ip", "end_ip", "netmask"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	if _, ok := resp.Schema.Attributes["uuid"]; !ok {
		t.Fatal("schema missing computed attribute uuid")
	}

	optional := []string{"gateway", "ip_range_type"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestSubnetIpRangeResource_IpRangeTypeValidator(t *testing.T) {
	var r subnetResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	attr, ok := resp.Schema.Attributes["ip_range_type"].(rschema.StringAttribute)
	if !ok {
		t.Fatal("ip_range_type should be a string attribute")
	}
	if !attr.Optional {
		t.Fatal("ip_range_type should be optional")
	}
	if len(attr.Validators) == 0 {
		t.Fatal("ip_range_type should have validators")
	}

	tests := []struct {
		value       string
		expectError bool
	}{
		{value: "Normal"},
		{value: "AddressPool"},
		{value: "Reserved", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			req := validator.StringRequest{ConfigValue: types.StringValue(tt.value)}
			res := validator.StringResponse{}
			attr.Validators[0].ValidateString(context.Background(), req, &res)

			if res.Diagnostics.HasError() != tt.expectError {
				t.Fatalf("unexpected diagnostics for %q: %v", tt.value, res.Diagnostics)
			}
		})
	}
}

func TestSubnetAddIpRangeParamOmitsNullIpRangeType(t *testing.T) {
	plan := subnetModel{
		Name:        types.StringValue("test-subnet"),
		StartIp:     types.StringValue("192.168.50.2"),
		EndIp:       types.StringValue("192.168.50.254"),
		Netmask:     types.StringValue("255.255.255.0"),
		Gateway:     types.StringValue("192.168.50.1"),
		IpRangeType: types.StringNull(),
	}

	p := subnetAddIpRangeParam(plan)
	if p.Params.IpRangeType != nil {
		t.Fatalf("expected nil ipRangeType, got %q", *p.Params.IpRangeType)
	}

	payload, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal param: %v", err)
	}
	if strings.Contains(string(payload), "ipRangeType") {
		t.Fatalf("expected ipRangeType to be omitted from JSON, got %s", payload)
	}
}

func TestSubnetAddIpRangeParamIncludesConfiguredIpRangeType(t *testing.T) {
	plan := subnetModel{
		Name:        types.StringValue("test-subnet"),
		StartIp:     types.StringValue("192.168.50.2"),
		EndIp:       types.StringValue("192.168.50.254"),
		Netmask:     types.StringValue("255.255.255.0"),
		Gateway:     types.StringValue("192.168.50.1"),
		IpRangeType: types.StringValue("AddressPool"),
	}

	p := subnetAddIpRangeParam(plan)
	if p.Params.IpRangeType == nil {
		t.Fatal("expected ipRangeType to be set")
	}
	if *p.Params.IpRangeType != "AddressPool" {
		t.Fatalf("unexpected ipRangeType: %q", *p.Params.IpRangeType)
	}
}

func TestSubnetIpRangeResource_Metadata(t *testing.T) {
	var r subnetResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_subnet_ip_range" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccSubnetIpRangeResource(t *testing.T) {
	env := loadEnvData(t)

	if len(env.L2Networks) == 0 {
		t.Skip("no l2_networks in env.json, skipping subnet ip range acceptance test")
	}

	// Use the L2NoVlanNetwork to create a fresh L3 for this test
	var l2UUID string
	for _, l2 := range env.L2Networks {
		if envStr(l2, "type") == "L2NoVlanNetwork" {
			l2UUID = envStr(l2, "uuid")
			break
		}
	}
	if l2UUID == "" {
		l2UUID = envStr(env.L2Networks[0], "uuid")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSubnetIpRangeDestroy,
		Steps: []tfresource.TestStep{
			// Step 1: Create (all attrs are ForceNew, no update step).
			// The backend may omit ipRangeType from QueryIpRange; the provider should keep
			// the configured value in state to avoid a replacement-only refresh plan.
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_l3network" "for_subnet" {
  name            = "acc-test-l3-for-subnet"
  l2_network_uuid = %q
  category        = "Private"
  ip_version      = 4
}

resource "zstack_subnet_ip_range" "test" {
  name            = "acc-test-subnet"
  l3_network_uuid = zstack_l3network.for_subnet.uuid
  start_ip        = "192.168.50.2"
  end_ip          = "192.168.50.254"
  netmask         = "255.255.255.0"
  gateway         = "192.168.50.1"
  ip_range_type   = "Normal"
}
`, l2UUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("name"), knownvalue.StringExact("acc-test-subnet")),
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("start_ip"), knownvalue.StringExact("192.168.50.2")),
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("end_ip"), knownvalue.StringExact("192.168.50.254")),
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("netmask"), knownvalue.StringExact("255.255.255.0")),
					statecheck.ExpectKnownValue("zstack_subnet_ip_range.test", tfjsonpath.New("gateway"), knownvalue.StringExact("192.168.50.1")),
				},
			},
			// Step 2: Import (ip_range_type is not returned by API, skip verify)
			{
				ResourceName:                         "zstack_subnet_ip_range.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_subnet_ip_range.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"ip_range_type"},
			},
		},
	})
}
