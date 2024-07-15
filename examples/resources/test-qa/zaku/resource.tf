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

data "zstack_zsclusters" "clusters" {
  name_regex = "Cluster-1"
}


data "zstack_images" "images" {
  name_regex = "zaku-3.1.0"
}

data "zstack_l3network" "networks" {
  name_regex = "L3-公有网络"
}

resource "zstack_vm" "vm" {

#  count          = 3
  name           = "zaku"
  description    = "zaku testi mo"
  imageuuid      = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  clusteruuid = data.zstack_zsclusters.clusters.clusters.0.uuid
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
  zoneuuid = data.zstack_zones.zones.zones.0.uuid
#  rootdiskofferinguuid = "2bddfda17a794fea8d9d66a8d1c7ec54"

#  rootdisksize = 60720000
  memorysize   = 24890375100
  cupnum       = 8
  #  cpumode = "host-passthrough"
}

#data "zstack_images" "images" {}


output "zstack_images" {
  value = zstack_vm.vm
}


