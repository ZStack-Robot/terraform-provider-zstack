# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "test-qa.zstack.io"
  accesskeyid     = "46Ca208N3yLLn1JfkQR3"
  accesskeysecret = "JjdhscQnoN1uqS7VmZKZGuUmYQz6rFu35fh3hvcu"
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


