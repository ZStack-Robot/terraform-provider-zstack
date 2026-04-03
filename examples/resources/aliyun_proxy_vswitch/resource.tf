# Copyright (c) ZStack.io, Inc.

resource "zstack_aliyun_proxy_vswitch" "example" {
  aliyun_proxy_vpc_uuid = "example-uuid-placeholder"
  vpc_l3_network_uuid   = "example-uuid-placeholder"
}

output "zstack_aliyun_proxy_vswitch" {
  value = zstack_aliyun_proxy_vswitch.example
}
