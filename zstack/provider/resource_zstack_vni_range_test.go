// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/view"
)

type vniRangePageQueryFunc func(context.Context, *param.QueryParam) ([]view.VniRangeInventoryView, int, error)

func (f vniRangePageQueryFunc) PageVniRanges(ctx context.Context, query *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
	return f(ctx, query)
}

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
			error:    regexp.MustCompile(`(?s)start_vni.*between 1 and 16777214`),
		},
		"end above maximum": {
			startVni: 100,
			endVni:   maxVni + 1,
			error:    regexp.MustCompile(`(?s)end_vni.*between 1 and 16777214`),
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

func TestVniRangeOverlapDetection(t *testing.T) {
	existing := view.VniRangeInventoryView{
		BaseInfoView:  view.BaseInfoView{UUID: "existing-uuid", Name: "existing"},
		StartVni:      100,
		EndVni:        200,
		L2NetworkUuid: "pool-uuid",
	}

	testCases := map[string]struct {
		startVni    int64
		endVni      int64
		currentUuid string
		wantOverlap bool
	}{
		"start inside existing": {
			startVni:    150,
			endVni:      250,
			wantOverlap: true,
		},
		"end inside existing": {
			startVni:    50,
			endVni:      150,
			wantOverlap: true,
		},
		"contains existing": {
			startVni:    50,
			endVni:      250,
			wantOverlap: true,
		},
		"contained by existing": {
			startVni:    120,
			endVni:      180,
			wantOverlap: true,
		},
		"equal start endpoint": {
			startVni:    100,
			endVni:      100,
			wantOverlap: true,
		},
		"equal end endpoint": {
			startVni:    200,
			endVni:      200,
			wantOverlap: true,
		},
		"adjacent before": {
			startVni:    50,
			endVni:      99,
			wantOverlap: false,
		},
		"adjacent after": {
			startVni:    201,
			endVni:      250,
			wantOverlap: false,
		},
		"update excludes itself": {
			startVni:    100,
			endVni:      200,
			currentUuid: "existing-uuid",
			wantOverlap: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			resource := &vniRangeResource{
				vniRangePageQuery: vniRangePageQueryFunc(func(_ context.Context, query *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
					if got := query.Get("q"); got != "l2NetworkUuid=pool-uuid" {
						t.Fatalf("unexpected pool query: %q", got)
					}
					return []view.VniRangeInventoryView{existing}, 1, nil
				}),
			}

			conflict, err := resource.findOverlappingVniRange(
				context.Background(),
				"pool-uuid",
				testCase.startVni,
				testCase.endVni,
				testCase.currentUuid,
			)
			if err != nil {
				t.Fatalf("find overlapping VNI range: %v", err)
			}
			if got := conflict != nil; got != testCase.wantOverlap {
				t.Fatalf("unexpected overlap result: got %t, want %t", got, testCase.wantOverlap)
			}
		})
	}
}

func TestVniRangeQueryAllPagination(t *testing.T) {
	call := 0
	resource := &vniRangeResource{
		vniRangePageQuery: vniRangePageQueryFunc(func(_ context.Context, query *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
			wantStart := fmt.Sprintf("%d", call)
			if got := query.Get("start"); got != wantStart {
				t.Fatalf("unexpected page start: got %q, want %q", got, wantStart)
			}
			if got := query.Get("limit"); got != fmt.Sprintf("%d", vniRangePageSize) {
				t.Fatalf("unexpected page limit: %q", got)
			}
			call++
			return []view.VniRangeInventoryView{{
				BaseInfoView:  view.BaseInfoView{UUID: fmt.Sprintf("range-%d", call)},
				StartVni:      call * 100,
				EndVni:        call*100 + 99,
				L2NetworkUuid: "pool-uuid",
			}}, 2, nil
		}),
	}

	ranges, err := resource.queryAllVniRanges(context.Background(), "pool-uuid")
	if err != nil {
		t.Fatalf("query all VNI ranges: %v", err)
	}
	if len(ranges) != 2 || call != 2 {
		t.Fatalf("unexpected pagination result: ranges=%d calls=%d", len(ranges), call)
	}
}

