# # Copyright (c) ZStack.io, Inc.

data "zstack_zones" "zones" {
  name_regex = "ZONE-1"
  # name_pattern = "zone name% Pattern"   # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.
}

output "zstack_zones" {
  value = data.zstack_zones.zones
}


