# Copyright (c) ZStack.io, Inc.

resource "zstack_guest_tool_attachment" "test" {
  instance_uuid = "uuid of vm instance"
}

output "zstack_guest_tool_attachment" {
  value = zstack_guest_tool_attachment.test
}