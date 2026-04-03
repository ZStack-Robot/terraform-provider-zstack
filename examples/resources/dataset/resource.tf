# Copyright (c) ZStack.io, Inc.

resource "zstack_dataset" "example" {
  name              = "example-dataset"
  url               = "https://example.com/datasets/training-data"
  model_center_uuid = "example-model-center-uuid"
}

output "zstack_dataset" {
  value = zstack_dataset.example
}
