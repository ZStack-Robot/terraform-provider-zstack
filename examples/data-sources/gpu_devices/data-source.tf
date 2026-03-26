# Copyright (c) ZStack.io, Inc.

data "zstack_gpu_devices" "example" {
  # name = "my-gpu-device"
  # name_pattern = "NVIDIA-%"  # Pattern for fuzzy name search, similar to MySQL LIKE. Use % for multiple characters and _ for exactly one character.

  filter {
    name   = "vendor"
    values = ["NVIDIA"]
  }
}

output "zstack_gpu_devices" {
  value = data.zstack_gpu_devices.example
}
