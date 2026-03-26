# Copyright (c) ZStack.io, Inc.

data "zstack_l2vlan_networks" "example" {
  name = "my-vlan-network"
  # name_pattern = "vlan-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "vlan"
    values = ["100"]
  }
}

output "zstack_l2vlan_networks" {
  value = data.zstack_l2vlan_networks.example
}
