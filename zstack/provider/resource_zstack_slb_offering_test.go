// Copyright (c) ZStack.io, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSlbOfferingLifecycle(t *testing.T) {
	var deleted atomic.Bool
	var unavailable atomic.Bool
	var creates atomic.Int32
	inv := map[string]any{"uuid": "slb-uuid", "name": "wuying-slb-8c8g", "description": "", "cpuNum": 8, "memorySize": int64(8589934592), "type": "SLB", "zoneUuid": "zone-uuid", "managementNetworkUuid": "management-uuid", "imageUuid": "shared-image-uuid", "state": "Enabled"}
	mux := http.NewServeMux()
	mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"inventory": map[string]string{"uuid": "session"}})
	})
	mux.HandleFunc("/zstack/v1/instance-offerings/slb", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{"inventories": []any{inv}})
			return
		}
		var payload struct {
			Params map[string]any `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		for k, v := range map[string]any{"cpuNum": float64(8), "memorySize": float64(8589934592), "zoneUuid": "zone-uuid", "managementNetworkUuid": "management-uuid", "imageUuid": "shared-image-uuid"} {
			if payload.Params[k] != v {
				t.Errorf("%s got %v want %v", k, payload.Params[k], v)
			}
		}
		inv["name"] = payload.Params["name"]
		deleted.Store(false)
		creates.Add(1)
		json.NewEncoder(w).Encode(map[string]any{"inventory": inv})
	})
	mux.HandleFunc("/zstack/v1/instance-offerings/slb/slb-uuid", func(w http.ResponseWriter, r *http.Request) {
		if unavailable.Load() {
			http.Error(w, "SLB read unavailable", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		items := []any{inv}
		if deleted.Load() {
			items = []any{}
		}
		json.NewEncoder(w).Encode(map[string]any{"inventories": items})
	})
	mux.HandleFunc("/zstack/v1/instance-offerings/slb-uuid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected %s", r.Method)
		}
		deleted.Store(true)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "{}")
	})
	server := httptest.NewServer(mux)
	defer server.Close()
	host, port, _ := net.SplitHostPort(server.Listener.Addr().String())
	cfg := fmt.Sprintf(`provider "zstack" {
 host = %q
 port = %s
 account_name = "admin"
 account_password = "test"
 }
 resource "zstack_slb_offering" "test" {
 name = "wuying-slb-8c8g"
 cpu_num = 8
 memory_size = 8192
 zone_uuid = "zone-uuid"
 management_network_uuid = "management-uuid"
 image_uuid = "shared-image-uuid"
 }
 data "zstack_slb_offerings" "selected" { uuid = zstack_slb_offering.test.uuid }
 `, host, port)
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: cfg, Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("zstack_slb_offering.test", "memory_size", "8192"), resource.TestCheckResourceAttr("zstack_slb_offering.test", "image_uuid", "shared-image-uuid"), resource.TestCheckResourceAttr("data.zstack_slb_offerings.selected", "slb_offers.#", "1"))},
		{ResourceName: "zstack_slb_offering.test", ImportState: true, ImportStateId: "slb-uuid", ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "uuid"},
		{Config: cfg, PlanOnly: true, PreConfig: func() { unavailable.Store(true) }, ExpectError: regexp.MustCompile("Error reading SLB Offering")},
		{Config: cfg, PlanOnly: true, PreConfig: func() { unavailable.Store(false) }},
		{Config: cfg, PreConfig: func() { deleted.Store(true) }},
		{Config: strings.ReplaceAll(cfg, "wuying-slb-8c8g", "renamed-slb"), Check: resource.TestCheckResourceAttr("zstack_slb_offering.test", "name", "renamed-slb")},
	}})
	if creates.Load() != 3 {
		t.Errorf("created %d times", creates.Load())
	}
}
