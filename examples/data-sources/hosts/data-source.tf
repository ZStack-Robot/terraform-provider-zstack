# Copyright (c) ZStack.io, Inc.

data "zstack_hosts" "example" {
  #  name = "hostname"
  #   name_pattern = "hostname%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
    filter = {  # option
      State = "Enabled"
      Status = "Connected"
      HypervisorType = "KVM"
      Architecture = "x86_64"
      TotalCpuCapacity = "480"
      ManagementIp = "172.30.3.4"      
   }
}

output "zstack_hosts" {
  value = data.zstack_hosts.example
}