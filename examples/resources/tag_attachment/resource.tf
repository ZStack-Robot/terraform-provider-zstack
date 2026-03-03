resource "zstack_tag_attachment" "test" {

  tag_uuid = "4c49107f43554d87bb4163aa7d9f205d"
  resource_uuids = [
    "c24847378787471f85887e875c3d9064"
  ]

  #   tokens = {
  #    performance = "high"
  #  } 
}

output "zstack_tag_attach" {
  value = zstack_tag_attachment.test
}