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
        host              = "172.30.3.3"
		port              = "8080"
        access_key_id     = "RB6frcmLjORqMyM5jhBz"
        access_key_secret = "qdt6QI9TUr7PxT3F8UU1BSrvfjryiRxSbrLWrzZ2"	
	}
	`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"zstack": providerserver.NewProtocol6WithError(New("test")()),
	}
)
