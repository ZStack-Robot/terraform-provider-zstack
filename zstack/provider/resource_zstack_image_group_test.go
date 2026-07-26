// Copyright (c) ZStack.io, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	tfresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/param"
)

const (
	testImageGroupUuid = "11111111111111111111111111111111"
	testRootImageUuid  = "22222222222222222222222222222222"
	testDataImageUuid  = "33333333333333333333333333333333"
)

type imageGroupMutationCheck struct {
	mutate func()
}

func (c imageGroupMutationCheck) CheckState(_ context.Context, _ statecheck.CheckStateRequest, _ *statecheck.CheckStateResponse) {
	c.mutate()
}

func TestImageGroupResourceSchemaAndMetadata(t *testing.T) {
	var imageGroup imageGroupResource
	schemaResp := &resource.SchemaResponse{}
	imageGroup.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)

	for _, name := range []string{"root_volume_template_uuid", "name"} {
		attribute, ok := schemaResp.Schema.Attributes[name]
		if !ok || !attribute.IsRequired() {
			t.Errorf("%s should be a required attribute", name)
		}
	}
	for _, name := range []string{"uuid", "image_count", "status"} {
		attribute, ok := schemaResp.Schema.Attributes[name]
		if !ok || !attribute.IsComputed() {
			t.Errorf("%s should be a computed attribute", name)
		}
	}

	metadataResp := &resource.MetadataResponse{}
	imageGroup.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "zstack"}, metadataResp)
	if metadataResp.TypeName != "zstack_image_group" {
		t.Fatalf("unexpected resource type name: %s", metadataResp.TypeName)
	}
}

func TestImageGroupCreateParam(t *testing.T) {
	createParam := imageGroupCreateParam(imageGroupResourceModel{
		Name:                    types.StringValue("group"),
		Description:             types.StringValue("description"),
		DataVolumeTemplateUuids: stringSliceToSet([]string{testDataImageUuid}),
		ResourceUuid:            types.StringValue(testImageGroupUuid),
		TagUuids:                stringSliceToSet([]string{"tag-uuid"}),
		SystemTags:              stringSliceToSet([]string{"system-tag"}),
	})

	if createParam.Params.Name != "group" ||
		createParam.Params.Description == nil ||
		*createParam.Params.Description != "description" ||
		createParam.Params.ResourceUuid == nil ||
		*createParam.Params.ResourceUuid != testImageGroupUuid {
		t.Fatalf("unexpected create parameters: %#v", createParam)
	}
	if !slices.Equal(createParam.Params.DataVolumeTemplateUuids, []string{testDataImageUuid}) {
		t.Fatalf("unexpected data templates: %#v", createParam.Params.DataVolumeTemplateUuids)
	}
	if !slices.Equal(createParam.Params.TagUuids, []string{"tag-uuid"}) {
		t.Fatalf("unexpected tag UUIDs: %#v", createParam.Params.TagUuids)
	}
	if !slices.Equal(createParam.SystemTags, []string{"system-tag"}) {
		t.Fatalf("unexpected system tags: %#v", createParam.SystemTags)
	}
}

func TestImageGroupResourceRejectsEmptyConfiguredValues(t *testing.T) {
	testCases := map[string]string{
		"description":                        `description = ""`,
		"resource_uuid":                      `resource_uuid = ""`,
		"data_volume_template_uuids element": `data_volume_template_uuids = [""]`,
		"tag_uuids element":                  `tag_uuids = [""]`,
		"system_tags element":                `system_tags = [""]`,
	}

	for name, invalidConfiguration := range testCases {
		t.Run(name, func(t *testing.T) {
			tfresource.UnitTest(t, tfresource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []tfresource.TestStep{{
					Config: fmt.Sprintf(`
provider "zstack" {
  host              = "127.0.0.1"
  port              = 1
  access_key_id     = "test-access-key"
  access_key_secret = "test-access-key-secret"
}

resource "zstack_image_group" "test" {
  name                      = "invalid-image-group"
  root_volume_template_uuid = %q
  %s
}
`, testRootImageUuid, invalidConfiguration),
					PlanOnly:    true,
					ExpectError: regexp.MustCompile(`(?s)at.*least 1`),
				}},
			})
		})
	}
}

