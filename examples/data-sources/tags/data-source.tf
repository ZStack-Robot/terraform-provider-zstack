data "zstack_tags" "all_tags" {
  #    name_pattern = "AI%" 
  tag_type = "tag"
  filter {
    name   = "color"
    values = ["purple", "yellow"]
  }
}

# Output examples
output "all_tags" {
  value = data.zstack_tags.all_tags
}