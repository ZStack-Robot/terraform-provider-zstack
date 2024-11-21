# Copyright (c) ZStack.io, Inc.


data "zstack_instances" "vminstances" {
  #   name = "name of vm instance"
  #    name_pattern = "L3 networks name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}


output "zstack_vminstances" {
  value = data.zstack_instances.vminstances
}



