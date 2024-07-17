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

data "zstack_l3network" "networks" {
    name_regex = "DPortGroup-1"
}



output "zstack_networks" {
  value = data.zstack_l3network.networks #.l3networks.*.uuid
}


