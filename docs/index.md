---
page_title: "ZStack Provider"
description: |-
  The ZStack provider is used to interact with the resources supported by ZStack Cloud, a powerful cloud management platform. 
  This provider allows you to manage various cloud resources such as virtual machines, networks, storage, and more. 
  It provides a seamless integration with Terraform, enabling you to define and manage your cloud infrastructure as code.
---

# zstack Provider



The ZStack provider is designed to manage resources in a ZStack Cloud environment. 
It provides a way to interact with ZStack's API to create, update, and delete resources such as virtual machines, networks, storage, and more. 
This provider is ideal for organizations looking to automate their cloud infrastructure management using Terraform.

## Example Usage

To use the ZStack provider, you need to configure it with the necessary credentials and endpoint information. Login ZStack Cloud, Operation Management -> Access Control -> AccessKey Management, Click Generate AccessKey
Below is an example configuration:

```terraform
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
data "zstack_instance_offerings" "offer" {
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
  instance_offering_uuid = data.zstack_instance_offerings.offer.instance_offers[0].uuid # Using Instance offering UUID or custom CPU and memory
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
```

## Automation / AI Usage Guide

When generating Terraform configurations programmatically (e.g., via LLMs or CI scripts), prefer the `uuid` argument over `name` / `name_pattern` on data sources:

- **Stable** — UUIDs survive renames, environment migrations, and name collisions.
- **Deterministic** — exactly 0 or 1 match; no fuzzy semantics.
- **Idempotent** — re-running an LLM-generated prompt produces an identical `.tf` file when the UUID is in scope.

The `uuid` argument is mutually exclusive with `name` / `name_pattern`. Mixing them produces a plan-time validation error.

```hcl
# Recommended for AI / automation
data "zstack_zone" "by_uuid" {
  uuid = "abc123def456..."
}

# Recommended for human-authored configs
data "zstack_zone" "by_name" {
  name = "ZONE-1"
}
```

Phase A coverage (UUID lookup support): `zstack_zone`, `zstack_clusters`, `zstack_hosts`, `zstack_l3networks`, `zstack_instances`. Remaining list data sources still accept `name` / `name_pattern` / `filter` only — Phase B will add `uuid` to the rest.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `access_key_id` (String) AccessKey ID for ZStack API. Create AccessKey ID from MN,  Operational Management->Access Control->AccessKey Management. May also be provided via ZSTACK_ACCESS_KEY_ID environment variable. Required if using AccessKey authentication. Mutually exclusive with `account_name` and `account_password`.
- `access_key_secret` (String, Sensitive) AccessKey Secret for ZStack API. May also be provided via ZSTACK_ACCESS_KEY_SECRET environment variable. Required if using AccessKey authentication. Mutually exclusive with `account_name` and `account_password`.
- `account_name` (String) Username for ZStack API. May also be provided via ZSTACK_ACCOUNT_NAME environment variable. Required if using Account authentication.  Only supports the platform administrator account (`admin`). Mutually exclusive with `access_key_id` and `access_key_secret`. Using `access_key_id` and `access_key_secret` is the recommended approach for authentication, as it provides more flexibility and security.
- `account_password` (String, Sensitive) Password for ZStack API. May also be provided via ZSTACK_ACCOUNT_PASSWORD environment variable.Required if using Account authentication.  Only supports the platform administrator account (`admin`). Mutually exclusive with `access_key_id` and `access_key_secret`. Using `access_key_id` and `access_key_secret` is the recommended approach for authentication, as it provides more flexibility and security.
- `host` (String) ZStack Cloud MN HOST ip address. May also be provided via ZSTACK_HOST environment variable.
- `port` (Number) ZStack Cloud MN API port. May also be provided via ZSTACK_PORT environment variable.


