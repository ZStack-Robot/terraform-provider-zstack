# Copyright (c) ZStack.io, Inc.

data "zstack_zone" "zones" {
  name_regex = "ZONE-1"
  # name_pattern = "zone name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_zone" {
  value = data.zstack_zone.zones
}


