# Copyright (c) ZStack.io, Inc.

data "zstack_auto_scaling_groups" "example" {
  name = "web-tier-scaling"
  # name_pattern = "web-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_auto_scaling_groups" {
  value = data.zstack_auto_scaling_groups.example
}
