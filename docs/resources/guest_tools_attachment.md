---
page_title: "zstack_guest_tools_attachment Resource - terraform-provider-zstack"
subcategory: ""
description: |-
    Attaches guest tools ISO to a ZStack VM instance.
---

# zstack_guest_tools_attachment (Resource)

Attaches guest tools ISO to a ZStack VM instance.

## Example Usage

```terraform
# Copyright (c) ZStack.io, Inc.

resource "zstack_guest_tools_attachment" "test" {
  instance_uuid = "uuid of vm instance"
}

output "zstack_guest_tools_attachment" {
  value = zstack_guest_tools_attachment.test
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `instance_uuid` (String) UUID of the ZStack VM instance.

### Read-Only

- `guest_tools_status` (String) Status of the ZStack guest tools on the VM (e.g., 'Connected', 'Disconnected').
- `guest_tools_version` (String) Version of the ZStack guest tools installed on the VM.
- `id` (String) Same as the vm_instance_uuid. Used for Terraform tracking.

