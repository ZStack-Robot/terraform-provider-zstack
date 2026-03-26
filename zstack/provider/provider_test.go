// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//
//	"go.testEnvVars": {
//	     "TF_ACC": "1"
//	},
//
// Required environment variables for test:
//
//	ZSTACK_HOST              - ZStack API host (e.g. "172.x.x.x")
//	ZSTACK_PORT              - ZStack API port (default: "8080")
//
// Authentication (choose one):
//
//	Option 1 - AccessKey (recommended):
//	  ZSTACK_ACCESS_KEY_ID     - AccessKey ID
//	  ZSTACK_ACCESS_KEY_SECRET - AccessKey Secret
//
//	Option 2 - Account/Password:
//	  ZSTACK_ACCOUNT_NAME      - ZStack account name
//	  ZSTACK_ACCOUNT_PASSWORD  - ZStack account password

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func providerConfig() string {
	host := getEnvOrDefault("ZSTACK_HOST", "172.30.3.2")
	port := getEnvOrDefault("ZSTACK_PORT", "8080")

	akID := os.Getenv("ZSTACK_ACCESS_KEY_ID")
	akSecret := os.Getenv("ZSTACK_ACCESS_KEY_SECRET")

	if akID != "" && akSecret != "" {
		return fmt.Sprintf(`
	provider "zstack" {
		host              = %q
		port              = %s
		access_key_id     = %q
		access_key_secret = %q
	}
	`, host, port, akID, akSecret)
	}

	return fmt.Sprintf(`
	provider "zstack" {
		host              = %q
		port              = %s
		account_name      = %q
		account_password  = %q
	}
	`,
		host, port,
		getEnvOrDefault("ZSTACK_ACCOUNT_NAME", "admin"),
		getEnvOrDefault("ZSTACK_ACCOUNT_PASSWORD", "password"),
	)
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"zstack": providerserver.NewProtocol6WithError(New("test")()),
	}
)
