# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.x.x.x"  #ip address for zstack cloud api endpoint
  accountname     = "admin"   
  accountpassword = "password"
  accesskeyid     = "accesskeyid"
  accesskeysecret = "accesskeysecret"
}


data "zstack_backupstorage" "imagestorage" {
  #  name_regex = "image"
}



output "zstack_imagestorage" {
  value = data.zstack_backupstorage.imagestorage
}


