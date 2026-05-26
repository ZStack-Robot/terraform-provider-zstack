// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestInstanceStateResource_Schema(t *testing.T) {
	var r instanceStateResource
	resp := &resource.SchemaResponse{}
	r.Schema(context.Background(), resource.SchemaRequest{}, resp)
	if len(resp.Schema.Attributes) == 0 {
		t.Fatal("schema should not be empty")
	}

	required := []string{"vm_instance_uuid", "state"}
	for _, attr := range required {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing required attribute %q", attr)
		}
		if !a.IsRequired() {
			t.Errorf("attribute %q should be required", attr)
		}
	}

	computed := []string{"id"}
	for _, attr := range computed {
		a, ok := resp.Schema.Attributes[attr]
		if !ok {
			t.Fatalf("schema missing computed attribute %q", attr)
		}
		if !a.IsComputed() {
			t.Errorf("attribute %q should be computed", attr)
		}
	}

	stopTypeAttr, ok := resp.Schema.Attributes["stop_type"].(rschema.StringAttribute)
	if !ok {
		t.Fatal("stop_type should be a string attribute")
	}
	if !stopTypeAttr.IsOptional() || !stopTypeAttr.IsComputed() {
		t.Fatal("stop_type should be optional and computed")
	}

	timeoutAttr, ok := resp.Schema.Attributes["operation_timeout"].(rschema.Int64Attribute)
	if !ok {
		t.Fatal("operation_timeout should be an int64 attribute")
	}
	if !timeoutAttr.IsOptional() || !timeoutAttr.IsComputed() {
		t.Fatal("operation_timeout should be optional and computed")
	}
}

func TestInstanceStateResource_Metadata(t *testing.T) {
	var r instanceStateResource
	resp := &resource.MetadataResponse{}
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, resp)
	if resp.TypeName != "zstack_instance_state" {
		t.Errorf("unexpected type name: %s", resp.TypeName)
	}
}

func TestInstanceStateForRead(t *testing.T) {
	tests := []struct {
		name           string
		actualState    string
		previousState  types.String
		wantState      string
		wantNormalized bool
		wantErr        bool
	}{
		{
			name:        "running terminal state",
			actualState: instanceStateRunning,
			wantState:   instanceStateRunning,
		},
		{
			name:        "stopped terminal state",
			actualState: instanceStateStopped,
			wantState:   instanceStateStopped,
		},
		{
			name:           "starting maps to running",
			actualState:    "Starting",
			wantState:      instanceStateRunning,
			wantNormalized: true,
		},
		{
			name:           "stopping maps to stopped",
			actualState:    "Stopping",
			wantState:      instanceStateStopped,
			wantNormalized: true,
		},
		{
			name:           "unknown state keeps previous running state",
			actualState:    "Unknown",
			previousState:  types.StringValue(instanceStateRunning),
			wantState:      instanceStateRunning,
			wantNormalized: true,
		},
		{
			name:           "unknown state keeps previous stopped state",
			actualState:    "Unknown",
			previousState:  types.StringValue(instanceStateStopped),
			wantState:      instanceStateStopped,
			wantNormalized: true,
		},
		{
			name:          "unknown state without previous terminal state returns error",
			actualState:   "Unknown",
			previousState: types.StringNull(),
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, gotNormalized, err := instanceStateForRead(tt.actualState, tt.previousState)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotState != tt.wantState {
				t.Fatalf("expected state %q, got %q", tt.wantState, gotState)
			}
			if gotNormalized != tt.wantNormalized {
				t.Fatalf("expected normalized %t, got %t", tt.wantNormalized, gotNormalized)
			}
		})
	}
}

