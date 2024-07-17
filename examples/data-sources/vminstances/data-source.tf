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

data "zstack_vminstances" "vminstances" {}



output "zstack_vminstances" {
  value = data.zstack_vminstances.vminstances
}



