resource "zstack_script_execution" "example" {
  script_uuid    = "The uuid of script"
  instance_uuid  = "The uuid of vm instance"
  script_timeout = 180
}


output "zstack_script_out" {
  value = zstack_script_execution.example
}