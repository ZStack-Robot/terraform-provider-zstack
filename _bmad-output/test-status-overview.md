# Test Status Overview — terraform-provider-zstack

**Date**: 2026-04-20
**Branch**: `test/progress`
**Base**: `fix/review-findings-2026-04-13`

---

## Executive Summary

| Metric | Count | Coverage |
|---|---|---|
| Resource implementations | 111 | — |
| Resources with ANY test file | 27 | **24%** |
| Resources with acceptance tests | 19 | **17%** |
| Resources with unit-only tests (no acc) | 8 | 7% |
| Resources with NO test file | 84 | **76% untested** |
| Data source implementations | 42 | — |
| Data sources with ANY test file | 29 | **69%** |
| Data sources with NO test file | 13 | **31% untested** |
| Total acceptance test functions | 78 | — |
| Total unit test functions | 53 | — |
| Build status | PASS | `go build ./...` |
| Unit test status | 52/53 PASS | 1 expected failure |

**Overall test coverage: ~24% of resources, ~69% of data sources.**

---

## Build & Unit Test Results

### Build
```
go build ./...  ✅ PASS (exit 0)
```

### Unit Tests (non-acceptance)
```
go test ./zstack/provider/ -run 'Test[^A]' -timeout 5m
Result: 52 PASS, 1 FAIL
```

| Test | Status | Notes |
|---|---|---|
| All `*_Schema` tests (26 tests) | ✅ PASS | Schema validation for resources with tests |
| All `*_Metadata` tests (25 tests) | ✅ PASS | Metadata/type name validation |
| `TestNoEmptyUUIDStateCorruption` | ✅ PASS | Anti-pattern scanner |
| `TestNoReadRemoveResourceOnTransientError` | ✅ PASS | Anti-pattern scanner |
| `TestNoDeleteEmptyUUIDGuard` | ✅ PASS | Anti-pattern scanner |
| `TestQueryEnvironment` | ❌ FAIL | **Expected**: requires `ZSTACK_ACCESS_KEY_ID`/`ZSTACK_ACCESS_KEY_SECRET` env vars. Should be guarded with `t.Skip`. |

### Known Issue: `TestQueryEnvironment`
This test fails when env vars are not set. Unlike other env-dependent tests which use `loadEnvData(t)` (which calls `t.Skip`), `TestQueryEnvironment` directly checks env vars and calls `t.Fatalf`. This is a test bug — it should use `t.Skip` for missing credentials.

---

## Resource Test Coverage Detail

### Tier 1: Resources WITH Acceptance Tests (19 resources) ✅

| Resource | Unit Tests | Acc Tests | Disappears | Destroy Check | Import |
|---|---|---|---|---|---|
| zone | 2 (Schema/Meta) | 2 (basic + disappears) | ✅ | ✅ | ✅ |
| account | 2* | 2 (basic + disappears) | ✅ | ✅ | — |
| affinity_group | 2* | 2 (basic + disappears) | ✅ | ✅ | — |
| cluster | 2 | 2 (basic + disappears) | ✅ | ✅ | — |
| ssh_key_pair | 2 | 2 (basic + disappears) | ✅ | ✅ | — |
| iam2_project | 2 | 2 (basic + disappears) | ✅ | ✅ | — |
| image | 0 | 1 (create) | — | ✅ | — |
| volume | 1 (Schema) | 1 | — | ✅ | — |
| l2vlan_network | 2 | 1 | — | ✅ | — |
| load_balancer | 2 | 1 | — | ✅ | — |
| load_balancer_listener | 2 | 1 | — | ✅ | — |
| port_forwarding_rule | 2 | 1 | — | ✅ | — |
| vip | 2 | 1 | — | ✅ | — |
| virtual_router_image | 2 | 1 | — | ✅ | — |
| virtual_router_offering | 2 | 1 | — | ✅ | — |
| virtual_router_instance | 0 | 1 | — | — | — |
| auto_scaling_group | 0 | 1 | — | ✅ | — |
| reserved_ip | 0 | 1 | — | ✅ | — |
| networking_secgroup_attachment | 0 | 1 | — | ✅ | — |

*Note: account/affinity_group unit count appears as 0 in grep due to test naming pattern (TestAccountResource_Schema matches TestAcc prefix)*

**Test quality observations for Tier 1:**
- Only `zone` has import state test
- Only 6 resources have `_disappears` tests
- No update/modify tests observed (only create + destroy)
- Good: destroy checks exist for most resources
- Good: modern `statecheck.StateCheck` pattern used (not deprecated `TestCheckFunc`)

### Tier 2: Resources with Unit Tests ONLY (8 resources) ⚠️

| Resource | Unit Tests | Acc Tests | Notes |
|---|---|---|---|
| backup_storage | 2 (Schema/Meta) | 0 | Needs acc tests |
| host | 2 | 0 | Needs acc tests |
| networking_secgroup_rule | 2 | 0 | Needs acc tests |
| primary_storage | 2 | 0 | Needs acc tests |
| subnet_ip_range | 2 | 0 | Needs acc tests |
| tag_attachment | 2 | 0 | Needs acc tests |
| volume_snapshot | 2 | 0 | Needs acc tests |
| vpc | 2 | 0 | Needs acc tests |

