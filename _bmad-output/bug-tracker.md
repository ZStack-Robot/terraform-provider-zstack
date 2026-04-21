# Bug Tracker — terraform-provider-zstack

> Generated: 2026-04-20  
> Updated: 2026-04-21  
> Branch: `test/progress`  
> Tools used: `golangci-lint run`, `go vet`, `go test -short`, manual code review, automated codebase scanning

---

## Fix Status

**Commit**: `5a52557` — `fix: resolve 17 bugs from QA audit (BUG-001 through BUG-021)`

| Verification | Before | After |
|---|---|---|
| `go build ./...` | ✅ clean | ✅ clean |
| `go vet ./...` | ✅ clean | ✅ clean |
| `golangci-lint run ./...` | 13 issues | ✅ 0 issues |
| `go test ./... -short` | 53 pass, 1 fail | ✅ 53 pass, 0 fail |

### Fixed (22 bugs)

| Bug | Priority | Status | Description |
|-----|----------|--------|-------------|
| BUG-001 | P0 | ✅ Fixed | Remove dead code loop in tag_attachment Delete |
| BUG-002 | P1 | ✅ Fixed | Remove empty if-branch in instance resource |
| BUG-003 | P0 | ✅ Fixed | Replace deprecated strings.Title with cases.Title |
| BUG-004 | P0 | ✅ Fixed | Fix nil pointer panic in testAccClientLoggedIn |
| BUG-005 | P0 | ✅ Fixed | Fix TestQueryEnvironment failure (t.Fatal → t.Skip) |
| BUG-006 | P1 | ✅ Fixed | Remove hardcoded default credentials from tests |
| BUG-007 | P1 | ✅ Fixed | Fix 6 unchecked error return values |
| BUG-008 | P1 | ✅ Fixed | Simplify embedded field selectors |
| BUG-009 | P2 | ✅ Fixed | Refactor if-chain to switch in filter.go |
| BUG-010 | P2 | ✅ Fixed | Remove unused copyStringValues function |
| BUG-011 | P2 | ✅ Fixed | Remove unused envInt function |
| BUG-012 | P2 | ✅ Fixed | Fix typo virtural → virtual in filename |
| BUG-013 | P2 | ✅ Fixed | Fix typo BackupStorges → BackupStorages |
| BUG-014 | P2 | ✅ Fixed | Fix typo gpuDeviceTyp → gpuDeviceType |
| BUG-015 | P2 | ✅ Fixed | Fix truncated description in clusters data source |
| BUG-016 | P2 | ✅ Fixed | Standardize factory function naming in provider.go |
| BUG-017 | P3 | ✅ Fixed | Align vmResource/InstanceResource naming |
| BUG-020 | P2 | ✅ Fixed | Convert vague TODO to actionable comment |
| BUG-021 | P1 | ✅ Fixed | Fix antipattern test glob to use absolute path |
| BUG-022 | P2 | ✅ Fixed | Replace string scanning with AST in antipattern tests |
| BUG-023 | P2 | ✅ Fixed | Randomize test resource names |
| BUG-019 | P2 | ✅ Fixed | Move disk state logic to Read() per TODO |

### Remaining (11 bugs — deferred / lower priority)

| Bug | Priority | Status | Description |
|-----|----------|--------|-------------|
| BUG-018 | P3 | 🔲 Open | Standardize acronym casing (UUID/Uuid/IP/Ip) |
| BUG-024 | P3 | 🔲 Open | Add update steps to acceptance tests |
| BUG-025 | P2 | 🔲 Open | Clean up commented-out code blocks in 19+ files |
| BUG-026–033 | P2–P3 | 🔲 Open | Various naming consistency and test improvements |

---

## Summary

| Category | Critical | High | Medium | Low | Total |
|----------|----------|------|--------|-----|-------|
| Logic Bugs | 1 | 1 | 1 | 0 | 3 |
| Error Handling | 0 | 1 | 2 | 0 | 3 |
| Linter / Static Analysis | 0 | 2 | 3 | 1 | 6 |
| Dead Code | 0 | 0 | 2 | 1 | 3 |
| Spelling / Naming | 0 | 2 | 4 | 2 | 8 |
| Test Infrastructure | 1 | 2 | 3 | 1 | 7 |
| TODO / Tech Debt | 0 | 1 | 1 | 0 | 2 |
| Deprecated API | 0 | 1 | 0 | 0 | 1 |
| **Total** | **2** | **10** | **16** | **5** | **33** |

---

