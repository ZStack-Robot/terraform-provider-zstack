# Copyright (c) ZStack.io, Inc.

# Create an L2 VLAN network with a specific VLAN ID and attach it to a cluster
resource "zstack_l2vlan_network" "example" {
  name               = "my-vlan-network"
  description        = "Production VLAN network for application tier"
  vlan               = 100
  zone_uuid          = "zone-uuid-placeholder"
  physical_interface = "eth0"

  # Optional: attach to one or more clusters
  attached_cluster_uuids = ["cluster-uuid-placeholder"]
}

output "zstack_l2vlan_network" {
  value = zstack_l2vlan_network.example
}
