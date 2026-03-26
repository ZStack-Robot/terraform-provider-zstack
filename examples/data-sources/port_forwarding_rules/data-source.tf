# Copyright (c) ZStack.io, Inc.

data "zstack_port_forwarding_rules" "example" {
  name = "web-server-forwarding"
  # name_pattern = "web-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "protocol_type"
    values = ["TCP"]
  }
}

output "zstack_port_forwarding_rules" {
  value = data.zstack_port_forwarding_rules.example
}
