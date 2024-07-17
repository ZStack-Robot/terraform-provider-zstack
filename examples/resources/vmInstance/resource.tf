# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "IP"
  accesskeyid     = "accesskeyid"
  accesskeysecret = "accesskeysecret"
}

provider "local" {

}

provider "template" {
  # Configuration options
}
data "zstack_images" "images" {
  name_regex = "C790123newname"
}

data "zstack_l3network" "networks" {
  name_regex = "public"
}

resource "zstack_vm" "vm" {
  count          = 2
  name           = "disk-test-${count.index + 1}"
  description    = "test"
  imageuuid      = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
  rootdisksize = 202400
  memorysize   = 1147483640
  cupnum       = 1
}

output "zstack_vm" {
  value = zstack_vm.vm
}


