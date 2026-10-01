package resources_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestJSONRefreshDetectsDrift(t *testing.T) {
	for _, kind := range []string{"flow", "widget", "styles"} {
		t.Run(kind, func(t *testing.T) {
			var mu sync.Mutex
			var remote map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodDelete {
					_, _ = w.Write([]byte(`{}`))
					return
				}
				if r.Method == http.MethodPost && r.URL.Path != "/v2/mgmt/theme/export" {
					if err := json.NewDecoder(r.Body).Decode(&remote); err != nil {
						t.Error(err)
						return
					}
					if kind == "styles" {
						theme, ok := remote["theme"].(map[string]any)
						require.True(t, ok)
						remote = theme
					}
					remote["serverDefault"] = true
					if kind == "styles" {
						remote["componentsVersion"] = "1"
					} else {
						metadata, ok := remote["metadata"].(map[string]any)
						require.True(t, ok)
						metadata["componentsVersion"] = "1"
					}
				}
				var result any = remote
				if kind == "styles" {
					result = map[string]any{"theme": remote}
				}
				_ = json.NewEncoder(w).Encode(result)
			}))
			defer server.Close()
			payload := `{"metadata":{"name":"before"},"contents":{}}`
			identity := `flow_id = "test"`
			if kind == "widget" {
				payload = `{"metadata":{"name":"before"},"screens":{}}`
				identity = `widget_id = "test"`
			}
			if kind == "styles" {
				payload = `{"styles":{"primary":"before"}}`
				identity = ""
			}
			importID := "Pmock/test"
			if kind == "styles" {
				importID = "Pmock"
			}
			config := fmt.Sprintf(`provider "descope" {
 management_key = "mock-management"
 base_url = %q
}
resource "descope_%s" "test" {
 project_id = "Pmock"
%s
 data = %q
}`, server.URL, kind, identity, payload)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testacc.ProviderFactories,
				Steps: []resource.TestStep{
					{Config: config},
					{Config: config, PlanOnly: true, PreConfig: func() {
						mu.Lock()
						defer mu.Unlock()
						if kind == "styles" {
							remote["componentsVersion"] = "2"
						} else {
							metadata, ok := remote["metadata"].(map[string]any)
							require.True(t, ok)
							metadata["componentsVersion"] = "2"
						}
					}},
					{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true, PreConfig: func() {
						mu.Lock()
						defer mu.Unlock()
						if kind == "styles" {
							styles, ok := remote["styles"].(map[string]any)
							require.True(t, ok)
							styles["primary"] = "after"
						} else {
							metadata, ok := remote["metadata"].(map[string]any)
							require.True(t, ok)
							metadata["name"] = "after"
						}
					}},
					{Config: config},
					{ResourceName: "descope_" + kind + ".test", ImportState: true, ImportStateId: importID, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"data"}},
				},
			})
		})
	}
}

func TestSSODomainsRefresh(t *testing.T) {
	var mu sync.Mutex
	domains := []string{"before.example.com"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		result := map[string]any{}
		if r.URL.Path == "/v1/mgmt/sso/settings/new" {
			result["ssoId"] = "S1"
		}
		if r.Method == http.MethodGet {
			result = map[string]any{"ssoId": "S1", "tenant": map[string]any{"domains": domains}, "oidc": map[string]any{"name": "test", "clientId": "mock-client"}}
		}
		_ = json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()
	config := fmt.Sprintf(`provider "descope" {
 management_key = "mock-management"
 base_url = %q
}
resource "descope_sso" "test" {
 project_id = "Pmock"
 tenant_id = "Tmock"
 display_name = "test"
 domains = ["before.example.com"]
 oidc = { name = "test", client_id = "mock-client" }
}`, server.URL)
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testacc.ProviderFactories, Steps: []resource.TestStep{
		{Config: config},
		{Config: config, PlanOnly: true, ExpectNonEmptyPlan: true, PreConfig: func() { mu.Lock(); domains = []string{}; mu.Unlock() }},
	}})
}
