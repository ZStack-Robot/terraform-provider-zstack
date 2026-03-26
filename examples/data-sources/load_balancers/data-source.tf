# Copyright (c) ZStack.io, Inc.

data "zstack_load_balancers" "example" {
  name = "my-load-balancer"
  # name_pattern = "web-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_load_balancers" {
  value = data.zstack_load_balancers.example
}
