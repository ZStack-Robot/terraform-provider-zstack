resource "zstack_l2vlan_network" "example" {
  name               = "example-l2-vlan"
  vlan               = 100
  zone_uuid          = "zone-uuid"
  physical_interface = "eth0"
}

resource "zstack_l2_network_cluster_attachment" "example" {
  l2_network_uuid = zstack_l2vlan_network.example.uuid
  cluster_uuid    = "cluster-uuid"
}

output "zstack_l2_network_cluster_attachment" {
  value = zstack_l2_network_cluster_attachment.example
}
