# Copyright (c) ZStack.io, Inc.

data "zstack_virtual_routers" "test" {
  #   name = "name of vm instance"
  #    name_pattern = "virtual router instances name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
    filter = {  # option
      State = "Running"
      Status = "Connected"
      HypervisorType = "KVM"
      Architecture = "x86_s64"
   }   
}


output "zstack_vrouters" {
  value = data.zstack_virtual_routers.test
}



