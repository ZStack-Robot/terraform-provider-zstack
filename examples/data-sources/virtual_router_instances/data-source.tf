# Copyright (c) ZStack.io, Inc.

data "zstack_virtual_router_instances" "test" {
  #   name = "name of vm instance"
  #    name_pattern = "virtual router instances name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "hypervisor_type"
    values = ["KVM"]
  }
  filter {
    name   = "memory_size"
    values = ["1024"] # in megabytes, MB
  }
  filter {
    name   = "cpu_num"
    values = [1]
  }
}

output "zstack_vrouters" {
  value = data.zstack_virtual_router_instances.test
}



