# QA 2026-04-22 P0 Quick Wins + Story-15 (IsNull HIGH risk)

## TL;DR

> **Quick Summary**: On a NEW branch off latest `origin/master` (decoupled from MR #28), ship two batches in one MR: (1) P0 = BUG-001 tag_attachment Delete fix + minimal `go/ast` Layer-2 scanner harness with checks 2d/2b/2a + drop QA shell scan scripts + tracker/docs update; (2) Story-15 = 12 IsNull/IsUnknown guard sites across 7 resource files (1 of which is an ordering bug requiring a reorder, not just a guard).
>
> **Deliverables**:
> - New branch `fix/qa-20260422-p0-plus-story15` off `origin/master` (HEAD `6ed8cdef`)
> - 1 BUG-001 Delete fix + spy/mock test (`resource_zstack_tag_attachment.go`)
> - 12 IsNull/IsUnknown guard fixes across 7 files + per-resource Unknown-value unit tests
> - 1 ordering-bug fix (`resource_zstack_instance_scripts_execution.go:157-160`)
> - New `go/ast` AST scanner harness `resource_antipattern_ast_test.go` + checks 2d → 2b → 2a + fixtures under `zstack/provider/testdata/antipatterns/`
> - 2 QA shell scripts dropped at `scripts/qa/scan_isnull_gap.sh`, `scripts/qa/scan_update_no_reread.sh`
> - Local `_bmad-output/bug-tracker.md` + `_bmad-output/test-status-overview.md` updates (force-added, gitignored)
> - Independent GitLab MR to `master`
>
> **Estimated Effort**: Large (4 waves, ~21 tasks + final verification wave)
> **Parallel Execution**: YES — 4 implementation waves + 1 verification wave
> **Critical Path**: T1 (AST harness skeleton) → T13/T14/T15 (scanner checks) → T19 (repo-sweep smoke) → F1-F4 → user okay

---

## Context

### Original Request
> "下一步我们还是创建分支修复bug http://dev.zstack.io:9080/zstackio/terraform-provider-zstack，继续分析qa结论。"
> Plus local QA report `/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md` (2026-04-22 architecture analysis).

### Interview Summary
**Key Decisions** (all 8 decisions locked, no open questions):
- **Scope**: P0 + Story-15 only — NOT Story-14, NOT Story-16
- **Branch**: NEW `fix/qa-20260422-p0-plus-story15` off latest `origin/master` (HEAD `6ed8cdef`); independent MR; **fully decoupled from MR #28**
- **Test strategy**: TDD (RED → GREEN → REFACTOR) + Agent QA (offline-only — `go build ./...`, targeted unit tests with Unknown-value injection, scanner fixture tests, `go test ./... -short`)
- **Pre-research**: Yes — parallel explore + oracle, results synthesized
- **P0-C scope**: Include — drop both QA shell scan scripts into `scripts/qa/`
- **`instance_scripts_execution`**: Include as Story-15e dedicated task (7 resources total, 8 fix tasks)
- **Unknown-value semantics**: **Per-field decision** — each resource's executor decides whether Unknown means omit, default, or warn-then-omit, based on existing IsNull behavior at that site. Documented per task.
- **Tracker sync**: Local tracker + docs only on this branch — NO QA cross-sync (separate follow-up after MR merges)

**Research Findings** (from explore + oracle agents — synthesized in draft):
- BUG-022 AST scanner is NOT in `master` — must build minimal `go/ast` harness on this branch
- Existing `resource_antipattern_test.go` is string/line-based — do NOT extend it; create sibling AST file
- Oracle: ship checks 2a/2b/2d in this branch; defer 2c (Schema↔Read drift) — would misfire on `resource_zstack_license.go` and `resource_zstack_host.go`
- Oracle order: 2d (lowest risk, hardens harness) → 2b → 2a
- BUG-001 real root cause: `tag_attachment.go:205` calls `DetachTagFromResources` without passing `state.ResourceUuids` (NOT "append result discarded" as peer tracker described)
- Story-15 verified hotspots: 12 sites across 7 files (table in Wave 2 tasks)
- `instance_scripts_execution.go:157-160` is an **ordering bug** (calls `ValueInt64()` BEFORE `IsNull()` check) — needs reorder + guard, not just guard
- `l2vlan_network` originally listed but explore found NO Int64/Bool gap — dropped from scope

### Metis Review (gaps caught and resolved before plan generation)
- **Resolved by user decisions**: P0-C scope, Story-15e classification, Unknown-value policy, tracker/docs scope
- **Locked into guardrails**: branch isolation, no `tools/tools.go` / `examples/` / `go generate` interaction, no Story-14/16 creep, no scanner check 2c
- **Locked into acceptance criteria**: Unknown-value injection unit tests required (compile/test green alone is insufficient), spy/mock client test for BUG-001, fixture-based positive + repo-based negative scanner expectations
- **Wave dependency fix**: Wave 3 scanner repo-sweeps wait for Wave 2 to complete; fixture self-tests can run in parallel with Wave 2
- **Type scope locked**: scanner 2b matches Int64/Bool scalar fields (incl. nested in structs); explicitly excludes List/Set/Map/Object containers and String/Float64

---

## Work Objectives

### Core Objective
Ship a single, independent GitLab MR to `master` that fixes BUG-001 + 12 Story-15 IsNull/IsUnknown sites + adds a minimal Layer-2 AST scanner with 3 checks, all guarded by TDD + agent-executed offline QA, with zero coupling to MR #28.

### Concrete Deliverables
- `zstack/provider/resource_zstack_tag_attachment.go` — Delete passes `state.ResourceUuids`
- `zstack/provider/resource_zstack_tag_attachment_test.go` — spy/mock client test asserting `DetachTagFromResources` receives expected UUIDs
- `zstack/provider/resource_zstack_instance.go` — IsUnknown guards at lines 628, 664-666, 708-714
- `zstack/provider/resource_zstack_port_forwarding_rule.go` — IsUnknown guards at lines 217-225
- `zstack/provider/resource_zstack_volume.go` — IsUnknown guards at lines 217, 344-345
- `zstack/provider/resource_zstack_instance_scripts.go` — IsUnknown guards at lines 164-167, 310-315
- `zstack/provider/resource_zstack_instance_scripts_execution.go` — REORDER lines 157-160 + IsUnknown guard
- `zstack/provider/resource_zstack_l3network.go` — IsUnknown guards at lines 185-187, 190-192
- `zstack/provider/resource_zstack_vip_qos.go` — IsUnknown guards at lines 119-130
- 7 corresponding `*_test.go` files with Unknown-value injection unit tests (one per resource)
- `zstack/provider/resource_antipattern_ast_test.go` — new minimal `go/ast` scanner harness (test-only)
- `zstack/provider/testdata/antipatterns/check_2a/{bad,good}/*.go` — fixtures
- `zstack/provider/testdata/antipatterns/check_2b/{bad,good}/*.go` — fixtures
- `zstack/provider/testdata/antipatterns/check_2d/{bad,good}/*.go` — fixtures
- `scripts/qa/scan_isnull_gap.sh` — QA report's first scan script
- `scripts/qa/scan_update_no_reread.sh` — QA report's second scan script
- `_bmad-output/bug-tracker.md` — entry added for BUG-001 (next free local ID) + Story-15 references
- `_bmad-output/test-status-overview.md` — updated with new tests + scanner status
- GitLab MR: `fix/qa-20260422-p0-plus-story15` → `master`

### Definition of Done
- [ ] `go build ./...` exits 0
- [ ] `go test ./... -short` exits 0 (all unit tests pass, including new Unknown-value tests + scanner fixture tests)
- [ ] `go test ./zstack/provider/ -run TestAntipatternScanner -short` exits 0 (scanner self-tests pass on fixtures + repo sweep returns zero findings post-fix)
- [ ] All 12 Story-15 sites have a unit test that constructs an Unknown value and asserts no panic + correct fallback per per-field policy
- [ ] BUG-001 has a spy/mock test asserting `DetachTagFromResources` receives `state.ResourceUuids`
- [ ] Branch pushed to `origin/fix/qa-20260422-p0-plus-story15`
- [ ] GitLab MR opened with `glab mr create` to target `master`
- [ ] All 4 final-verification reviewers (F1-F4) report APPROVE
- [ ] User explicitly says okay before marking work complete

### Must Have
- TDD discipline: every fix lands as failing test first, then code, then green test
- Per-field Unknown-value policy DOCUMENTED in each fix task before implementation (which path: omit / default / warn+omit)
- Branch isolation: fresh worktree from `origin/master` only; no cherry-pick / copy / reuse from MR #28's branch or worktrees
- Scanner harness uses `go/parser` + `go/ast` (NOT `ast-grep`, NOT string matching)
- AST scanner is test-only Go (`*_test.go` file in `zstack/provider/` package)
- All git commands prefixed with `GIT_MASTER=1`
- 3+ files changed → at least 2 commits (per `git-master` skill)
- `_bmad-output/` files added with `git add -f` (gitignored)
- All subagent invocations include `load_skills=[]` and `run_in_background`
- Independent MR — does NOT depend on MR #28 being merged first

### Must NOT Have (Guardrails)
- ❌ NO touching `fix/bug-tracker-open-items` branch or `.worktrees/bug-tracker-fixes`
- ❌ NO cherry-picking commits or copying files from MR #28
- ❌ NO modifying `tools/tools.go`, `examples/`, generated `docs/`, or running `go generate`
- ❌ NO extending the existing `zstack/provider/resource_antipattern_test.go` (string-based) — create sibling AST file
- ❌ NO scanner check 2c (Schema↔Read drift) — deferred per oracle
- ❌ NO Story-14 work (56 batch Update read-after-write)
- ❌ NO Story-16 work (Read field alignment beyond what scanner immediately surfaces)
- ❌ NO refactoring beyond minimal IsNull/IsUnknown guards + reorder for the ordering bug
- ❌ NO new resources, no new data sources
- ❌ NO live ZStack acceptance tests required (`TF_ACC=1` is OPTIONAL/non-blocking)
- ❌ NO QA tracker cross-sync on this branch (separate follow-up)
- ❌ NO scanner 2b false-positive types: List/Set/Map/Object/String/Float64 are EXPLICITLY out of scope
- ❌ NO touching `resource_zstack_license.go` or `resource_zstack_host.go` for scanner work (false-positive sources for deferred check 2c)
- ❌ NO l2vlan_network changes (no gap found)

---

## Verification Strategy (MANDATORY)

> **ZERO HUMAN INTERVENTION** — ALL verification is agent-executed. No "user manually tests" criteria.
> Evidence saved to `.sisyphus/evidence/task-{N}-{scenario-slug}.{ext}`.

### Test Decision
- **Infrastructure exists**: YES (Go test framework, `go test ./...`, existing `*_test.go` files)
- **Automated tests**: YES — TDD (RED → GREEN → REFACTOR) per fix task
- **Framework**: Go stdlib `testing` + `terraform-plugin-framework/types` + `go/parser` + `go/ast` for scanner
- **Live ZStack**: NOT required — all verification is offline. `TF_ACC=1` is optional/non-blocking.

### QA Policy
Every task includes agent-executed QA scenarios using:
- **Bash + go test**: `go test ./zstack/provider/ -run TestX -short -v` for unit tests + scanner fixture tests + Unknown-value injection
- **Bash + spy/mock**: BUG-001 uses an in-test fake client capturing `DetachTagFromResources` arguments
- **Bash + go build**: `go build ./...` for compile verification
- **Bash + go vet**: `go vet ./...` for static checks
- **Bash + scanner harness**: `go test ./zstack/provider/ -run TestAntipatternScanner -short -v` for scanner self-tests on fixtures + repo sweep

### Per-Wave Verification Gate
- **Wave 1 → Wave 2**: harness skeleton compiles + empty fixture run passes
- **Wave 2 → Wave 3 repo sweep**: all Wave 2 unit tests green + `go build ./...` clean
- **Wave 3 fixture work**: can run in parallel with Wave 2 (no dependency)
- **Wave 3 repo-sweep tasks (T19)**: BLOCKED until Wave 2 complete (otherwise sweeps will see pre-fix offenders)
- **Wave 4**: ALL prior waves complete + tracker entries added + docs updated

---

## Execution Strategy

### Parallel Execution Waves

> Maximize throughput. Each wave completes before the next begins (except scanner fixture work which overlaps Wave 2).
> Target: 5-8 tasks per wave. Final verification = 4 parallel reviewers.

```
Wave 1 (Foundation — start immediately, all parallel):
├── T1:  AST scanner harness skeleton (go/parser + go/ast) [quick]
├── T2:  testdata/antipatterns/ directory + README [quick]
├── T3:  Local tracker entry for BUG-001 [quick]
└── T4:  Drop QA shell scripts into scripts/qa/ (2 files) [quick]

Wave 2 (TDD fixes — start after T1+T2 done; all parallel):
├── T5:  BUG-001 tag_attachment.Delete fix + spy/mock test [unspecified-high]
├── T6:  Story-15a instance.go (3 sites) [unspecified-high]
├── T7:  Story-15b port_forwarding_rule.go (3 sites) [quick]
├── T8:  Story-15c volume.go (2 sites) [quick]
├── T9:  Story-15d instance_scripts.go (2 sites) [quick]
├── T10: Story-15e instance_scripts_execution.go REORDER + guard [unspecified-high]
├── T11: Story-15f l3network.go (2 sites) [quick]
└── T12: Story-15g vip_qos.go (3 sites) [quick]

Wave 3 (Scanner checks — fixtures can overlap Wave 2; repo sweeps wait):
├── T13: Check 2d (append result discarded) — fixtures + scanner code [deep]
├── T14: Check 2b (IsNull w/o IsUnknown Int64/Bool) — fixtures + scanner code [deep]
├── T15: Check 2a (Update missing read-after-write) — fixtures + scanner code [deep]
├── T16: Check 2d fixture self-tests pass [quick]
├── T17: Check 2b fixture self-tests pass [quick]
└── T18: Check 2a fixture self-tests pass [quick]

Wave 3.5 (Repo sweeps — BLOCKED until Wave 2 complete):
└── T19: Run all 3 scanners against repo, assert ZERO findings post-fix [unspecified-high]

Wave 4 (Integration & docs):
├── T20: Update _bmad-output/bug-tracker.md with all fixes [quick]
├── T21: Update _bmad-output/test-status-overview.md [quick]
└── T22: Push branch + open MR via glab [quick]

Final Verification Wave (4 parallel reviewers — ALL must approve, then user okay):
├── F1: Plan compliance audit [oracle]
├── F2: Code quality review [unspecified-high]
├── F3: Real manual QA execution [unspecified-high]
└── F4: Scope fidelity check [deep]
→ Present consolidated results → wait for user's explicit "okay" → done

Critical Path: T1 → T13/T14/T15 → T19 → F1-F4 → user okay
Parallel Speedup: ~70% vs sequential (max 8 concurrent in Wave 2)
Max Concurrent: 8 (Wave 2)
```

### Dependency Matrix

| Task | Depends On | Blocks |
|------|-----------|--------|
| T1   | —         | T13, T14, T15, T16, T17, T18 |
| T2   | —         | T13, T14, T15 |
| T3   | —         | T20 |
| T4   | —         | T20 |
| T5   | T1 (harness compiles) | T19, T20 |
| T6-T12 | T1 (harness compiles) | T19, T20 |
| T13  | T1, T2    | T16, T19 |
| T14  | T1, T2    | T17, T19 |
| T15  | T1, T2    | T18, T19 |
| T16-T18 | T13/T14/T15 respectively | T19 |
| T19  | T5-T18 ALL | T20, T21 |
| T20  | T3, T5-T19 | T22 |
| T21  | T19, T20  | T22 |
| T22  | T20, T21  | F1-F4 |
| F1-F4 | T22      | user okay |

### Agent Dispatch Summary

| Wave | Tasks | Agents |
|------|-------|--------|
| 1    | 4     | T1 → quick, T2 → quick, T3 → quick, T4 → quick |
| 2    | 8     | T5 → unspecified-high, T6 → unspecified-high, T7-T9 → quick, T10 → unspecified-high, T11-T12 → quick |
| 3    | 6     | T13-T15 → deep, T16-T18 → quick |
| 3.5  | 1     | T19 → unspecified-high |
| 4    | 3     | T20-T22 → quick |
| FINAL | 4    | F1 → oracle, F2 → unspecified-high, F3 → unspecified-high, F4 → deep |

---

## TODOs

> Implementation + Test = ONE Task (TDD). Never separate.
> EVERY task MUST have: Recommended Agent Profile + Parallelization info + QA Scenarios.
> A task WITHOUT QA Scenarios is INCOMPLETE.

- [x] 1. **AST Scanner Harness Skeleton** (Wave 1)

  **What to do**:
  - Create new test-only Go file `zstack/provider/resource_antipattern_ast_test.go`
  - Use `go/parser`, `go/ast`, `go/token` from stdlib (NOT `ast-grep`, NOT external deps)
  - Implement reusable scanning primitives:
    - `parsePackage(dir string) ([]*ast.File, *token.FileSet, error)` — parse all `.go` files in a directory
    - `walkFunctions(file *ast.File, visit func(*ast.FuncDecl))` — iterate function decls
    - `findCallExpr(node ast.Node, pkg, name string) []*ast.CallExpr` — locate specific calls
    - `Finding struct { File string; Line int; Check string; Message string }` — uniform finding type
    - Top-level `TestAntipatternScanner(t *testing.T)` shell that registers sub-tests for each check (initially empty placeholder for 2d/2b/2a)
  - Compile-only at this stage; no checks implemented yet
  - Add a no-op smoke sub-test that confirms harness can parse `zstack/provider/resource_zstack_volume.go` without error

  **Must NOT do**:
  - Do NOT extend `zstack/provider/resource_antipattern_test.go` (string-based; leave untouched)
  - Do NOT add ast-grep, go-ast-tool, or any third-party AST library
  - Do NOT pull code from `fix/bug-tracker-open-items` branch
  - Do NOT add `go/types` (full type-checking) — keep it lightweight, syntactic-only

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Skeleton-only, well-defined Go stdlib usage, no business logic
  - **Skills**: [`test-driven-development`]
    - `test-driven-development`: TDD discipline — write the harness skeleton with at least one passing smoke sub-test before considering it done
  - **Skills Evaluated but Omitted**:
    - `systematic-debugging`: No bug to debug yet
    - `frontend-design`: N/A

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T2, T3, T4)
  - **Blocks**: T13, T14, T15, T16, T17, T18 (all scanner check tasks need the harness)
  - **Blocked By**: None — can start immediately

  **References**:

  **Pattern References**:
  - `zstack/provider/resource_antipattern_test.go` — DO NOT extend, but read to understand the existing string-based scanner's structure (file iteration pattern, error reporting style) and use a CONSISTENT reporting format

  **API/Type References**:
  - Go stdlib `go/parser` — `parser.ParseDir`, `parser.ParseFile` for syntactic parsing
  - Go stdlib `go/ast` — `ast.Inspect`, `ast.FuncDecl`, `ast.CallExpr`, `ast.SelectorExpr`
  - Go stdlib `go/token` — `token.FileSet`, `token.Position` for file:line reporting

  **Test References**:
  - `zstack/provider/resource_zstack_volume_test.go` — existing test file structure for the package

  **External References**:
  - `https://pkg.go.dev/go/parser` — official parser docs
  - `https://pkg.go.dev/go/ast` — official ast docs

  **WHY Each Reference Matters**:
  - The existing string-based scanner shows the team's preferred error reporting format (file:line + message). Match it for consistency, but use real `token.Position` from go/ast (more accurate than string regex line counting).
  - `volume_test.go` shows the test naming convention used in this package — use the same.

  **Acceptance Criteria** (AGENT-EXECUTABLE ONLY):
  - [ ] File `zstack/provider/resource_antipattern_ast_test.go` exists
  - [ ] `go build ./zstack/provider/` exits 0
  - [ ] `go test ./zstack/provider/ -run TestAntipatternScanner/smoke -short -v` exits 0 with at least one PASS line
  - [ ] `go vet ./zstack/provider/` exits 0
  - [ ] `grep -E '"github.com/.*ast"' zstack/provider/resource_antipattern_ast_test.go` returns nothing (no third-party AST libs)

  **QA Scenarios**:

  ```
  Scenario: Harness parses provider package without panic
    Tool: Bash (go test)
    Preconditions: Branch checked out, T1 changes applied
    Steps:
      1. Run `cd /Users/michelvaillant/Developer/zstack.io/terraform-provider-zstack && go test ./zstack/provider/ -run TestAntipatternScanner -short -v 2>&1 | tee .sisyphus/evidence/task-1-harness-smoke.txt`
      2. Inspect output for `--- PASS: TestAntipatternScanner/smoke`
    Expected Result: Exit 0, smoke sub-test PASS
    Failure Indicators: panic, parse error, exit non-zero, no PASS line
    Evidence: .sisyphus/evidence/task-1-harness-smoke.txt

  Scenario: Compile fails cleanly if harness API misused
    Tool: Bash (go build with intentional misuse)
    Preconditions: Harness exists
    Steps:
      1. Create temp file `/tmp/misuse_test.go` that calls scanner with nil arg (won't be committed)
      2. Run `go build` — expect compile error or runtime nil-check
      3. Delete temp file
    Expected Result: Either compile error OR test fails fast with clear message (NOT silent nil deref panic)
    Evidence: .sisyphus/evidence/task-1-harness-error-handling.txt
  ```

  **Evidence to Capture**:
  - [ ] `.sisyphus/evidence/task-1-harness-smoke.txt` — go test output
  - [ ] `.sisyphus/evidence/task-1-harness-error-handling.txt` — error path output

  **Commit**: YES (own commit, T1 only)
  - Message: `test(provider): add AST scanner harness skeleton for antipattern checks`
  - Files: `zstack/provider/resource_antipattern_ast_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner -short`

