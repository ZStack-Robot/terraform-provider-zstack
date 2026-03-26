# Copyright (c) ZStack.io, Inc.

# Add a local primary storage to a zone
resource "zstack_primary_storage" "local_example" {
  name      = "local-ps-01"
  zone_uuid = "zone-uuid-placeholder"
  type      = "LocalStorage"
  url       = "/zstack_ps"

  attached_cluster_uuids = ["cluster-uuid-placeholder"]
}

# Add a Ceph primary storage
resource "zstack_primary_storage" "ceph_example" {
  name      = "ceph-ps-01"
  zone_uuid = "zone-uuid-placeholder"
  type      = "Ceph"
  mon_urls  = ["10.0.0.1", "10.0.0.2", "10.0.0.3"]
}

output "zstack_primary_storage" {
  value = zstack_primary_storage.local_example
}