func TestImageGroupResourceMockContract(t *testing.T) {
	mux := http.NewServeMux()
	var mu sync.Mutex
	deleted := false
	includeDataRef := true
	sawCreate := false
	refQueryCount := 0
	sawExpunge := false

	writeInventory := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventory": map[string]any{
				"uuid":        testImageGroupUuid,
				"name":        "contract-image-group",
				"description": "contract test",
				"imageCount":  2,
				"status":      "Ready",
			},
		})
	}
	writeInventories := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{
				"uuid":        testImageGroupUuid,
				"name":        "contract-image-group",
				"description": "contract test",
				"imageCount":  2,
				"status":      "Ready",
			}},
		})
	}

	mux.HandleFunc("/zstack/v1/imagegroup/from/image/"+testRootImageUuid, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			SystemTags []string `json:"systemTags"`
			Params     struct {
				Name                    string   `json:"name"`
				Description             string   `json:"description"`
				DataVolumeTemplateUuids []string `json:"dataVolumeTemplateUuids"`
				ResourceUuid            string   `json:"resourceUuid"`
				TagUuids                []string `json:"tagUuids"`
			} `json:"params"`
		}
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if payload.Params.Name != "contract-image-group" ||
			payload.Params.Description != "contract test" ||
			payload.Params.ResourceUuid != testImageGroupUuid ||
			!slices.Equal(payload.Params.DataVolumeTemplateUuids, []string{testDataImageUuid}) ||
			!slices.Equal(payload.Params.TagUuids, []string{"tag-uuid"}) ||
			!slices.Equal(payload.SystemTags, []string{"system-tag"}) {
			http.Error(w, fmt.Sprintf("unexpected create payload: %#v", payload), http.StatusBadRequest)
			return
		}
		mu.Lock()
		sawCreate = true
		deleted = false
		includeDataRef = true
		mu.Unlock()
		writeInventory(w)
	})

	mux.HandleFunc("/zstack/v1/imagegroups/"+testImageGroupUuid, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		isDeleted := deleted
		mu.Unlock()
		if req.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		if isDeleted {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"inventories": []any{}})
			return
		}
		writeInventories(w)
	})

	mux.HandleFunc("/zstack/v1/imagegrouprefs", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		queries := req.URL.Query()["q"]
		if !slices.Contains(queries, "imageGroupUuid="+testImageGroupUuid) {
			http.Error(w, fmt.Sprintf("missing image group filter: %#v", queries), http.StatusBadRequest)
			return
		}
		mu.Lock()
		refQueryCount++
		withDataRef := includeDataRef
		mu.Unlock()
		refs := []map[string]any{
			{"uuid": "root-ref", "imageUuid": testRootImageUuid, "imageGroupUuid": testImageGroupUuid},
		}
		if withDataRef {
			refs = append(refs, map[string]any{
				"uuid": "data-ref", "imageUuid": testDataImageUuid, "imageGroupUuid": testImageGroupUuid,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": refs,
		})
	})

	mux.HandleFunc("/zstack/v1/images/"+testRootImageUuid, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{"uuid": testRootImageUuid, "mediaType": "RootVolumeTemplate"}},
		})
	})
	mux.HandleFunc("/zstack/v1/images/"+testDataImageUuid, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{"uuid": testDataImageUuid, "mediaType": "DataVolumeTemplate"}},
		})
	})

	mux.HandleFunc("/zstack/v1/imagegroups/"+testImageGroupUuid+"/actions", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := payload["expungeImageGroup"]; !ok {
			http.Error(w, fmt.Sprintf("unexpected expunge payload: %#v", payload), http.StatusBadRequest)
			return
		}
		mu.Lock()
		deleted = true
		sawExpunge = true
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ZSTACK_HOST", host)
	t.Setenv("ZSTACK_PORT", port)
	t.Setenv("ZSTACK_ACCESS_KEY_ID", "test-access-key")
	t.Setenv("ZSTACK_ACCESS_KEY_SECRET", "test-access-key-secret")

	config := fmt.Sprintf(`
provider "zstack" {
  host              = %q
  port              = %s
  access_key_id     = "test-access-key"
  access_key_secret = "test-access-key-secret"
}

resource "zstack_image_group" "test" {
  name                       = "contract-image-group"
  description                = "contract test"
  root_volume_template_uuid  = %q
  data_volume_template_uuids = [%q]
  resource_uuid              = %q
  tag_uuids                  = ["tag-uuid"]
  system_tags                = ["system-tag"]
}
`, host, port, testRootImageUuid, testDataImageUuid, testImageGroupUuid)

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckImageGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("uuid"), knownvalue.StringExact(testImageGroupUuid)),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("image_count"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("status"), knownvalue.StringExact("Ready")),
				},
			},
			{
				ResourceName:                         "zstack_image_group.test",
				ImportState:                          true,
				ImportStateId:                        testImageGroupUuid,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"tag_uuids", "system_tags"},
			},
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					imageGroupMutationCheck{mutate: func() {
						mu.Lock()
						includeDataRef = false
						mu.Unlock()
					}},
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

	mu.Lock()
	defer mu.Unlock()
	if !sawCreate || refQueryCount < 2 || !sawExpunge {
		t.Fatalf(
			"incomplete contract coverage: create=%t ref-queries=%d expunge=%t",
			sawCreate,
			refQueryCount,
			sawExpunge,
		)
	}
}

