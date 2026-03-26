# Copyright (c) ZStack.io, Inc.

# Add an ImageStore backup storage
resource "zstack_backup_storage" "imagestore_example" {
  name     = "imagestore-bs-01"
  type     = "ImageStoreBackupStorage"
  hostname = "192.168.1.200"
  username = "root"
  password = "password"
  url      = "/zstack_bs"

  attached_zone_uuids = ["zone-uuid-placeholder"]
}

# Add a Ceph backup storage
resource "zstack_backup_storage" "ceph_example" {
  name     = "ceph-bs-01"
  type     = "CephBackupStorage"
  mon_urls = ["10.0.0.1", "10.0.0.2", "10.0.0.3"]

  attached_zone_uuids = ["zone-uuid-placeholder"]
}

output "zstack_backup_storage" {
  value = zstack_backup_storage.imagestore_example
}
