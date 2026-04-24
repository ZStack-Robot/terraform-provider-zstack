# Copyright (c) ZStack.io, Inc.

data "zstack_mn_nodes" "hosts" {

}

output "zstack_mn_nodes" {
  value = data.zstack_mn_nodes.hosts
}