### Tier 3: Resources with NO Test File (84 resources) ❌

#### Core / High Priority (frequently used by customers)
- `instance` — **CRITICAL**: Most important resource, has TODO in code
- `disk_offering` — Core compute resource
- `instance_offering` — Core compute resource
- `instance_scripts` / `instance_scripts_execution` — Operations
- `eip` — Core networking
- `l3network` — Core networking
- `networking_secgroup` — Core security
- `vm_nic` — Core networking
- `tag` — Universal tagging

#### IAM & Access Control
- `iam2_virtual_id`
- `iam2_organization`
- `user`
- `role`
- `policy`
- `access_key`
- `access_control_list`

#### Networking (Advanced)
- `vpc_firewall`
- `vpc_ha_group`
- `vpc_shared_qos`
- `l2vxlan_network`
- `vip_qos`
- `ipsec_connection`
- `policy_route_rule_set` / `policy_route_rule`
- `vrouter_route_table` / `vrouter_route_entry`
- `flow_meter` / `flow_collector`
- `multicast_router`
- `sdn_controller`
- `port_mirror` / `port_mirror_session`

#### Storage
- `ceph_primary_storage`
- `ceph_backup_storage`
- `ceph_pool`
- `image_store_backup_storage`
- `volume_backup`
- `database_backup`
- `cdp_policy` / `cdp_task`
- `nvme_server`
- `iscsi_server`

#### Monitoring / Alerting
- `alarm`
- `sns_topic`
- `sns_email_endpoint` / `sns_http_endpoint`
- `email_media`
- `monitor_template` / `monitor_group`
- `webhook`
- `log_server`
- `snmp_agent`

#### Scheduling
- `scheduler_job`
- `scheduler_trigger`

#### Misc / Specialized
- `global_config`
- `certificate`
- `lb_server_group`
- `guest_tool_attachment`
- `vm_cdrom`
- `preconfiguration_template`
- `resource_stack` / `stack_template`
- `price_table`
- `ldap_server`
- `directory`
- `vcenter`
- `v2v_conversion_host`
- `container_management_endpoint`
- `dataset`
- `zbox_backup`
- `license`
- `pci_device_offering`

#### Security Machines (specialized hardware)
- `baremetal_chassis` / `baremetal_instance` / `baremetal_pxe_server`
- `jit_security_machine`
- `san_sec_security_machine`
- `info_sec_security_machine`
- `fi_sec_security_machine`
- `flk_sec_security_machine`
- `aliyun_proxy_vpc` / `aliyun_proxy_vswitch` / `aliyun_nas_access_group`

---

## Data Source Test Coverage Detail

### Data Sources WITH Tests (29) ✅

| Data Source | Unit Tests | Acc Tests |
|---|---|---|
| zone | 0 | 1 |
| accounts | 0 | 1 |
| affinity_groups | 0 | 1 |
| auto_scaling_groups | 0 | 1 |
| backup_storages | 0 | 3 (basic + filterByName + filterByNamePattern) |
| clusters | 0 | 3 |
| disk_offers | 0 | 3 |
| gpu_devices | 2 (Schema/Meta) | 1 |
| hosts | 0 | 3 |
| iam2_projects | 0 | 1 |
| images | 0 | 3 |
| instance_guest_tools | 0 | 1 |
| instance_offers | 0 | 3 |
| instance_scripts | 0 | 3 |
| instances | 0 | 3 |
| l2networks | 0 | 3 |
| l2vlan_networks | 2 | 1 |
| l3networks | 0 | 3 |
| load_balancer_listeners | 2 | 1 |
| load_balancers | 2 | 1 |
| mn_nodes | 0 | 1 |
| networking_secgroup_rules | 0 | 1 |
| networking_secgroups | 0 | 1 |
| port_forwarding_rules | 2 | 1 |
| sdn_controllers | 0 | 1 |
| ssh_key_pairs | 0 | 1 |
| virtual_router_offers | 0 | 1 |
| virtual_routers | 0 | 1 |
| virtural_router_images | 0 | 3 |

*Note: `virtural_router_images` has a typo in filename (should be `virtual`)*

### Data Sources WITHOUT Tests (13) ❌

| Data Source | Notes |
|---|---|
| disks | No test file |
| eips | No test file |
| hook_scripts | No test file |
| license_authorized_capacity | No test file |
| license_authorized_nodes | No test file |
| primary_storages | No test file |
| reserved_ips | No test file |
| subnet_ip_ranges | No test file |
| tags | No test file |
| user_tags | No test file |
| vips | No test file |
| volume_snapshots | No test file |
| volumes | No test file |

---

## Test Infrastructure Assessment