## BUG-001 — `append` result discarded in tag_attachment Delete (Logic Bug)

- **Severity**: Critical
- **File**: `zstack/provider/resource_zstack_tag_attachment.go:205-210`
- **Category**: Logic Bug / SA4010
- **Detected by**: golangci-lint (staticcheck SA4010)

```go
// Line 205-210
resourceUuids := make([]string, 0, len(state.ResourceUuids.Elements()))
for _, v := range state.ResourceUuids.Elements() {
    resourceUuids = append(resourceUuids, v.(types.String).ValueString())
}

err := r.client.DetachTagFromResources(state.TagUuid.ValueString(), param.DeleteModePermissive)
```

**Problem**: `resourceUuids` is populated but **never passed** to `DetachTagFromResources`. The function call on line 210 only passes `TagUuid` and `DeleteModePermissive` — the specific resource UUIDs to detach from are ignored. This means **Delete detaches the tag from ALL resources** instead of just the ones in state, or the SDK ignores the missing parameter entirely.

**Fix**: Pass `resourceUuids` to the API call, or verify the SDK's `DetachTagFromResources` signature. The `resourceUuids` variable is computed but never used — this is almost certainly a bug where the intended behavior was to detach from specific resources.

---

## BUG-002 — Empty branch in instance resource (Logic Bug)

- **Severity**: Medium
- **File**: `zstack/provider/resource_zstack_instance.go:1143-1146`
- **Category**: Logic Bug / SA9003
- **Detected by**: golangci-lint (staticcheck SA9003)

```go
if dataDiskCephPoolName != "" {
    // Note: Pools field no longer available in PrimaryStorageInventoryView in SDK v2
    // Pool name validation is skipped; the API will validate on the server side
}
```

**Problem**: Empty `if` branch does nothing. The validation was removed during SDK v2 migration but the conditional structure was left behind.

**Fix**: Remove the empty `if` block entirely, keeping only the comment if documentation is needed.

---

## BUG-003 — `strings.Title` deprecated since Go 1.18 (Deprecated API)

- **Severity**: High
- **File**: `zstack/utils/filter.go:38`
- **Category**: Deprecated API / SA1019
- **Detected by**: golangci-lint (staticcheck SA1019)

```go
fieldName := strings.Title(apiFieldName)
```

**Problem**: `strings.Title` has been deprecated since Go 1.18 because it doesn't handle Unicode punctuation correctly. The Go team recommends `golang.org/x/text/cases` instead.

**Fix**: Replace with:
```go
import "golang.org/x/text/cases"
import "golang.org/x/text/language"

caser := cases.Title(language.Und)
fieldName := caser.String(apiFieldName)
```

---

## BUG-004 — Test helper panics on login failure (Test Infrastructure)

- **Severity**: Critical
- **File**: `zstack/provider/provider_test.go:104-113`
- **Category**: Test Infrastructure
- **Detected by**: Code review

```go
func testAccClientLoggedIn() *client.ZSClient {
    cli := testAccClient()
    if os.Getenv("ZSTACK_ACCESS_KEY_ID") == "" {
        if _, err := cli.Login(context.Background()); err != nil {
            panic(fmt.Sprintf("testAccClientLoggedIn: login failed: %v", err))
        }
    }
    return cli
}
```

**Problem**: `panic()` in test helpers crashes the **entire test binary** — not just the current test. This makes CI debugging extremely difficult: all subsequent tests are killed, stack traces are unhelpful, and the Go test harness cannot produce proper failure reports.

**Fix**: Change signature to accept `*testing.T` and use `t.Fatalf` or `t.Skip`:
```go
func testAccClientLoggedIn(t *testing.T) *client.ZSClient {
    t.Helper()
    cli := testAccClient()
    if os.Getenv("ZSTACK_ACCESS_KEY_ID") == "" {
        if _, err := cli.Login(context.Background()); err != nil {
            t.Skipf("testAccClientLoggedIn: login failed (skipping): %v", err)
        }
    }
    return cli
}
```

---

## BUG-005 — `TestQueryEnvironment` uses `t.Fatal` for missing env vars (Test Infrastructure)

- **Severity**: High
- **File**: `zstack/provider/query_env_test.go:28-30`
- **Category**: Test Infrastructure
- **Detected by**: `go test` (confirmed FAIL)

```go
if akID == "" || akSecret == "" {
    t.Fatal("ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
}
```

**Problem**: `t.Fatal` causes CI to report FAIL when env vars are not set. This is the **only test that fails** in the current suite (`go test -short` reports 1 FAIL / 53 PASS). Acceptance tests should `t.Skip` when environment is unavailable.

