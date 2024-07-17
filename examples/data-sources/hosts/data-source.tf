# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.x.x.x"
  port   = "8080"  #optional, default is 8080
  accesskeyid     = "accesskeyid"
  accesskeysecret = "accesskeysecret"
}

data "zstack_hosts" "hosts" {
 #  name_regex = "lc-2"
}



output "zstack_hosts" {
  value = data.zstack_hosts.hosts
}


