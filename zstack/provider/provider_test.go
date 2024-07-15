// Copyright (c) HashiCorp, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the ZStack client is properly configured.
	// It is also possible to use the ZSTACK_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.

	providerConfig = `
	provider "zstack" {
		host = "172.25.16.104"
		port = "8080"
		accountname = "admin"
		accountpassword = "password"
		accesskeyid = "mO6W9gzCQxsfK6OsE7dg"
		accesskeysecret = "Z1B3KVQlGqeaxcpeP55M3WxpRPDUyLqsppp1aLms"		
	}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"zstack": providerserver.NewProtocol6WithError(New("test")()),
	}
)
