data "zstack_networking_sdn_controllers" "exact_match" {
  name = "172.30.3.152"
}

data "zstack_networking_sdn_controllers" "pattern_match" {
  name_pattern = "172%"
}

data "zstack_networking_sdn_controllers" "filtered" {
  filter {
    name   = "status"
    values = ["Disconnected"]
  }
  filter {
    name   = "ip"
    values = ["172.30.3.152"]
  }
}

output "sdn_controller_exact" {
  value = data.zstack_networking_sdn_controllers.exact_match
}

output "sdn_controller_pattern" {
  value = data.zstack_networking_sdn_controllers.pattern_match
}

output "sdn_controller_filtered" {
  value = data.zstack_networking_sdn_controllers.filtered
}