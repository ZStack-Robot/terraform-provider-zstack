---
page_title: "zstack_tags Data Source - terraform-provider-zstack"
subcategory: ""
description: |-
    Data source to retrieve User, System, or regular Tags from the ZStack environment based on name, pattern, or filters.
---

# zstack_tags (Data Source)

Data source to retrieve User, System, or regular Tags from the ZStack environment based on name, pattern, or filters.

## Example Usage

```terraform
data "zstack_tags" "all_tags" {
  #    name_pattern = "AI%" 
  tag_type = "tag"
  filter {
    name   = "color"
    values = ["purple", "yellow"]
  }
}

# Output examples
output "all_tags" {
  value = data.zstack_tags.all_tags
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `tag_type` (String) Specifies which type of tag to query: 'user', 'system', or 'tag'.

### Optional

- `filter` (Block List) Filter resources based on any field in the schema. For example, to filter by status, use `name = "status"` and `values = ["Ready"]`. (see [below for nested schema](#nestedblock--filter))
- `name` (String) Exact name of the tag to match.
- `name_pattern` (String) Pattern for fuzzy matching the tag name (supports % as wildcard, _ as single character).

### Read-Only

- `system_tags` (Attributes List) List of system tags matching the query. (see [below for nested schema](#nestedatt--system_tags))
- `tags` (Attributes List) List of regular tags matching the query. (see [below for nested schema](#nestedatt--tags))
- `user_tags` (Attributes List) List of user tags matching the query. (see [below for nested schema](#nestedatt--user_tags))

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `name` (String) Name of the field to filter by (e.g., status, state).
- `values` (Set of String) Values to filter by. Multiple values will be treated as an OR condition.


<a id="nestedatt--system_tags"></a>
### Nested Schema for `system_tags`

Read-Only:

- `inherent` (Boolean) Indicates if the tag is inherent (built-in) or user-defined.
- `resource_type` (String) Type of the resource the system tag is attached to.
- `resource_uuid` (String) UUID of the resource the system tag is attached to.
- `tag` (String) System tag value.
- `uuid` (String) Unique identifier of the system tag.


<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Read-Only:

- `color` (String) Color label assigned to the tag.
- `description` (String) Description of the tag.
- `name` (String) Name of the tag.
- `type` (String) Type of the tag.
- `uuid` (String) Unique identifier of the tag.


<a id="nestedatt--user_tags"></a>
### Nested Schema for `user_tags`

Read-Only:

- `resource_type` (String) Type of the resource the tag is attached to.
- `resource_uuid` (String) UUID of the resource the tag is attached to.
- `tag` (String) Tag value.
- `type` (String) Tag category, typically 'User'.
- `uuid` (String) Unique identifier of the user tag.



