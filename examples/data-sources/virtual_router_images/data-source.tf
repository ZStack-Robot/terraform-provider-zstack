# Copyright (c) ZStack.io, Inc.

data "zstack_virtual_router_images" "test" {
  #   name = "name of virtual router images"
  #    name_pattern = "virtual router images name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name = "guest_os_type"
    values = ["VyOS 1.1.7"]
  } 
  filter {
    name = "status"
    values = ["Ready"]
  } 
}

output "zstack_vrimages" {
  value = data.zstack_virtual_router_images.test
}