**Fix**: Change `t.Fatal` to `t.Skip`:
```go
if akID == "" || akSecret == "" {
    t.Skip("ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
}
```

---

## BUG-006 — Hardcoded default credentials in test helper (Test Infrastructure)

- **Severity**: Medium
- **File**: `zstack/provider/provider_test.go:98-101`
- **Category**: Test Infrastructure / Security

```go
return client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(
    getEnvOrDefault("ZSTACK_ACCOUNT_NAME", "admin"),
    getEnvOrDefault("ZSTACK_ACCOUNT_PASSWORD", "password"),
).ReadOnly(true).Debug(false))
```

**Problem**: Falls back to `admin/password` when env vars are missing, which may cause unintended authentication attempts against real environments. Tests should fail or skip explicitly when credentials are not provided.

**Fix**: Remove default credentials; require explicit env vars or skip.

---

## BUG-007 — Unchecked error returns (Error Handling)

- **Severity**: High
- **File**: Multiple files
- **Category**: Error Handling / errcheck
- **Detected by**: golangci-lint (errcheck)

| File | Line | Unchecked Call |
|------|------|----------------|
| `resource_antipattern_test.go` | 62 | `defer f.Close()` |
| `resource_antipattern_test.go` | 233 | `defer f.Close()` |
| `resource_zstack_virtual_router_instance_test.go` | 30 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 46 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 64 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 72 | `w.Write([]byte(\`{}\`))` |

**Problem**: Return values of `f.Close()`, `json.Encode()`, and `w.Write()` are silently discarded. While often benign in tests, `f.Close()` errors can indicate data loss and `Encode` errors can cause silent test failures.

**Fix**: For `f.Close()` in production-like code, check error. For test mock handlers, either ignore with `_ =` to make intent explicit, or log the error.

---

## BUG-008 — Redundant embedded field selector (Error Handling)

- **Severity**: Low
- **File**: `zstack/provider/resource_zstack_networking_secgroup.go:167`
- **Category**: Code Quality / QF1008
- **Detected by**: golangci-lint (staticcheck QF1008)

```go
p.BaseParam.SystemTags = []string{"SdnControllerUuid::" + u}
```

**Problem**: `BaseParam` is an embedded field — the selector `p.BaseParam.SystemTags` can be simplified to `p.SystemTags`.

**Fix**: `p.SystemTags = []string{"SdnControllerUuid::" + u}`

---

## BUG-009 — Could use tagged switch (Code Quality)

- **Severity**: Low
- **File**: `zstack/utils/filter.go:65`
- **Category**: Code Quality / QF1003
- **Detected by**: golangci-lint (staticcheck QF1003)

```go
if key == "memory_size" {
    fieldValue = fmt.Sprintf("%d", BytesToMB(field.Int()))
} else if key == "disk_size" {
    ...
```

**Problem**: Chain of `if key == "..."` should be a `switch key {` for readability.

**Fix**: Refactor to `switch key { case "memory_size": ... case "disk_size": ... }`.

---

## BUG-010 — Unused function `copyStringValues` (Dead Code)

- **Severity**: Medium
- **File**: `zstack/provider/state_helpers.go:52`
- **Category**: Dead Code
- **Detected by**: golangci-lint (unused)

```go
func copyStringValues(values []types.String) []types.String {
    if len(values) == 0 {
        return nil
    }
    copied := make([]types.String, len(values))
    copy(copied, values)
    return copied
}
```

**Problem**: Function is defined but never called anywhere in the codebase.

**Fix**: Delete the function or add a usage. If it was intended for a future feature, convert to a TODO with issue reference.

---

## BUG-011 — Unused function `envInt` (Dead Code)

- **Severity**: Medium
- **File**: `zstack/provider/test_env_loader_test.go:175`
- **Category**: Dead Code
- **Detected by**: golangci-lint (unused)

```go
func envInt(m map[string]interface{}, key string) int {
```

**Problem**: Function is defined but never called.

**Fix**: Delete or use.

---

## BUG-012 — Misspelled file name `virtural` → `virtual` (Spelling)

- **Severity**: High
- **File**: `zstack/provider/data_source_zstack_virtural_router_images.go`
- **Category**: Spelling / File Name

**Problem**: File name contains `virtural` (typo for `virtual`). This affects maintainability, grep-ability, and developer confidence.

**Fix**: Rename to `data_source_zstack_virtual_router_images.go`. Update any references in `provider.go` registration or imports. Also rename the corresponding test file if one exists.

