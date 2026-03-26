# Copyright (c) ZStack.io, Inc.

# Create a cluster within a zone
resource "zstack_cluster" "example" {
  name            = "production-cluster"
  description     = "KVM compute cluster for production workloads"
  zone_uuid       = "zone-uuid-placeholder"
  hypervisor_type = "KVM"
}

output "zstack_cluster" {
  value = zstack_cluster.example
}
