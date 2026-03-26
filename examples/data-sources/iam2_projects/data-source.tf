# Copyright (c) ZStack.io, Inc.

data "zstack_iam2_projects" "example" {
  name = "my-project"
  # name_pattern = "my-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_iam2_projects" {
  value = data.zstack_iam2_projects.example
}
