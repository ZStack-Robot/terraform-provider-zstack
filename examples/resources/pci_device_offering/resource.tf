# Copyright (c) ZStack.io, Inc.

resource "zstack_pci_device_offering" "example" {
  name      = "example-pci-device-offering"
  vendor_id = "10de"
  device_id = "1234"
}

output "zstack_pci_device_offering" {
  value = zstack_pci_device_offering.example
}
