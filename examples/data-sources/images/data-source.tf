# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "ip"
  accesskeyid     = "accesskeyid"
  accesskeysecret = "accesskeysecret"
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
  value = data.zstack_images.images # images.0.uuid # "${data.zstack_images.images.uuid}"  # "${data.zstack_images.images.images.*.name}"    #{data.zstack_images.images.images.*.name} data.zstack_images.image

}


