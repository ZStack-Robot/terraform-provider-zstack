# Copyright (c) ZStack.io, Inc.

resource "zstack_volume_backup" "example" {
  name                = "example-volume-backup"
  volume_uuid         = "5d09ee200ebf450b9c962ce2082a64f8"
  backup_storage_uuid = "a1b2c3d4e5f6789a0b1c2d3e4f5a6b7c"
  description         = "Example Volume Backup"
}

output "volume_backup" {
  value = zstack_volume_backup.example
}
