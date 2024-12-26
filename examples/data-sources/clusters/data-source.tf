#  Copyright (c) ZStack.io, Inc.

data "zstack_clusters" "example" {
  #name = "cluster1"
  #name_pattern = "clu%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter = {  # option
      State = "Enabled"
      HypervisorType = "KVM"
      Architecture = "x86_64"
   }
}

output "zstack_clusters" {
  value = data.zstack_clusters.example
}