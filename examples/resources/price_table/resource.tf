# Copyright (c) ZStack.io, Inc.

resource "zstack_price_table" "example" {
  name = "example-price-table"

  prices = [
    {
      resource_name = "cpu"
      resource_unit = "Core"
      time_unit     = "s"
      price         = 0.01
    }
  ]
}

output "zstack_price_table" {
  value = zstack_price_table.example
}