---

## BUG-013 — Misspelled struct field `BackupStorges` → `BackupStorages` (Spelling)

- **Severity**: High
- **File**: `zstack/provider/data_source_zstack_backup_storages.go:40,144`
- **Category**: Spelling / Identifier

```go
BackupStorges []backupStorage `tfsdk:"backup_storages"`
// ...
state.BackupStorges = append(state.BackupStorges, backupStorageState)
```

**Problem**: Go field name `BackupStorges` is misspelled (missing "a"). The `tfsdk` tag is correct (`backup_storages`), so Terraform users won't see the typo, but it affects code readability and internal consistency.

**Fix**: Rename `BackupStorges` → `BackupStorages` across both occurrences.

---

## BUG-014 — Misspelled type name `gpuDeviceTyp` → `gpuDeviceType` (Spelling)

- **Severity**: Medium
- **File**: `zstack/provider/resource_zstack_instance.go:30,50-51`
- **Category**: Spelling / Identifier

```go
type gpuDeviceTyp string
// ...
mdevDevice gpuDeviceTyp = "mdevDevice"
pciDevice  gpuDeviceTyp = "pciDevice"
```

**Problem**: Type name `gpuDeviceTyp` is missing the trailing "e" — should be `gpuDeviceType`.

**Fix**: Rename to `gpuDeviceType` (3 occurrences in the same file).

---

## BUG-015 — Truncated schema description `"ype of the cluster"` (Spelling)

- **Severity**: Medium
- **File**: `zstack/provider/data_source_zstack_clusters.go:110`
- **Category**: Spelling / User-facing

```go
Description: "ype of the cluster",
```

**Problem**: Missing leading "T" — should be `"Type of the cluster"`. This is **user-facing** in Terraform registry documentation and CLI output.

**Fix**: Change to `"Type of the cluster"`.

---

## BUG-016 — Inconsistent factory/type naming in provider registration (Naming)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: `zstack/provider/provider.go:221-265`
- **Category**: Naming Consistency

```go
// Examples of inconsistent capitalization:
ZStackvmsDataSource          // lowercase "vms"
ZStackl3NetworkDataSource    // lowercase "l3"
ZStackmnNodeDataSource       // lowercase "mn"
ZStackvrouterDataSource      // "vrouter" not "VirtualRouter"
```

Compare with properly capitalized names:
```go
ZStackVirtualRouterImageDataSource  // correct
ZStackAutoScalingGroupDataSource    // correct
```

**Problem**: Exported factory function names had inconsistent capitalization of acronyms and abbreviations.

**Current fix**: Standardized to consistent PascalCase in both declarations and provider registration: `ZStackVMsDataSource`, `ZStackL3NetworkDataSource`, `ZStackMNNodeDataSource`, `ZStackVRouterDataSource`.

---

## BUG-017 — Type/factory name mismatch: `vmResource` vs `InstanceResource()` (Naming)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Low
- **File**: `zstack/provider/resource_zstack_instance.go:32,73`
- **Category**: Naming Consistency

```go
type vmResource struct { ... }           // internal type name
func InstanceResource() resource.Resource { return &vmResource{} }  // factory name
```

**Problem**: Factory returned `&vmResource{}` but was called `InstanceResource()`. Other resources are consistent (e.g., `imageResource` → `ImageResource()`). This made it harder to navigate the codebase.

**Current fix**: Renamed `vmResource` → `instanceResource` so the internal type now matches the exported factory naming pattern.

---

## BUG-018 — Mixed acronym casing: UUID vs Uuid, IP vs Ip (Naming)

- **Severity**: Low
- **File**: Multiple files across `zstack/provider/`
- **Category**: Naming Consistency

```go
// Examples found in various files:
L3NetworkUUID types.String  // "UUID" all-caps
PeerL3NetworkUuids types.String  // "Uuids" mixed case
Iprange []ipRangeModel  // "Ip" not "IP"
AttachedClusterUuids  // "Uuids" not "UUIDs"
```

**Problem**: Go convention recommends consistent acronym casing (all-caps `UUID`, `IP`, `VM`, or all-lower in unexported). The codebase mixes both styles.

**Fix**: Decide on convention (recommended: `UUID`, `IP`, `VM` for Go fields) and apply repo-wide. Note: this is a large refactor — lower priority than functional bugs.

---

## BUG-019 — TODO: Delete re-queries VM instance (Tech Debt)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: High
- **File**: `zstack/provider/resource_zstack_instance.go:1053`
- **Category**: TODO / Design Debt

