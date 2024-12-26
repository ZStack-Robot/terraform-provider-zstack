# Copyright (c) ZStack.io, Inc.

data "zstack_vips" "test" {
  filter = {  # option
      State = "Enabled"
      Ip = "172.30.9.47"
      UseFor = "LoadBalancer"
   }  
}

output "zstack_vips" {
  value = data.zstack_vips.test
}




