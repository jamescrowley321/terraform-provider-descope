package project_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jamescrowley321/terraform-provider-descope/tools/testacc"
)

func TestProject(t *testing.T) {
	p := testacc.Project(t)
	testacc.Run(t,
		resource.TestStep{
			Config: p.Config(`
				environment = "foo"
			`),
			ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = "production"
			`),
			Check: p.Check(map[string]any{
				"id":          testacc.AttributeIsSet,
				"name":        p.Name,
				"environment": "production",
				"tags":        []string{},
			}),
		},
		resource.TestStep{
			ResourceName: p.Path(),
			ImportState:  true,
		},
		resource.TestStep{
			PreConfig: func() {
				p.Name += "bar"
			},
			Config: p.Config(),
			Check: p.Check(map[string]any{
				"name":        p.Name,
				"environment": "production",
			}),
		},
		resource.TestStep{
			Config: p.Config(`
				environment = ""
				tags = ["foo", "bar"]
			`),
			Check: p.Check(map[string]any{
				"name":        p.Name,
				"tags":        []string{"foo", "bar"},
				"environment": "",
			}),
		},
	)
}
