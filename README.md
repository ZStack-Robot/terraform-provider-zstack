Terraform Provider For ZStack Cloud
==================

- [Getting Started](https://registry.terraform.io/providers/ZStack-Robot/zstack/latest)
- Usage
  - Documentation
  - [Examples](https://github.com/ZStack-Robot/terraform-provider-zstack/blob/main/docs/index.md)

The ZStack provider is used to interact with the resources supported by ZStack Cloud, a powerful cloud management platform. 
This provider allows you to manage various cloud resources such as virtual machines, networks, storage, and more. 
It provides a seamless integration with Terraform, enabling you to define and manage your cloud infrastructure as code.

Supported Versions
------------------

| Terraform version | minimum provider version |maximum provider version
| ---- | ---- | ----| 
| >= 1.5.x	| 1.0.0	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.5.x
-	[Go](https://golang.org/doc/install) 1.22 (to build the provider plugin)


Building The Provider
---------------------

```bash
git clone https://github.com/ZStack-Robot/terraform-provider-zstack.git
cd terraform-provider-zstack
go build -o terraform-provider-zstack
```

Using the provider
----------------------
Please see [instructions](https://www.zstack.io) on how to configure the ZStack Cloud Provider.


## Contributing to the provider

The ZStack Provider for Terraform is the work of many contributors. We appreciate your help!

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.20+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

You may also [report an issue](https://github.com/ZStack-Robot/terraform-provider-zstack/issues/new).

---

## Developing a New Resource or Data Source

When adding a new resource or data source, follow this workflow:

### 1. Implement

- **Resource:** create `zstack/provider/resource_zstack_<name>.go`, implement CRUD, register in `provider.go`
- **Data source:** create `zstack/provider/data_source_zstack_<name>.go`, implement Read, register in `provider.go`

### 2. Write acceptance tests

Create `zstack/provider/resource_zstack_<name>_test.go` (or `data_source_zstack_<name>_test.go`).

### 3. Create the example `.tf` file (required for doc generation)

The [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs) tool reads example files from the `examples/` directory. **You must create the example file before generating docs.**

| Type | Path | Filename |
|------|------|----------|
| Resource | `examples/resources/<name>/` | `resource.tf` |
| Data source | `examples/data-sources/<name>/` | `data-source.tf` |

Where `<name>` is the Terraform type name **without** the `zstack_` prefix. For example:

| Terraform type | Example file path |
|---|---|
| `zstack_disk_offer` (resource) | `examples/resources/disk_offer/resource.tf` |
| `zstack_disk_offers` (data source) | `examples/data-sources/disk_offers/data-source.tf` |
| `zstack_virtual_router_image` (resource) | `examples/resources/virtual_router_image/resource.tf` |

Each example file should start with `# Copyright (c) ZStack.io, Inc.` and contain a realistic, working HCL configuration.

### 4. Generate documentation

```bash
cd tools && go generate ./...
```

This pipeline (defined in `tools/tools.go`) runs three steps:

1. **Copyright headers** — adds copyright headers to all source files
2. **Terraform fmt** — formats all `.tf` files under `examples/`
3. **tfplugindocs generate** — reads the provider schema + example files and generates Markdown docs

Generated docs are written to:
- `docs/resources/<name>.md`
- `docs/data-sources/<name>.md`

### 5. Verify

```bash
go build ./...                                                           # compiles
TF_ACC=1 go test ./zstack/provider/ -v -run TestAcc<Name> -timeout 10m  # acceptance test
ls docs/resources/<name>.md docs/data-sources/<name>.md                  # doc generated
```

---

## Testing

This project provides three testing approaches: Go acceptance tests, Terraform batch integration tests, and ad-hoc single-resource tests.

### Prerequisites: Environment Variables

All testing approaches require the following environment variables to connect to a real ZStack environment:

```bash
export ZSTACK_HOST=172.24.227.46          # ZStack management node address
export ZSTACK_PORT=8080                   # Port (default: 8080)
export ZSTACK_ACCESS_KEY_ID=<your-key>    # Access Key ID
export ZSTACK_ACCESS_KEY_SECRET=<secret>  # Access Key Secret
```

> **Tip:** Store these in a `.env.test` file at the project root and load them with `source .env.test`.

---

### Approach 1: Go Acceptance Tests

Acceptance tests live in `zstack/provider/*_test.go` and use the Go testing framework with the Terraform SDK test harness. They create, read, and destroy real resources against a live ZStack environment.

#### 1. Generate the environment snapshot

Tests rely on `testdata/env.json`, which contains real resource UUIDs and names from your environment. Generate (or regenerate) it whenever the environment changes:

```bash
source .env.test
go run ./zstack/provider/testdata/generate_env.go
```

This writes to `zstack/provider/testdata/env.json`.

#### 2. Run all acceptance tests

```bash
source .env.test
TF_ACC=1 go test ./zstack/provider/ -v -timeout 30m 2>&1 | tee gotest.out
```

#### 3. Run tests for a single resource

```bash
# Example: test disk offering only
TF_ACC=1 go test ./zstack/provider/ -v -run TestAccDiskOffer -timeout 10m

# Example: test images data source only
TF_ACC=1 go test ./zstack/provider/ -v -run TestAccDataSourceImages -timeout 10m
```

#### 4. Understanding test results

- **PASS** — Test passed.
- **SKIP** — Skipped because the environment lacks the required resources (e.g., no VM instances → instance-related tests are skipped).
- **FAIL** — Test failed; investigate the output.

---

### Approach 2: Terraform Batch Integration Tests

This approach bypasses the Go test framework entirely. It generates standalone, runnable `.tf` files and tests each resource type with `terraform apply/destroy` directly. This is useful for:

- Quickly verifying provider compatibility against a real environment
- Debugging apply/destroy behavior for a specific resource
- CI/CD batch regression testing

#### 1. Generate the environment snapshot (same as above)

```bash
source .env.test
go run ./zstack/provider/testdata/generate_env.go
```

#### 2. Build the provider binary

```bash
go build -o $(go env GOPATH)/bin/terraform-provider-zstack
```

#### 3. Generate the .tf files

```bash
go run ./zstack/provider/testdata/generate_tf.go
```

Output is written to `zstack/provider/testdata/terraform/` with this structure:

```
zstack/provider/testdata/terraform/
├── dev.tfrc                                 # Terraform CLI config (dev_overrides → local binary)
├── provider.tf                              # Shared provider configuration (reads env vars)
├── data-images/
│   ├── main.tf                              # Data source query
│   └── provider.tf -> ../provider.tf        # Symlink
├── res-disk_offer/
│   ├── main.tf                              # Resource create/destroy
│   └── provider.tf -> ../provider.tf
├── res-instance/
│   ├── main.tf
│   └── provider.tf -> ../provider.tf
└── ...  (30+ subdirectories)
```

The generated `dev.tfrc` uses [`dev_overrides`](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers) to tell Terraform to use the locally compiled binary at `$(go env GOPATH)/bin/` instead of downloading from a registry. With `dev_overrides`, `terraform init` is not required.

The provider reads credentials directly from environment variables (`ZSTACK_HOST`, `ZSTACK_ACCESS_KEY_ID`, `ZSTACK_ACCESS_KEY_SECRET`), so no HCL variables or `TF_VAR_*` exports are needed.

#### 4. Test a single resource

```bash
source .env.test
export TF_CLI_CONFIG_FILE=$(pwd)/zstack/provider/testdata/terraform/dev.tfrc

cd zstack/provider/testdata/terraform/res-disk_offer
terraform apply -auto-approve
terraform destroy -auto-approve
```

#### 5. Batch-test all resources

```bash
source .env.test
export TF_CLI_CONFIG_FILE=$(pwd)/zstack/provider/testdata/terraform/dev.tfrc

for dir in zstack/provider/testdata/terraform/*/; do
  echo "=== Testing $dir ==="
  (cd "$dir" && terraform apply -auto-approve && terraform destroy -auto-approve)
done
```

#### 6. Generator categories

| Category | Prefix | Examples | Description |
|----------|--------|----------|-------------|
| Data sources | `data-` | `data-images`, `data-zones` | Read-only queries against existing resources |
| Self-contained resources | `res-` | `res-disk_offer`, `res-tag` | No environment UUIDs needed; works on any environment |
| Environment-dependent resources | `res-` | `res-image`, `res-instance` | Uses real UUIDs from `env.json` |

**Naming convention:** All generated resources use the `tf-batch-test-<type>` name prefix and `[batch-test]` description prefix for easy identification and cleanup in the ZStack console.

#### 7. Skipped resources

Some resources are skipped because they are too complex or potentially destructive:

- `vpc` — Requires VR + L2 and complex prerequisite setup
- `virtual_router_instance` — Creates a real virtual router (heavyweight)
- `volume_snapshot` — Requires an existing volume
- `guest_tool_attachment` / `script_execution` — Requires a running VM
- `tag_attachment` — Requires an existing resource UUID

---

### Approach 3: Ad-hoc Single-Resource Testing

To quickly test a specific resource, create a `.tf` file manually:

```bash
mkdir -p /tmp/zstack-test && cd /tmp/zstack-test

cat > main.tf <<'EOF'
terraform {
  required_providers {
    zstack = { source = "zstack.io/cloud/zstack" }
  }
}

provider "zstack" {}

# --- Add the resource or data source you want to test ---
data "zstack_zones" "all" {}
output "zones" { value = data.zstack_zones.all }
EOF

# Build (if not already done) and run:
go build -o $(go env GOPATH)/bin/terraform-provider-zstack
source .env.test
export TF_CLI_CONFIG_FILE=$(pwd)/zstack/provider/testdata/terraform/dev.tfrc
terraform apply -auto-approve
```

> **Note:** The provider reads `ZSTACK_HOST`, `ZSTACK_ACCESS_KEY_ID`, and `ZSTACK_ACCESS_KEY_SECRET` directly from environment variables, so the `provider "zstack" {}` block can be empty. The `dev.tfrc` tells Terraform to use the local binary.

---

### Test File Reference

| File | Purpose |
|------|---------|
| `zstack/provider/testdata/generate_env.go` | Queries the ZStack environment and writes an `env.json` snapshot |
| `zstack/provider/testdata/env.json` | Environment data snapshot (UUIDs, names, etc.) |
| `zstack/provider/testdata/generate_tf.go` | Reads `env.json` and generates batch-test `.tf` files |
| `zstack/provider/testdata/terraform/` | Generated `.tf` file directory (git-ignored) |
| `zstack/provider/*_test.go` | Go acceptance test source code |