- [x] 2. **Testdata Fixtures Directory** (Wave 1)

  **What to do**:
  - Create directory tree:
    - `zstack/provider/testdata/antipatterns/check_2a/{bad,good}/` (empty placeholders, README only)
    - `zstack/provider/testdata/antipatterns/check_2b/{bad,good}/` (empty placeholders, README only)
    - `zstack/provider/testdata/antipatterns/check_2d/{bad,good}/` (empty placeholders, README only)
  - Add `zstack/provider/testdata/antipatterns/README.md` explaining:
    - Purpose: fixtures for AST scanner self-tests
    - Layout: each check has `bad/` (must be flagged) and `good/` (must NOT be flagged) subdirs
    - File naming: `*.go.fixture` extension to AVOID being compiled by `go build ./...` (renamed by scanner before parsing in-memory)
  - Create one `.gitkeep` per `bad/` and `good/` directory so git tracks empty dirs
  - DO NOT add actual fixture content here — that's the responsibility of T13/T14/T15

  **Must NOT do**:
  - Do NOT add files with `.go` extension (would be picked up by `go build ./...`)
  - Do NOT add fixture content (defer to scanner check tasks)
  - Do NOT modify `tools/tools.go` or trigger any code generation

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Pure directory + README scaffolding
  - **Skills**: []
    - No skills needed — trivial scaffolding

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T3, T4)
  - **Blocks**: T13, T14, T15
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `zstack/provider/testdata/` — existing testdata layout (env.json, terraform/) for naming conventions
  - `zstack/provider/testdata/generate_env.go` — note: existing testdata uses `.go` extension for code; we deliberately use `.go.fixture` to avoid compilation

  **WHY**:
  - The `.go.fixture` extension is critical: if fixtures use `.go`, then `go build ./...` will try to compile invalid example code and fail. Renaming in-memory in the scanner is the simplest workaround.

  **Acceptance Criteria**:
  - [ ] Directory `zstack/provider/testdata/antipatterns/` exists with subdirs `check_2a/{bad,good}`, `check_2b/{bad,good}`, `check_2d/{bad,good}`
  - [ ] Each leaf dir contains `.gitkeep`
  - [ ] `zstack/provider/testdata/antipatterns/README.md` exists with documented layout
  - [ ] `find zstack/provider/testdata/antipatterns -name '*.go'` returns nothing (no real .go files yet)
  - [ ] `go build ./...` exits 0 (no compilation interference)

  **QA Scenarios**:

  ```
  Scenario: Directory tree present and correct
    Tool: Bash
    Steps:
      1. Run `find zstack/provider/testdata/antipatterns -type d | sort | tee .sisyphus/evidence/task-2-dirs.txt`
      2. Verify output contains all 7 expected dirs (root + 6 leaves)
    Expected Result: Exact match with expected layout
    Evidence: .sisyphus/evidence/task-2-dirs.txt

  Scenario: No accidental .go files
    Tool: Bash
    Steps:
      1. Run `find zstack/provider/testdata/antipatterns -name '*.go' | tee .sisyphus/evidence/task-2-no-go.txt`
    Expected Result: Empty output
    Failure Indicators: Any .go file present
    Evidence: .sisyphus/evidence/task-2-no-go.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-2-dirs.txt`
  - [ ] `.sisyphus/evidence/task-2-no-go.txt`

  **Commit**: YES (own commit)
  - Message: `test(provider): scaffold testdata/antipatterns/ fixture layout`
  - Files: `zstack/provider/testdata/antipatterns/**`
  - Pre-commit: `GIT_MASTER=1 go build ./...`

- [x] 3. **Local Tracker Entry for BUG-001** (Wave 1)

  **What to do**:
  - Read current `_bmad-output/bug-tracker.md` (use `cat` since it's gitignored — `git show HEAD:_bmad-output/bug-tracker.md` may fail)
  - Determine next free local BUG ID by scanning existing entries
  - Add new entry for BUG-001 (peer tracker name: tag_attachment Delete) using next free ID
  - Entry must include: ID, peer-tracker reference (BUG-001), title, file:line (`zstack/provider/resource_zstack_tag_attachment.go:197-211`), root cause (Delete ignores `state.ResourceUuids`), planned fix (pass ResourceUuids as 3rd arg to `DetachTagFromResources`), status (PLANNED → will become FIXED in T5/T20)
  - Use `git add -f _bmad-output/bug-tracker.md` to stage (gitignored)

  **Must NOT do**:
  - Do NOT modify QA tracker in `.worktrees/qa-tracker-repo` (out of scope per user decision)
  - Do NOT touch any other entry's status field
  - Do NOT use `git add` without `-f` (will silently no-op)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Small markdown edit
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T2, T4)
  - **Blocks**: T20 (final tracker sync depends on initial entry)
  - **Blocked By**: None

  **References**:
  - `_bmad-output/bug-tracker.md` — existing entries for ID format and status vocabulary

  **Acceptance Criteria**:
  - [ ] `_bmad-output/bug-tracker.md` contains a new entry referencing peer BUG-001 with status PLANNED
  - [ ] `git status --ignored _bmad-output/` shows the file as modified (or use `git diff --no-index` against a snapshot)
  - [ ] Entry includes file path `zstack/provider/resource_zstack_tag_attachment.go` and line range `197-211`
  - [ ] Entry includes root cause text mentioning `state.ResourceUuids`

  **QA Scenarios**:

  ```
  Scenario: New BUG entry parseable
    Tool: Bash (grep)
    Steps:
      1. Run `grep -A5 'BUG-001' _bmad-output/bug-tracker.md | tee .sisyphus/evidence/task-3-entry.txt`
    Expected Result: Output contains "tag_attachment" AND "ResourceUuids" AND "PLANNED"
    Evidence: .sisyphus/evidence/task-3-entry.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-3-entry.txt`

  **Commit**: NO (will be committed together with T20 final tracker sync to keep tracker churn minimal)

