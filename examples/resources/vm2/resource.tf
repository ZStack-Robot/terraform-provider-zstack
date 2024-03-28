terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host  =  "172.25.16.104"
  accountname = "admin"
  accountpassword = "password"
  accesskeyid = "mO6W9gzCQxsfK6OsE7dg"
  accesskeysecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms"
}

provider "local" {

}

provider "template" {
  # Configuration options
}
data "zstack_images" "images" {
    name_regex = "C79"
}

data "zstack_l3network" "networks" {
  name_regex = "public"
}

resource "zstack_vm" "vm" {
  count = 1
  name = "temp-modift-${count.index + 1}"
  description = "chi test"
  imageuuid = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
#  l3networkuuids = "de7f26a7304d45aea9e9871a1ba7dbae"
#  rootdiskofferinguuid = "e6ed934030244c7c8465975f7a23ae79"
  rootdisksize = 202400
  memorysize = 11474836480
  cupnum = 2
#  cpumode = "host-passthrough"
}

#data "zstack_images" "images" {}


output "zstack_vm" {
   value =  zstack_vm.vm
}


