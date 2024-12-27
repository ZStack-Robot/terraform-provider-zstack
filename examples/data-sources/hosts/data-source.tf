# Copyright (c) ZStack.io, Inc.

data "zstack_hosts" "example" {
  #  name = "hostname"
  #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "architecture"
    values = ["aarch64", "x86_64"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "status"
    values = ["Disconnected"]
  } 
  filter {
    name   = "cluster_uuid"
    values = ["37c25209578c495ca176f60ad0cd97fa"]
  } 
}

output "zstack_hosts" {
  value = data.zstack_hosts.example
}