# Copyright (c) ZStack.io, Inc.

data "zstack_instances" "vminstances" {
  #   name = "name of vm instance"
  #    name_pattern = "virtual machine instances name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "architecture"
    values = ["aarch64", "x86_64"]
  }
  filter {
    name   = "state"
    values = ["Running", "Stopped"]
  }
  filter {
    name   = "hypervisor_type"
    values = ["KVM"]
  }
}


output "zstack_vminstances" {
  value = data.zstack_instances.vminstances
}



