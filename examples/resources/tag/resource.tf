resource "zstack_tag" "test" {
  name        = "tag"
  description = "tag for test"
  value       = "performance1::{performance1}" # if use withToken typ, value format is token::{key}, only key can update
  color       = "#57D355"
  type        = "withToken" # type support simple and withToken
}

output "zstack_tag" {
  value = zstack_tag.test
}
