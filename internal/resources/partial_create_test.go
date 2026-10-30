package resources_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestPartialCreateRetainsOwnership(t *testing.T) {
	for _, kind := range []string{"tenant", "sso"} {
		t.Run(kind, func(t *testing.T) {
			var lock sync.Mutex
			fail := true
			created, deleted := 0, 0
			id := ""
			jit := false
			var settingsBodies []map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				lock.Lock()
				defer lock.Unlock()
				w.Header().Set("Content-Type", "application/json")
				var response any = map[string]any{}
				switch r.URL.Path {
				case "/v1/mgmt/tenant/create", "/v1/mgmt/sso/settings/new":
					created++
					id = fmt.Sprintf("created-%d", created)
					response = map[string]any{"id": id, "ssoId": id}
				case "/v1/mgmt/tenant/settings", "/v1/mgmt/sso/oidc":
					if r.Method == http.MethodPost && fail {
						http.Error(w, `{"errorCode":"E999999","errorDescription":"injected settings failure"}`, http.StatusBadRequest)
						return
					}
					if r.URL.Path == "/v1/mgmt/tenant/settings" && r.Method == http.MethodPost {
						var body map[string]any
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						settingsBodies = append(settingsBodies, body)
						jit, _ = body["JITDisabled"].(bool)
					}
					if r.URL.Path == "/v1/mgmt/tenant/settings" {
						response = map[string]any{"sessionTokenExpiration": 0, "JITDisabled": jit}
					}
				case "/v1/mgmt/tenant", "/v1/mgmt/tenant/delete":
					if r.Method == http.MethodDelete || r.URL.Path == "/v1/mgmt/tenant/delete" {
						deleted++
						id = ""
					} else {
						response = map[string]any{"id": id, "name": "test", "createdTime": 1, "domains": []string{"sso.example.com"}, "authType": "oidc"}
					}
				case "/v2/mgmt/sso/settings", "/v1/mgmt/sso/settings":
					if r.Method == http.MethodDelete {
						deleted++
						id = ""
					} else {
						response = map[string]any{"ssoId": id, "oidc": map[string]any{"name": "test", "clientId": "mock-client"}}
					}
				}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()
			config := fmt.Sprintf(`provider "descope" {
 management_key = "mock-management"
 base_url = %q
}

resource "descope_%s" "test" {
 project_id = "Pmock"
`, server.URL, kind)
			if kind == "tenant" {
				config += `name = "test"
 settings = {}
}`
			} else {
				config += `tenant_id = "Tmock"
 display_name = "test"
 oidc = { name = "test", client_id = "mock-client" }
}`
			}
			nextConfig := strings.ReplaceAll(config, "settings = {}", "settings = { jit_disabled = true }")
			nextID := "created-2"
			wantCreated := 2
			if kind == "sso" {
				nextConfig = strings.ReplaceAll(config, `display_name = "test"`, `display_name = "renamed"`)
				nextID = "Tmock/created-3"
				wantCreated = 3
			}
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testacc.ProviderFactories,
				Steps: []resource.TestStep{
					{Config: config, ExpectError: regexp.MustCompile("injected settings failure")},
					{Config: config, PreConfig: func() { lock.Lock(); fail = false; lock.Unlock() }, Check: resource.TestCheckResourceAttr("descope_"+kind+".test", "id", map[string]string{"tenant": "created-2", "sso": "Tmock/created-2"}[kind])},
					{Config: nextConfig, Check: resource.TestCheckResourceAttr("descope_"+kind+".test", "id", nextID)},
				},
			})
			lock.Lock()
			defer lock.Unlock()
			require.Equal(t, wantCreated, created)
			require.Equal(t, wantCreated, deleted)
			if kind == "tenant" {
				require.Len(t, settingsBodies, 2)
				for _, body := range settingsBodies {
					require.Equal(t, []any{"sso.example.com"}, body["domains"])
					require.Equal(t, "oidc", body["authType"])
				}
			}
			require.Empty(t, id)
		})
	}
}

func TestSSORejectsConflictingProtocols(t *testing.T) {
	config := `provider "descope" { management_key = "mock-management" }
resource "descope_sso" "test" {
 project_id = "Pmock"
 tenant_id = "Tmock"
 display_name = "test"
 oidc = { name = "test", client_id = "mock-client" }
 saml_metadata = { idp_metadata_url = "https://example.com/metadata" }
}`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProviderFactories,
		Steps:                    []resource.TestStep{{Config: config, PlanOnly: true, ExpectError: regexp.MustCompile("Invalid Attribute Combination")}},
	})
}

func TestSSOAppPartialCreateRetainsOwnership(t *testing.T) {
	for _, kind := range []string{"oidc", "saml", "wsfed"} {
		for _, failAt := range []string{"load", "secret"} {
			if failAt == "secret" && kind != "oidc" {
				continue
			}
			t.Run(kind+"/"+failAt, func(t *testing.T) {
				var mu sync.Mutex
				fail := true
				created, deleted := 0, 0
				id := ""
				var data map[string]any
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					mu.Lock()
					defer mu.Unlock()
					w.Header().Set("Content-Type", "application/json")
					var result any = map[string]any{}
					switch {
					case strings.HasSuffix(r.URL.Path, "/create"):
						created++
						id = fmt.Sprintf("app-%d", created)
						if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
							t.Error(err)
							return
						}
						result = map[string]any{"id": id}
					case strings.HasSuffix(r.URL.Path, "/"+failAt) && fail:
						http.Error(w, `{"errorCode":"E999999","errorDescription":"injected app read failure"}`, http.StatusBadRequest)
						return
					case strings.HasSuffix(r.URL.Path, "/load"):
						result = map[string]any{"id": id, "name": "test", "appType": kind, kind + "Settings": data}
					case strings.HasSuffix(r.URL.Path, "/secret"):
						result = map[string]any{"cleartext": "mock-secret"}
					case strings.HasSuffix(r.URL.Path, "/delete"):
						deleted++
						id = ""
					}
					_ = json.NewEncoder(w).Encode(result)
				}))
				defer server.Close()
				config := fmt.Sprintf(`provider "descope" {
 management_key = "mock-management"
 base_url = %q
}
resource "descope_%s_app" "test" {
 project_id = "Pmock"
 name = "test"
 deletion_protection = false
}`, server.URL, kind)
				if kind == "saml" {
					config = strings.TrimSuffix(config, "}") + `manual_configuration = { acs_url = "https://sp.example.com/acs", entity_id = "sp-entity" }
}`
				}
				resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testacc.ProviderFactories, Steps: []resource.TestStep{
					{Config: config, ExpectError: regexp.MustCompile("injected app read failure")},
					{Config: config, PreConfig: func() { mu.Lock(); fail = false; mu.Unlock() }, Check: resource.TestCheckResourceAttr("descope_"+kind+"_app.test", "id", "app-2")},
				}})
				mu.Lock()
				defer mu.Unlock()
				require.Equal(t, 2, created)
				require.Equal(t, 2, deleted)
				require.Empty(t, id)
			})
		}
	}
}
