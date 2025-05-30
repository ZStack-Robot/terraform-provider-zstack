---
page_title: "zstack_subnet_ip_range Resource - terraform-provider-zstack"
subcategory: ""
description: |-
    This resource allows you to manage subnets in ZStack. A subnet is a logical subdivision of an IP network, defined by a range of IP addresses. You can define the subnet's properties, such as its name, IP range, netmask, gateway, and the L3 network it belongs to.
---

# zstack_subnet_ip_range (Resource)

This resource allows you to manage subnets in ZStack. A subnet is a logical subdivision of an IP network, defined by a range of IP addresses. You can define the subnet's properties, such as its name, IP range, netmask, gateway, and the L3 network it belongs to.

## Example Usage

```terraform
# Copyright (c) ZStack.io, Inc.

resource "zstack_subnet_ip_range" "test" {
  l3_network_uuid = "6a7c9dd9d6e449f992a59df8c102b3ba"
  name            = "net1"
  start_ip        = "192.168.1.1"
  end_ip          = "192.168.1.100"
  netmask         = "255.255.0.0"
  gateway         = "192.168.100.1"
}

output "zstack_subnet_ip_range" {
  value = zstack_subnet_ip_range.test
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `end_ip` (String) The ending IP address of the subnet range.
- `gateway` (String) The default gateway for the subnet.
- `l3_network_uuid` (String) The UUID of the L3 network to which the subnet belongs.
- `name` (String) The name of the subnet. This is a user-defined identifier for the subnet.
- `netmask` (String) The subnet mask, used to define the network portion of an IP address.
- `start_ip` (String) The starting IP address of the subnet range.

### Optional

- `ip_range_type` (String) The type of IP range. Possible values depend on the ZStack configuration (e.g., 'Normal' or 'Reserved').

### Read-Only

- `uuid` (String) The unique identifier of the subnet.