- [x] 4. **Drop QA Shell Scan Scripts** (Wave 1)

  **What to do**:
  - Create directory `scripts/qa/`
  - Extract the two shell scan scripts from the QA report (`/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md`):
    - `scripts/qa/scan_isnull_gap.sh` — scans for `IsNull()` without subsequent `IsUnknown()`
    - `scripts/qa/scan_update_no_reread.sh` — scans for Update methods without read-after-write
  - Both scripts must:
    - Start with `#!/usr/bin/env bash` and `set -euo pipefail`
    - Be `chmod +x` (executable)
    - Print findings as `file:line: <message>` (consistent with Go scanner output)
    - Exit 0 if no findings, exit 1 if findings (so CI can use them)
  - Add `scripts/qa/README.md` explaining purpose, usage, relationship to Go scanner (the Go scanner is the source of truth for CI; shell scripts are for ad-hoc review and reproducibility of the QA report)

  **Must NOT do**:
  - Do NOT add scripts to `tools/tools.go` or `go generate` chain
  - Do NOT make scripts depend on Go (must be pure shell + grep/awk for portability)
  - Do NOT use BSD-only or GNU-only flags exclusively — use POSIX-portable subset where possible

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Mechanical script transcription from spec
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with T1, T2, T3)
  - **Blocks**: T20 (mentioned in tracker)
  - **Blocked By**: None

  **References**:
  - `/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md` — QA report contains the two scripts verbatim
  - `scripts/run_tests.sh` — existing shell script for style/header conventions

  **Acceptance Criteria**:
  - [ ] `scripts/qa/scan_isnull_gap.sh` exists and is executable (`test -x scripts/qa/scan_isnull_gap.sh`)
  - [ ] `scripts/qa/scan_update_no_reread.sh` exists and is executable
  - [ ] `scripts/qa/README.md` exists
  - [ ] `bash -n scripts/qa/scan_isnull_gap.sh` exits 0 (syntax check)
  - [ ] `bash -n scripts/qa/scan_update_no_reread.sh` exits 0
  - [ ] Running `scripts/qa/scan_isnull_gap.sh zstack/provider/` after Wave 2 returns exit 0 (no findings remain)

  **QA Scenarios**:

  ```
  Scenario: Scripts pass shellcheck syntax
    Tool: Bash
    Steps:
      1. Run `bash -n scripts/qa/scan_isnull_gap.sh && bash -n scripts/qa/scan_update_no_reread.sh && echo OK | tee .sisyphus/evidence/task-4-syntax.txt`
    Expected Result: "OK" printed
    Evidence: .sisyphus/evidence/task-4-syntax.txt

  Scenario: Pre-Wave-2 sanity (scripts find known offenders)
    Tool: Bash
    Preconditions: T4 done, Wave 2 NOT yet done (or run on commit before Wave 2)
    Steps:
      1. Run `scripts/qa/scan_isnull_gap.sh zstack/provider/ | tee .sisyphus/evidence/task-4-pre-wave2-findings.txt`
    Expected Result: Non-empty output listing at least the 12 known Story-15 hotspots
    Evidence: .sisyphus/evidence/task-4-pre-wave2-findings.txt
    Note: This scenario is informational; pass criteria is "scripts run without error", not "exact count match" (shell heuristics differ from Go AST)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-4-syntax.txt`
  - [ ] `.sisyphus/evidence/task-4-pre-wave2-findings.txt`

  **Commit**: YES (own commit)
  - Message: `chore(qa): drop QA report scan scripts into scripts/qa/`
  - Files: `scripts/qa/scan_isnull_gap.sh`, `scripts/qa/scan_update_no_reread.sh`, `scripts/qa/README.md`
  - Pre-commit: `bash -n scripts/qa/*.sh`

- [x] 5. **BUG-001 Fix: tag_attachment.Delete passes ResourceUuids** (Wave 2)

  **What to do**:
  - **RED**: Add to `zstack/provider/resource_zstack_tag_attachment_test.go` a unit test `TestTagAttachmentDeletePassesResourceUuids` using a fake/spy client that records `DetachTagFromResources` arguments. Construct a `tagAttachmentResource{client: spy}` and a state with non-empty `ResourceUuids = ["uuid-a","uuid-b"]`. Call `Delete`. Assert spy recorded ResourceUuids exactly `["uuid-a","uuid-b"]`. Test MUST FAIL initially.
  - **VERIFY SDK SIGNATURE FIRST**: Before writing the fix, grep `vendor/` and `go.sum` to confirm `client.DetachTagFromResources` real signature (1st arg = tag UUID, 2nd arg = resource type, 3rd arg = `[]string ResourceUuids`?). Document the confirmed signature in test file as a comment.
  - **GREEN**: Modify `zstack/provider/resource_zstack_tag_attachment.go:197-211` `Delete` so it passes `state.ResourceUuids.ElementsAs(...)` (or equivalent slice extraction) as the 3rd arg to `DetachTagFromResources`. Use the same slice-extraction idiom already in `Create` for symmetry.
  - **REFACTOR**: Only if `Create`/`Delete` slice extraction is duplicated, extract a tiny `extractResourceUuids(ctx, types.List)` helper IN THE SAME FILE (no cross-file refactoring). Otherwise inline.

  **Per-field Unknown semantics for this site**:
  - `state.ResourceUuids` in Delete: state values come from refresh, so Unknown is impossible here. NO IsUnknown guard needed (state values are always Known). Document this in test comment.

  **Must NOT do**:
  - Do NOT modify `Create`, `Read`, or `Update` paths
  - Do NOT change `DetachTagFromResources` signature in vendor/SDK
  - Do NOT add ZStack live API call (must use spy)
  - Do NOT introduce a new mock framework — handcrafted spy struct is sufficient

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Real bug fix; needs careful SDK signature verification + spy construction + diagnostics handling
  - **Skills**: [`test-driven-development`, `systematic-debugging`]
    - `test-driven-development`: Strict RED → GREEN → REFACTOR
    - `systematic-debugging`: To verify SDK signature and existing Create call pattern

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T6-T12)
  - **Parallel Group**: Wave 2
  - **Blocks**: T19 (repo sweep), T20 (tracker entry update)
  - **Blocked By**: T1 (harness compiles — sanity baseline)

  **References**:

  **Pattern References**:
  - `zstack/provider/resource_zstack_tag_attachment.go:Create` — existing slice extraction pattern from `plan.ResourceUuids` (mirror this for Delete)
  - Any existing `_test.go` in `zstack/provider/` that builds a fake client — search for `type fake` or `type spy` in `zstack/provider/*_test.go`

  **API/Type References**:
  - `vendor/.../client.go` (or wherever ZStack SDK client lives) — exact signature of `DetachTagFromResources`
  - `terraform-plugin-framework/types.List.ElementsAs(ctx, &slice, false)` — canonical extraction call

  **Test References**:
  - `zstack/provider/resource_zstack_tag_attachment_test.go` (existing) — extend, do not replace; mirror its setup style

  **WHY**:
  - The peer tracker said "append result discarded" but oracle confirmed real cause is missing 3rd arg. Test must demonstrate the 3rd-arg semantic explicitly so future readers understand the bug.

  **Acceptance Criteria** (AGENT-EXECUTABLE ONLY):
  - [ ] `git log -p zstack/provider/resource_zstack_tag_attachment_test.go` shows new test `TestTagAttachmentDeletePassesResourceUuids` added
  - [ ] On the commit BEFORE the fix, running `go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short` exits NON-ZERO (proves RED)
  - [ ] On the fix commit, running `go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short -v` exits 0 with PASS
  - [ ] `go build ./...` exits 0
  - [ ] `go vet ./zstack/provider/` exits 0
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_tag_attachment.go` shows ONLY the Delete-path change (no Create/Read/Update edits)

  **QA Scenarios**:

  ```
  Scenario: Spy receives ResourceUuids on Delete (happy path)
    Tool: Bash (go test)
    Preconditions: T5 fix + test applied
    Steps:
      1. Run `go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short -v 2>&1 | tee .sisyphus/evidence/task-5-spy-pass.txt`
      2. Confirm output shows "--- PASS: TestTagAttachmentDeletePassesResourceUuids"
    Expected Result: Exit 0, PASS line, spy recorded ["uuid-a","uuid-b"]
    Failure Indicators: FAIL line, spy got nil/empty slice, panic
    Evidence: .sisyphus/evidence/task-5-spy-pass.txt

  Scenario: RED state proven before fix (TDD discipline)
    Tool: Bash (git + go test)
    Preconditions: T5 commit landed; previous commit (test only, no fix) exists
    Steps:
      1. Run `git stash && git checkout HEAD~1 -- zstack/provider/resource_zstack_tag_attachment.go`
      2. Run `go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short -v 2>&1 | tee .sisyphus/evidence/task-5-red-proof.txt`
      3. Restore: `git checkout HEAD -- zstack/provider/resource_zstack_tag_attachment.go && git stash pop`
    Expected Result: Test FAILS (exit non-zero) on pre-fix code
    Evidence: .sisyphus/evidence/task-5-red-proof.txt
    Note: If TDD was done as ONE commit (not RED then GREEN), use `git show HEAD -- zstack/provider/resource_zstack_tag_attachment.go | grep '^-' | grep -v '^---'` to prove the test would have failed against the removed line.
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-5-spy-pass.txt`
  - [ ] `.sisyphus/evidence/task-5-red-proof.txt`

  **Commit**: YES (own commit, T5 only)
  - Message: `fix(tag_attachment): pass state.ResourceUuids to DetachTagFromResources`
  - Files: `zstack/provider/resource_zstack_tag_attachment.go`, `zstack/provider/resource_zstack_tag_attachment_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short`

- [ ] 6. **Story-15a: instance.go IsUnknown guards (3 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_instance.go:628` (NeverStop), `:664-666` (MemorySize), `:708-714` (CPUNum)
  - **Per-field Unknown semantics** (DOCUMENT IN TEST FILE COMMENTS BEFORE FIX):
    - `:628 NeverStop` (Bool): Unknown → **omit** from update payload (fall through, don't send field)
    - `:664-666 MemorySize` (Int64): Unknown → **omit** (memory is provider-computed when scaled; sending Unknown=0 would shrink VM)
    - `:708-714 CPUNum` (Int64): Unknown → **omit** (same reasoning as MemorySize)
  - **RED**: Add to `zstack/provider/resource_zstack_instance_test.go` (or sibling `_unknown_test.go` if file too large) a unit test `TestInstanceUpdateGuardsUnknownValues` with subtests for each of 3 fields. Each subtest constructs a plan with that field = `types.Int64Unknown()` / `types.BoolUnknown()` and asserts the produced API payload OMITS the field. Tests MUST FAIL initially with panic/wrong value.
  - **GREEN**: At each of the 3 sites, change `if !plan.X.IsNull()` → `if !plan.X.IsNull() && !plan.X.IsUnknown()`.
  - **REFACTOR**: None (minimal guard only).

  **Must NOT do**:
  - Do NOT touch any other line in `instance.go`
  - Do NOT change handler signature, return value, or error wrapping
  - Do NOT add Unknown handling to Read or Create paths (Unknown only appears in plan/Update)
  - Do NOT add live ZStack call

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Largest file in batch, 3 sites, payload-shape assertions need careful test construction
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T5, T7-T12)
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:

  **Pattern References**:
  - Existing `if !plan.<Field>.IsNull()` calls in same file (search for pattern in `instance.go`) — match exact style
  - `zstack/provider/resource_zstack_volume.go` Update path — uses similar IsNull pattern at `:217`

  **API/Type References**:
  - `types.Int64Unknown()`, `types.BoolUnknown()` from `github.com/hashicorp/terraform-plugin-framework/types`
  - `(types.Int64).IsUnknown() bool`, `(types.Bool).IsUnknown() bool`

  **Test References**:
  - `zstack/provider/resource_zstack_volume_test.go` — test naming convention reference
  - Existing `instance_test.go` — extend, mirror setup

  **WHY**:
  - Unknown values arise during plan when an upstream resource attribute isn't yet computed. Calling `.ValueInt64()` on Unknown returns 0 → silently shrinks VM. The IsNull-only check misses Unknown.

  **Acceptance Criteria**:
  - [ ] All 3 sites at `:628`, `:664-666`, `:708-714` show `&& !plan.X.IsUnknown()` added
  - [ ] `TestInstanceUpdateGuardsUnknownValues` exists with 3 subtests (one per field)
  - [ ] `go test ./zstack/provider/ -run TestInstanceUpdateGuardsUnknownValues -short -v` exits 0
  - [ ] `go build ./...` exits 0
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_instance.go | grep '^+' | grep -v '^+++' | grep -c IsUnknown` returns exactly 3

  **QA Scenarios**:

  ```
  Scenario: NeverStop Unknown is omitted from payload
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceUpdateGuardsUnknownValues/NeverStop -short -v 2>&1 | tee .sisyphus/evidence/task-6-neverstop.txt`
    Expected Result: PASS, payload assertion confirms field absent
    Evidence: .sisyphus/evidence/task-6-neverstop.txt

  Scenario: MemorySize Unknown is omitted (regression: prevents VM shrink to 0)
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceUpdateGuardsUnknownValues/MemorySize -short -v 2>&1 | tee .sisyphus/evidence/task-6-memorysize.txt`
    Expected Result: PASS; assertion: `assert.NotContains(payload, "memorySize")` AND `assert.NotEqual(payload["memorySize"], int64(0))`
    Failure Indicators: PASS but payload contains memorySize=0 (silent shrink bug)
    Evidence: .sisyphus/evidence/task-6-memorysize.txt

  Scenario: CPUNum Unknown is omitted
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceUpdateGuardsUnknownValues/CPUNum -short -v 2>&1 | tee .sisyphus/evidence/task-6-cpunum.txt`
    Expected Result: PASS
    Evidence: .sisyphus/evidence/task-6-cpunum.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-6-neverstop.txt`
  - [ ] `.sisyphus/evidence/task-6-memorysize.txt`
  - [ ] `.sisyphus/evidence/task-6-cpunum.txt`

  **Commit**: YES (own commit, T6 only)
  - Message: `fix(instance): guard NeverStop/MemorySize/CPUNum against Unknown values`
  - Files: `zstack/provider/resource_zstack_instance.go`, `zstack/provider/resource_zstack_instance_test.go` (or `_unknown_test.go`)
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestInstanceUpdateGuardsUnknownValues -short`

- [ ] 7. **Story-15b: port_forwarding_rule.go IsUnknown guards (3 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_port_forwarding_rule.go:217-225` (3 adjacent fields — confirm exact field names by reading the file: typically `VipPortStart`, `VipPortEnd`, `PrivatePortStart`/`PrivatePortEnd` ranges)
  - **Per-field Unknown semantics** (DOCUMENT IN TEST FILE BEFORE FIX):
    - All 3 port-range Int64 fields: Unknown → **omit** (zero-port = invalid; sending 0 would create invalid rule)
  - **RED**: Add `TestPortForwardingRuleUpdateGuardsUnknownValues` with 3 subtests. Each constructs plan with that port field = `types.Int64Unknown()`, asserts payload omits the field. MUST FAIL initially.
  - **GREEN**: Add `&& !plan.X.IsUnknown()` at each site.
  - **REFACTOR**: None.

  **Must NOT do**:
  - Do NOT touch any line outside :217-225
  - Do NOT modify VIP-related lookups
  - Do NOT add live API call

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 3 adjacent identical-shape fixes, mechanical
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_port_forwarding_rule.go:217-225` — exact site
  - T6 instance.go pattern — same fix shape

  **Acceptance Criteria**:
  - [ ] All 3 sites have `&& !plan.X.IsUnknown()` added
  - [ ] `TestPortForwardingRuleUpdateGuardsUnknownValues` exists with 3 subtests
  - [ ] `go test ./zstack/provider/ -run TestPortForwardingRuleUpdateGuardsUnknownValues -short -v` exits 0
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_port_forwarding_rule.go | grep -c IsUnknown` returns exactly 3

  **QA Scenarios**:

  ```
  Scenario: Each port-range field Unknown is omitted
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestPortForwardingRuleUpdateGuardsUnknownValues -short -v 2>&1 | tee .sisyphus/evidence/task-7-pf-rule.txt`
      2. Confirm 3 PASS lines (one per subtest)
    Expected Result: All 3 PASS
    Evidence: .sisyphus/evidence/task-7-pf-rule.txt

  Scenario: Negative — invalid port 0 NOT sent
    Tool: Bash (go test)
    Steps:
      1. Within same test, assert payload does NOT contain `port: 0`
    Expected Result: assertion holds
    Evidence: .sisyphus/evidence/task-7-pf-rule.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-7-pf-rule.txt`

  **Commit**: YES (own commit)
  - Message: `fix(port_forwarding_rule): guard port-range fields against Unknown values`
  - Files: `zstack/provider/resource_zstack_port_forwarding_rule.go`, `zstack/provider/resource_zstack_port_forwarding_rule_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestPortForwardingRuleUpdateGuardsUnknownValues -short`

- [ ] 8. **Story-15c: volume.go IsUnknown guards (2 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_volume.go:217` (likely `Size` Int64), `:344-345` (likely IOPS or shareable Bool/Int64 — confirm by reading)
  - **Per-field Unknown semantics**:
    - `:217 Size` (Int64): Unknown → **omit** (zero-size resize would destroy data; never send Unknown=0)
    - `:344-345` field (read file to confirm type): if Int64 → **omit**; if Bool → **omit**
  - **RED**: Add `TestVolumeUpdateGuardsUnknownValues` with 2 subtests. MUST FAIL initially.
  - **GREEN**: Add `&& !plan.X.IsUnknown()` at both sites.
  - **REFACTOR**: None.

  **Must NOT do**:
  - Do NOT touch volume snapshot logic
  - Do NOT change resize semantics
  - Do NOT add live API call

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_volume.go:217, 344-345` — exact sites
  - `zstack/provider/resource_zstack_volume_test.go` — extend

  **Acceptance Criteria**:
  - [ ] Both sites have `&& !plan.X.IsUnknown()` added
  - [ ] `TestVolumeUpdateGuardsUnknownValues` exists with 2 subtests, both PASS
  - [ ] `go build ./...` exits 0
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_volume.go | grep -c IsUnknown` returns exactly 2

  **QA Scenarios**:

  ```
  Scenario: Size Unknown is omitted (regression: prevents data-destroying resize-to-0)
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestVolumeUpdateGuardsUnknownValues/Size -short -v 2>&1 | tee .sisyphus/evidence/task-8-volume-size.txt`
      2. Assertion: payload does NOT contain `size: 0`
    Expected Result: PASS; payload omits size key
    Failure Indicators: PASS but payload has size=0 (data-loss bug)
    Evidence: .sisyphus/evidence/task-8-volume-size.txt

  Scenario: Second field Unknown handled
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestVolumeUpdateGuardsUnknownValues -short -v 2>&1 | tee .sisyphus/evidence/task-8-volume-all.txt`
    Expected Result: All subtests PASS
    Evidence: .sisyphus/evidence/task-8-volume-all.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-8-volume-size.txt`
  - [ ] `.sisyphus/evidence/task-8-volume-all.txt`

  **Commit**: YES (own commit)
  - Message: `fix(volume): guard Size and update fields against Unknown values`
  - Files: `zstack/provider/resource_zstack_volume.go`, `zstack/provider/resource_zstack_volume_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestVolumeUpdateGuardsUnknownValues -short`

