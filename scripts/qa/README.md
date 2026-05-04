# QA Scan Scripts

Shell scripts for detecting systematic bug patterns in the Terraform ZStack Provider codebase.

## Scripts

### `scan_isnull_gap.sh`

**Purpose**: Detect `IsNull()` checks missing `IsUnknown()` guards in resource Create/Update methods.

**Problem**: Optional+Computed fields have `Unknown` (not `Null`) when users don't set them. Code checking only `IsNull()` misses `Unknown`, sending zero values (`0`, `false`, `""`) to the API.

**Usage**:
```bash
./scripts/qa/scan_isnull_gap.sh
```

**Output Format**:
- `file:line: HIGH — <code snippet>` — Int64/Bool fields (high risk)
- `file: MEDIUM — gap=N` — String fields (medium risk)

**Exit Code**:
- `0` — No findings
- `1` — Findings detected

**Related Issues**: BUG-5, BUG-12b, BUG-13

---

### `scan_update_no_reread.sh`

**Purpose**: Detect Update methods that trust SDK return values without read-after-write verification.

**Problem**: SDK's `PutWithRespKey` may return empty structs. Provider Update methods using these values directly zero out state fields, causing "inconsistent result" errors.

**Usage**:
```bash
./scripts/qa/scan_update_no_reread.sh
```

**Output Format**:
- `file:line: SAFE` — Has read-after-write (Query/Get after Update)
- `file:line: DANGER — no read-after-write` — Trusts SDK return value

**Exit Code**:
- `0` — No dangerous patterns
- `1` — Dangerous patterns detected

**Related Issues**: BUG-6, BUG-9, BUG-11, BUG-14

---

## Relationship to Go Scanner (Layer 2)

These shell scripts provide **quick, portable detection** for two of the three systematic bug patterns identified in the QA report.

The Go-based AST scanner (`TestAntipattern_*` in `internal/provider/antipattern_test.go`) provides:
- **Deeper analysis**: AST-based pattern matching vs. grep
- **Automated enforcement**: Runs in CI/CD pipeline
- **Precise categorization**: Distinguishes field types and risk levels

**Complementary use**:
1. Run shell scripts for **quick local verification** during development
2. Run Go scanner in **CI/CD** for automated enforcement
3. Use both for **cross-validation** of findings

---

## QA Report Reference

Source: `/Users/michelvaillant/Downloads/terraform-provider-zstack-测试架构分析-20260422.md`

- **Section VI**: Full script definitions and usage
- **Section IV**: Three systematic bug patterns (模式 1-3)
- **Appendix I**: Detailed pattern analysis with affected resources
