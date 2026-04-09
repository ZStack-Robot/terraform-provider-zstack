// Copyright (c) ZStack.io, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccZStackIAM2ProjectsDataSource(t *testing.T) {
	env := loadEnvData(t)
	if len(env.IAM2Projects) == 0 {
		t.Skip("no IAM2 projects in env data")
	}
	proj := env.IAM2Projects[0]
	name := envStr(proj, "name")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
data "zstack_iam2_projects" "test" {
	name = %q
}`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.zstack_iam2_projects.test", tfjsonpath.New("iam2_projects").AtSliceIndex(0).AtMapKey("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue("data.zstack_iam2_projects.test", tfjsonpath.New("iam2_projects").AtSliceIndex(0).AtMapKey("uuid"), knownvalue.StringExact(envStr(proj, "uuid"))),
				},
			},
		},
	})
}
