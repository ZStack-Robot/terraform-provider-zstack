# Copyright (c) ZStack.io, Inc.

resource "zstack_price_table" "example" {
  name = "example-price-table"
}

output "zstack_price_table" {
  value = zstack_price_table.example
}
