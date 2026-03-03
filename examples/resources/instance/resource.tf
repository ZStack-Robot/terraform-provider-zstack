data "zstack_images" "centos" {
  name = "centos8"
}

data "zstack_l3networks" "l3networks" {
  name = "public-net"
}

data "zstack_instance_offers" "offer" {
  name = "min"
}

resource "zstack_instance" "example_vm" {
  name       = "example-v"
  image_uuid = data.zstack_images.centos.images[0].uuid
  # l3_network_uuids = [data.zstack_l3networks.l3networks.l3networks[0].uuid] # Removed use of deprecated `l3_network_uuids` in favor of `network_interfaces`
  description = "jumper server"
  #  instance_offering_uuid = data.zstack_instance_offers.offer.instance_offers[0].uuid #using Instance offering uuid or custom cpu and memory 
  network_interfaces = [
    {
      l3_network_uuid = data.zstack_l3networks.networks.l3networks.0.uuid
      # default_l3      = true
      static_ip = "172.30.3.154"
    },
    {
      l3_network_uuid = data.zstack_l3networks.vpc_net.l3networks.0.uuid
      default_l3      = true
      static_ip       = "192.168.2.20"
    }
  ]
  memory_size = 4096
  cpu_num     = 4
  expunge     = true
}

output "zstack_instance" {
  value = zstack_instance.example_vm
}