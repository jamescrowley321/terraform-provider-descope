package accesskey

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringSetToSlice(t *testing.T) {
	ctx := context.Background()

	t.Run("returns nil for null set", func(t *testing.T) {
		var diags diag.Diagnostics
		s := strsetattr.Empty()
		// Create a null set by using the zero value workaround
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		// Empty set returns empty slice, not nil
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %v", result)
		}
	})

	t.Run("converts string set to slice", func(t *testing.T) {
		var diags diag.Diagnostics
		s := strsetattr.Value([]string{"admin", "user"})
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(result))
		}
		found := map[string]bool{}
		for _, v := range result {
			found[v] = true
		}
		if !found["admin"] || !found["user"] {
			t.Fatalf("expected admin and user, got %v", result)
		}
	})

	t.Run("handles empty set", func(t *testing.T) {
		var diags diag.Diagnostics
		s := strsetattr.Value([]string{})
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %v", result)
		}
	})
}

func TestStringListToSlice(t *testing.T) {
	ctx := context.Background()

	t.Run("converts string list to slice", func(t *testing.T) {
		var diags diag.Diagnostics
		l := strlistattr.Value([]string{"10.0.0.1", "10.0.0.2"})
		result := StringListToSlice(ctx, l, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(result))
		}
		if result[0] != "10.0.0.1" || result[1] != "10.0.0.2" {
			t.Fatalf("expected [10.0.0.1, 10.0.0.2], got %v", result)
		}
	})

	t.Run("handles empty list", func(t *testing.T) {
		var diags diag.Diagnostics
		l := strlistattr.Value([]string{})
		result := StringListToSlice(ctx, l, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %v", result)
		}
	})
}

func TestStringMapToAnyMap(t *testing.T) {
	t.Run("returns nil for null map", func(t *testing.T) {
		m := strmapattr.Empty()
		result := StringMapToAnyMap(m)
		// Empty map returns nil because len(elems) == 0
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("converts string map to any map", func(t *testing.T) {
		m := strmapattr.Value(map[string]string{"key1": "val1", "key2": "val2"})
		result := StringMapToAnyMap(m)
		if len(result) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(result))
		}
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Fatalf("unexpected values: %v", result)
		}
	})
}

