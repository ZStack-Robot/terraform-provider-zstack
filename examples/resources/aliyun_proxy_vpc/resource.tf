# Copyright (c) ZStack.io, Inc.

resource "zstack_aliyun_proxy_vpc" "example" {
  name         = "example-aliyun-proxy-vpc"
  cidr_block   = "10.0.0.0/16"
  vrouter_uuid = "example-uuid-placeholder"
}

output "zstack_aliyun_proxy_vpc" {
  value = zstack_aliyun_proxy_vpc.example
}
