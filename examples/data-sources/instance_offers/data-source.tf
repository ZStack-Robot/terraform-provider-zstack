# Copyright (c) ZStack.io, Inc.

data "zstack_instance_offers" "example" {
  # name = "InstanceOffering-1"
  # name_pattern = "clu%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "allocator_strategy"
    values = ["LeastVmPreferredHostAllocatorStrategy"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "memory_size"
    values = ["1073741824"]
  }
  filter {
    name   = "cpu_num"
    values = [1]
  }
}

output "zstack_instance_offers" {
  value = data.zstack_instance_offers.example
}


