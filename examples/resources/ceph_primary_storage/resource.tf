# Copyright (c) ZStack.io, Inc.

resource "zstack_ceph_primary_storage" "example" {
  name        = "ceph-ps-1"
  description = "Ceph primary storage"
  zone_uuid   = "example-zone-uuid"
  mon_urls    = ["10.0.0.1:6789", "10.0.0.2:6789", "10.0.0.3:6789"]
}

output "ceph_primary_storage" {
  value = zstack_ceph_primary_storage.example
}
