# Copyright (c) ZStack.io, Inc.

resource "zstack_ceph_backup_storage" "example" {
  name        = "ceph-bs-1"
  description = "Ceph backup storage"
  mon_urls    = ["10.0.0.1:6789", "10.0.0.2:6789", "10.0.0.3:6789"]
}

output "ceph_backup_storage" {
  value = zstack_ceph_backup_storage.example
}