func TestAccInstanceStateResource_StartStopDeleteNoop(t *testing.T) {
	mux := http.NewServeMux()

	var mu sync.Mutex
	currentState := instanceStateStopped
	actions := make([]string, 0, 2)
	stopType := ""

	mux.HandleFunc("/zstack/v1/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventory": map[string]string{
				"uuid": "mock-session-uuid",
			},
		})
	})

	mux.HandleFunc("/zstack/v1/vm-instances/mock-vm-uuid", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		mu.Lock()
		state := currentState
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []any{
				map[string]any{
					"uuid":  "mock-vm-uuid",
					"name":  "mock-vm",
					"state": state,
				},
			},
		})
	})

	mux.HandleFunc("/zstack/v1/vm-instances/mock-vm-uuid/actions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}

		var payload map[string]map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		switch {
		case payload["startVmInstance"] != nil:
			actions = append(actions, "start")
			currentState = instanceStateRunning
		case payload["stopVmInstance"] != nil:
			actions = append(actions, "stop")
			stopType, _ = payload["stopVmInstance"]["type"].(string)
			currentState = instanceStateStopped
		default:
			mu.Unlock()
			http.Error(w, "unexpected action", http.StatusBadRequest)
			return
		}
		state := currentState
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventory": map[string]any{
				"uuid":  "mock-vm-uuid",
				"name":  "mock-vm",
				"state": state,
			},
		})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	hostParts := strings.Split(server.Listener.Addr().String(), ":")
	config := func(state string) string {
		return fmt.Sprintf(`
provider "zstack" {
  host             = "%s"
  port             = %s
  account_name     = "admin"
  account_password = "password"
}

resource "zstack_instance_state" "test" {
  vm_instance_uuid   = "mock-vm-uuid"
  state              = %q
  operation_timeout  = 5
}
`, hostParts[0], hostParts[1], state)
	}

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: config(instanceStateRunning),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("id"), knownvalue.StringExact("mock-vm-uuid")),
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("state"), knownvalue.StringExact(instanceStateRunning)),
				},
			},
			{
				Config: config(instanceStateStopped),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("state"), knownvalue.StringExact(instanceStateStopped)),
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("stop_type"), knownvalue.StringExact(defaultInstanceStateStopType)),
				},
			},
		},
	})

	mu.Lock()
	gotActions := append([]string(nil), actions...)
	gotStopType := stopType
	mu.Unlock()

	wantActions := []string{"start", "stop"}
	if strings.Join(gotActions, ",") != strings.Join(wantActions, ",") {
		t.Fatalf("expected actions %v, got %v", wantActions, gotActions)
	}
	if gotStopType != defaultInstanceStateStopType {
		t.Fatalf("expected stop type %q, got %q", defaultInstanceStateStopType, gotStopType)
	}
}

func TestAccInstanceStateResource_RealEnvExistingVM(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("acceptance test skipped unless TF_ACC=1")
	}

	cli := testAccClientLoggedIn()

	targetUUID := os.Getenv("ZSTACK_INSTANCE_STATE_TEST_VM_UUID")
	if targetUUID == "" {
		t.Skip("set ZSTACK_INSTANCE_STATE_TEST_VM_UUID to a test VM UUID before running this acceptance test")
	}

	vm, err := cli.GetVmInstance(targetUUID)
	if err != nil {
		t.Fatalf("get VM instance %s: %v", targetUUID, err)
	}
	if vm.Type != "UserVm" {
		t.Fatalf("ZSTACK_INSTANCE_STATE_TEST_VM_UUID must refer to a UserVm, got type %q", vm.Type)
	}
	if vm.State != instanceStateRunning {
		t.Fatalf("ZSTACK_INSTANCE_STATE_TEST_VM_UUID must start in Running state, got %q", vm.State)
	}

	config := func(state string) string {
		return providerConfig() + fmt.Sprintf(`
resource "zstack_instance_state" "test" {
  vm_instance_uuid  = %q
  state             = %q
  operation_timeout = 300
}
`, targetUUID, state)
	}

	t.Logf("using existing VM %q (%s) for start/stop acceptance test", vm.Name, targetUUID)

	tfresource.Test(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []tfresource.TestStep{
			{
				Config: config(instanceStateStopped),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("state"), knownvalue.StringExact(instanceStateStopped)),
				},
			},
			{
				Config: config(instanceStateRunning),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_instance_state.test", tfjsonpath.New("state"), knownvalue.StringExact(instanceStateRunning)),
				},
			},
			{
				Config:   config(instanceStateRunning),
				PlanOnly: true,
			},
		},
	})
}
