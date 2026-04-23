# Antipatterns Test Fixtures

This directory contains test fixtures for the AST scanner self-tests, validating that the provider's static analysis checks correctly identify and flag antipatterns.

## Directory Layout

```
antipatterns/
├── check_2a/
│   ├── bad/      # Fixtures that MUST be flagged by check_2a
│   └── good/     # Fixtures that MUST NOT be flagged by check_2a
├── check_2b/
│   ├── bad/      # Fixtures that MUST be flagged by check_2b
│   └── good/     # Fixtures that MUST NOT be flagged by check_2b
└── check_2d/
    ├── bad/      # Fixtures that MUST be flagged by check_2d
    └── good/     # Fixtures that MUST NOT be flagged by check_2d
```

## File Naming Convention

All fixture files use the `.go.fixture` extension (NOT `.go`) to prevent compilation by `go build ./...`.

Example: `bad_example.go.fixture`

This ensures:
- Fixtures are not compiled as part of the provider build
- Fixtures can contain intentionally invalid or incomplete code
- The test harness can parse and analyze fixtures independently

## Purpose

These fixtures enable the AST scanner to validate that:
1. **Bad fixtures** are correctly identified and flagged by their respective checks
2. **Good fixtures** pass validation without false positives
3. The scanner maintains accuracy across diverse code patterns
