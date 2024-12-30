# Copyright (c) ZStack.io, Inc.

data "zstack_l2networks" "networks" {
  #   name = "L2 networks name"
  #   name_pattern = "L2 networks name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "vlan"
    values = [36]
  }
  filter {
    name   = "physical_interface"
    values = ["ens29f1"]
  }
}

output "zstack_l2networks" {
  value = data.zstack_l2networks.networks
}

