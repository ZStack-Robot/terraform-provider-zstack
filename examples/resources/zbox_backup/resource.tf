# Copyright (c) ZStack.io, Inc.

resource "zstack_zbox_backup" "example" {
  name      = "example-zbox-backup"
  zbox_uuid = "example-zbox-uuid"
}

output "zstack_zbox_backup" {
  value = zstack_zbox_backup.example
}
