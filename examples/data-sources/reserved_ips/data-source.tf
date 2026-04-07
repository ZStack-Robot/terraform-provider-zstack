# Copyright (c) ZStack.io, Inc.

data "zstack_reserved_ips" "all" {
}

output "reserved_ips" {
  value = data.zstack_reserved_ips.all
}
