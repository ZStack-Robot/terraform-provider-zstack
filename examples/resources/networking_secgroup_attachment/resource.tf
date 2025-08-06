resource "zstack_networking_secgroup_attachment" "example" {
  secgroup_uuid = "f450b20497c34397977091bc1c8f87f9"
  nic_uuid      = "a8aa88c413704717b138190832864b54"
}

output "secgroup_attachment" {
  value = zstack_networking_secgroup_attachment.example
}