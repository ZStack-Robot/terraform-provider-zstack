terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host  =  "172.25.16.104"
  accountname = "admin"
  accountpassword = "password"
  accesskeyid = "mO6W9gzCQxsfK6OsE7dg"
  accesskeysecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms"
}

data "zstack_backupstorage" "imagestorage" {
  name_regex = "image"
}

resource "zstack_image" "image" {
  count = 3
  name = "C790123newname"
  description = "chi test"
  url = "file:///opt/zstack-dvd/zstack-image-1.4.qcow2"
  #url = "http://172.20.1.130:8001/imagestore/download/C79-image-a1c72a746afa809ff4a2214b6a5595f4da59ce97.raw"
  # url = "http://storage.zstack.io/mirror/nightly/diskimages/CentOS7.9.qcow2"
  guestostype = "Linux"
  platform = "Linux"
  format = "qcow2"
#  type = ""
  architecture = "x86_64"
  virtio = false
   backupstorageuuid = data.zstack_backupstorage.imagestorage.backupstorages.0.uuid
}

#data "zstack_images" "images" {}



output "zstack_image" {
   value =  zstack_image.image
}


