// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type filterTestBase struct {
	UUID string `json:"uuid"`
}

type filterTestAffinityGroup struct {
	filterTestBase
	Name     string `json:"name,omitempty"`
	ZoneUuid string `json:"zoneUuid,omitempty"`
}

type filterTestInstance struct {
	CPUNum int64
}

type filterTestVolume struct {
	PrimaryStorageUUID string `json:"primaryStorageUUID,omitempty"`
}

type filterTestTerraformModel struct {
	ZoneUuid types.String `tfsdk:"zone_uuid"`
}

func TestFilterResourceResolvesMappedCamelCaseFieldByJSONTag(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
		{filterTestBase: filterTestBase{UUID: "ag-2"}, Name: "ag-b", ZoneUuid: "zone-b"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"zone_uuid": {"zone-a"},
	}, "affinity_group")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].UUID != "ag-1" {
		t.Fatalf("expected ag-1, got %s", filtered[0].UUID)
	}
}

func TestFilterResourceResolvesEmbeddedFieldByJSONTag(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
		{filterTestBase: filterTestBase{UUID: "ag-2"}, Name: "ag-b", ZoneUuid: "zone-b"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"uuid": {"ag-2"},
	}, "affinity_group")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].UUID != "ag-2" {
		t.Fatalf("expected ag-2, got %s", filtered[0].UUID)
	}
}

func TestFilterResourceResolvesExactGoFieldMapping(t *testing.T) {
	resources := []filterTestInstance{
		{CPUNum: 2},
		{CPUNum: 4},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"cpu_num": {"4"},
	}, "instance")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].CPUNum != 4 {
		t.Fatalf("expected cpu_num 4, got %d", filtered[0].CPUNum)
	}
}

func TestFilterResourceResolvesAcronymFieldByJSONTag(t *testing.T) {
	resources := []filterTestVolume{
		{PrimaryStorageUUID: "ps-1"},
		{PrimaryStorageUUID: "ps-2"},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"primary_storage_uuid": {"ps-2"},
	}, "volume")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].PrimaryStorageUUID != "ps-2" {
		t.Fatalf("expected ps-2, got %s", filtered[0].PrimaryStorageUUID)
	}
}

func TestFilterResourceResolvesTerraformSDKTag(t *testing.T) {
	resources := []filterTestTerraformModel{
		{ZoneUuid: types.StringValue("zone-a")},
		{ZoneUuid: types.StringValue("zone-b")},
	}

	filtered, diags := FilterResource(context.Background(), resources, map[string][]string{
		"zone_uuid": {"zone-b"},
	}, "unknown")
	if diags.HasError() {
		t.Fatalf("FilterResource returned diagnostics: %v", diags)
	}

	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered resource, got %d", len(filtered))
	}
	if filtered[0].ZoneUuid.ValueString() != "zone-b" {
		t.Fatalf("expected zone-b, got %s", filtered[0].ZoneUuid.ValueString())
	}
}

func TestFilterResourceInvalidFilterKey(t *testing.T) {
	resources := []filterTestAffinityGroup{
		{filterTestBase: filterTestBase{UUID: "ag-1"}, Name: "ag-a", ZoneUuid: "zone-a"},
	}

	_, diags := FilterResource(context.Background(), resources, map[string][]string{
		"missing": {"value"},
	}, "affinity_group")

	if !diags.HasError() {
		t.Fatal("expected invalid filter key diagnostic")
	}
}
