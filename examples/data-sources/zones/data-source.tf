# Copyright (c) ZStack.io Inc.

data "zstack_zones" "zones" {
   name_regex = "ZONE-1"
}



output "zstack_zones" {
  value = data.zstack_zones.zones
}


