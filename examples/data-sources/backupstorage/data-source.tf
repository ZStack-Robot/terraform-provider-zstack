# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.25.16.104"
  accountname     = "admin"
  accountpassword = "password"
  accesskeyid     = "mO6W9gzCQxsfK6OsE7dg"
  accesskeysecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms"
}


data "zstack_backupstorage" "imagestorage" {
  #  name_regex = "image"
}



output "zstack_imagestorage" {
  value = data.zstack_backupstorage.imagestorage
}


