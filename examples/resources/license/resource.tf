# Copyright (c) ZStack.io, Inc.

resource "zstack_license" "example" {
  license              = "your-license-text-here"
  management_node_uuid = "management-node-uuid"
}