### Strengths ✅
1. **Modern test patterns** — Uses `statecheck.StateCheck` (not deprecated `TestCheckFunc`)
2. **Disappears test infrastructure** — Generic `disappearsCheck` struct in `check_disappears_test.go`
3. **Destroy check infrastructure** — Generic `testAccCheckResourceDestroyByGet` helper
4. **Anti-pattern scanners** — `resource_antipattern_test.go` catches common bugs:
   - Empty UUID state corruption
   - RemoveResource on transient errors
   - Delete empty UUID guard (dead code detection)
5. **Environment data loader** — `testdata/env.json` + `generate_env.go` for real UUIDs
6. **Batch TF generator** — `generate_tf.go` creates standalone `.tf` files per resource
7. **Provider test helpers** — `provider_test.go` with `testAccProtoV6ProviderFactories`

### Weaknesses ❌
1. **No CI test stage** — `.gitlab-ci.yml` only has `build` stage, no `test` stage
2. **No unit tests for 84 resources** — Not even Schema/Metadata validation
3. **No update/modify acceptance tests** — All acc tests only test create+destroy, not update-in-place
4. **No import tests** — Only `zone` has ImportState test; all other resources lack import testing
5. **TestQueryEnvironment bug** — Uses `t.Fatalf` instead of `t.Skip` for missing env vars
6. **Filename typo** — `data_source_zstack_virtural_router_images.go` (should be `virtual`)
7. **Missing disappears tests** — Only 6/19 tested resources have `_disappears` variant
8. **No error/edge-case unit tests** — No tests for API error handling, invalid inputs, etc.
9. **instance resource has TODO** — Delete/Read behavior flagged as incomplete

---

## CI/CD Pipeline Status

### Current `.gitlab-ci.yml`
```yaml
stages:
  - build   # Only stage — no test stage!

build_job:
  script:
    - GOARCH=amd64 go build -o dist/terraform-provider-zstack
    - GOARCH=arm64 go build -o dist/terraform-provider-zstack-aarch64
    # Upload to Minio
```

### Missing
- No `test` stage for unit tests (`go test ./... -run 'Test[^A]'`)
- No `acceptance-test` stage (even optional/manual)
- No lint stage (`golangci-lint`, `staticcheck`)
- No doc generation verification

---

## Recommended Improvement Plan

### Phase 1: Quick Wins (1-2 days)
1. **Fix `TestQueryEnvironment`** — Change `t.Fatalf` to `t.Skip` for missing env vars
2. **Fix filename typo** — `virtural` → `virtual` in data source filename
3. **Add CI test stage** — Unit tests (`Test[^A]`) in `.gitlab-ci.yml`
4. **Add Schema/Metadata unit tests** for Tier 2 resources that lack them (8 resources already have test files but some are missing Schema/Meta)

### Phase 2: Unit Test Coverage (1-2 weeks)
5. **Add Schema/Metadata unit tests** for ALL 84 untested resources — Mechanical, can be generated
6. **Add Schema/Metadata unit tests** for 13 untested data sources
7. **Add anti-pattern test coverage** — Extend `resource_antipattern_test.go` to scan for additional patterns

### Phase 3: Acceptance Test Coverage — Core Resources (2-3 weeks)
8. **Add acceptance tests for core resources** (highest customer impact):
   - `instance` (most critical — also fix the TODO)
   - `disk_offering`, `instance_offering`
   - `eip`, `l3network`, `networking_secgroup`
   - `tag`, `vm_nic`
9. **Add acceptance tests for IAM resources**:
   - `user`, `role`, `policy`, `iam2_virtual_id`, `iam2_organization`
10. **Add acceptance tests for Tier 2 resources** (have unit tests, need acc tests):
    - `backup_storage`, `host`, `primary_storage`, `vpc`, etc.

### Phase 4: Test Quality Improvements (1-2 weeks)
11. **Add update/modify tests** — Test attribute changes (not just create+destroy)
12. **Add ImportState tests** — Currently only `zone` has import test
13. **Add `_disappears` tests** — Extend to all resources (infrastructure exists, just needs configuration)
14. **Add error handling tests** — Invalid inputs, API failure simulation

### Phase 5: Advanced (ongoing)
15. **Add acceptance tests for networking resources** — VPC, firewall, routes, etc.
16. **Add acceptance tests for storage resources** — Ceph, backups, snapshots
17. **Add acceptance tests for monitoring resources** — Alarms, SNS, webhooks
18. **CI acceptance test stage** — Manual/scheduled, against test environment

---

## Appendix: Test File Inventory

### Test Infrastructure Files
| File | Purpose |
|---|---|
| `provider_test.go` | Provider factories, config helpers |
| `test_env_loader_test.go` | Load `testdata/env.json`, helper functions |
| `check_destroy_test.go` | Destroy verification helpers (15 resources) |
| `check_disappears_test.go` | Disappears simulation helpers (6 resources) |
| `resource_antipattern_test.go` | Anti-pattern scanner (3 checks) |
| `query_env_test.go` | Environment query test (has bug) |
| `testdata/generate_env.go` | Generate env.json from live environment |
| `testdata/generate_tf.go` | Generate batch `.tf` test files |
| `scripts/run_tests.sh` | Test execution script |
