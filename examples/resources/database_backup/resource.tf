# Copyright (c) ZStack.io, Inc.

resource "zstack_database_backup" "example" {
  name                = "example-database-backup"
  backup_storage_uuid = "example-backup-storage-uuid"
}

output "zstack_database_backup" {
  value = zstack_database_backup.example
}
