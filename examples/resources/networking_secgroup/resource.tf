resource "zstack_networking_secgroup" "linuxbridge_test" {
  name         = "tf-secgroup-linuxbridge"
  description  = "Test Security Group with LinuxBridge"
  vswitch_type = "LinuxBridge"
  ip_version   = 4
}

resource "zstack_networking_secgroup" "ovndpdk_test" {
  name                = "tf-secgroup-ovndpdk"
  description         = "Test Security Group with OvnDpdk"
  vswitch_type        = "OvnDpdk"
  sdn_controller_uuid = "7f6ae65fb89d468ea66dd772301b6bbb" # snd controller UUID
  ip_version          = 4
}

resource "zstack_networking_secgroup" "linuxbridge_defautvstype_test" {
  name        = "tf-secgroup-linuxbridge"
  description = "Test Security Group with LinuxBridge"
  ip_version  = 4
}

output "linuxbridge_uuid" {
  value = zstack_networking_secgroup.linuxbridge_test.uuid
}

output "ovndpdk_uuid" {
  value = zstack_networking_secgroup.ovndpdk_test.uuid
}

output "defaut_vstype_ovndpdk_uuid" {
  value = zstack_networking_secgroup.linuxbridge_defautvstype_test.uuid
}