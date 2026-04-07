# Copyright (c) ZStack.io, Inc.

resource "zstack_ceph_pool" "example" {
  pool_name            = "example-ceph-pool"
  primary_storage_uuid = "5d09ee200ebf450b9c962ce2082a64f8"
  type                 = "replicated"
  alias_name           = "example-alias"
  description          = "Example Ceph Pool"
  is_create            = true
}

output "ceph_pool" {
  value = zstack_ceph_pool.example
}
