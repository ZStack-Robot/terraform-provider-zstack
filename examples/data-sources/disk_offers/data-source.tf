# Copyright (c) ZStack.io, Inc.

data "zstack_disk_offers" "example" {
  name = "smallDiskOffering"
  # name_pattern = "sm%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "name"
    values = ["name1", "name2"]
  }
  filter {
    name   = "state"
    values = ["Enabled"]
  }
  filter {
    name   = "disk_size"
    values = ["20"] # in gigabytes (GB)
  }
}

output "zstack_disk_offers" {
  value = data.zstack_disk_offers.example
}



