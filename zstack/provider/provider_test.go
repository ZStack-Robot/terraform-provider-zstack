// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/zstackio/zstack-sdk-go-v2/pkg/client"
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

// testAccClient creates a ZStack SDK client for use in CheckDestroy and other
// test helper functions. It reads connection details from environment variables,
// matching the same variables used by providerConfig().
func testAccClient() *client.ZSClient {
	host := getEnvOrDefault("ZSTACK_HOST", "172.30.3.2")
	port, _ := strconv.Atoi(getEnvOrDefault("ZSTACK_PORT", "8080"))

	akID := os.Getenv("ZSTACK_ACCESS_KEY_ID")
	akSecret := os.Getenv("ZSTACK_ACCESS_KEY_SECRET")

	if akID != "" && akSecret != "" {
		return client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(akID, akSecret).ReadOnly(true).Debug(false))
	}

	return client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(
		getEnvOrDefault("ZSTACK_ACCOUNT_NAME", "admin"),
		getEnvOrDefault("ZSTACK_ACCOUNT_PASSWORD", "password"),
	).ReadOnly(true).Debug(false))
}

// testAccCheckResourceDestroyByGet returns a CheckDestroy function for resources
// that can be verified via a Get<Resource>(uuid) SDK call.
// The getFunc should attempt to fetch the resource by ID and return an error if
// not found.
func testAccCheckResourceDestroyByGet(resourceType string, getFunc func(cli *client.ZSClient, id string) error) func(*terraform.State) error {
	return func(s *terraform.State) error {
		cli := testAccClient()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			id := rs.Primary.Attributes["uuid"]
			if id == "" {
				id = rs.Primary.ID
			}
			if id == "" {
				continue
			}
			err := getFunc(cli, id)
			if err == nil {
				return fmt.Errorf("%s %s still exists", resourceType, id)
			}
			if !isZStackNotFoundError(err) {
				return fmt.Errorf("error checking %s %s destroyed: %w", resourceType, id, err)
			}
		}
		return nil
	}
}
