// Copyright (c) ZStack.io, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackVirtualRouterInstance_Create(t *testing.T) {
	mux := http.NewServeMux()

	// Setup generic login route
	mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := map[string]interface{}{
			"inventory": map[string]string{
				"uuid": "mock-session-uuid-123",
			},
		}
		json.NewEncoder(w).Encode(res)
	})

	// Setup mock route for creating Virtual Router Instance
	mux.HandleFunc("/zstack/v1/vpc/virtual-routers", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := map[string]interface{}{
			"inventory": map[string]interface{}{
				"uuid":                 "mock-vr-uuid",
				"name":                 "test-mock-vr",
				"description":          "mock vr created via test",
				"instanceOfferingUuid": "mock-offering-uuid",
				"state":                "Running",
				"status":               "Connected",
			},
		}
		json.NewEncoder(w).Encode(res)
	})

	// Setup GET Virtual Router Instance (read path)
	mux.HandleFunc("/zstack/v1/vm-instances/appliances/virtual-routers/mock-vr-uuid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		res := map[string]interface{}{
			"inventories": []interface{}{
				map[string]interface{}{
					"uuid":                 "mock-vr-uuid",
					"name":                 "test-mock-vr",
					"description":          "mock vr created via test",
					"instanceOfferingUuid": "mock-offering-uuid",
					"state":                "Running",
					"status":               "Connected",
				},
			},
		}
		json.NewEncoder(w).Encode(res)
	})

	// Setup DELETE Virtual Router Instance
	mux.HandleFunc("/zstack/v1/vm-instances/mock-vr-uuid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		}
	})

	// Spin up generic mock server
	server := httptest.NewServer(mux)
	defer server.Close()

	// Override ZSTACK environment variables used by the provider initialization block
	hostParts := strings.Split(server.Listener.Addr().String(), ":")
	t.Setenv("ZSTACK_HOST", hostParts[0])
	t.Setenv("ZSTACK_PORT", hostParts[1])
	t.Setenv("ZSTACK_ACCOUNT_NAME", "admin")
	t.Setenv("ZSTACK_ACCOUNT_PASSWORD", "password")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "zstack" {
    host              = "%s"
    port              = %s
    account_name      = "admin" 
    account_password  = "password" 
}

resource "zstack_virtual_router_instance" "foo" {
  name                        = "test-mock-vr"
  description                 = "mock vr created via test"
  virtual_router_offering_uuid = "mock-offering-uuid"
}
`, hostParts[0], hostParts[1]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zstack_virtual_router_instance.foo", "name", "test-mock-vr"),
					resource.TestCheckResourceAttr("zstack_virtual_router_instance.foo", "uuid", "mock-vr-uuid"),
					resource.TestCheckResourceAttr("zstack_virtual_router_instance.foo", "state", "Running"),
					resource.TestCheckResourceAttr("zstack_virtual_router_instance.foo", "description", "mock vr created via test"),
					resource.TestCheckResourceAttr("zstack_virtual_router_instance.foo", "virtual_router_offering_uuid", "mock-offering-uuid"),
				),
			},
		},
	})
}
