# Create a security group first (assuming you have a security group resource)
resource "zstack_networking_secgroup" "linuxbridge_test" {
  name         = "tf-secgroup-linuxbridge-for-rules"
  description  = "Test Security Group with LinuxBridge"
  vswitch_type = "LinuxBridge"
  ip_version   = 4
}

# Create an ingress rule
resource "zstack_networking_secgroup_rule" "ingress_http" {
  name                    = "allow-http-ingress"
  security_group_uuid     = zstack_networking_secgroup.linuxbridge_test.uuid
  direction               = "Ingress"
  action                  = "ACCEPT"
  protocol                = "TCP"
  priority                = 2
  ip_version              = 4
  ip_ranges               = "10.0.1.0/8,30.1.1.0/24"
  destination_port_ranges = "80"
  description             = "Allow HTTP traffic from anywhere"
  state                   = "Disabled"
}

# Create an egress rule with a specific port range
resource "zstack_networking_secgroup_rule" "egress_custom" {
  name                    = "allow-custom-egress"
  security_group_uuid     = zstack_networking_secgroup.linuxbridge_test.uuid
  direction               = "Egress"
  action                  = "ACCEPT"
  protocol                = "TCP"
  priority                = 1
  ip_version              = 4
  ip_ranges               = "192.168.0.0/16"
  destination_port_ranges = "8000-9000"
  description             = "Allow custom port range traffic to internal network"
  state                   = "Enabled"
}

# Create a rule for ICMP protocol
resource "zstack_networking_secgroup_rule" "icmp_rule" {
  depends_on          = [zstack_networking_secgroup_rule.ingress_http]
  name                = "allow-icmp"
  security_group_uuid = zstack_networking_secgroup.linuxbridge_test.uuid
  direction           = "Ingress"
  action              = "ACCEPT"
  protocol            = "ICMP"
  priority            = 1
  ip_version          = 4
  ip_ranges           = "20.0.0.0/0"
  description         = "Allow ICMP traffic (ping) from anywhere"
  state               = "Enabled"
}


# Output the created resources
output "ingress_http_rule_id" {
  value = zstack_networking_secgroup_rule.ingress_http
}

output "egress_custom_rule_id" {
  value = zstack_networking_secgroup_rule.egress_custom
}

output "icmp_rule_id" {
  value = zstack_networking_secgroup_rule.icmp_rule
}