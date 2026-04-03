# Copyright (c) ZStack.io, Inc.

resource "zstack_ldap_server" "example" {
  name     = "example-ldap-server"
  url      = "ldap://ldap.example.com:389"
  base_dn  = "dc=example,dc=com"
  username = "cn=admin,dc=example,dc=com"
  password = "example-password"
}

output "zstack_ldap_server" {
  value     = zstack_ldap_server.example
  sensitive = true
}
