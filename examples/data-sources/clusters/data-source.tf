#  Copyright (c) ZStack.io, Inc.

data "zstack_clusters" "example" {
  #name = "cluster1"
  #name_pattern = "clu%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "architecture"
    values = ["aarch64", "x86_64"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "hypervisor_type"
    values = ["KVM"]
  }
}

output "zstack_clusters" {
  value = data.zstack_clusters.example
}