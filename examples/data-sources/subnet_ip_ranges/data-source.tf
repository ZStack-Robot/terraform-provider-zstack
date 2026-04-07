# Copyright (c) ZStack.io, Inc.

data "zstack_subnet_ip_ranges" "all" {
}

data "zstack_subnet_ip_ranges" "filtered" {
  name_pattern = "subnet%"
}

output "subnet_ip_ranges" {
  value = data.zstack_subnet_ip_ranges.all
}