```go
//TODO: query vm instance again in delete function is not smart.
// Update vm instance's data disk state in read function is a better way
```

**Problem**: The Delete function re-queried the VM to get data disk state. This was acknowledged as suboptimal — data disk state should be populated in Read() so Delete() can use the state directly.

**Current fix**: `data_disks` state now persists the data volume UUIDs during Create/Read, and Delete uses those persisted state values instead of issuing a fresh `GetVmInstance()` call.

---

## BUG-020 — TODO: modify mapping tools (Tech Debt)

- **Severity**: Medium
- **File**: `zstack/provider/data_source_zstack_hook_scripts.go:44`
- **Category**: TODO / Tech Debt

```go
// todo modify mapping tools
```

**Problem**: Vague TODO without owner, issue link, or action plan.

**Fix**: Either implement the mapping change, or convert to `// TODO(owner): description — see issue #N`.

---

## BUG-021 — Antipattern test uses relative glob (Test Infrastructure)

- **Severity**: Medium
- **File**: `zstack/provider/resource_antipattern_test.go:27-35`
- **Category**: Test Infrastructure

```go
files, err := filepath.Glob("resource_zstack_*.go")
if err != nil {
    t.Fatalf("failed to glob resource files: %v", err)
}
if len(files) == 0 {
    t.Fatal("no resource files found — test may be running from wrong directory")
}
```

**Problem**: Relative `filepath.Glob` depends on the working directory. Running `go test` from the module root vs package dir produces different results. The test `t.Fatal`s instead of skipping.

**Fix**: Use `runtime.Caller(0)` to compute absolute path (like `test_env_loader_test.go` does):
```go
_, filename, _, _ := runtime.Caller(0)
dir := filepath.Dir(filename)
files, err := filepath.Glob(filepath.Join(dir, "resource_zstack_*.go"))
```

---

## BUG-022 — Antipattern test uses naive string scanning (Test Infrastructure)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: `zstack/provider/resource_antipattern_test.go`
- **Category**: Test Infrastructure

**Problem**: The antipattern scanner searched for code patterns using simple substring matching (`strings.Contains`) on raw source text. This approach:
- Produces **false positives** from comments, string literals, or multiline formatting
- Produces **false negatives** when code is formatted differently than expected

**Current fix**: Replaced string scanning with AST-based scanning using `go/ast` and `go/parser`, plus regression tests covering commented code vs real violations.

---

## BUG-023 — Hardcoded test resource names (Test Infrastructure)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: Multiple acceptance test files (e.g., `resource_zstack_volume_test.go:48`)
- **Category**: Test Infrastructure

```go
name = "acc-test-volume"
```

**Problem**: Fixed resource names can cause collisions when tests run concurrently or in shared environments.

**Current fix**: Added shared `testAccName()` helper in `provider_test.go` and switched the affected acceptance tests to randomized `acc-test-<name>-<suffix>` resource names.

Example helper pattern:
```go
name := fmt.Sprintf("acc-test-volume-%s", acctest.RandString(8))
```

---

## BUG-024 — Acceptance tests lack update steps (Test Infrastructure)

- **Severity**: Low
- **File**: All acceptance tests
- **Category**: Test Coverage Gap

**Problem**: All existing acceptance tests only exercise Create → (optional Import) → Destroy. None test attribute updates. If Update logic has bugs, they go undetected.

**Fix**: Add TestStep with modified attributes to acceptance tests (start with high-value resources: instance, disk_offering, image).

---

## BUG-025 — Commented-out code blocks across many files (Dead Code)

- **Severity**: Low
- **File**: 19+ files in `zstack/provider/`
- **Category**: Dead Code

**Files with significant `/* ... */` commented blocks**:
- `data_source_zstack_instances.go`
- `data_source_zstack_tags.go`
- `data_source_zstack_backup_storages.go`
- `data_source_zstack_disk_offers.go`
- `data_source_zstack_zone.go`
- `data_source_zstack_virtural_router_images.go`
- `data_source_zstack_l3networks.go`
- `data_source_zstack_hosts.go`
- `data_source_zstack_images.go`
- `data_source_zstack_l2networks.go`
- `data_source_zstack_vips.go`
- `data_source_zstack_clusters.go`
- `resource_zstack_image.go` (commented `removeStringFromSlice` function)
- `resource_zstack_disk_offering.go`
- `provider.go`

**Problem**: Large commented code blocks add noise, confuse diffs, and suggest incomplete cleanups.

