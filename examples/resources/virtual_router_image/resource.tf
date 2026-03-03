# Copyright (c) ZStack.io, Inc.

data "zstack_backupstorages" "example" {
  name = "bs"
}

resource "zstack_virtual_router_image" "test" {
  name                 = "virtual router name"
  description          = "An example virtual router image"
  url                  = "http://zstack-image-server.zstack.io/zstack_iso/cloud/5.2.0/zstack-vrouter-240904-14.qcow2"
  architecture         = "x86_64"
  virtio               = true
  backup_storage_uuids = [data.zstack_backupstorages.example.backup_storages.0.uuid]
  boot_mode            = "legacy"
  guest_os_type        = "VyOS 1.1.7" #Attribute guest_os_type value must be one of: ["VyOS 1.1.7" "openEuler 22.03"]
}

output "zstack_vrimage" {
  value = zstack_virtual_router_image.test
}