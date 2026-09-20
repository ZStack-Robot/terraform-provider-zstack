// Copyright (c) ZStack.io, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSlbOfferingsQueries(t *testing.T) {
	for _, tc := range []struct {
		name, config, query, count string
		empty                      bool
	}{
		{name: "empty", count: "0", empty: true},
		{name: "filtered_empty", config: `filter {
 name = "state"
 values = ["Disabled"]
}`, count: "0"},
		{name: "memory_filter_mib", config: `filter {
 name = "memory_size"
 values = ["8192"]
}`, count: "1"},
		{name: "uuid", config: `uuid = "slb-uuid"`, query: "uuid=slb-uuid", count: "1"},
		{name: "name", config: `name = "slb-test"`, query: "name=slb-test", count: "1"},
		{name: "name_pattern", config: `name_pattern = "slb-%"`, query: "name~=slb-%", count: "1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"inventory": map[string]string{"uuid": "session"}})
			})
			mux.HandleFunc("/zstack/v1/instance-offerings/slb", func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", r.Method)
				}
				if got := r.URL.Query().Get("q"); got != tc.query {
					t.Errorf("query = %q, want %q", got, tc.query)
				}
				items := []any{}
				if !tc.empty {
					items = append(items, map[string]any{"uuid": "slb-uuid", "name": "slb-test", "memorySize": int64(8589934592), "cpuNum": 8, "state": "Enabled", "type": "SLB"})
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{"inventories": items})
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
data "zstack_slb_offerings" "test" {
 %s
}
output "count" {
 value = length(data.zstack_slb_offerings.test.slb_offers)
}
output "names" {
 value = [for offer in data.zstack_slb_offerings.test.slb_offers : offer.name]
}
`, host, port, tc.config)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{{Config: cfg, Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_slb_offerings.test", "slb_offers.#", tc.count),
					resource.TestCheckOutput("count", tc.count),
				)}},
			})
		})
	}
}
