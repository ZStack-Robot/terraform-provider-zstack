// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZStackIAM2ProjectsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.IAM2Projects) == 0 {
		t.Skip("no IAM2 projects in env data")
	}
	proj := env.IAM2Projects[0]
	name := envStr(proj, "name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_iam2_projects" "test" {
	name = %q
}`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.zstack_iam2_projects.test", "iam2_projects.0.name", name),
					resource.TestCheckResourceAttr("data.zstack_iam2_projects.test", "iam2_projects.0.uuid", envStr(proj, "uuid")),
				),
			},
		},
	})
}
