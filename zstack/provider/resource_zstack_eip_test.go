// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestEIPResource_Schema(t *testing.T) {
	var r eipResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "vip_uuid", "vm_nic_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "description"}
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

func TestEIPResource_Metadata(t *testing.T) {
	var r eipResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_eip" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccEIPResource(t *testing.T) {
	env := loadEnvData(t)

	vipUUID, vmNicUUID, skipReason := selectEIPFixture(env)
	if skipReason != "" {
		t.Skip(skipReason)
	}

	name := testAccName("eip")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEipDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_eip" "test" {
  name        = %q
  description = "acceptance EIP"
  vip_uuid    = %q
  vm_nic_uuid = %q
}
`, name, vipUUID, vmNicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance EIP")),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vip_uuid"), knownvalue.StringExact(vipUUID)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vm_nic_uuid"), knownvalue.StringExact(vmNicUUID)),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_eip" "test" {
  name        = %q
  description = "acceptance EIP updated"
  vip_uuid    = %q
  vm_nic_uuid = %q
}
`, updatedName, vipUUID, vmNicUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance EIP updated")),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vip_uuid"), knownvalue.StringExact(vipUUID)),
					statecheck.ExpectKnownValue("zstack_eip.test", tfjsonpath.New("vm_nic_uuid"), knownvalue.StringExact(vmNicUUID)),
				},
			},
			{
				ResourceName:                         "zstack_eip.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_eip.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func selectEIPFixture(env *EnvData) (string, string, string) {
	for _, vip := range env.Vips {
		if !isAvailableEIPVip(env, vip) {
			continue
		}
		vipL3UUID := envStr(vip, "l3_network_uuid")
		if vmNicUUID := selectEIPVmNic(env, vipL3UUID); vmNicUUID != "" {
			return envStr(vip, "uuid"), vmNicUUID, ""
		}
	}

	return "", "", "no available EIP fixture in env.json: need an unused VIP on an L3 with Eip service and an unbound VM NIC on an Eip-enabled compatible non-Public L3 (Public VIP -> non-Public NIC; non-Public VIP -> different non-Public NIC)"
}

func isAvailableEIPVip(env *EnvData, vip map[string]interface{}) bool {
	l3UUID := envStr(vip, "l3_network_uuid")
	return envStr(vip, "use_for") == "" && l3HasNetworkService(env, l3UUID, "Eip")
}

func selectEIPVmNic(env *EnvData, vipL3UUID string) string {
	var fallback string
	for _, nic := range env.VmNics {
		nicUUID := envStr(nic, "uuid")
		nicL3UUID := envStr(nic, "l3_network_uuid")
		if !isCompatibleEIPVmNic(env, vipL3UUID, nicL3UUID, nicUUID) {
			continue
		}
		if l3Category(env, nicL3UUID) == "Private" {
			return nicUUID
		}
		if fallback == "" {
			fallback = nicUUID
		}
	}
	return fallback
}

func isCompatibleEIPVmNic(env *EnvData, vipL3UUID, nicL3UUID, nicUUID string) bool {
	if nicUUID == "" || nicL3UUID == "" {
		return false
	}
	if !l3HasNetworkService(env, nicL3UUID, "Eip") || vmNicHasEip(env, nicUUID) {
		return false
	}
	if l3Category(env, nicL3UUID) == "Public" {
		return false
	}
	if l3Category(env, vipL3UUID) != "Public" && nicL3UUID == vipL3UUID {
		return false
	}
	return true
}

func l3Category(env *EnvData, l3UUID string) string {
	for _, l3 := range env.L3Networks {
		if envStr(l3, "uuid") == l3UUID {
			return envStr(l3, "category")
		}
	}
	return ""
}

func l3HasNetworkService(env *EnvData, l3UUID, serviceType string) bool {
	if l3UUID == "" {
		return false
	}
	for _, l3 := range env.L3Networks {
		if envStr(l3, "uuid") != l3UUID {
			continue
		}
		services, ok := l3["network_services"].([]interface{})
		if !ok {
			return false
		}
		for _, service := range services {
			serviceMap, ok := service.(map[string]interface{})
			if ok && envStr(serviceMap, "network_service_type") == serviceType {
				return true
			}
		}
	}
	return false
}

func vmNicHasEip(env *EnvData, vmNicUUID string) bool {
	if vmNicUUID == "" {
		return false
	}
	for _, eip := range env.Eips {
		if envStr(eip, "vm_nic_uuid") == vmNicUUID {
			return true
		}
	}
	return false
}
