# Copyright (c) ZStack.io, Inc.

data "zstack_virtual_router_offers" "test" {
  #   name = "name of virtual router offers"
  #    name_pattern = "virtual router offers name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}

output "zstack_offers" {
  value = data.zstack_virtual_router_offers.test
}



