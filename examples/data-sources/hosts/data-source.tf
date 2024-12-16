# # Copyright (c) ZStack.io, Inc.

data "zstack_hosts" "example" {
  #  name = "hostname"
  #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}

output "zstack_hosts" {
  value = data.zstack_hosts.example
}