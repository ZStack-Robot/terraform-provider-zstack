---
page_title: "zstack_primary_storages Data Source - terraform-provider-zstack"
subcategory: ""
description: |-
    List all primary storages, or query primary storages by exact name match, or query primary storages by name pattern fuzzy match.
---

# zstack_primary_storages (Data Source)

List all primary storages, or query primary storages by exact name match, or query primary storages by name pattern fuzzy match.

## Example Usage

```terraform
data "zstack_primary_storages" "primary_storages_test" {
  filter {
    name   = "status"
    values = ["Connected"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "total_capacity"
    values = ["7937400246271"]
  }
}

output "zstack_primary_storages_test" {
  value = data.zstack_primary_storages.primary_storages_test
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `filter` (Block List) Filter resources based on any field in the schema. For example, to filter by status, use `name = "status"` and `values = ["Ready"]`. (see [below for nested schema](#nestedblock--filter))
- `name` (String) Exact name for searching primary storage.
- `name_pattern` (String) Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

### Read-Only

- `primary_storages` (Attributes List) List of primary storage entries (see [below for nested schema](#nestedatt--primary_storages))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Name of the field to filter by (e.g., status, state).
- `values` (Set of String) Values to filter by. Multiple values will be treated as an OR condition.


<a id="nestedatt--primary_storages"></a>
### Nested Schema for `primary_storages`

Read-Only:

- `available_capacity` (Number) Available capacity of the primary storage in bytes
- `available_physical_capacity` (Number) Available physical capacity of the primary storage in bytes
- `name` (String) Name of the primary storage
- `state` (String) State of the primary storage (Enabled or Disabled)
- `status` (String) Readiness status of the primary storage
- `system_used_capacity` (Number) System used capacity of the primary storage in bytes
- `total_capacity` (Number) Total capacity of the primary storage in bytes
- `total_physical_capacity` (Number) Total physical capacity of the primary storage in bytes
- `uuid` (String) UUID identifier of the primary storage