func TestImageGroupResourceMockNoOpAndDisappears(t *testing.T) {
	mux := http.NewServeMux()
	var mu sync.Mutex
	deleted := false

	mux.HandleFunc("/zstack/v1/imagegroup/from/image/"+testRootImageUuid, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		deleted = false
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventory": map[string]any{
				"uuid": testImageGroupUuid, "name": "minimal-image-group", "imageCount": 1, "status": "Ready",
			},
		})
	})
	mux.HandleFunc("/zstack/v1/imagegroups/"+testImageGroupUuid, func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		isDeleted := deleted
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if isDeleted {
			_ = json.NewEncoder(w).Encode(map[string]any{"inventories": []any{}})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{
				"uuid": testImageGroupUuid, "name": "minimal-image-group", "imageCount": 1, "status": "Ready",
			}},
		})
	})
	mux.HandleFunc("/zstack/v1/imagegrouprefs", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{
				"uuid": "root-ref", "imageUuid": testRootImageUuid, "imageGroupUuid": testImageGroupUuid,
			}},
		})
	})
	mux.HandleFunc("/zstack/v1/images/"+testRootImageUuid, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"inventories": []map[string]any{{"uuid": testRootImageUuid, "mediaType": "RootVolumeTemplate"}},
		})
	})
	mux.HandleFunc("/zstack/v1/imagegroups/"+testImageGroupUuid+"/actions", func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		deleted = true
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ZSTACK_HOST", host)
	t.Setenv("ZSTACK_PORT", port)
	t.Setenv("ZSTACK_ACCESS_KEY_ID", "test-access-key")
	t.Setenv("ZSTACK_ACCESS_KEY_SECRET", "test-access-key-secret")

	config := fmt.Sprintf(`
provider "zstack" {
  host              = %q
  port              = %s
  access_key_id     = "test-access-key"
  access_key_secret = "test-access-key-secret"
}

resource "zstack_image_group" "test" {
  name                      = "minimal-image-group"
  root_volume_template_uuid = %q
}
`, host, port, testRootImageUuid)

	tfresource.UnitTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckImageGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
			},
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckImageGroupDisappears("zstack_image_group.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccImageGroupResource(t *testing.T) {
	rootImageUuid := imageGroupAcceptanceRootTemplate(t)
	name := testAccName("image-group")
	config := imageGroupAcceptanceConfig(name, rootImageUuid)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckImageGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: config,
				ConfigPlanChecks: tfresource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("root_volume_template_uuid"), knownvalue.StringExact(rootImageUuid)),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("image_count"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue("zstack_image_group.test", tfjsonpath.New("status"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:                         "zstack_image_group.test",
				ImportState:                          true,
				ImportStateIdFunc:                    importStateIdFromUUID("zstack_image_group.test"),
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"tag_uuids", "system_tags"},
			},
		},
	})
}

func TestAccImageGroupResourceDisappears(t *testing.T) {
	rootImageUuid := imageGroupAcceptanceRootTemplate(t)

	tfresource.ParallelTest(t, tfresource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckImageGroupDestroy,
		Steps: []tfresource.TestStep{
			{
				Config: imageGroupAcceptanceConfig(testAccName("image-group-disappears"), rootImageUuid),
				ConfigStateChecks: []statecheck.StateCheck{
					stateCheckImageGroupDisappears("zstack_image_group.test"),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func imageGroupAcceptanceConfig(name, rootImageUuid string) string {
	return providerConfig() + fmt.Sprintf(`
resource "zstack_image_group" "test" {
  name                      = %q
  description               = "Terraform image group acceptance test"
  root_volume_template_uuid = %q
}
`, name, rootImageUuid)
}

func imageGroupAcceptanceRootTemplate(t *testing.T) string {
	t.Helper()
	if rootImageUuid := os.Getenv("ZSTACK_TEST_IMAGE_GROUP_ROOT_IMAGE_UUID"); rootImageUuid != "" {
		return rootImageUuid
	}
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}

	query := param.NewQueryParam()
	query.AddQ("mediaType=RootVolumeTemplate")
	query.AddQ("status=Ready")
	images, err := testAccClientLoggedIn().QueryImage(&query)
	if err != nil {
		t.Fatalf("query a root volume template for image group acceptance test: %v", err)
	}
	if len(images) == 0 {
		t.Skip("no Ready RootVolumeTemplate image is available")
	}
	return images[0].UUID
}
