# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.24.244.20"
  accountname     = "admin"
  accountpassword = "password"
  accesskeyid     = "a4vS7kyX2TIaWCF1tq5S"
  accesskeysecret = "tlBIu7e7TvsYdO33r0AHVy7OowoftDMUrlhb2ALQ"
}


data "zstack_images" "images" {
  name_regex = "Image-1"
}

data "zstack_l3network" "networks" {
  name_regex = "q"
}

resource "zstack_vm" "vm" {
  count          = 5
  name           = "apilog-${count.index + 1}"
  description    = "5.1.8 test"
  imageuuid      = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
  rootdisksize = 202400
  memorysize   = 148903751
  cupnum       = 1
  #  cpumode = "host-passthrough"
}

#data "zstack_images" "images" {}


output "zstack_vm" {
  value = zstack_vm.vm
}