func TestVniRangeQueryAllFailures(t *testing.T) {
	t.Run("API timeout", func(t *testing.T) {
		resource := &vniRangeResource{
			vniRangePageQuery: vniRangePageQueryFunc(func(context.Context, *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
				return nil, 0, context.DeadlineExceeded
			}),
		}

		_, err := resource.queryAllVniRanges(context.Background(), "pool-uuid")
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected deadline exceeded, got %v", err)
		}
	})

	t.Run("incomplete pagination", func(t *testing.T) {
		call := 0
		resource := &vniRangeResource{
			vniRangePageQuery: vniRangePageQueryFunc(func(context.Context, *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
				call++
				if call == 1 {
					return []view.VniRangeInventoryView{{
						BaseInfoView:  view.BaseInfoView{UUID: "range-1"},
						StartVni:      100,
						EndVni:        199,
						L2NetworkUuid: "pool-uuid",
					}}, 2, nil
				}
				return nil, 2, nil
			}),
		}

		_, err := resource.queryAllVniRanges(context.Background(), "pool-uuid")
		if err == nil || !strings.Contains(err.Error(), "incomplete VNI range query") {
			t.Fatalf("expected incomplete pagination error, got %v", err)
		}
	})

	t.Run("cancellation", func(t *testing.T) {
		requestStarted := make(chan struct{})
		resource := &vniRangeResource{
			vniRangePageQuery: vniRangePageQueryFunc(func(ctx context.Context, _ *param.QueryParam) ([]view.VniRangeInventoryView, int, error) {
				close(requestStarted)
				<-ctx.Done()
				return nil, 0, ctx.Err()
			}),
		}

		ctx, cancel := context.WithCancel(context.Background())
		errs := make(chan error, 1)
		go func() {
			_, err := resource.queryAllVniRanges(ctx, "pool-uuid")
			errs <- err
		}()
		<-requestStarted
		cancel()

		if err := <-errs; !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context cancellation, got %v", err)
		}
	})
}

func TestVniRangeResourcePlanRejectsOverlap(t *testing.T) {
	server := newVniRangePlanServer(t, []view.VniRangeInventoryView{{
		BaseInfoView:  view.BaseInfoView{UUID: "existing-uuid", Name: "existing"},
		StartVni:      100,
		EndVni:        200,
		L2NetworkUuid: "pool-uuid",
	}})
	defer server.Close()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config:      vniRangePlanTestConfig(t, server, 150, 250),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)Overlapping VNI range.*existing-uuid.*100, 200`),
			},
		},
	})
}

func TestVniRangeResourcePlanAllowsAdjacentRange(t *testing.T) {
	server := newVniRangePlanServer(t, []view.VniRangeInventoryView{{
		BaseInfoView:  view.BaseInfoView{UUID: "existing-uuid", Name: "existing"},
		StartVni:      100,
		EndVni:        200,
		L2NetworkUuid: "pool-uuid",
	}})
	defer server.Close()

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config:             vniRangePlanTestConfig(t, server, 201, 250),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func newVniRangePlanServer(t *testing.T, ranges []view.VniRangeInventoryView) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/zstack/v1/l2-networks/vxlan-pool/vni-range", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		if got := req.URL.Query().Get("q"); got != "l2NetworkUuid=pool-uuid" {
			t.Fatalf("unexpected pool query: %q", got)
		}
		if got := req.URL.Query().Get("replyWithCount"); got != "true" {
			t.Fatalf("replyWithCount was not enabled: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"inventories": ranges,
			"total":       len(ranges),
		}); err != nil {
			t.Fatalf("write mock response: %v", err)
		}
	})
	return httptest.NewServer(mux)
}

func vniRangePlanTestConfig(t *testing.T, server *httptest.Server, startVni, endVni int) string {
	t.Helper()
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatalf("split mock server address: %v", err)
	}
	return fmt.Sprintf(`
provider "zstack" {
  host              = %q
  port              = %s
  access_key_id     = "test-access-key"
  access_key_secret = "test-access-key-secret"
}

resource "zstack_vni_range" "test" {
  name       = "plan-test"
  start_vni  = %d
  end_vni    = %d
  pool_uuid  = "pool-uuid"
}
`, host, port, startVni, endVni)
}
