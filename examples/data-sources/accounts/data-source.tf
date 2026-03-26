# Copyright (c) ZStack.io, Inc.

data "zstack_accounts" "example" {
  name = "admin"
  # name_pattern = "admin%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "type"
    values = ["SystemAdmin"]
  }
}

output "zstack_accounts" {
  value = data.zstack_accounts.example
}
