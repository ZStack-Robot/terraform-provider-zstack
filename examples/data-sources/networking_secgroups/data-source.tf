data "zstack_networking_secgroups" "test" {
  priority = 8
  filter {
    name   = "state"
    values = ["Enabled"]
  }
}

output "zstack_secs" {
  value = data.zstack_networking_secgroups.test
}