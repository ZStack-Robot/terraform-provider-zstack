# Copyright (c) ZStack.io, Inc.

data "zstack_l3networks" "networks" {
  #   name = "L3 networks name"
  #    name_pattern = "L3 networks name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter = {
    Category  = "Private"
    IpVersion = 6
  }
}

output "zstack_networks" {
  value = data.zstack_l3networks.networks
}

