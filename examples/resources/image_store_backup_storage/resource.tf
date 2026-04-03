# Copyright (c) ZStack.io, Inc.

resource "zstack_image_store_backup_storage" "example" {
  name        = "imagestore-bs-1"
  description = "Image store backup storage"
  hostname    = "192.168.1.50"
  username    = "root"
  password    = "password"
  url         = "/zstack_bs"
}

output "image_store_backup_storage" {
  value = zstack_image_store_backup_storage.example
}
