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


data "zstack_zones" "zones" {
   name_regex = "ZONE-1"
}



output "zstack_zones" {
  value = data.zstack_zones.zones
}


