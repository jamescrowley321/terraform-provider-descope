package oauthprovider_test

import (
	"strings"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOAuthClaimMappingLifecycle(t *testing.T) {
	projectID := testacc.ProjectID(t)
	o := &testacc.Resource{Type: "oauth_provider", ID: "claims"}
	config := o.Block(`
 project_id = "` + projectID + `"
 id = "custom_claims"
 client_id = "my-client-id"
 client_secret = "my-client-secret"
 allowed_grant_types = ["authorization_code"]
 authorization_endpoint = "https://auth.example.com"
 token_endpoint = "https://token.example.com"
 user_info_endpoint = "https://userinfo.example.com"
 claim_mapping = { loginId = "sub", email = null, department = null }
 `)
	update := strings.ReplaceAll(config, `{ loginId = "sub", email = null, department = null }`, `{ loginId = "subject", department = "dept" }`)
	cleared := strings.ReplaceAll(update, `{ loginId = "subject", department = "dept" }`, `{}`)
	testacc.Run(t,
		resource.TestStep{Config: config, Check: resource.TestCheckResourceAttr(o.Path(), "claim_mapping.loginId", "sub")},
		resource.TestStep{Config: update, Check: resource.ComposeAggregateTestCheckFunc(resource.TestCheckResourceAttr(o.Path(), "claim_mapping.loginId", "subject"), resource.TestCheckResourceAttr(o.Path(), "claim_mapping.department", "dept"))},
		resource.TestStep{ResourceName: o.Path(), ImportState: true, ImportStateIdFunc: testacc.GenerateImportStateID(o.Path(), "project_id", "id"), ImportStateVerify: true, ImportStateVerifyIgnore: []string{"client_secret"}},
		resource.TestStep{Config: cleared, Check: resource.TestCheckResourceAttr(o.Path(), "claim_mapping.%", "0")},
	)
}
