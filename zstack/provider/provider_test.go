// Copyright (c) ZStack.io, Inc.

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
        host              = "172.24.192.246"
		port              = "8080"
		account_name      = "admin" 
        account_password  = "password" 
	}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"zstack": providerserver.NewProtocol6WithError(New("test")()),
	}
)
