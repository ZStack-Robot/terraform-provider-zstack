# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.20.14.21"
  port   = "8080"
  accountname     = "admin"
  accountpassword = "password"
 # accesskeyid     = "fOrgOdYb1j871Jmn6kWf"
 # accesskeysecret = "SOH8LzoCGEBrbyH4YTRTE8qppgqYT43XKbkjx03I"
}

data "zstack_hosts" "hosts" {
}



output "zstack_hosts" {
  value = data.zstack_hosts.hosts
}