- [ ] 9. **Story-15d: instance_scripts.go IsUnknown guards (2 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_instance_scripts.go:164-167`, `:310-315` (read file to confirm exact field names — likely `ScriptTimeout` Int64 + `Privilege`/similar)
  - **Per-field Unknown semantics** (DOCUMENT):
    - `:164-167` field (Int64): Unknown → **omit** (script timeout default applied server-side; sending Unknown=0 = immediate timeout)
    - `:310-315` field: Unknown → **omit** (same default-server-side reasoning)
  - **RED**: Add `TestInstanceScriptsUpdateGuardsUnknownValues` with 2 subtests. MUST FAIL initially.
  - **GREEN**: Add `&& !plan.X.IsUnknown()` at both sites.
  - **REFACTOR**: None.

  **Must NOT do**:
  - Do NOT touch script content rendering
  - Do NOT change script storage path
  - Do NOT confuse this resource with `instance_scripts_execution` (T10)

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_instance_scripts.go:164-167, 310-315` — exact sites
  - `zstack/provider/resource_zstack_instance_scripts_test.go` — extend (or create if absent)

  **Acceptance Criteria**:
  - [ ] Both sites have `&& !plan.X.IsUnknown()` added
  - [ ] `TestInstanceScriptsUpdateGuardsUnknownValues` exists with 2 subtests, both PASS
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_instance_scripts.go | grep -c IsUnknown` returns exactly 2

  **QA Scenarios**:

  ```
  Scenario: Each script field Unknown is omitted
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceScriptsUpdateGuardsUnknownValues -short -v 2>&1 | tee .sisyphus/evidence/task-9-instance-scripts.txt`
    Expected Result: 2 PASS lines
    Evidence: .sisyphus/evidence/task-9-instance-scripts.txt

  Scenario: Negative — Unknown timeout NOT sent as 0
    Tool: Bash (go test)
    Steps:
      1. Same test asserts payload omits the field (does NOT contain `timeout: 0`)
    Evidence: .sisyphus/evidence/task-9-instance-scripts.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-9-instance-scripts.txt`

  **Commit**: YES (own commit)
  - Message: `fix(instance_scripts): guard timeout and update fields against Unknown values`
  - Files: `zstack/provider/resource_zstack_instance_scripts.go`, `zstack/provider/resource_zstack_instance_scripts_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestInstanceScriptsUpdateGuardsUnknownValues -short`

- [ ] 10. **Story-15e: instance_scripts_execution.go REORDER + guard (1 ordering bug)** (Wave 2)

  **What to do**:
  - Site: `zstack/provider/resource_zstack_instance_scripts_execution.go:157-160`
  - **THIS IS NOT A SIMPLE GUARD — IT'S AN ORDERING BUG**: current code calls `.ValueInt64()` BEFORE checking `IsNull()`/`IsUnknown()`. Calling `.ValueInt64()` on Unknown returns 0 silently; calling on Null returns 0 silently. Either way the timeout becomes 0 = immediate timeout.
  - **Per-field Unknown semantics** (DOCUMENT):
    - `ScriptTimeout` (Int64): Null OR Unknown → **omit** from execution payload (server applies default timeout). Add explicit comment in code: `// fixes ordering bug 2026-04-22 QA report — IsNull/IsUnknown must precede ValueInt64`
  - **RED**: Add `TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout` with 2 subtests:
    - subtest A: plan.ScriptTimeout = `types.Int64Null()` → assert payload omits timeout
    - subtest B: plan.ScriptTimeout = `types.Int64Unknown()` → assert payload omits timeout AND no panic
  - Both MUST FAIL initially (because pre-fix code calls `.ValueInt64()` first, returning 0, sending `timeout: 0`).
  - **GREEN**: Restructure the code block at :157-160 from:
    ```go
    timeout := plan.ScriptTimeout.ValueInt64()
    if !plan.ScriptTimeout.IsNull() { payload.Timeout = timeout }
    ```
    to:
    ```go
    if !plan.ScriptTimeout.IsNull() && !plan.ScriptTimeout.IsUnknown() {
        payload.Timeout = plan.ScriptTimeout.ValueInt64()
    }
    ```
    (exact existing variable names confirmed by reading file first)
  - Add code comment: `// fixes ordering bug 2026-04-22 QA report — IsNull/IsUnknown must precede ValueInt64`

  **Must NOT do**:
  - Do NOT modify any other line in the file
  - Do NOT change execution payload struct
  - Do NOT change handler signature
  - Do NOT add live API call

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Subtle ordering bug — code restructure (not just guard insertion); needs careful verification both Null AND Unknown paths covered
  - **Skills**: [`test-driven-development`, `systematic-debugging`]
    - `systematic-debugging`: To trace the exact value flow and ensure no regression in Null path

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_instance_scripts_execution.go:157-160` — exact site
  - QA report `/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md` — Story-15 section explicitly calls this out as ordering bug, not guard gap
  - T9 (instance_scripts.go) for sibling resource pattern

  **Acceptance Criteria**:
  - [ ] Site at :157-160 shows: IsNull/IsUnknown check BEFORE ValueInt64 call
  - [ ] Inline comment present: `fixes ordering bug 2026-04-22 QA report`
  - [ ] `TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout` exists with 2 subtests (Null + Unknown), both PASS
  - [ ] `go test ./zstack/provider/ -run TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout -short -v` exits 0
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_instance_scripts_execution.go` shows ordering swap (not just `&& IsUnknown()` append)

  **QA Scenarios**:

  ```
  Scenario: Null timeout — payload omits field (regression: prevents immediate timeout)
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout/Null -short -v 2>&1 | tee .sisyphus/evidence/task-10-null.txt`
      2. Assert: payload does NOT contain `timeout: 0`
    Expected Result: PASS; payload omits timeout
    Failure Indicators: PASS but payload has timeout=0
    Evidence: .sisyphus/evidence/task-10-null.txt

  Scenario: Unknown timeout — no panic, payload omits field
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout/Unknown -short -v 2>&1 | tee .sisyphus/evidence/task-10-unknown.txt`
    Expected Result: PASS, no panic
    Failure Indicators: panic message in output
    Evidence: .sisyphus/evidence/task-10-unknown.txt

  Scenario: Ordering verified syntactically
    Tool: Bash (grep)
    Steps:
      1. Run `awk 'NR>=155 && NR<=165' zstack/provider/resource_zstack_instance_scripts_execution.go | tee .sisyphus/evidence/task-10-source.txt`
      2. Confirm IsNull AND IsUnknown appear on a line BEFORE any line with ValueInt64
    Expected Result: ordering correct; comment present
    Evidence: .sisyphus/evidence/task-10-source.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-10-null.txt`
  - [ ] `.sisyphus/evidence/task-10-unknown.txt`
  - [ ] `.sisyphus/evidence/task-10-source.txt`

  **Commit**: YES (own commit)
  - Message: `fix(instance_scripts_execution): reorder ScriptTimeout IsNull/IsUnknown check before ValueInt64`
  - Files: `zstack/provider/resource_zstack_instance_scripts_execution.go`, `zstack/provider/resource_zstack_instance_scripts_execution_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestInstanceScriptsExecutionGuardsUnknownAndNullTimeout -short`

- [ ] 11. **Story-15f: l3network.go IsUnknown guards (2 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_l3network.go:185-187`, `:190-192` (read file to confirm exact field names — likely IP-range or DHCP-related Int64/Bool)
  - **Per-field Unknown semantics** (DOCUMENT):
    - `:185-187` field: Unknown → **omit** (network-config defaults applied server-side)
    - `:190-192` field: Unknown → **omit** (same)
  - **RED**: Add `TestL3NetworkUpdateGuardsUnknownValues` with 2 subtests. MUST FAIL initially.
  - **GREEN**: Add `&& !plan.X.IsUnknown()` at both sites.
  - **REFACTOR**: None.

  **Must NOT do**:
  - Do NOT touch IP-range computation
  - Do NOT change DNS server logic
  - Do NOT add live API call

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_l3network.go:185-187, 190-192` — exact sites
  - `zstack/provider/resource_zstack_l3network_test.go` — extend (or create)

  **Acceptance Criteria**:
  - [ ] Both sites have `&& !plan.X.IsUnknown()` added
  - [ ] `TestL3NetworkUpdateGuardsUnknownValues` exists with 2 subtests, both PASS
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_l3network.go | grep -c IsUnknown` returns exactly 2

  **QA Scenarios**:

  ```
  Scenario: Each L3 network field Unknown is omitted
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestL3NetworkUpdateGuardsUnknownValues -short -v 2>&1 | tee .sisyphus/evidence/task-11-l3network.txt`
    Expected Result: 2 PASS lines
    Evidence: .sisyphus/evidence/task-11-l3network.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-11-l3network.txt`

  **Commit**: YES (own commit)
  - Message: `fix(l3network): guard update fields against Unknown values`
  - Files: `zstack/provider/resource_zstack_l3network.go`, `zstack/provider/resource_zstack_l3network_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestL3NetworkUpdateGuardsUnknownValues -short`

- [ ] 12. **Story-15g: vip_qos.go IsUnknown guards (3 sites)** (Wave 2)

  **What to do**:
  - Sites: `zstack/provider/resource_zstack_vip_qos.go:119-130` (3 adjacent QoS rate fields — confirm by reading)
  - **Per-field Unknown semantics** (DOCUMENT):
    - All 3 QoS rate Int64 fields: Unknown → **omit** (zero rate = unlimited or invalid depending on semantic; sending Unknown=0 changes QoS unintentionally)
  - **RED**: Add `TestVipQosUpdateGuardsUnknownValues` with 3 subtests. MUST FAIL initially.
  - **GREEN**: Add `&& !plan.X.IsUnknown()` at each site.
  - **REFACTOR**: None.

  **Must NOT do**:
  - Do NOT touch VIP base resource (`resource_zstack_vip.go`)
  - Do NOT change QoS validation rules
  - Do NOT add live API call

  **Recommended Agent Profile**:
  - **Category**: `quick`
  - **Skills**: [`test-driven-development`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2
  - **Blocks**: T19, T20
  - **Blocked By**: T1

  **References**:
  - `zstack/provider/resource_zstack_vip_qos.go:119-130` — exact site
  - `zstack/provider/resource_zstack_vip_qos_test.go` — extend (or create)

  **Acceptance Criteria**:
  - [ ] All 3 sites have `&& !plan.X.IsUnknown()` added
  - [ ] `TestVipQosUpdateGuardsUnknownValues` exists with 3 subtests, all PASS
  - [ ] `git diff origin/master -- zstack/provider/resource_zstack_vip_qos.go | grep -c IsUnknown` returns exactly 3

  **QA Scenarios**:

  ```
  Scenario: Each QoS rate field Unknown is omitted
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestVipQosUpdateGuardsUnknownValues -short -v 2>&1 | tee .sisyphus/evidence/task-12-vip-qos.txt`
    Expected Result: 3 PASS lines
    Evidence: .sisyphus/evidence/task-12-vip-qos.txt

  Scenario: Negative — Unknown rate NOT sent as 0
    Tool: Bash (go test)
    Steps:
      1. Same test asserts payload omits the field (does NOT contain `rate: 0`)
    Evidence: .sisyphus/evidence/task-12-vip-qos.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-12-vip-qos.txt`

  **Commit**: YES (own commit)
  - Message: `fix(vip_qos): guard rate fields against Unknown values`
  - Files: `zstack/provider/resource_zstack_vip_qos.go`, `zstack/provider/resource_zstack_vip_qos_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestVipQosUpdateGuardsUnknownValues -short`

- [ ] 13. **Scanner Check 2d: Append result discarded** (Wave 3 — implement first per oracle order)

  **What to do**:
  - Add fixtures (within `zstack/provider/testdata/antipatterns/check_2d/`):
    - `bad/append_discarded.go.fixture` — function with `append(s, x)` whose result is discarded (e.g., `append(diags, diag)` not assigned). Add 2-3 distinct bad patterns.
    - `good/append_assigned.go.fixture` — function with `s = append(s, x)` properly reassigned. Add 2-3 distinct good patterns including method-call-result `r.Diags = append(r.Diags, ...)`.
  - In `resource_antipattern_ast_test.go`, add subtest `TestAntipatternScanner/check_2d`:
    - Walks all `*.go.fixture` files under `testdata/antipatterns/check_2d/` (rename to `.go` in-memory before `parser.ParseFile`)
    - For each `ast.CallExpr` where `Fun` is identifier `append`, check whether enclosing `ast.AssignStmt` exists where this call is the RHS and result is assigned
    - Findings: file:line for any unassigned `append(...)` call
  - Acceptance: subtest passes when scanner finds N bad patterns in `bad/` and 0 in `good/`
  - Add a final repo-sweep call inside the subtest that runs the scanner on `zstack/provider/` and prints findings (assertion deferred to T19)

  **Must NOT do**:
  - Do NOT use `go/types` (full type-checking) — keep syntactic
  - Do NOT extend `resource_antipattern_test.go` (the legacy string scanner)
  - Do NOT add fixtures with `.go` extension (would break `go build ./...`)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Real AST traversal logic; needs careful enclosing-AssignStmt detection (parent tracking); easy to false-positive
  - **Skills**: [`test-driven-development`, `systematic-debugging`]

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T14, T15 — different checks; with Wave 2 since fixtures are isolated)
  - **Parallel Group**: Wave 3 (with T14, T15)
  - **Blocks**: T16 (fixture self-test), T19 (repo sweep)
  - **Blocked By**: T1 (harness skeleton), T2 (fixture dirs)

  **References**:

  **Pattern References**:
  - `zstack/provider/resource_antipattern_test.go` — note its check structure for consistency, but DO NOT extend it
  - `zstack/provider/resource_zstack_tag_attachment.go` — historical example of bug pattern (now fixed in T5)

  **API/Type References**:
  - `go/ast.CallExpr.Fun` — for identifier check
  - `go/ast.AssignStmt.Rhs` — for assignment detection
  - `go/ast.Inspect` with parent tracking via stack (Go ast doesn't carry parent pointers; use closure-based stack)

  **External References**:
  - `https://pkg.go.dev/go/ast#Inspect` — official traversal docs

  **WHY**:
  - Oracle ordered 2d first because it's the most localized check (single function scope) — successfully landing it builds confidence in the harness API before tackling 2b (whole-file pattern) and 2a (whole-resource semantic).

  **Acceptance Criteria**:
  - [ ] Fixtures: `bad/append_discarded.go.fixture` exists with at least 2 distinct bad patterns
  - [ ] Fixtures: `good/append_assigned.go.fixture` exists with at least 2 distinct good patterns
  - [ ] `go test ./zstack/provider/ -run TestAntipatternScanner/check_2d -short -v` exits 0
  - [ ] Subtest output shows `bad/` flagged (N findings >= 2) and `good/` flagged 0 times
  - [ ] `go build ./...` exits 0 (fixtures don't compile)

  **QA Scenarios**:

  ```
  Scenario: Bad fixtures correctly flagged
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2d/bad -short -v 2>&1 | tee .sisyphus/evidence/task-13-bad.txt`
    Expected Result: PASS, output lists each bad fixture file:line
    Evidence: .sisyphus/evidence/task-13-bad.txt

  Scenario: Good fixtures NOT flagged
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2d/good -short -v 2>&1 | tee .sisyphus/evidence/task-13-good.txt`
    Expected Result: PASS, zero findings on good/
    Failure Indicators: any finding listed → false-positive
    Evidence: .sisyphus/evidence/task-13-good.txt

  Scenario: Build still clean (fixtures excluded from compilation)
    Tool: Bash (go build)
    Steps:
      1. Run `go build ./... 2>&1 | tee .sisyphus/evidence/task-13-build.txt`
    Expected Result: exit 0
    Evidence: .sisyphus/evidence/task-13-build.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-13-bad.txt`
  - [ ] `.sisyphus/evidence/task-13-good.txt`
  - [ ] `.sisyphus/evidence/task-13-build.txt`

  **Commit**: YES (own commit)
  - Message: `test(antipattern): add check 2d (append result discarded) with fixtures`
  - Files: `zstack/provider/resource_antipattern_ast_test.go`, `zstack/provider/testdata/antipatterns/check_2d/**`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2d -short`

- [ ] 14. **Scanner Check 2b: IsNull without IsUnknown for Int64/Bool fields** (Wave 3)

  **What to do**:
  - Add fixtures (`zstack/provider/testdata/antipatterns/check_2b/`):
    - `bad/isnull_no_isunknown.go.fixture` — function with `if !plan.X.IsNull() { ... plan.X.ValueInt64() }` where X is Int64 or Bool, no IsUnknown check (3+ patterns)
    - `good/isnull_with_isunknown.go.fixture` — `if !plan.X.IsNull() && !plan.X.IsUnknown() { ... }` (3+ patterns)
    - `good/string_field_isnull.go.fixture` — same pattern but with `.IsNull()` on `types.String` — MUST NOT be flagged (out of scope)
    - `good/list_field_isnull.go.fixture` — same with `types.List` — MUST NOT be flagged
  - In `resource_antipattern_ast_test.go`, add subtest `TestAntipatternScanner/check_2b`:
    - Find `ast.IfStmt` whose Cond contains `<expr>.IsNull()` selector call AND DOES NOT contain `<same expr>.IsUnknown()`
    - Heuristic for type filter: examine the receiver expr — if its identifier is followed by a call to `.ValueInt64()` or `.ValueBool()` inside the IfStmt body, treat as in-scope (Int64/Bool); skip if `.ValueString()` / `.Elements()` / `.ElementsAs()` (List/Set/Map/String)
    - Findings: file:line of the IfStmt
  - Acceptance: scanner flags ALL bad patterns; flags ZERO good patterns (including string and list goods)

  **Must NOT do**:
  - Do NOT match List/Set/Map/Object/String/Float64 (see Must NOT Have section)
  - Do NOT use `go/types` for type resolution (syntactic heuristic only)
  - Do NOT report findings for `if !plan.X.IsNull() && !plan.X.IsUnknown()` (compound check is the GOOD pattern)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Most complex check — requires AST traversal + receiver-type heuristic + boolean expression decomposition
  - **Skills**: [`test-driven-development`, `systematic-debugging`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T13, T15)
  - **Blocks**: T17, T19
  - **Blocked By**: T1, T2

  **References**:
  - QA report `/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md` — Story-15 section defines the bad pattern
  - `zstack/provider/resource_zstack_volume.go:217` (post-fix) — canonical good example
  - `zstack/provider/resource_zstack_volume.go:217` (pre-fix, via `git show origin/master:`) — canonical bad example (use as reference for fixture content)

  **API/Type References**:
  - `go/ast.IfStmt.Cond`, `go/ast.BinaryExpr.Op == token.LAND` — to walk `&&` clauses
  - `go/ast.SelectorExpr.Sel.Name` — to compare method names

  **WHY**:
  - This is the largest false-positive surface. The receiver-type heuristic via co-located `.ValueInt64()`/`.ValueBool()` call is the agreed compromise (vs. full type-checking with `go/types` which oracle deferred).

  **Acceptance Criteria**:
  - [ ] All 4 fixture files exist (1 bad + 3 good — including string and list goods)
  - [ ] `go test ./zstack/provider/ -run TestAntipatternScanner/check_2b -short -v` exits 0
  - [ ] Bad fixtures flagged with file:line for each bad pattern
  - [ ] Good fixtures (including string and list) flagged zero times
  - [ ] Subtest report explicitly mentions "string fixture not flagged" and "list fixture not flagged" in test output (use `t.Logf`)

  **QA Scenarios**:

  ```
  Scenario: Bad fixtures flagged with correct count
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2b/bad -short -v 2>&1 | tee .sisyphus/evidence/task-14-bad.txt`
    Expected Result: PASS, finding count >= 3
    Evidence: .sisyphus/evidence/task-14-bad.txt

  Scenario: Good fixtures NOT flagged (including type-out-of-scope cases)
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2b/good -short -v 2>&1 | tee .sisyphus/evidence/task-14-good.txt`
      2. Confirm output mentions string and list fixtures explicitly NOT flagged
    Expected Result: PASS, zero findings, explicit type-exclusion log lines
    Failure Indicators: any string/list pattern flagged → false-positive
    Evidence: .sisyphus/evidence/task-14-good.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-14-bad.txt`
  - [ ] `.sisyphus/evidence/task-14-good.txt`

  **Commit**: YES (own commit)
  - Message: `test(antipattern): add check 2b (IsNull without IsUnknown Int64/Bool) with fixtures`
  - Files: `zstack/provider/resource_antipattern_ast_test.go`, `zstack/provider/testdata/antipatterns/check_2b/**`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2b -short`

- [ ] 15. **Scanner Check 2a: Update missing read-after-write** (Wave 3)

  **What to do**:
  - Add fixtures (`zstack/provider/testdata/antipatterns/check_2a/`):
    - `bad/update_no_reread.go.fixture` — Update method that calls some Update API but does NOT call any subsequent Read/Get/Query API and does NOT call `setStateFromAPIResponse`-style helper. (2+ patterns)
    - `good/update_with_reread.go.fixture` — Update method ending with read call. (2+ patterns)
  - In `resource_antipattern_ast_test.go`, add subtest `TestAntipatternScanner/check_2a`:
    - Find `ast.FuncDecl` named `Update` (receiver is a resource type)
    - Walk body for any `.Update*` API call
    - Then check if any subsequent statement in same function calls `.Get*`, `.Query*`, or assigns to state from a non-plan source
    - Findings: file:line of Update funcdecl when Update API call exists but no subsequent read
  - Acceptance: scanner flags `bad/`; ZERO findings on `good/`

  **Must NOT do**:
  - Do NOT flag Create methods (out of scope for 2a)
  - Do NOT flag Delete methods
  - Do NOT use `go/types`
  - Do NOT touch `resource_zstack_license.go` or `resource_zstack_host.go` for fixture references (oracle: false-positive sources for the deferred 2c, but 2a may also see noise — verify by reading these files first; if scanner WOULD flag them spuriously, refine the heuristic before T19 repo sweep)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Whole-function semantic check; needs careful "subsequent statement" analysis
  - **Skills**: [`test-driven-development`, `systematic-debugging`]

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T13, T14)
  - **Blocks**: T18, T19
  - **Blocked By**: T1, T2

  **References**:
  - `zstack/provider/resource_zstack_volume.go` Update method — likely a good example (read after write)
  - QA report Story-14 section — Update read-after-write context (we're NOT fixing Story-14, just scanning for it)

  **API/Type References**:
  - `go/ast.FuncDecl.Name.Name == "Update"`, `Recv != nil`
  - Iterate `FuncDecl.Body.List` for ordered statement walk

  **WHY**:
  - This is the cheapest static signal for Story-14 (which we're explicitly NOT fixing this batch). Surfacing the count gives the next batch hard data without requiring fix work.

  **Acceptance Criteria**:
  - [ ] Both fixture pairs exist
  - [ ] `go test ./zstack/provider/ -run TestAntipatternScanner/check_2a -short -v` exits 0
  - [ ] Bad fixtures flagged
  - [ ] Good fixtures flagged 0 times
  - [ ] Repo sweep on `zstack/provider/` (called inside subtest) prints finding count to test log (assertion deferred to T19; T15 only logs)

  **QA Scenarios**:

  ```
  Scenario: Bad Update fixtures flagged
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2a/bad -short -v 2>&1 | tee .sisyphus/evidence/task-15-bad.txt`
    Expected Result: PASS, findings >= 2
    Evidence: .sisyphus/evidence/task-15-bad.txt

  Scenario: Good Update fixtures NOT flagged
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2a/good -short -v 2>&1 | tee .sisyphus/evidence/task-15-good.txt`
    Expected Result: PASS, zero findings
    Evidence: .sisyphus/evidence/task-15-good.txt

  Scenario: Repo finding count logged (informational; assertion in T19)
    Tool: Bash (go test)
    Steps:
      1. Run `go test ./zstack/provider/ -run TestAntipatternScanner/check_2a/repo_sweep -short -v 2>&1 | tee .sisyphus/evidence/task-15-repo-count.txt`
    Expected Result: PASS, log line printing finding count (Story-14 backlog visibility)
    Evidence: .sisyphus/evidence/task-15-repo-count.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-15-bad.txt`
  - [ ] `.sisyphus/evidence/task-15-good.txt`
  - [ ] `.sisyphus/evidence/task-15-repo-count.txt`

  **Commit**: YES (own commit)
  - Message: `test(antipattern): add check 2a (Update missing read-after-write) with fixtures`
  - Files: `zstack/provider/resource_antipattern_ast_test.go`, `zstack/provider/testdata/antipatterns/check_2a/**`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2a -short`

- [ ] 16. **T16 — Self-test pass: check_2d (BUG-001 family — Delete ignores ResourceUuids)**

  **What to do**:
  - Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2d -short -v` against the fixtures landed in T13.
  - Confirm `bad/*.go.fixture` files trigger findings (subtest `TestAntipatternScanner/check_2d/bad` PASS = scanner reported finding).
  - Confirm `good/*.go.fixture` files trigger ZERO findings (subtest `TestAntipatternScanner/check_2d/good` PASS = scanner stayed silent).
  - If a subtest fails: this task **rolls back** to a micro-fix on T13 scanner logic OR the fixture (whichever is wrong per oracle review). Do NOT add a third fixture to "make it pass" — fix the root cause.
  - Capture stdout to evidence file.

  **Must NOT do**:
  - Skip subtest failures with `t.Skip` or `// TODO`
  - Add new fixtures beyond what T13 declared
  - Touch scanner check_2b or check_2a logic (different tasks)
  - Modify `bad/*.go.fixture` to silence findings — bad MUST stay bad

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Pure verification task — run a test, read output, decide pass/fail. No new code unless rollback to T13.
  - **Skills**: []
    - Reason: No skill overlap; `go test` invocation is trivial.

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T17, T18) — all three are independent self-test verifications
  - **Blocks**: T19 (repo-sweep cannot start until all three checks self-test green)
  - **Blocked By**: T13 (check_2d implementation + fixtures)

  **References**:
  - **Pattern References**:
    - `.sisyphus/plans/qa-20260422-p0-plus-story15.md#T13` — fixture filenames and expected finding counts T13 declared
  - **Test References**:
    - `zstack/provider/resource_antipattern_ast_test.go:TestAntipatternScanner/check_2d` (created in T13) — subtest names to invoke
  - **WHY Each Reference Matters**:
    - T13 owns the contract; T16 is the gate that proves the contract holds before T19 runs the scanner against the real codebase

  **Acceptance Criteria**:
  - [ ] `go test ./zstack/provider/ -run TestAntipatternScanner/check_2d -short -v` exit code = 0
  - [ ] Output contains both `--- PASS: TestAntipatternScanner/check_2d/bad` and `--- PASS: TestAntipatternScanner/check_2d/good`
  - [ ] Zero `FAIL` lines in scanner stdout
  - [ ] `.sisyphus/evidence/task-16-self-test.txt` exists and is non-empty

  **QA Scenarios**:
  ```
  Scenario: Happy path — both bad and good subtests pass
    Tool: Bash (go test)
    Preconditions: T13 commit landed; testdata/antipatterns/check_2d/{bad,good}/ populated
    Steps:
      1. Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2d -short -v 2>&1 | tee .sisyphus/evidence/task-16-self-test.txt`
      2. Assert exit code = 0
      3. Assert grep -c "^--- PASS: TestAntipatternScanner/check_2d/" = 2 (bad + good)
      4. Assert grep -c "^--- FAIL" = 0
    Expected Result: PASS, both subtests green, evidence captured
    Failure Indicators: Any FAIL line; missing PASS line for bad or good
    Evidence: .sisyphus/evidence/task-16-self-test.txt

  Scenario: Failure case — scanner regression detected
    Tool: Bash (go test)
    Preconditions: Hypothetical — assume T13 scanner has bug missing bad fixture
    Steps:
      1. Run the same command as above
      2. Observe `--- FAIL: TestAntipatternScanner/check_2d/bad` in stdout
      3. STOP — do NOT proceed to T19. Rollback: open `resource_antipattern_ast_test.go`, fix scanner logic, re-run T16
    Expected Result: Failure surfaces immediately; T19 stays blocked until green
    Evidence: .sisyphus/evidence/task-16-self-test-fail.txt (only if rollback path triggered)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-16-self-test.txt`

  **Commit**: NO (folded into T13 commit if no rollback needed; only commits if T13 micro-fix required, in which case amend T13 with `GIT_MASTER=1 git commit --amend --no-edit`)

- [ ] 17. **T17 — Self-test pass: check_2b (Story-15 IsNull-only HIGH gap on Int64/Bool scalars)**

  **What to do**:
  - Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2b -short -v` against fixtures from T14.
  - Confirm `bad/` fixtures trigger findings on Int64/Bool scalar fields with `IsNull()` only (no `IsUnknown()` companion).
  - Confirm `good/` fixtures (Int64/Bool with proper `!IsNull() && !IsUnknown()` chain, plus String/Float64/List/Set/Map/Object/nested-Object as exclusion controls) trigger ZERO findings.
  - Per oracle review: T14 explicitly EXCLUDES String/Float64/List/Set/Map/Object — verify those are silent in `good/` fixtures (they should never be reported).
  - If subtest fails: rollback to T14 micro-fix.

  **Must NOT do**:
  - Extend scanner to flag String/Float64/Collection types — out of scope per oracle, would create false positives across the codebase
  - Add fixtures for Story-14 (BUG-022 / overall scanner) — separate scope
  - Touch check_2d or check_2a logic

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Verification only.
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T16, T18)
  - **Blocks**: T19
  - **Blocked By**: T14

  **References**:
  - **Pattern References**:
    - `.sisyphus/plans/qa-20260422-p0-plus-story15.md#T14` — exact fixture inventory
  - **Test References**:
    - `zstack/provider/resource_antipattern_ast_test.go:TestAntipatternScanner/check_2b` (from T14)
  - **WHY**: T14 contract enforcement gate — ensures only Int64/Bool scalars are flagged, eliminating false-positive risk on String/Float64/Collections

  **Acceptance Criteria**:
  - [ ] Exit code = 0
  - [ ] Both bad and good subtests PASS
  - [ ] Specifically: scanner reports findings ONLY for Int64/Bool fixtures in bad/; reports ZERO findings for String/Float64/List/Set/Map/Object exclusion-control fixtures in good/
  - [ ] `.sisyphus/evidence/task-17-self-test.txt` exists

  **QA Scenarios**:
  ```
  Scenario: Happy path — Int64/Bool flagged, exclusion types silent
    Tool: Bash (go test)
    Preconditions: T14 landed; testdata/antipatterns/check_2b/{bad,good}/ populated with both target types AND exclusion controls
    Steps:
      1. Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2b -short -v 2>&1 | tee .sisyphus/evidence/task-17-self-test.txt`
      2. Assert exit code = 0
      3. Assert both `--- PASS: TestAntipatternScanner/check_2b/bad` and `--- PASS: TestAntipatternScanner/check_2b/good` present
      4. Manually grep finding output: confirm reported file:line matches Int64/Bool fixtures only, never String/Float64/Collection fixtures
    Expected Result: PASS; scanner discrimination between target types and exclusion types verified
    Evidence: .sisyphus/evidence/task-17-self-test.txt

  Scenario: Negative — exclusion type false positive caught
    Tool: Bash (go test)
    Preconditions: Hypothetical — scanner over-broad and flags a String field
    Steps:
      1. Run the same command
      2. Observe `--- FAIL: TestAntipatternScanner/check_2b/good` because a String exclusion-control fixture was wrongly flagged
      3. STOP — rollback to T14, tighten type filter (must check `types.Int64` or `types.Bool` AST node only)
    Expected Result: Over-broad scanner caught immediately
    Evidence: .sisyphus/evidence/task-17-self-test-fail.txt (rollback path only)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-17-self-test.txt`

  **Commit**: NO (folded into T14; amend if micro-fix required)

- [ ] 18. **T18 — Self-test pass: check_2a (Update without read-after-write — informational only)**

  **What to do**:
  - Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2a -short -v` against fixtures from T15.
  - Confirm `bad/` (Update writes via SDK then immediately calls `resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)` without re-reading the resource) → finding reported.
  - Confirm `good/` (Update writes, then re-reads via Read function or explicit GET, then sets state from re-read) → ZERO findings.
  - Per oracle review: T15 includes a `repo_sweep` subtest that LOGS finding count for the real `zstack/provider/` codebase but does NOT assert (Story-14 backlog visibility, NOT in this batch's scope).
  - Verify the repo_sweep subtest logs a count line and does NOT fail the test even if N > 0.

  **Must NOT do**:
  - Add an assertion on `repo_sweep` finding count — this is informational only this batch (Story-14 is out of scope)
  - Auto-fix any finding from `repo_sweep` in this batch — out of scope; only log
  - Touch check_2d or check_2b logic

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Verification only.
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 3 (with T16, T17)
  - **Blocks**: T19
  - **Blocked By**: T15

  **References**:
  - **Pattern References**:
    - `.sisyphus/plans/qa-20260422-p0-plus-story15.md#T15` — fixture inventory and repo_sweep behavior contract
  - **Test References**:
    - `zstack/provider/resource_antipattern_ast_test.go:TestAntipatternScanner/check_2a` (from T15)
  - **WHY**: Confirms informational-only contract holds; ensures repo_sweep does NOT fail the build despite Story-14 backlog being open

  **Acceptance Criteria**:
  - [ ] Exit code = 0
  - [ ] Subtests `bad`, `good`, AND `repo_sweep` all PASS (PASS for repo_sweep means "ran without panic"; finding count is logged, not asserted)
  - [ ] Stdout contains a line matching pattern `check_2a repo finding count: \d+` (for Story-14 backlog visibility)
  - [ ] `.sisyphus/evidence/task-18-self-test.txt` exists

  **QA Scenarios**:
  ```
  Scenario: Happy path — fixtures discriminate, repo_sweep logs count without asserting
    Tool: Bash (go test)
    Preconditions: T15 landed
    Steps:
      1. Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/check_2a -short -v 2>&1 | tee .sisyphus/evidence/task-18-self-test.txt`
      2. Assert exit code = 0
      3. Assert all three PASS lines: bad / good / repo_sweep
      4. Assert grep -E "check_2a repo finding count: [0-9]+" matches exactly one line
    Expected Result: PASS; informational count logged
    Evidence: .sisyphus/evidence/task-18-self-test.txt

  Scenario: Negative — repo_sweep wrongly fails the build
    Tool: Bash (go test)
    Preconditions: Hypothetical — T15 scanner accidentally calls `t.Errorf` instead of `t.Logf` in repo_sweep
    Steps:
      1. Run the same command
      2. Observe `--- FAIL: TestAntipatternScanner/check_2a/repo_sweep` because repo has known Story-14 backlog findings
      3. STOP — rollback to T15, change `t.Errorf` → `t.Logf` for repo_sweep
    Expected Result: Failure caught; informational contract restored
    Evidence: .sisyphus/evidence/task-18-self-test-fail.txt (rollback path only)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-18-self-test.txt`

  **Commit**: NO (folded into T15; amend if micro-fix required)

- [ ] 19. **T19 — Wave 3.5: Repo-sweep — assert ZERO findings for check_2d and check_2b on `zstack/provider/`**

  **What to do**:
  - Add a new subtest `TestAntipatternScanner/repo_sweep_postfix` (separate from the per-check subtests in T13-T15) that runs check_2d and check_2b against the **real** `zstack/provider/` directory (not testdata fixtures) and **asserts ZERO findings**.
  - check_2a is excluded from this assertion (Story-14 backlog visibility — informational count only, see T18).
  - Walk `zstack/provider/*.go` (skip `*_test.go`, `testdata/`, `examples/`, `tools/`); for each file: parse via `go/parser`, run check_2d and check_2b visitors, collect findings.
  - Assert `len(findings_2d) == 0`. If non-zero: fail with full file:line list — this means a Story-15 hotspot was missed in T5-T12 OR a new BUG-001-class regression was introduced.
  - Assert `len(findings_2b) == 0`. If non-zero: same — Story-15 hotspot still has IsNull-only guard missing the IsUnknown companion.
  - Log finding count for check_2a (informational; do NOT assert).
  - Capture full output to evidence.

  **Must NOT do**:
  - Run scanner against directories outside `zstack/provider/` (out of scope; vendor/SDK/tools/examples are explicitly excluded)
  - Auto-fix findings — if findings exist, FAIL the test; rollback to whichever Wave 2 task missed the hotspot
  - Add `t.Skip` to dodge findings
  - Run before all Wave 2 (T5-T12) AND all Wave 3 self-tests (T16-T18) are green — this is the integration gate

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Integration verification spanning all Wave 2 fixes + scanner; needs careful failure analysis if any finding surfaces. Not trivial like a self-test (`quick`); not deep architecture (`deep`). Mid-tier integration scope.
  - **Skills**: []
    - Reason: Pure Go test execution + diff inspection; no domain-specific skill applies.

  **Parallelization**:
  - **Can Run In Parallel**: NO — gate task; must run alone after Wave 2 + Wave 3
  - **Parallel Group**: Wave 3.5 (solo)
  - **Blocks**: Wave 4 (T20, T21, T22) — tracker/docs/MR cannot ship until repo-sweep proves zero findings
  - **Blocked By**: T5, T6, T7, T8, T9, T10, T11, T12 (ALL Wave 2 fixes), AND T16, T17, T18 (ALL Wave 3 self-tests). This is the **integration gate**.

  **References**:
  - **Pattern References**:
    - `.sisyphus/plans/qa-20260422-p0-plus-story15.md` — `Concrete Deliverables` section lists all 12 Story-15 hotspots that MUST be guarded post-fix
    - `zstack/provider/resource_antipattern_ast_test.go:TestAntipatternScanner/check_2{d,b}` — visitor functions to reuse
  - **Test References**:
    - T16 / T17 / T18 self-tests prove the scanner discriminates correctly; T19 applies that proven scanner to real code
  - **External References**:
    - Go `path/filepath.Walk` — directory traversal for `zstack/provider/*.go`
    - Go `go/parser.ParseFile` + `go/ast.Inspect` — AST walking
  - **WHY Each Reference Matters**:
    - This is the **single most important gate in the plan** — it proves the 12 fixes are complete AND no regression slipped in. If T19 fails, Wave 4 cannot start; the plan recurses to whichever task missed a hotspot
    - The visitor functions are reused (not reimplemented) — single source of scanner truth from T13/T14

  **Acceptance Criteria**:
  - [ ] New subtest `TestAntipatternScanner/repo_sweep_postfix` exists in `resource_antipattern_ast_test.go`
  - [ ] `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/repo_sweep_postfix -short -v` exit code = 0
  - [ ] Stdout contains line matching `repo_sweep_postfix: check_2d findings = 0`
  - [ ] Stdout contains line matching `repo_sweep_postfix: check_2b findings = 0`
  - [ ] Stdout contains line matching `repo_sweep_postfix: check_2a findings = \d+` (informational only — count logged, no assertion)
  - [ ] If ANY finding surfaces for check_2d or check_2b: this task FAILS; rollback to corresponding Wave 2 task by file:line mapping in the failure output
  - [ ] `.sisyphus/evidence/task-19-repo-sweep.txt` exists

  **QA Scenarios**:
  ```
  Scenario: Happy path — all 12 hotspots guarded, scanner reports zero findings on real code
    Tool: Bash (go test)
    Preconditions: T5-T12 ALL committed and merged into branch; T16-T18 self-tests green
    Steps:
      1. Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/repo_sweep_postfix -short -v 2>&1 | tee .sisyphus/evidence/task-19-repo-sweep.txt`
      2. Assert exit code = 0
      3. Assert grep "repo_sweep_postfix: check_2d findings = 0" matches exactly one line
      4. Assert grep "repo_sweep_postfix: check_2b findings = 0" matches exactly one line
      5. Assert grep -E "repo_sweep_postfix: check_2a findings = [0-9]+" matches exactly one line (count is logged, not asserted)
    Expected Result: PASS; the 12 Story-15 hotspots and BUG-001 are all confirmed guarded
    Failure Indicators:
      - `--- FAIL: TestAntipatternScanner/repo_sweep_postfix` with file:line list of unguarded sites
      - Indicates a Wave 2 task missed its target file or a regression was introduced
    Evidence: .sisyphus/evidence/task-19-repo-sweep.txt

  Scenario: Negative — Wave 2 missed a hotspot (e.g., T6 forgot one of the 3 instance.go sites)
    Tool: Bash (go test)
    Preconditions: Hypothetical — T6 only patched 2 of 3 instance.go sites
    Steps:
      1. Run the same command
      2. Observe `--- FAIL: TestAntipatternScanner/repo_sweep_postfix` with output `check_2b findings = 1` and file:line `resource_zstack_instance.go:708`
      3. STOP — do NOT proceed to T20. Rollback: open T6, add the missing guard at line 708, re-run T17 self-test, re-run T19
    Expected Result: Missed hotspot caught at integration gate; root cause traced via file:line
    Evidence: .sisyphus/evidence/task-19-repo-sweep-fail.txt (rollback path only)

  Scenario: Negative — regression introduced (Wave 2 fix accidentally introduced a new IsNull-only site)
    Tool: Bash (go test)
    Preconditions: Hypothetical — T8 volume.go fix accidentally copy-pasted an IsNull-only check at a new location
    Steps:
      1. Run the same command
      2. Observe FAIL with file:line pointing at a NEW location not in the original 12-hotspot list
      3. STOP — rollback the offending edit; this is exactly the regression-detection role of T19
    Expected Result: Regression caught at integration gate
    Evidence: .sisyphus/evidence/task-19-repo-sweep-regression.txt (rollback path only)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-19-repo-sweep.txt`

  **Commit**: YES (own commit; the new `repo_sweep_postfix` subtest is a permanent regression guard)
  - Message: `test(antipattern): repo_sweep_postfix asserts zero IsNull-only / Delete-ignores-uuids findings`
  - Files: `zstack/provider/resource_antipattern_ast_test.go`
  - Pre-commit: `GIT_MASTER=1 go build ./... && GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner -short`

- [ ] 20. **T20 — Wave 4: Update `_bmad-output/bug-tracker.md` with BUG-001 closure + Story-15 references**

  **What to do**:
  - Open `_bmad-output/bug-tracker.md` (gitignored — must use `git add -f`).
  - Add new entry for BUG-001 (tag_attachment.Delete ignoring `state.ResourceUuids`) under the Closed section:
    - Status: `CLOSED — fixed in branch fix/qa-20260422-p0-plus-story15`
    - Root cause: per Metis gap analysis, Delete dispatch loop ignores the SDK's per-uuid signature and passes empty/nil ResourceUuids
    - Fix: T5 commit SHA (replace placeholder during execution)
    - Verification: spy/mock client test in `resource_zstack_tag_attachment_test.go` (added in T5)
  - Add aggregate Story-15 entry referencing all 12 hotspots:
    - Status: `CLOSED — IsUnknown guards added across 7 resources, 12 sites`
    - Sites table: file:line | resource | scalar type | per-field semantics doc'd in commit
    - Verification: scanner check_2b on `zstack/provider/` reports 0 findings (T19 evidence)
  - Add Story-15e (instance_scripts_execution ordering bug) entry:
    - Status: `CLOSED — IsNull/IsUnknown checked before ValueInt64; ordering bug fixed`
    - Note: distinct task because root cause is wrong call order, not just missing guard
  - Add Story-14 backlog visibility note:
    - Status: `OPEN — out of scope this batch; check_2a scanner reports informational count via T19 evidence; tracked for next batch`
  - DO NOT touch any other entries in the tracker.
  - Source for entries: this plan (Concrete Deliverables section + per-task acceptance) and Metis gap analysis findings already incorporated.

  **Must NOT do**:
  - Modify entries for unrelated bugs already in the tracker
  - Sync to QA cross-repo tracker (per user decision: Local tracker + docs only — NO cross-repo sync this batch)
  - Add Story-14 / Story-16 closure entries — out of scope
  - Skip `git add -f` (file is gitignored; default `git add` would silently no-op)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Doc edit; structured append; references commit SHAs from prior tasks. No analysis.
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T21 — both edit independent docs)
  - **Parallel Group**: Wave 4 (with T21)
  - **Blocks**: T22 (push + MR cannot ship without tracker + docs in commit)
  - **Blocked By**: T19 (must have repo-sweep evidence file referenced in BUG-001 + Story-15 closure entries)

  **References**:
  - **Pattern References**:
    - `_bmad-output/bug-tracker.md` (existing entries — match formatting/heading style of e.g. BUG-019, BUG-024 closure entries from prior batch)
  - **WHY**: Local tracker is the single source of truth for what's been fixed; without it, future `/start-work` sessions cannot see Story-15 + BUG-001 closure status

  **Acceptance Criteria**:
  - [ ] `_bmad-output/bug-tracker.md` contains new BUG-001 entry under Closed section with status, root cause, fix commit SHA, verification reference
  - [ ] New Story-15 aggregate entry with 12-site table (file:line | resource | scalar type)
  - [ ] New Story-15e entry for ordering bug
  - [ ] New Story-14 OPEN backlog note referencing T19 evidence as informational source
  - [ ] `GIT_MASTER=1 git add -f _bmad-output/bug-tracker.md` succeeds (file added despite gitignore)
  - [ ] `GIT_MASTER=1 git diff --cached _bmad-output/bug-tracker.md` shows ONLY appended entries (no modifications to existing entries)

  **QA Scenarios**:
  ```
  Scenario: Happy path — tracker updated cleanly with all four entries
    Tool: Bash (cat / grep)
    Preconditions: T19 complete; commit SHAs from T5, T6-T12, T10 known
    Steps:
      1. Open `_bmad-output/bug-tracker.md` and append entries per spec
      2. Run `GIT_MASTER=1 git add -f _bmad-output/bug-tracker.md`
      3. Assert `git diff --cached _bmad-output/bug-tracker.md | grep "^+" | grep -c "BUG-001\|Story-15\|Story-15e\|Story-14"` >= 4
      4. Capture diff to evidence: `GIT_MASTER=1 git diff --cached _bmad-output/bug-tracker.md > .sisyphus/evidence/task-20-tracker-diff.txt`
    Expected Result: All four entries present in cached diff; existing entries unmodified
    Evidence: .sisyphus/evidence/task-20-tracker-diff.txt

  Scenario: Negative — accidentally modified existing entry
    Tool: Bash (git diff)
    Preconditions: Hypothetical — fat-finger edit to BUG-019 line
    Steps:
      1. Run `GIT_MASTER=1 git diff --cached _bmad-output/bug-tracker.md`
      2. Observe modifications to existing BUG-019 lines
      3. STOP — `GIT_MASTER=1 git restore --staged _bmad-output/bug-tracker.md`, redo edits append-only
    Expected Result: Caught before commit; only appends allowed
    Evidence: .sisyphus/evidence/task-20-tracker-diff-fail.txt (rollback path only)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-20-tracker-diff.txt`

  **Commit**: YES (own commit; can be combined with T21 if landed in same execution wave per git-master 3+ files rule check — but tracker and overview are 2 files, allowed in single commit)
  - Message: `docs(bug-tracker): close BUG-001 + Story-15 (12 sites) + Story-15e; note Story-14 backlog`
  - Files: `_bmad-output/bug-tracker.md`
  - Pre-commit: none (doc only)

- [ ] 21. **T21 — Wave 4: Update `_bmad-output/test-status-overview.md` with new tests + scanner status**

  **What to do**:
  - Open `_bmad-output/test-status-overview.md` (gitignored — `git add -f` required).
  - Append a new dated section `## 2026-04-22 — P0 + Story-15 batch` containing:
    - **New tests added** (per Wave 2): list each `*_test.go` test name added in T5-T12, with the resource it covers and the Unknown-value injection assertion shape (per-field semantics doc'd inline)
    - **Scanner status**: 3 new checks active (2d, 2b, 2a); check 2c deferred (per oracle review); finding counts on `zstack/provider/` from T19 (0/0/N for 2d/2b/2a)
    - **BUG-001 spy/mock test**: noted with file path
    - **Coverage delta**: short note on which of the 7 affected resources now have Unknown-value injection coverage (vs none before this batch)
  - DO NOT touch entries for prior dated sections.
  - Source: this plan + T19 evidence + commit history of T5-T12 / T13-T15 / T19.

  **Must NOT do**:
  - Modify pre-existing dated sections
  - Add coverage claims for resources NOT in the 7 affected list (out of scope)
  - Reference QA cross-repo (per user decision)
  - Skip `git add -f`

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Doc append.
  - **Skills**: []

  **Parallelization**:
  - **Can Run In Parallel**: YES (with T20)
  - **Parallel Group**: Wave 4
  - **Blocks**: T22
  - **Blocked By**: T19

  **References**:
  - **Pattern References**:
    - `_bmad-output/test-status-overview.md` (existing dated sections — match formatting)
  - **WHY**: Future planning sessions read this to know test coverage state; without it, BUG-001 and Story-15 coverage gains are invisible to next `/start-work`

  **Acceptance Criteria**:
  - [ ] New section `## 2026-04-22 — P0 + Story-15 batch` exists
  - [ ] Section lists every test added in T5-T12 by exact test function name
  - [ ] Scanner status block shows `check_2d: 0`, `check_2b: 0`, `check_2a: <N> (informational)`
  - [ ] `GIT_MASTER=1 git add -f` succeeds
  - [ ] Diff is append-only

  **QA Scenarios**:
  ```
  Scenario: Happy path — overview appended with new section
    Tool: Bash (grep / git diff)
    Preconditions: T19 complete; T5-T12 test function names known
    Steps:
      1. Append section per spec
      2. Run `GIT_MASTER=1 git add -f _bmad-output/test-status-overview.md`
      3. Assert `grep "## 2026-04-22 — P0 + Story-15 batch" _bmad-output/test-status-overview.md` matches one line
      4. Assert `grep -c "TestAccTagAttachment\|TestAccInstance\|TestAccPortForwarding\|TestAccVolume\|TestAccInstanceScripts\|TestAccL3Network\|TestAccVipQos" _bmad-output/test-status-overview.md` >= 7 (one per affected resource family)
      5. Capture diff: `GIT_MASTER=1 git diff --cached _bmad-output/test-status-overview.md > .sisyphus/evidence/task-21-overview-diff.txt`
    Expected Result: Section present; per-resource test references present
    Evidence: .sisyphus/evidence/task-21-overview-diff.txt

  Scenario: Negative — modified pre-existing dated section
    Tool: Bash (git diff)
    Preconditions: Hypothetical — accidental edit to prior date
    Steps:
      1. Run git diff
      2. Observe modifications to pre-existing date heading or content
      3. STOP — restore staged, redo append-only
    Expected Result: Caught pre-commit
    Evidence: .sisyphus/evidence/task-21-overview-diff-fail.txt (rollback path only)
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-21-overview-diff.txt`

  **Commit**: YES (combined with T20 OK per git-master rules — 2 files, both `_bmad-output/`, single docs concern)
  - Message (if combined): `docs: close BUG-001 + Story-15 + record test/scanner coverage delta`
  - Message (if separate): `docs(test-status): record P0+Story-15 batch test additions and scanner counts`
  - Files: `_bmad-output/test-status-overview.md` (and `_bmad-output/bug-tracker.md` if combined)
  - Pre-commit: none

- [ ] 22. **T22 — Wave 4: Push branch + `glab mr create` to `master`**

  **What to do**:
  - Pre-flight checks (in order; each must pass before proceeding):
    1. `GIT_MASTER=1 git status` — assert clean working tree (all T1-T21 commits landed, no untracked files outside `.sisyphus/evidence/` or expected scope)
    2. `GIT_MASTER=1 git log --oneline origin/master..HEAD` — assert commit count matches plan (one per Wave 1-3 commit + Wave 4 docs commit; ~12-16 commits expected per git-master 3+ files rule)
    3. `GIT_MASTER=1 git fetch origin master` — required because earlier session attempt timed out; rebase if origin moved
    4. If `origin/master` advanced: `GIT_MASTER=1 git rebase origin/master` and re-run T19 self-test
  - Push: `GIT_MASTER=1 git push -u origin fix/qa-20260422-p0-plus-story15`
  - Open MR via `glab mr create`:
    - `--target-branch master`
    - `--source-branch fix/qa-20260422-p0-plus-story15`
    - `--title`: `fix(qa-20260422): P0 quick wins + Story-15 IsUnknown guards (12 sites) + BUG-001`
    - `--description`: include sections — Summary / Closes (BUG-001, Story-15, Story-15e) / Out of scope (Story-14, Story-16, cross-repo sync) / Test plan (T16-T19 evidence summary) / Independent of MR #28
    - `--remove-source-branch` (so branch is cleaned up after merge)
  - Capture MR URL to evidence file.
  - Verify MR mergeable status: `glab mr view <number>` — assert no conflicts; if conflicts: stop and report (do not auto-resolve in this task).

  **Must NOT do**:
  - Force-push (`git push -f`) — branch is fresh, no need
  - Touch `test/progress` branch or MR #28 in any way (independent MR per user decision)
  - Use `--draft` flag — MR opens ready for review
  - Auto-merge (`glab mr merge`) — that's reviewer's call
  - Skip the rebase if `origin/master` advanced — must re-verify on top of latest master
  - Push if T19 evidence shows non-zero check_2d or check_2b findings — gate violation

  **Recommended Agent Profile**:
  - **Category**: `unspecified-high`
    - Reason: Multi-step git + glab orchestration with conditional rebase + conflict-handling decision tree. Beyond `quick`. Not `deep` (no architecture). Mid-tier integration.
  - **Skills**: [`git-master`]
    - Reason: Git push + MR creation + conflict-handling = exactly the git-master skill domain (per skill description: "MUST USE for ANY git operations"). git-master will enforce `GIT_MASTER=1` prefix and 3+ files / 2+ commits rule on any final cleanup commits.

  **Parallelization**:
  - **Can Run In Parallel**: NO — terminal task; only one MR per branch
  - **Parallel Group**: Wave 4 (solo after T20 + T21)
  - **Blocks**: Final verification wave (F1-F4 review the open MR)
  - **Blocked By**: T20, T21 (docs commits must be in branch before push)

  **References**:
  - **Pattern References**:
    - Prior MR #28 description structure (use as template; MR is at `http://dev.zstack.io:9080/zstackio/terraform-provider-zstack/-/merge_requests/28`)
  - **External References**:
    - `glab mr create --help` — flag reference (note: this `glab` version 1.57 does NOT support `--state`; do not use that flag)
    - `git rebase` workflow — required if origin/master advanced during this session
  - **WHY**: Final delivery to GitLab; must be independent of MR #28; must surface clean diff vs current `master`

  **Acceptance Criteria**:
  - [ ] Pre-flight: `git status` shows clean working tree
  - [ ] Pre-flight: `git log --oneline origin/master..HEAD` count matches expected commit count
  - [ ] Pre-flight: `git fetch origin master` succeeds (no timeout); if it fails, this task BLOCKS until network restored
  - [ ] If `origin/master` advanced: rebase succeeds AND T19 self-test re-passes
  - [ ] `git push -u origin fix/qa-20260422-p0-plus-story15` exit code = 0
  - [ ] `glab mr create` exit code = 0; MR URL captured
  - [ ] `glab mr view <number>` shows `Mergeable: Yes` (no conflicts)
  - [ ] MR description includes all required sections (Summary / Closes / Out of scope / Test plan / Independence note)
  - [ ] `.sisyphus/evidence/task-22-mr-url.txt` contains MR URL
  - [ ] `.sisyphus/evidence/task-22-mr-status.txt` contains `glab mr view` output

  **QA Scenarios**:
  ```
  Scenario: Happy path — clean push, MR opened, mergeable
    Tool: Bash (git + glab)
    Preconditions: All Wave 1-4 commits landed; T19 evidence shows zero findings; network to dev.zstack.io:9022 working
    Steps:
      1. Run `GIT_MASTER=1 git status` and assert "nothing to commit, working tree clean"
      2. Run `GIT_MASTER=1 git fetch origin master` and assert exit 0
      3. Run `GIT_MASTER=1 git rebase origin/master` (no-op if already up-to-date) and assert success
      4. Run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/repo_sweep_postfix -short` post-rebase; assert pass
      5. Run `GIT_MASTER=1 git push -u origin fix/qa-20260422-p0-plus-story15`
      6. Run `glab mr create --target-branch master --source-branch fix/qa-20260422-p0-plus-story15 --title "fix(qa-20260422): P0 quick wins + Story-15 IsUnknown guards (12 sites) + BUG-001" --description "$(cat .sisyphus/evidence/mr-description.md)" --remove-source-branch | tee .sisyphus/evidence/task-22-mr-url.txt`
      7. Extract MR number from URL; run `glab mr view <number> | tee .sisyphus/evidence/task-22-mr-status.txt`
      8. Assert grep -i "mergeable.*yes" .sisyphus/evidence/task-22-mr-status.txt
    Expected Result: MR open, mergeable, URL captured
    Evidence: .sisyphus/evidence/task-22-mr-url.txt, .sisyphus/evidence/task-22-mr-status.txt

  Scenario: Negative — origin/master advanced, rebase introduces conflicts
    Tool: Bash (git)
    Preconditions: Hypothetical — another developer landed conflicting changes to one of the 7 affected resources during this session
    Steps:
      1. `git fetch origin master` succeeds, advancing remote tracking ref
      2. `git rebase origin/master` reports conflicts in (e.g.) `resource_zstack_volume.go`
      3. STOP — `GIT_MASTER=1 git rebase --abort`; do NOT auto-resolve. Report to user with conflict file list and ask whether to merge-resolve or replan
      4. Do NOT push; do NOT open MR
    Expected Result: Conflict surfaces; user asked for direction; T22 stays incomplete
    Evidence: .sisyphus/evidence/task-22-conflict.txt

  Scenario: Negative — network to dev.zstack.io:9022 still timing out
    Tool: Bash (git fetch)
    Preconditions: Same as earlier session — ssh timeout
    Steps:
      1. `GIT_MASTER=1 git fetch origin master` times out
      2. STOP — report to user; T22 BLOCKS until network restored
      3. Do NOT push to a stale base
    Expected Result: Network failure surfaces; user can address infra; task resumable
    Evidence: .sisyphus/evidence/task-22-network-fail.txt

  Scenario: Negative — T19 evidence missing or shows non-zero findings (regression sneak-in between T19 and T22)
    Tool: Bash
    Preconditions: Hypothetical — between T19 commit and T22, a stray edit re-introduced an IsNull-only site
    Steps:
      1. Re-run `GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner/repo_sweep_postfix -short -v` as last pre-flight
      2. Observe finding count > 0
      3. STOP — do NOT push; rollback the offending edit; re-verify; resume T22
    Expected Result: Pre-flight catches regression
    Evidence: .sisyphus/evidence/task-22-regression-fail.txt
  ```

  **Evidence**:
  - [ ] `.sisyphus/evidence/task-22-mr-url.txt`
  - [ ] `.sisyphus/evidence/task-22-mr-status.txt`
  - [ ] `.sisyphus/evidence/mr-description.md` (the description text used for `glab mr create`, kept as artifact)

  **Commit**: NO (T22 produces no source-tree changes; only push + MR creation)

---

## Final Verification Wave (MANDATORY — after ALL implementation tasks)

> 4 review agents run in PARALLEL. ALL must APPROVE. Present consolidated results to user and get explicit "okay" before completing.
>
> **Do NOT auto-proceed after verification. Wait for user's explicit approval before marking work complete.**

- [ ] F1. **Plan Compliance Audit** — `oracle`
  Read this plan end-to-end. For each "Must Have": verify implementation exists (read file, run command). For each "Must NOT Have": search codebase for forbidden patterns — REJECT with file:line if found. Verify all 12 Story-15 hotspots received guards at the documented line ranges. Verify BUG-001 fix passes ResourceUuids. Verify scanner harness uses go/parser + go/ast (not ast-grep). Verify scanner check 2c is NOT present. Verify NO files outside scope were touched (`git diff origin/master --stat` should show only files in Concrete Deliverables list). Check evidence files exist in `.sisyphus/evidence/`.
  Output: `Must Have [N/N] | Must NOT Have [N/N] | Tasks [N/N] | VERDICT: APPROVE/REJECT`

- [ ] F2. **Code Quality Review** — `unspecified-high`
  Run `go build ./...` + `go vet ./...` + `go test ./... -short`. Review all changed files for: `interface{}`/`any` casts hiding errors, swallowed errors, panic instead of diagnostics, console prints in production paths, commented-out code, unused imports, magic numbers without comments. Check AI slop: excessive comments stating the obvious, over-abstraction, generic names (data/result/item/temp). Verify each fix is MINIMAL — only adding `&& !plan.Field.IsUnknown()` guards (or reorder for T10), NOT refactoring surrounding code.
  Output: `Build [PASS/FAIL] | Vet [PASS/FAIL] | Tests [N pass/N fail] | Files [N clean/N issues] | VERDICT`

- [ ] F3. **Real Manual QA** — `unspecified-high`
  Start from clean state (fresh checkout of branch). Execute EVERY QA scenario from EVERY task — follow exact steps, capture evidence to `.sisyphus/evidence/final-qa/`. Run scanner repo sweep — assert ZERO findings. Run all Unknown-value injection unit tests — assert no panics. Run BUG-001 spy test — assert ResourceUuids passed. Verify shell scripts in `scripts/qa/` are executable (`chmod +x`) and runnable. Test cross-task integration: do all fixed resources still compile and unit-test together? Check edge cases: empty Unknown, mixed Null+Unknown, repeated Unknown across multiple fields.
  Output: `Scenarios [N/N pass] | Integration [N/N] | Edge Cases [N tested] | VERDICT`

- [ ] F4. **Scope Fidelity Check** — `deep`
  For each task: read "What to do", read actual diff (`git log origin/master..HEAD --stat` and `git diff origin/master`). Verify 1:1 — everything in spec was built (no missing), nothing beyond spec (no creep). Check "Must NOT do" compliance per task. Detect cross-task contamination: Task N touching Task M's files. Flag unaccounted changes. Specifically verify: NO commits touch `fix/bug-tracker-open-items` files, NO modifications to `tools/tools.go` / `examples/` / `docs/`, NO new resources or data sources, NO l2vlan_network changes, NO license/host resource changes for scanner.
  Output: `Tasks [N/N compliant] | Contamination [CLEAN/N issues] | Unaccounted [CLEAN/N files] | VERDICT`

---

## Commit Strategy

> Per `git-master` skill: 3+ files changed → 2+ commits (NO EXCEPTIONS). Every git command prefixed with `GIT_MASTER=1`.

### Commit Granularity Rules
- **One task = one logical commit** (with its test + fix together for TDD tasks)
- **Foundation tasks** (T1-T4): each gets its own commit
- **Per-resource Story-15 fixes** (T6-T12): each resource = its own commit (test + fix in same commit per TDD)
- **BUG-001** (T5): own commit (Delete fix + spy test)
- **Scanner checks** (T13-T18): each check (2d/2b/2a) = its own commit (fixtures + scanner code + self-test pass together)
- **Repo sweep** (T19): own commit if any harness tweaks needed
- **Tracker/docs** (T20-T21): combined commit using `git add -f` for `_bmad-output/`

### Commit Message Style (English, semantic)
```
type(scope): description

Examples:
test(provider): add AST scanner harness skeleton for antipattern checks
fix(tag_attachment): pass state.ResourceUuids to DetachTagFromResources
fix(instance): guard NeverStop/MemorySize/CPUNum against Unknown values
fix(instance_scripts_execution): reorder ScriptTimeout IsNull check before ValueInt64
test(antipattern): add check 2d (append result discarded) with fixtures
chore(qa): drop QA report scan scripts into scripts/qa/
docs(tracker): add BUG-001 and Story-15 fix records
```

### Pre-Commit Verification (per commit)
- `GIT_MASTER=1 go build ./...` exits 0
- `GIT_MASTER=1 go test ./zstack/provider/ -run <relevant test> -short -v` passes
- `GIT_MASTER=1 git status` shows expected files only

---

## Success Criteria

### Verification Commands
```bash
# Compile
GIT_MASTER=1 go build ./...
# Expected: exit 0, no output

# All short tests
GIT_MASTER=1 go test ./... -short
# Expected: ok zstack/provider ... PASS

# Scanner self-tests on fixtures + repo sweep
GIT_MASTER=1 go test ./zstack/provider/ -run TestAntipatternScanner -short -v
# Expected: PASS — fixtures flagged correctly, repo sweep returns 0 findings

# Story-15 unit tests (Unknown-value injection)
GIT_MASTER=1 go test ./zstack/provider/ -run TestUnknownValueGuards -short -v
# Expected: PASS — no panics on Unknown plan values

# BUG-001 spy test
GIT_MASTER=1 go test ./zstack/provider/ -run TestTagAttachmentDeletePassesResourceUuids -short -v
# Expected: PASS — spy receives expected UUID slice

# Branch + MR exists
GIT_MASTER=1 git log origin/master..HEAD --oneline
# Expected: list of commits per Commit Strategy

glab mr list --source-branch fix/qa-20260422-p0-plus-story15
# Expected: 1 open MR targeting master
```

### Final Checklist
- [ ] All "Must Have" present (verify via F1)
- [ ] All "Must NOT Have" absent (verify via F1 + F4)
- [ ] All 12 Story-15 sites guarded
- [ ] BUG-001 fix verified by spy/mock test
- [ ] Scanner harness compiles + 3 checks active + repo sweep clean
- [ ] QA shell scripts dropped + executable
- [ ] Local tracker + docs updated
- [ ] Branch pushed, MR open
- [ ] F1-F4 all APPROVE
- [ ] User explicitly says "okay"
