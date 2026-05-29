# Copyright (c) ZStack.io, Inc.

resource "zstack_instance_state" "test" {
  vm_instance_uuid = zstack_instance.test.uuid
  state            = "Running"
}
