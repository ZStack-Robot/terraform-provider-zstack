# Copyright (c) ZStack.io, Inc.

resource "zstack_vpc_ha_group" "example" {
  name        = "example-vpc-ha-group"
  description = "Example VPC HA Group"
  monitor_ips = ["192.168.1.1", "192.168.1.2"]
}

output "vpc_ha_group" {
  value = zstack_vpc_ha_group.example
}