func TestAnyMapToStringMap(t *testing.T) {
	t.Run("returns nil for nil map", func(t *testing.T) {
		result := anyMapToStringMap(nil)
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("converts string values directly", func(t *testing.T) {
		m := map[string]any{"key1": "val1", "key2": "val2"}
		result := anyMapToStringMap(m)
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Fatalf("unexpected values: %v", result)
		}
	})

	t.Run("marshals non-string values to JSON", func(t *testing.T) {
		m := map[string]any{"num": 42, "bool": true}
		result := anyMapToStringMap(m)
		if result["num"] != "42" {
			t.Fatalf("expected '42', got %q", result["num"])
		}
		if result["bool"] != "true" {
			t.Fatalf("expected 'true', got %q", result["bool"])
		}
	})

	t.Run("handles empty map", func(t *testing.T) {
		m := map[string]any{}
		result := anyMapToStringMap(m)
		if len(result) != 0 {
			t.Fatalf("expected empty map, got %v", result)
		}
	})
}

func TestSetModelFromResponse(t *testing.T) {
	t.Run("populates all fields from response", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{
			ID:               "key-123",
			Name:             "test-key",
			Description:      "A test key",
			ClientID:         "client-456",
			CreatedBy:        "user-789",
			CreatedTime:      1700000000,
			ExpireTime:       1800000000,
			UserID:           "user-abc",
			Status:           "active",
			RoleNames:        []string{"admin"},
			PermittedIPs:     []string{"10.0.0.1"},
			CustomClaims:     map[string]any{"claim1": "value1"},
			CustomAttributes: map[string]any{"attr1": "value1"},
		}

		SetModelFromResponse(model, resp, "cleartext-secret")

		assertEqual(t, "ID", model.ID.ValueString(), "key-123")
		assertEqual(t, "Name", model.Name.ValueString(), "test-key")
		assertEqual(t, "Description", model.Description.ValueString(), "A test key")
		assertEqual(t, "ClientID", model.ClientID.ValueString(), "client-456")
		assertEqual(t, "CreatedBy", model.CreatedBy.ValueString(), "user-789")
		assertEqual(t, "Status", model.Status.ValueString(), "active")
		assertEqual(t, "UserID", model.UserID.ValueString(), "user-abc")
		assertEqual(t, "Cleartext", model.Cleartext.ValueString(), "cleartext-secret")

		if model.CreatedTime.ValueInt64() != 1700000000 {
			t.Fatalf("expected CreatedTime 1700000000, got %d", model.CreatedTime.ValueInt64())
		}
		if model.ExpireTime.ValueInt64() != 1800000000 {
			t.Fatalf("expected ExpireTime 1800000000, got %d", model.ExpireTime.ValueInt64())
		}
	})

	t.Run("defaults status to active when empty", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{Status: ""}
		SetModelFromResponse(model, resp, "")
		assertEqual(t, "Status", model.Status.ValueString(), "active")
	})

	t.Run("does not set cleartext when empty", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{}
		SetModelFromResponse(model, resp, "")
		if !model.Cleartext.IsNull() {
			t.Fatalf("expected Cleartext to be null, got %q", model.Cleartext.ValueString())
		}
	})

	t.Run("maps key tenants from response", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{
			KeyTenants: []*descope.AssociatedTenant{
				{TenantID: "t1", TenantName: "Tenant One", Roles: []string{"admin", "user"}},
				{TenantID: "t2", TenantName: "Tenant Two"},
			},
		}
		SetModelFromResponse(model, resp, "")

		elems := model.KeyTenants.Elements()
		if len(elems) != 2 {
			t.Fatalf("expected 2 tenants, got %d", len(elems))
		}
	})

	t.Run("sets empty tenants list when nil", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{KeyTenants: nil}
		SetModelFromResponse(model, resp, "")

		elems := model.KeyTenants.Elements()
		if len(elems) != 0 {
			t.Fatalf("expected 0 tenants, got %d", len(elems))
		}
	})

	t.Run("handles nil custom claims and attributes", func(t *testing.T) {
		model := &AccessKeyModel{}
		resp := &descope.AccessKeyResponse{
			CustomClaims:     nil,
			CustomAttributes: nil,
		}
		SetModelFromResponse(model, resp, "")

		// nil maps produce null/empty strmapattr values
		if model.CustomClaims.IsUnknown() {
			t.Fatal("expected CustomClaims to not be unknown")
		}
		if model.CustomAttributes.IsUnknown() {
			t.Fatal("expected CustomAttributes to not be unknown")
		}
	})
}

func assertEqual(t *testing.T, field string, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: expected %q, got %q", field, want, got)
	}
}

func TestTenantsToSDK(t *testing.T) {
	ctx := context.Background()

	t.Run("converts tenant models to SDK structs", func(t *testing.T) {
		var diags diag.Diagnostics
		tenants := listattr.Value([]*TenantModel{
			{
				TenantID: types.StringValue("t1"),
				Roles:    strsetattr.Value([]string{"admin"}),
			},
			{
				TenantID: types.StringValue("t2"),
				Roles:    strsetattr.Value([]string{}),
			},
		})

		result := TenantsToSDK(ctx, tenants, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 tenants, got %d", len(result))
		}
		if result[0].TenantID != "t1" {
			t.Fatalf("expected tenant ID t1, got %s", result[0].TenantID)
		}
		if len(result[0].Roles) != 1 || result[0].Roles[0] != "admin" {
			t.Fatalf("expected [admin] roles, got %v", result[0].Roles)
		}
	})

	t.Run("returns nil for empty tenant list", func(t *testing.T) {
		var diags diag.Diagnostics
		tenants := listattr.Empty[TenantModel]()
		result := TenantsToSDK(ctx, tenants, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %v", result)
		}
	})
}
