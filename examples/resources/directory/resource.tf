# Copyright (c) ZStack.io, Inc.

resource "zstack_directory" "example" {
  name      = "example-directory"
  zone_uuid = "example-zone-uuid"
  type      = "VmInstanceGroup"
}

output "zstack_directory" {
  value = zstack_directory.example
}
