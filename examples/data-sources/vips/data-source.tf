# Copyright (c) ZStack.io, Inc.

data "zstack_vips" "test" {
  filter {
    name = "use_for"
    values = ["LoadBalancer"]
  } 
}

output "zstack_vips" {
  value = data.zstack_vips.test
}




