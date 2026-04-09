# Copyright (c) ZStack.io, Inc.

resource "zstack_ipsec_connection" "example" {
  name            = "ipsec-conn-1"
  description     = "IPsec VPN connection"
  l3_network_uuid = "example-l3-network-uuid"
  peer_address    = "203.0.113.1"
  auth_key        = "my-secret-auth-key"
  vip_uuid        = "example-vip-uuid"
  peer_cidrs      = ["10.0.0.0/24", "10.0.1.0/24"]
}

output "ipsec_connection" {
  value = zstack_ipsec_connection.example
}
