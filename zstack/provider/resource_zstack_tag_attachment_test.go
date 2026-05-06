// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

func TestTagAttachmentResource_Schema(t *testing.T) {
	var r tagAttachmentResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"tag_uuid", "resource_uuids"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	if _, ok := resp.Schema.Attributes["tokens"]; !ok {
		t.Fatal("schema missing optional attribute tokens")
	}
}

func TestTagAttachmentResource_Metadata(t *testing.T) {
	var r tagAttachmentResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_tag_attachment" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccTagAttachmentResource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.Zones) == 0 {
		t.Skip("no zones in env data, required for tag attachment target")
	}
	zoneUUID := envStr(env.Zones[0], "uuid")

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckTagAttachmentDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "zstack_tag" "test_for_attach" {
  name  = "acc-test-tag-for-attach"
  value = "attach-value"
  type  = "simple"
  color = "#FF0000"
}

resource "zstack_tag_attachment" "test" {
  tag_uuid       = zstack_tag.test_for_attach.uuid
  resource_uuids = [%q]
}
`, zoneUUID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_tag_attachment.test", tfjsonpath.New("tag_uuid"), knownvalue.NotNull()),
				},
			},
		},
	})
}

// TestTagAttachmentDeletePassesResourceUuids verifies that Delete passes
// state.ResourceUuids to the detach operation. This is a unit test using a spy client.
//
// SDK SIGNATURE INVESTIGATION (via go doc):
// Current SDK v0.0.4 has:
//
//	func (cli *ZSClient) DetachTagFromResources(uuid string, deleteMode param.DeleteMode) error
//
// This signature only accepts 2 args (uuid, deleteMode), NOT resourceUuids.
//
// However, there exists param.DetachTagFromResourcesParamDetail with:
//
//	type DetachTagFromResourcesParamDetail struct {
//	    ResourceUuids []string `json:"resourceUuids" validate:"required"`
//	}
//
// The current implementation extracts resourceUuids from state but never uses them,
// causing DetachTagFromResources to detach the tag from ALL resources instead of
// just the ones specified in the attachment state.
//
// This test uses a spy to verify that resourceUuids SHOULD be passed to the detach operation.
// State values come from refresh, so Unknown is impossible here; no IsUnknown guard needed.
func TestTagAttachmentDeletePassesResourceUuids(t *testing.T) {
	ctx := context.Background()

	spy := &spyTagClient{
		detachCalls: []detachCall{},
	}

	r := &tagAttachmentResource{
		client: spy,
	}

	req := resource.DeleteRequest{
		State: tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"tag_uuid":       tftypes.String,
					"resource_uuids": tftypes.List{ElementType: tftypes.String},
					"tokens":         tftypes.Map{ElementType: tftypes.String},
				},
			}, map[string]tftypes.Value{
				"tag_uuid": tftypes.NewValue(tftypes.String, "tag-uuid-123"),
				"resource_uuids": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
					tftypes.NewValue(tftypes.String, "resource-uuid-a"),
					tftypes.NewValue(tftypes.String, "resource-uuid-b"),
				}),
				"tokens": tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
			}),
			Schema: schema.Schema{
				Attributes: map[string]schema.Attribute{
					"tag_uuid": schema.StringAttribute{
						Required: true,
					},
					"resource_uuids": schema.ListAttribute{
						Required:    true,
						ElementType: types.StringType,
					},
					"tokens": schema.MapAttribute{
						Optional:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	}
	resp := &resource.DeleteResponse{}

	r.Delete(ctx, req, resp)

	// Verify no diagnostics errors
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}

	// Verify DetachTagFromResources was called exactly once
	if len(spy.detachCalls) != 1 {
		t.Fatalf("expected 1 detach call, got %d", len(spy.detachCalls))
	}

	// Verify the call received the correct arguments
	call := spy.detachCalls[0]
	if call.tagUuid != "tag-uuid-123" {
		t.Errorf("expected tagUuid=tag-uuid-123, got %s", call.tagUuid)
	}

	// CRITICAL ASSERTION: Verify resourceUuids were passed
	expectedUuids := []string{"resource-uuid-a", "resource-uuid-b"}
	if len(call.resourceUuids) != len(expectedUuids) {
		t.Fatalf("expected %d resourceUuids, got %d", len(expectedUuids), len(call.resourceUuids))
	}
	for i, expected := range expectedUuids {
		if call.resourceUuids[i] != expected {
			t.Errorf("resourceUuids[%d]: expected %s, got %s", i, expected, call.resourceUuids[i])
		}
	}
}

func TestTagAttachmentReadQueriesByTagPatternUuid(t *testing.T) {
	ctx := context.Background()
	spy := &spyTagClient{
		queryResults: []view.UserTagInventoryView{
			{
				BaseInfoView:   view.BaseInfoView{UUID: "user-tag-uuid-1"},
				TagPatternUuid: "tag-pattern-uuid-123",
				ResourceUuid:   "resource-uuid-a",
			},
		},
	}

	r := &tagAttachmentResource{client: spy}
	req := resource.ReadRequest{
		State: tagAttachmentTestState("tag-pattern-uuid-123", []string{"resource-uuid-a"}),
	}
	resp := &resource.ReadResponse{State: req.State}

	r.Read(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected errors: %v", resp.Diagnostics.Errors())
	}
	if len(spy.queryCalls) != 1 {
		t.Fatalf("expected 1 QueryUserTag call, got %d", len(spy.queryCalls))
	}
	qs := spy.queryCalls[0].Values["q"]
	if len(qs) == 0 || qs[0] != "tagPatternUuid=tag-pattern-uuid-123" {
		t.Fatalf("expected tagPatternUuid query, got %v", qs)
	}
}

// detachCall records arguments to DetachTagFromResources
type detachCall struct {
	tagUuid       string
	deleteMode    param.DeleteMode
	resourceUuids []string
}

type spyTagClient struct {
	detachCalls  []detachCall
	queryCalls   []param.QueryParam
	queryResults []view.UserTagInventoryView
}

func (s *spyTagClient) AttachTagToResources(tagUuid string, params param.AttachTagToResourcesParam) (*view.AttachTagToResourcesEventView, error) {
	return nil, nil
}

func (s *spyTagClient) QueryUserTag(params *param.QueryParam) ([]view.UserTagInventoryView, error) {
	if params != nil {
		s.queryCalls = append(s.queryCalls, *params)
	}
	return s.queryResults, nil
}

func (s *spyTagClient) DetachTagFromResources(tagUuid string, deleteMode param.DeleteMode, resourceUuids []string) error {
	s.detachCalls = append(s.detachCalls, detachCall{
		tagUuid:       tagUuid,
		deleteMode:    deleteMode,
		resourceUuids: resourceUuids,
	})
	return nil
}

func tagAttachmentTestState(tagUUID string, resourceUUIDs []string) tfsdk.State {
	values := make([]tftypes.Value, 0, len(resourceUUIDs))
	for _, uuid := range resourceUUIDs {
		values = append(values, tftypes.NewValue(tftypes.String, uuid))
	}

	return tfsdk.State{
		Raw: tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"tag_uuid":       tftypes.String,
				"resource_uuids": tftypes.List{ElementType: tftypes.String},
				"tokens":         tftypes.Map{ElementType: tftypes.String},
			},
		}, map[string]tftypes.Value{
			"tag_uuid":       tftypes.NewValue(tftypes.String, tagUUID),
			"resource_uuids": tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, values),
			"tokens":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
		}),
		Schema: schema.Schema{
			Attributes: map[string]schema.Attribute{
				"tag_uuid": schema.StringAttribute{
					Required: true,
				},
				"resource_uuids": schema.ListAttribute{
					Required:    true,
					ElementType: types.StringType,
				},
				"tokens": schema.MapAttribute{
					Optional:    true,
					ElementType: types.StringType,
				},
			},
		},
	}
}
