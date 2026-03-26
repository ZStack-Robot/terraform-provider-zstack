# Sprint 2 Test Results — STORY-TF-001-03 + STORY-TF-001-07

**Date**: 2026-03-26
**Environment**: ZStack @ 172.24.227.46:8080

## Unit Tests (Schema + Metadata)

All 8 unit tests pass:

```
=== RUN   TestAccountResource_Schema
--- PASS: TestAccountResource_Schema (0.00s)
=== RUN   TestAccountResource_Metadata
--- PASS: TestAccountResource_Metadata (0.00s)
=== RUN   TestAffinityGroupResource_Schema
--- PASS: TestAffinityGroupResource_Schema (0.00s)
=== RUN   TestAffinityGroupResource_Metadata
--- PASS: TestAffinityGroupResource_Metadata (0.00s)
=== RUN   TestIAM2ProjectResource_Schema
--- PASS: TestIAM2ProjectResource_Schema (0.00s)
=== RUN   TestIAM2ProjectResource_Metadata
--- PASS: TestIAM2ProjectResource_Metadata (0.00s)
=== RUN   TestSshKeyPairResource_Schema
--- PASS: TestSshKeyPairResource_Schema (0.00s)
=== RUN   TestSshKeyPairResource_Metadata
--- PASS: TestSshKeyPairResource_Metadata (0.00s)
PASS
```

## Terraform Apply / Destroy (Live Environment)

| Resource | Apply | Destroy | Notes |
|---|---|---|---|
| `zstack_affinity_group` | ✅ PASS | ✅ PASS | policy=antiSoft, type=HOST, state=Enabled |
| `zstack_ssh_key_pair` | ✅ PASS | ✅ PASS | public_key import works |
| `zstack_account` | ✅ PASS | ✅ PASS | password sensitive field, type=Normal |
| `zstack_iam2_project` | ✅ PASS | ✅ PASS | state=Enabled |

### Resource Apply Details

**zstack_affinity_group**
- UUID: aa9b6afea74b4555a8d5a1225aa2766e
- Created with: name="tf-batch-test-affinity-group", policy="antiSoft"
- Returned: type="HOST", state="Enabled"
- Fixed: API uppercases policy ("antiSoft" → "ANTISOFT"), resolved with case-insensitive comparison

**zstack_ssh_key_pair**
- UUID: 651b30f84fa04b53a6afc43fb8bad8b2
- Created with: name="tf-batch-test-ssh-key-pair", public_key="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7 test@batch-test"

**zstack_account**
- UUID: edb9bba28d4444dd80cc7a7e2a86d65f
- Created with: name="tf-batch-test-account", password=(sensitive)
- Returned: type="Normal"

**zstack_iam2_project**
- UUID: 0ce5f6d2ba0f486db9fc4d25108ec03a
- Created with: name="tf-batch-test-iam2-project"
- Returned: state="Enabled"

## Data Source Queries (Live Environment)

| Data Source | Status | Records |
|---|---|---|
| `zstack_accounts` | ✅ PASS | 1 (admin / SystemAdmin) |
| `zstack_affinity_groups` | ✅ PASS | 1 (zstack.affinity.group.for.virtual.router) |
| `zstack_iam2_projects` | ✅ PASS | Query returned results |
| `zstack_ssh_key_pairs` | ✅ PASS | 0 (empty list, expected) |

## Documentation Generation

```bash
cd tools && go generate ./...
```

Generated docs:
- `docs/resources/account.md`
- `docs/resources/affinity_group.md`
- `docs/resources/iam2_project.md`
- `docs/resources/ssh_key_pair.md`
- `docs/data-sources/accounts.md`
- `docs/data-sources/affinity_groups.md`
- `docs/data-sources/iam2_projects.md`
- `docs/data-sources/ssh_key_pairs.md`

## Issues Found & Fixed

1. **Affinity Group Policy Case Mismatch**: ZStack API returns policy in uppercase (e.g. `ANTISOFT`), but users input camelCase (e.g. `antiSoft`). Terraform's `ProducedInconsistentResult` error was triggered. Fixed by using `strings.EqualFold` to preserve user-provided case.

## Build Verification

```
go build ./...  ✅ PASS
```
