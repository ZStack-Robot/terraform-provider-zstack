data "zstack_networking_secgroup_rules" "example" {
  # priority = 8
  filter {
    name   = "action"
    values = ["DROP"]
  }
  filter {
    name   = "protocol"
    values = ["TCP"]
  }
}

output "secgroup_rules" {
  value = data.zstack_networking_secgroup_rules.example.rules
}