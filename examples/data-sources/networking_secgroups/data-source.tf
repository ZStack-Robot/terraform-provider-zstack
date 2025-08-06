data "zstack_networking_secgroups" "test" {
  #name = "p1"
  name_pattern = "p%"
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_secs" {
  value = data.zstack_networking_secgroups.test
}