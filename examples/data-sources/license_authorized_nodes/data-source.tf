# Copyright (c) ZStack.io, Inc.

data "zstack_license_authorized_nodes" "all" {
}

output "license_authorized_nodes" {
  value = data.zstack_license_authorized_nodes.all.nodes
}
