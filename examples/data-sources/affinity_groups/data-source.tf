# Copyright (c) ZStack.io, Inc.

data "zstack_affinity_groups" "example" {
  name = "my-affinity-group"
  # name_pattern = "my-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "policy"
    values = ["antiSoft"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_affinity_groups" {
  value = data.zstack_affinity_groups.example
}
