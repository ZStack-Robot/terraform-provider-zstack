terraform {
  required_providers {
    zstack = {
      source = "zstack.io/terraform-provider-zstack/zstack"
    }
  }
}

provider "zstack" {
  host  =  "172.24.197.190"
  accountname = "admin"
  accountpassword = "password"
  accesskeyid = "EstXFYGGv1tWj4e9l9vx"
  accesskeysecret = "0sutxoYYCLrWFqxhsIAu7nCmohzpASvouu6pfn3m"
}


data "zstack_images" "images" {
    name_regex = "image_for_sg_test1"
}

data "zstack_l3network" "networks" {
  name_regex = "l3VlanNetwork19"
}

resource "zstack_vm" "vm" {
  count = 5
  name = "temp-modift-${count.index + 1}"
  description = "5.0 test"
  imageuuid = data.zstack_images.images.images.0.uuid #"${data.zstack_images.images.images[0].uuid}" #"9b26312501614ec0b6dc731e6977dfb2"
  l3networkuuids = data.zstack_l3network.networks.l3networks.0.uuid
#  l3networkuuids = "de7f26a7304d45aea9e9871a1ba7dbae"
#  rootdiskofferinguuid = "e6ed934030244c7c8465975f7a23ae79"
  rootdisksize = 202400
  memorysize = 5147483640
  cupnum = 1
#  cpumode = "host-passthrough"
}

#data "zstack_images" "images" {}


output "zstack_vm" {
   value =  zstack_vm.vm
}


