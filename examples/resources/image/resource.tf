# Copyright (c) ZStack.io, Inc.

data "zstack_backupstorages" "example" {
  name = "name of image Storage"
}

resource "zstack_image" "image" {
  count       = 3
  name        = "name of image"
  description = "description of image"
  url         = "file:///opt/zstack-dvd/zstack-image-1.4.qcow2"
  #url = "http://url/image.raw"
  # url = "http://url/image.qcow2"
  guest_os_type        = "Linux"
  platform             = "Linux"
  format               = "qcow2"
  architecture         = "x86_64"
  virtio               = false
  backup_storage_uuids = [data.zstack_backupstorages.example.backup_storages.0.uuid]
  boot_mode            = "legacy"
  expunge              = true
}

output "zstack_image" {
  value = zstack_image.image
}


