package resources_test

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/descope/terraform-provider-descope/tools/testacc"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyClaimsLifecycle(t *testing.T) {
	var lock sync.Mutex
	var stored map[string]any
	updates := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		if r.URL.Path != "/v1/mgmt/infra" || r.Header.Get("Authorization") != "Bearer Pmock:mock-management" {
			http.Error(w, "unexpected request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var response map[string]any
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			var request struct {
				Entity string         `json:"entity"`
				Data   map[string]any `json:"data"`
			}
			decoder := json.NewDecoder(r.Body)
			decoder.UseNumber()
			if err := decoder.Decode(&request); err != nil || request.Entity != "access_key" {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			stored = request.Data
			stored["clientId"] = "mock-client"
			stored["createdTime"] = json.Number("1700000000")
			stored["createdBy"] = "test"
			response = maps.Clone(stored)
			if r.Method == http.MethodPost {
				response["cleartext"] = "mock-secret"
			} else {
				updates++
			}
		case http.MethodGet:
			if stored == nil {
				http.Error(w, `{"errorCode":"E084004","errorDescription":"not found"}`, http.StatusNotFound)
				return
			}
			response = maps.Clone(stored)
		case http.MethodDelete:
			stored = nil
			_, _ = w.Write([]byte(`{}`))
			return
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "Kmock", "data": response})
	}))
	defer server.Close()
	config := func(claims, attributes string) string {
		return fmt.Sprintf(`provider "descope" {
 management_key = "mock-management"
 base_url = %q
 }
 resource "descope_access_key" "test" {
 project_id = "Pmock"
 name = "service"
 custom_claims = %q
 custom_attributes = %q
 }`, server.URL, claims, attributes)
	}
	claims := "{\n  \"account\": 9007199254740993,\n  \"nested\": {\"enabled\": true, \"scopes\": [\"read\", \"write\"]}\n}"
	attributes := `{"owner":"platform"}`
	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("descope_access_key.test", "custom_claims", claims),
		resource.TestCheckResourceAttr("descope_access_key.test", "custom_attributes", attributes),
		resource.TestCheckResourceAttr("descope_access_key.test", "cleartext", "mock-secret"),
		resource.TestCheckResourceAttr("descope_access_key.test", "created_time", "1700000000"),
	)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testacc.ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			lock.Lock()
			defer lock.Unlock()
			if stored != nil {
				return fmt.Errorf("key still exists after destroy")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: config(claims, attributes), Check: check},
			{Config: config(claims, attributes), PlanOnly: true},
			{Config: config(claims, attributes), PlanOnly: true, ExpectNonEmptyPlan: true, PreConfig: func() {
				lock.Lock()
				defer lock.Unlock()
				stored["customClaims"] = map[string]any{"account": json.Number("9007199254740992")}
				delete(stored, "customAttributes")
			}},
			{Config: config(claims, attributes), Check: check},
			{Config: config("{}", "{}"), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("descope_access_key.test", "custom_claims", "{}"),
				resource.TestCheckResourceAttr("descope_access_key.test", "custom_attributes", "{}"),
				resource.TestCheckResourceAttr("descope_access_key.test", "cleartext", "mock-secret"),
			)},
		},
	})
	lock.Lock()
	defer lock.Unlock()
	require.Equal(t, 2, updates)
}
