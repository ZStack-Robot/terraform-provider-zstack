# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host            = "172.24.248.201"
  accountname     = "admin"
  accountpassword = "password"
  accesskeyid     = "t8CRrQMcddB1jLO0l45r"
  accesskeysecret = "DIwdwk8JjxZrDhMsRnoF5PhcgYR0mcu2841DWPSb"
}

provider "local" {

}

provider "template" {
  # Configuration options
}
data "zstack_images" "images" {
  name_regex = "centos7-test.qcow2.no-qemu-ga"
}

data "zstack_l3network" "networks" {
  name_regex = "L3-1"
}

resource "zstack_vm" "vm" {
  count          = 2
  name           = "disk-test-${count.index + 1}"
  description    = "chi test"
  imageuuid      = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
  #  l3networkuuids = "de7f26a7304d45aea9e9871a1ba7dbae"
  #  rootdiskofferinguuid = "e6ed934030244c7c8465975f7a23ae79"
  rootdisksize = 202400
  memorysize   = 1147483640
  cupnum       = 1
  #  cpumode = "host-passthrough"
}

#data "zstack_images" "images" {}


output "zstack_vm" {
  value = zstack_vm.vm
}

