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
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_virtual_router_instance.foo", tfjsonpath.New("name"), knownvalue.StringExact("test-mock-vr")),
					statecheck.ExpectKnownValue("zstack_virtual_router_instance.foo", tfjsonpath.New("uuid"), knownvalue.StringExact("mock-vr-uuid")),
					statecheck.ExpectKnownValue("zstack_virtual_router_instance.foo", tfjsonpath.New("state"), knownvalue.StringExact("Running")),
					statecheck.ExpectKnownValue("zstack_virtual_router_instance.foo", tfjsonpath.New("description"), knownvalue.StringExact("mock vr created via test")),
					statecheck.ExpectKnownValue("zstack_virtual_router_instance.foo", tfjsonpath.New("virtual_router_offering_uuid"), knownvalue.StringExact("mock-offering-uuid")),
				},
			},
		},
	})
}
