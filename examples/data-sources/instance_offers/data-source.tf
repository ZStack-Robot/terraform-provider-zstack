# Copyright (c) ZStack.io, Inc.

data "zstack_instance_offers" "example" {
  # name = "InstanceOffering-1"
  # name_pattern = "clu%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter = {  # option
    State = "Enabled"
  }
}

output "zstack_instance_offers" {
  value = data.zstack_instance_offers.example
}