**Fix**: Delete stale blocks. If code is needed later, it's in git history. For planned-but-deferred work, replace with a TODO linking to an issue.

---

## Static Analysis Summary (golangci-lint)

Full output from `golangci-lint run ./...`:

| # | Rule | File | Line | Message |
|---|------|------|------|---------|
| 1 | errcheck | resource_antipattern_test.go | 62 | `f.Close()` unchecked |
| 2 | errcheck | resource_antipattern_test.go | 233 | `f.Close()` unchecked |
| 3 | errcheck | resource_zstack_virtual_router_instance_test.go | 30 | `Encode()` unchecked |
| 4 | errcheck | resource_zstack_virtual_router_instance_test.go | 46 | `Encode()` unchecked |
| 5 | errcheck | resource_zstack_virtual_router_instance_test.go | 64 | `Encode()` unchecked |
| 6 | errcheck | resource_zstack_virtual_router_instance_test.go | 72 | `w.Write()` unchecked |
| 7 | SA9003 | resource_zstack_instance.go | 1143 | Empty branch |
| 8 | QF1008 | resource_zstack_networking_secgroup.go | 167 | Redundant embedded field selector |
| 9 | **SA4010** | **resource_zstack_tag_attachment.go** | **207** | **Append result never used** |
| 10 | SA1019 | utils/filter.go | 38 | `strings.Title` deprecated |
| 11 | QF1003 | utils/filter.go | 65 | Could use tagged switch |
| 12 | unused | state_helpers.go | 52 | `copyStringValues` never used |
| 13 | unused | test_env_loader_test.go | 175 | `envInt` never used |

**`go vet ./...`**: Clean (0 issues)  
**`go build ./...`**: Clean (0 errors)  
**`go test -short`**: 53 PASS, 1 FAIL (`TestQueryEnvironment`), 76 SKIP (acceptance tests)

---

## Test Results Summary

| Status | Count | Notes |
|--------|-------|-------|
| PASS | 53 | All schema/metadata unit tests + 1 mock test + 3 antipattern tests |
| FAIL | 1 | `TestQueryEnvironment` — uses `t.Fatal` for missing env vars (BUG-005) |
| SKIP | 76 | Acceptance tests — skipped without `TF_ACC` or env data |
| **Total** | **130** | |

---

## Priority Fix Order

### Immediate (P0) — Fix before next release

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-001 | `append` result discarded in tag_attachment Delete — likely functional bug | 30 min |
| BUG-004 | Test helper `panic()` crashes test binary | 15 min |
| BUG-005 | `TestQueryEnvironment` fails CI with `t.Fatal` | 5 min |
| BUG-003 | `strings.Title` deprecated — may break on future Go versions | 30 min |

### Short-term (P1) — Fix within 1 sprint

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-012 | Rename `virtural` → `virtual` in filename | 15 min |
| BUG-013 | Fix `BackupStorges` → `BackupStorages` | 10 min |
| BUG-014 | Fix `gpuDeviceTyp` → `gpuDeviceType` | 10 min |
| BUG-015 | Fix truncated description `"ype of the cluster"` | 5 min |
| BUG-007 | Unchecked error returns (6 occurrences) | 30 min |
| BUG-010 | Remove unused `copyStringValues` | 5 min |
| BUG-011 | Remove unused `envInt` | 5 min |
| BUG-021 | Fix antipattern test relative glob | 15 min |

### Medium-term (P2) — Technical debt cleanup

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-002 | Remove empty `if` branch in instance resource | 5 min |
| BUG-006 | Remove hardcoded default credentials in tests | 15 min |
| BUG-008 | Simplify embedded field selector | 5 min |
| BUG-009 | Refactor if-chain to switch | 10 min |
| BUG-016 | Standardize factory naming in provider.go | 1 hr |
| BUG-019 | Implement TODO: move disk state to Read() | 2-3 hrs |
| BUG-020 | Resolve or track TODO in hook_scripts | 15 min |
| BUG-022 | Replace naive string scanning with AST | 2-3 hrs |
| BUG-023 | Randomize test resource names | 1 hr |
| BUG-025 | Clean up commented-out code blocks (19+ files) | 2-3 hrs |

### Low Priority (P3) — Style consistency

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-017 | Align vmResource/InstanceResource naming | 30 min |
| BUG-018 | Standardize acronym casing (UUID/Uuid/IP/Ip) | 3-4 hrs |
| BUG-024 | Add update steps to acceptance tests | Ongoing |
