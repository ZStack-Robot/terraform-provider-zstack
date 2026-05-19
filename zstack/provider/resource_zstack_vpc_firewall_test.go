// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestVpcFirewallResource_Schema(t *testing.T) {
	var r vpcFirewallResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}
	// Check required attributes
	required := []string{"name", "vpc_uuid"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	vpcUuidAttr, ok := resp.Schema.Attributes["vpc_uuid"].(rschema.StringAttribute)
	if !ok {
		t.Fatal("vpc_uuid should be a string attribute")
	}
	if !strings.Contains(vpcUuidAttr.Description, "virtual router instance") {
		t.Fatalf("vpc_uuid description should explain virtual router instance UUID semantics, got %q", vpcUuidAttr.Description)
	}

	// Check computed attributes
	computed := []string{"uuid"}
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

func TestVpcFirewallResource_Metadata(t *testing.T) {
	var r vpcFirewallResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_vpc_firewall" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestAccVpcFirewallResource_CreateUsesVirtualRouterUuid(t *testing.T) {
	mux := http.NewServeMux()

	var mu sync.Mutex
	virtualRouterValidated := false
	var createdVpcUuid string

	mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"inventory": map[string]string{
				"uuid": "mock-session-uuid",
			},
		})
	})

	mux.HandleFunc("/zstack/v1/vm-instances/appliances/virtual-routers/mock-vrouter-uuid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		mu.Lock()
		virtualRouterValidated = true
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"inventories": []interface{}{
				map[string]interface{}{
					"uuid": "mock-vrouter-uuid",
					"name": "mock-vrouter",
				},
			},
		})
	})

	mux.HandleFunc("/zstack/v1/vpcfirewalls", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			var payload struct {
				Params map[string]interface{} `json:"params"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mu.Lock()
			createdVpcUuid, _ = payload.Params["vpcUuid"].(string)
			mu.Unlock()

			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"inventory": map[string]interface{}{
					"uuid":        "mock-firewall-uuid",
					"name":        "mock-firewall",
					"description": "mock firewall",
				},
			})
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"inventories": []interface{}{
					map[string]interface{}{
						"uuid":        "mock-firewall-uuid",
						"name":        "mock-firewall",
						"description": "mock firewall",
					},
				},
			})
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	hostParts := strings.Split(server.Listener.Addr().String(), ":")

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "zstack" {
  host             = "%s"
  port             = %s
  account_name     = "admin"
  account_password = "password"
}

resource "zstack_vpc_firewall" "test" {
  name        = "mock-firewall"
  description = "mock firewall"
  vpc_uuid    = "mock-vrouter-uuid"
}
`, hostParts[0], hostParts[1]),
			},
		},
	})

	mu.Lock()
	validated := virtualRouterValidated
	vpcUuid := createdVpcUuid
	mu.Unlock()

	if !validated {
		t.Fatal("expected create to validate vpc_uuid with GetVirtualRouterVm before creating firewall")
	}
	if vpcUuid != "mock-vrouter-uuid" {
		t.Fatalf("expected create request params.vpcUuid to be virtual router UUID, got %q", vpcUuid)
	}
}

func TestAccVpcFirewallResource_RejectsNonVirtualRouterUuid(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"inventory": map[string]string{
				"uuid": "mock-session-uuid",
			},
		})
	})

	mux.HandleFunc("/zstack/v1/vm-instances/appliances/virtual-routers/mock-l3-vpc-uuid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"inventories": []interface{}{},
		})
	})

	mux.HandleFunc("/zstack/v1/vpcfirewalls", func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("CreateVpcFirewall should not be called when vpc_uuid is not a virtual router UUID")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	hostParts := strings.Split(server.Listener.Addr().String(), ":")

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "zstack" {
  host             = "%s"
  port             = %s
  account_name     = "admin"
  account_password = "password"
}

resource "zstack_vpc_firewall" "test" {
  name     = "mock-firewall"
  vpc_uuid = "mock-l3-vpc-uuid"
}
`, hostParts[0], hostParts[1]),
				ExpectError: regexp.MustCompile(`(?s)Invalid VPC Firewall vpc_uuid.*virtual router instance UUID`),
			},
		},
	})
}
