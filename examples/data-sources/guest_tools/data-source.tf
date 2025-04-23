data "zstack_guest_tools" "example" {
  instance_uuid = "uuid of vm instance"
}

output "guest_tools" {
  value = data.zstack_guest_tools.example
}