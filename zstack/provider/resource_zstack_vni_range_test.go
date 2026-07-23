// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

func TestVniRangeResourceSchema(t *testing.T) {
	var r vniRangeResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)

	for _, name := range []string{"name", "start_vni", "end_vni", "pool_uuid"} {
		attribute, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Fatalf("schema missing required attribute %q", name)
		}
		if !attribute.IsRequired() {
			t.Errorf("attribute %q should be required", name)
		}
	}

	uuid, ok := resp.Schema.Attributes["uuid"]
	if !ok || !uuid.IsComputed() {
		t.Fatal("uuid should be a computed attribute")
	}
}

func TestVniRangeResourceMetadata(t *testing.T) {
	var r vniRangeResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vni_range" {
		t.Fatalf("unexpected type name: %s", resp.TypeName)
	}
}

func TestVniRangeModelFromViewPreservesCreateOnlyFields(t *testing.T) {
	prior := vniRangeResourceModel{
		ResourceUuid: types.StringValue("custom-uuid"),
		TagUuids:     stringSliceToList([]string{"tag-1"}),
		SystemTags:   stringSliceToList([]string{"system-tag"}),
	}

	state := vniRangeModelFromView(&view.VniRangeInventoryView{
		BaseInfoView:  view.BaseInfoView{UUID: "range-uuid", Name: "range"},
		Description:   "test range",
		StartVni:      100,
		EndVni:        200,
		L2NetworkUuid: "pool-uuid",
	}, prior)

	if state.PoolUuid.ValueString() != "pool-uuid" || state.StartVni.ValueInt64() != 100 || state.EndVni.ValueInt64() != 200 {
		t.Fatalf("inventory fields were not mapped: %#v", state)
	}
	if state.ResourceUuid.ValueString() != "custom-uuid" {
		t.Fatalf("resource_uuid was not preserved: %s", state.ResourceUuid.ValueString())
	}
	if got := listToStringSlice(state.TagUuids); len(got) != 1 || got[0] != "tag-1" {
		t.Fatalf("tag_uuids were not preserved: %#v", got)
	}
	if got := listToStringSlice(state.SystemTags); len(got) != 1 || got[0] != "system-tag" {
		t.Fatalf("system_tags were not preserved: %#v", got)
	}
}

func TestVniRangeCreateParam(t *testing.T) {
	createParam := vniRangeCreateParam(vniRangeResourceModel{
		Name:       types.StringValue("range"),
		StartVni:   types.Int64Value(100),
		EndVni:     types.Int64Value(200),
		TagUuids:   stringSliceToList([]string{"tag-1"}),
		SystemTags: stringSliceToList([]string{"system-tag"}),
	})

	if createParam.Params.StartVni != 100 || createParam.Params.EndVni != 200 {
		t.Fatalf("unexpected VNI bounds: %d-%d", createParam.Params.StartVni, createParam.Params.EndVni)
	}
	if len(createParam.Params.TagUuids) != 1 || createParam.Params.TagUuids[0] != "tag-1" {
		t.Fatalf("tag UUIDs were not passed: %#v", createParam.Params.TagUuids)
	}
	if len(createParam.SystemTags) != 1 || createParam.SystemTags[0] != "system-tag" {
		t.Fatalf("system tags were not passed: %#v", createParam.SystemTags)
	}
}

func TestVniRangeResourceRejectsInvalidBounds(t *testing.T) {
	testCases := map[string]struct {
		startVni int
		endVni   int
		error    *regexp.Regexp
	}{
		"start below minimum": {
			startVni: 0,
			endVni:   100,
			error:    regexp.MustCompile(`(?s)start_vni.*between 1 and 16777215`),
		},
		"end above maximum": {
			startVni: 100,
			endVni:   maxVni + 1,
			error:    regexp.MustCompile(`(?s)end_vni.*between 1 and 16777215`),
		},
		"start after end": {
			startVni: 200,
			endVni:   100,
			error:    regexp.MustCompile(`(?s)start_vni must be less than or equal to end_vni`),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []tfresource.TestStep{
					{
						Config: providerConfig() + fmt.Sprintf(`
resource "zstack_vni_range" "test" {
  name       = "invalid-vni-range"
  start_vni  = %d
  end_vni    = %d
  pool_uuid  = "pool-uuid"
}
`, testCase.startVni, testCase.endVni),
						PlanOnly:    true,
						ExpectError: testCase.error,
					},
				},
			})
		})
	}
}
