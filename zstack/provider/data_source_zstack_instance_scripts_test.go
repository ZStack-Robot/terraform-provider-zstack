// Copyright (c) ZStack.io, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Run go testing with TF_ACC environment variable set. Edit vscode settings.json and insert
//   "go.testEnvVars": {
//        "TF_ACC": "1"
//   },

func TestAccZStackInstanceScriptsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_scripts" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.#", "5"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.name", "ps"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.uuid", "159051a78fbb49d3b72436a0de6c3c2b"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.platform", "Windows"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.render_params", "[]"),
					//	resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_content", `$path = \"C:\\Users\\Administrator\\Documents\\hello.txt\"\r\n\"Hello, World!\" | Out-File -FilePath $path -Encoding UTF8`),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_timeout", "60"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_type", "Powershell"),
				),
			},
		},
	})
}

func TestAccZStackInstanceScriptsDataSourceFilterByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_scripts" "test" { name = "test-script"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.#", "1"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.name", "test-script"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.uuid", "43a9dcf3860540bf998f07e48fa143cf"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.platform", "Linux"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.render_params", ""),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_timeout", "50"),
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.0.script_type", "Shell"),
				),
			},
		},
	})
}

func TestAccZStackInstanceScriptsDataSourceFilterByNamePattern(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "zstack_scripts" "test" { name_pattern = "test%"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_scripts.test", "scripts.#", "2"),
				),
			},
		},
	})
}
