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

data "zstack_images" "images" {
#   name_regex = "RDS-3.13.10"
#  images = [
#
#    {
#    name_regx = "C790123newname"
#    }
#  ]
}



output "zstack_images" {
   value =  data.zstack_images.images.images.0.uuid # "${data.zstack_images.images.uuid}"  # "${data.zstack_images.images.images.*.name}"    #{data.zstack_images.images.images.*.name} data.zstack_images.image

}


