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
)

func TestAffinityGroupResource_Schema(t *testing.T) {
	var r affinityGroupResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"name", "policy"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"uuid", "state"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	optional := []string{"description", "type", "zone_uuid"}
	for _, attr := range optional {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema missing optional attribute %q", attr)
		}
	}
}

func TestAffinityGroupResource_Metadata(t *testing.T) {
	var r affinityGroupResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_affinity_group" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestNormalizeAffinityGroupType(t *testing.T) {
	tests := []struct {
		name     string
		config   types.String
		apiValue string
		want     types.String
	}{
		{
			name:     "preserves configured host when API returns uppercase",
			config:   types.StringValue("host"),
			apiValue: "HOST",
			want:     types.StringValue("host"),
		},
		{
			name:     "normalizes API uppercase host for computed value",
			config:   types.StringNull(),
			apiValue: "HOST",
			want:     types.StringValue("host"),
		},
		{
			name:     "normalizes API uppercase host when config is unknown",
			config:   types.StringUnknown(),
			apiValue: "HOST",
			want:     types.StringValue("host"),
		},
		{
			name:     "preserves configured value when API value is empty",
			config:   types.StringValue("host"),
			apiValue: "",
			want:     types.StringValue("host"),
		},
		{
			name:     "uses null when config is unknown and API value is empty",
			config:   types.StringUnknown(),
			apiValue: "",
			want:     types.StringNull(),
		},
		{
			name:     "keeps unknown server value unchanged",
			config:   types.StringNull(),
			apiValue: "OTHER",
			want:     types.StringValue("OTHER"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeAffinityGroupType(tt.config, tt.apiValue)
			if !got.Equal(tt.want) {
				t.Fatalf("normalizeAffinityGroupType() = %s, want %s", got.String(), tt.want.String())
			}
		})
	}
}

func TestNormalizeAffinityGroupPolicy(t *testing.T) {
	tests := []struct {
		name     string
		config   types.String
		apiValue string
		want     types.String
	}{
		{
			name:     "preserves configured antiSoft when API returns uppercase",
			config:   types.StringValue("antiSoft"),
			apiValue: "ANTISOFT",
			want:     types.StringValue("antiSoft"),
		},
		{
			name:     "normalizes imported uppercase antiSoft",
			config:   types.StringNull(),
			apiValue: "ANTISOFT",
			want:     types.StringValue("antiSoft"),
		},
		{
			name:     "normalizes uppercase antiHard",
			config:   types.StringNull(),
			apiValue: "ANTIHARD",
			want:     types.StringValue("antiHard"),
		},
		{
			name:     "normalizes uppercase proSoft",
			config:   types.StringNull(),
			apiValue: "PROSOFT",
			want:     types.StringValue("proSoft"),
		},
		{
			name:     "normalizes uppercase proHard",
			config:   types.StringNull(),
			apiValue: "PROHARD",
			want:     types.StringValue("proHard"),
		},
		{
			name:     "uses null when config is unknown and API value is empty",
			config:   types.StringUnknown(),
			apiValue: "",
			want:     types.StringNull(),
		},
		{
			name:     "keeps unknown server value unchanged",
			config:   types.StringNull(),
			apiValue: "OTHER",
			want:     types.StringValue("OTHER"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeAffinityGroupPolicy(tt.config, tt.apiValue)
			if !got.Equal(tt.want) {
				t.Fatalf("normalizeAffinityGroupPolicy() = %s, want %s", got.String(), tt.want.String())
			}
		})
	}
}

func TestAccAffinityGroupResource(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("affinity-group")
	updatedName := name + "-updated"

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAffinityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_affinity_group" "test" {
  name        = %q
  description = "acceptance affinity group"
  policy      = "antiSoft"
  type        = "host"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance affinity group")),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("policy"), knownvalue.StringExact("antiSoft")),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("type"), knownvalue.StringExact("host")),
				},
			},
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_affinity_group" "test" {
  name        = %q
  description = "acceptance affinity group updated"
  policy      = "antiSoft"
  type        = "host"
}
`, updatedName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("name"), knownvalue.StringExact(updatedName)),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("description"), knownvalue.StringExact("acceptance affinity group updated")),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("policy"), knownvalue.StringExact("antiSoft")),
					statecheck.ExpectKnownValue("zstack_affinity_group.test", tfjsonpath.New("type"), knownvalue.StringExact("host")),
				},
			},
			{
				ResourceName:                         "zstack_affinity_group.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_affinity_group.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
		},
	})
}

func TestAccAffinityGroupResource_disappears(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("ag-disappears")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAffinityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_affinity_group" "test_disappears" {
  name   = %q
  policy = "antiSoft"
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_affinity_group.test_disappears", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					stateCheckAffinityGroupDisappears("zstack_affinity_group.test_disappears"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAffinityGroupsDataSourceFilterByZoneUuid(t *testing.T) {
	_ = loadEnvData(t)
	name := testAccName("ag-filter-zone")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckAffinityGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_zone" "enabled" {
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

resource "zstack_affinity_group" "test" {
  name        = %q
  description = "acceptance affinity group for zone filter"
  policy      = "antiSoft"
  type        = "host"
  zone_uuid   = data.zstack_zone.enabled.zones[0].uuid
}

data "zstack_affinity_groups" "created" {
  depends_on = [zstack_affinity_group.test]

  filter {
    name   = "zone_uuid"
    values = [data.zstack_zone.enabled.zones[0].uuid]
  }

  filter {
    name   = "uuid"
    values = [zstack_affinity_group.test.uuid]
  }
}
`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_affinity_groups.created", tfjsonpath.New("affinity_groups"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue("data.zstack_affinity_groups.created", tfjsonpath.New("affinity_groups").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_affinity_groups.created", tfjsonpath.New("affinity_groups").AtSliceIndex(0).AtMapKey("zone_uuid"), knownvalue.NotNull()),
				},
			},
		},
	})
}
