# Copyright (c) ZStack.io, Inc.

# Create a port forwarding rule mapping VIP port 8080 to a VM's private port 80
resource "zstack_port_forwarding_rule" "example" {
  name               = "web-server-forwarding"
  description        = "Forward external port 8080 to internal web server port 80"
  vip_uuid           = "vip-uuid-placeholder"
  vip_port_start     = 8080
  vip_port_end       = 8080
  private_port_start = 80
  private_port_end   = 80
  protocol_type      = "TCP" # TCP or UDP
  allowed_cidr       = "0.0.0.0/0"

  # Optional: attach to a specific VM NIC
  # vm_nic_uuid = "vm-nic-uuid-placeholder"
}

output "zstack_port_forwarding_rule" {
  value = zstack_port_forwarding_rule.example
}
