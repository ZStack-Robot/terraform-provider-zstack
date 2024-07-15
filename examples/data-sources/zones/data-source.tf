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


data "zstack_zones" "zones" {
   name_regex = "ZONE-1"
}



output "zstack_zones" {
  value = data.zstack_zones.zones
}


