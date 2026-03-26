# Copyright (c) ZStack.io, Inc.

data "zstack_load_balancer_listeners" "example" {
  name = "http-listener"
  # name_pattern = "http-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "protocol"
    values = ["tcp"]
  }
}

output "zstack_load_balancer_listeners" {
  value = data.zstack_load_balancer_listeners.example
}
