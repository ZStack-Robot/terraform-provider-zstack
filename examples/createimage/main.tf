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

resource "zstack_image" "image" {
  count = 2
  name = "C790123newname"
#  description = "chi test"
   url = "http://storage.zstack.io/mirror/nightly/diskimages/CentOS7.9.qcow2"
#  guestostype = "Linux"
#  platform = "Linux"
#  type = ""
#  architecture = "x86_64"
#  virtio = false
}

#data "zstack_images" "images" {}



output "zstack_image" {
   value =  zstack_image.image
}


