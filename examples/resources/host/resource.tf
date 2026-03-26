# Copyright (c) ZStack.io, Inc.

# Add a KVM host to a cluster
resource "zstack_host" "example" {
  name          = "kvm-host-01"
  description   = "KVM compute host in production cluster"
  management_ip = "192.168.1.100"
  cluster_uuid  = "cluster-uuid-placeholder"
  username      = "root"
  password      = "password"
  ssh_port      = 22
}

output "zstack_host" {
  value = zstack_host.example
}
