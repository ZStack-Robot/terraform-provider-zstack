# Copyright (c) ZStack.io, Inc.

# Configure the ZStack provider with the necessary credentials and endpoint information.
# - `host`: The IP address or domain name of the ZStack Cloud API endpoint.
# - `access_key_id`: The Access Key ID for authenticating with ZStack Cloud.
# - `access_key_secret`: The Access Key Secret for authenticating with ZStack Cloud.

provider "zstack" {
  host              = "ip address of zstack cloud api endpoint"
  access_key_id     = "access_key_id of zstack cloud"
  access_key_secret = "access_key_secret of zstack cloud"
}

# Fetch the details of an image from ZStack Cloud by its name.
# - `name`: The name of the image to retrieve.
data "zstack_images" "centos" {
  name = "Image-1"
}

# Fetch the details of an L3 network from ZStack Cloud by its name.
# - `name`: The name of the L3 network to retrieve.
data "zstack_l3networks" "l3networks" {
  name = "L3Network-1"
}

# Fetch the details of an instance offering from ZStack Cloud by its name.
# - `name`: The name of the instance offering to retrieve.
data "zstack_instance_offers" "offer" {
  name = "InstanceOffering-1"
}

# Create a new virtual machine instance in ZStack Cloud.
# - `name`: The name of the virtual machine.
# - `image_uuid`: The UUID of the image to use for the virtual machine.
# - `l3_network_uuids`: A list of L3 network UUIDs to attach to the virtual machine.
# - `description`: A description of the virtual machine.
# - `instance_offering_uuid`: The UUID of the instance offering to use for the virtual machine.
# - `memory_size`: (Optional) The memory size in bytes. If not specified, the instance offering's memory size will be used.
# - `cpu_num`: (Optional) The number of CPUs. If not specified, the instance offering's CPU count will be used.
# - `never_stop`: If set to `true`, the virtual machine will never be stopped.
# - `root_disk`: Configuration for the root disk of the virtual machine.
resource "zstack_instance" "example_vm" {
  name       = "moexample-v"
  image_uuid = data.zstack_images.centos.images[0].uuid
  #  l3_network_uuids       = [data.zstack_l3networks.l3networks.l3networks[0].uuid]
  description            = "Example VM instance"
  instance_offering_uuid = data.zstack_instance_offers.offer.instance_offers[0].uuid # Using Instance offering UUID or custom CPU and memory
  # memory_size            = 1024000000
  # cpu_num                = 1

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

  never_stop = true
  root_disk = {
    size = 123456788
  }
}

# Output the details of the created virtual machine instance.
output "zstack_instance" {
  value = zstack_instance.example_vm
}