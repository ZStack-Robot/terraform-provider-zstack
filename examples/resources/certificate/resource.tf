# Copyright (c) ZStack.io, Inc.

resource "zstack_certificate" "example" {
  name        = "example-certificate"
  certificate = "-----BEGIN CERTIFICATE-----\nMIICljCCAX4CCQD3J4MsvW1fJzANBgkqhkiG9w0BAQsFADANMQswCQYDVQQGEwJ\nVUzAeFw0yMzAxMDExMjAwMDBaFw0yNDAxMDExMjAwMDBaMA0xCzAJBgNVBAYTAlVT\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2Z3J...\n-----END CERTIFICATE-----"
  description = "Example SSL Certificate"
}

output "certificate" {
  value = zstack_certificate.example
}
