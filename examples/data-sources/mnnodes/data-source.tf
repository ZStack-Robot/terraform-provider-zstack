# Copyright (c) ZStack.io Inc.

data "zstack_mnnodes" "hosts" {

}

output "zstack_mnnodes" {
  value = data.zstack_mnnodes.hosts
}


