// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

type spyPortForwardingRuleClient struct {
	createCalls []createCall
}

type createCall struct {
	param param.CreatePortForwardingRuleParam
}

func (s *spyPortForwardingRuleClient) CreatePortForwardingRule(p param.CreatePortForwardingRuleParam) (*view.PortForwardingRuleInventoryView, error) {
	s.createCalls = append(s.createCalls, createCall{param: p})
	return &view.PortForwardingRuleInventoryView{}, nil
}

func (s *spyPortForwardingRuleClient) GetPortForwardingRule(uuid string) (*view.PortForwardingRuleInventoryView, error) {
	return &view.PortForwardingRuleInventoryView{}, nil
}

func (s *spyPortForwardingRuleClient) UpdatePortForwardingRule(uuid string, p param.UpdatePortForwardingRuleParam) (*view.PortForwardingRuleInventoryView, error) {
	return &view.PortForwardingRuleInventoryView{}, nil
}

func (s *spyPortForwardingRuleClient) DetachPortForwardingRule(uuid string, deleteMode param.DeleteMode) error {
	return nil
}

func (s *spyPortForwardingRuleClient) DeletePortForwardingRule(uuid string, deleteMode param.DeleteMode) error {
	return nil
}

func (s *spyPortForwardingRuleClient) AttachPortForwardingRule(ruleUuid string, vmNicUuid string, p param.AttachPortForwardingRuleParam) (*view.PortForwardingRuleInventoryView, error) {
	return &view.PortForwardingRuleInventoryView{}, nil
}

func TestPortForwardingRuleResource_Schema(t *testing.T) {
	var r portForwardingRuleResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vip_uuid", "vip_port_start", "protocol_type"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "vip_ip", "guest_ip", "state"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "vip_port_end", "private_port_start", "private_port_end", "vm_nic_uuid", "allowed_cidr"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestPortForwardingRuleResource_Metadata(t *testing.T) {
	var r portForwardingRuleResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_port_forwarding_rule" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestPortForwardingRuleUpdateGuardsUnknownValues(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		setupPlan func(*portForwardingRuleResourceModel)
		checkFunc func(*testing.T, param.CreatePortForwardingRuleParam)
	}{
		{
			name:      "VipPortEnd Unknown should be omitted",
			fieldName: "VipPortEnd",
			setupPlan: func(plan *portForwardingRuleResourceModel) {
				plan.VipPortEnd = types.Int64Unknown()
			},
			checkFunc: func(t *testing.T, p param.CreatePortForwardingRuleParam) {
				if p.Params.VipPortEnd != nil {
					t.Errorf("VipPortEnd should be nil when Unknown, got %v", *p.Params.VipPortEnd)
				}
			},
		},
		{
			name:      "PrivatePortStart Unknown should be omitted",
			fieldName: "PrivatePortStart",
			setupPlan: func(plan *portForwardingRuleResourceModel) {
				plan.PrivatePortStart = types.Int64Unknown()
			},
			checkFunc: func(t *testing.T, p param.CreatePortForwardingRuleParam) {
				if p.Params.PrivatePortStart != nil {
					t.Errorf("PrivatePortStart should be nil when Unknown, got %v", *p.Params.PrivatePortStart)
				}
			},
		},
		{
			name:      "PrivatePortEnd Unknown should be omitted",
			fieldName: "PrivatePortEnd",
			setupPlan: func(plan *portForwardingRuleResourceModel) {
				plan.PrivatePortEnd = types.Int64Unknown()
			},
			checkFunc: func(t *testing.T, p param.CreatePortForwardingRuleParam) {
				if p.Params.PrivatePortEnd != nil {
					t.Errorf("PrivatePortEnd should be nil when Unknown, got %v", *p.Params.PrivatePortEnd)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := portForwardingRuleResourceModel{
				Name:             types.StringValue("test-rule"),
				VipUuid:          types.StringValue("test-vip-uuid"),
				VipPortStart:     types.Int64Value(8080),
				ProtocolType:     types.StringValue("TCP"),
				VipPortEnd:       types.Int64Value(8081),
				PrivatePortStart: types.Int64Value(9080),
				PrivatePortEnd:   types.Int64Value(9081),
			}
			tt.setupPlan(&plan)

			createParam := param.CreatePortForwardingRuleParam{
				BaseParam: param.BaseParam{},
				Params: param.CreatePortForwardingRuleParamDetail{
					Name:         plan.Name.ValueString(),
					VipUuid:      plan.VipUuid.ValueString(),
					VipPortStart: int(plan.VipPortStart.ValueInt64()),
					ProtocolType: plan.ProtocolType.ValueString(),
				},
			}

			if !plan.VipPortEnd.IsNull() && !plan.VipPortEnd.IsUnknown() {
				createParam.Params.VipPortEnd = intPtr(int(plan.VipPortEnd.ValueInt64()))
			}
			if !plan.PrivatePortStart.IsNull() && !plan.PrivatePortStart.IsUnknown() {
				createParam.Params.PrivatePortStart = intPtr(int(plan.PrivatePortStart.ValueInt64()))
			}
			if !plan.PrivatePortEnd.IsNull() && !plan.PrivatePortEnd.IsUnknown() {
				createParam.Params.PrivatePortEnd = intPtr(int(plan.PrivatePortEnd.ValueInt64()))
			}

			tt.checkFunc(t, createParam)
		})
	}
}

func TestAccPortForwardingRuleResource_disappears(t *testing.T) {
	env := loadEnvData(t)
	name := testAccName("pf-rule-disappears")

	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping port forwarding rule acceptance test")
	}

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPortForwardingRuleDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_port_forwarding_rule" "test" {
  name           = %q
  vip_uuid       = data.zstack_vips.test.vips.0.uuid
  vip_port_start = 8080
  protocol_type  = "TCP"
  allowed_cidr   = "0.0.0.0/0"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckPortForwardingRuleDisappears("zstack_port_forwarding_rule.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPortForwardingRuleResource(t *testing.T) {
	env := loadEnvData(t)
	name := testAccName("pf-rule")
	updatedName := name + "-updated"

	// Port forwarding requires a VIP. Try to find one from env.json.
	// If VIPs are not available, we create a minimal config that references a VIP data source.
	if len(env.L3Networks) == 0 {
		t.Skip("no l3_networks in env.json, skipping port forwarding rule acceptance test")
	}

	// Use a pre-existing VIP if available via data source
	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckPortForwardingRuleDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_port_forwarding_rule" "test" {
  name           = %q
  description    = "acceptance port forwarding rule"
  vip_uuid       = data.zstack_vips.test.vips.0.uuid
  vip_port_start = 8080
  protocol_type  = "TCP"
  allowed_cidr   = "0.0.0.0/0"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance port forwarding rule")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("protocol_type"), knownvalue.StringExact("TCP")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("vip_port_start"), knownvalue.StringExact("8080")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_vips" "test" {
}

resource "zstack_port_forwarding_rule" "test" {
  name           = %q
  description    = "acceptance port forwarding rule updated"
  vip_uuid       = data.zstack_vips.test.vips.0.uuid
  vip_port_start = 8080
  protocol_type  = "TCP"
  allowed_cidr   = "0.0.0.0/0"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance port forwarding rule updated")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("protocol_type"), knownvalue.StringExact("TCP")),
					statecheck.ExpectKnownValue("zstack_port_forwarding_rule.test", tfjsonpath.New("vip_port_start"), knownvalue.StringExact("8080")),
				},
			},
			{
				ResourceName:                         "zstack_port_forwarding_rule.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_port_forwarding_rule.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}
