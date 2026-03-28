package convert

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/types/valuemaptype"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/types/valuesettype"
)

func TestStringSetToSlice(t *testing.T) {
	ctx := context.Background()

	t.Run("returns nil for null set", func(t *testing.T) {
		var diags diag.Diagnostics
		s := valuesettype.NewNullValue[types.String](ctx)
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if result != nil {
			t.Fatalf("expected nil for null set, got %v", result)
		}
	})

	t.Run("returns nil for unknown set", func(t *testing.T) {
		var diags diag.Diagnostics
		s := valuesettype.NewUnknownValue[types.String](ctx)
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if result != nil {
			t.Fatalf("expected nil for unknown set, got %v", result)
		}
	})

	t.Run("extracts all string values preserving content", func(t *testing.T) {
		var diags diag.Diagnostics
		s := strsetattr.Value([]string{"admin", "user"}) //nolint:contextcheck // Value uses context.Background() by design
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

	t.Run("returns empty slice for empty set", func(t *testing.T) {
		var diags diag.Diagnostics
		s := strsetattr.Empty() //nolint:contextcheck // Empty uses context.Background() by design
		result := StringSetToSlice(ctx, s, &diags)
		if diags.HasError() {
			t.Fatalf("unexpected error: %v", diags.Errors())
		}
		if len(result) != 0 {
			t.Fatalf("expected empty slice, got %v", result)
		}
	})
}

func TestStringMapToAnyMap(t *testing.T) {
	ctx := context.Background()

	t.Run("returns nil for null map", func(t *testing.T) {
		m := valuemaptype.NewMapValue[types.String](ctx)
		result := StringMapToAnyMap(m)
		if result != nil {
			t.Fatalf("expected nil for null map, got %v", result)
		}
	})

	t.Run("returns nil for unknown map", func(t *testing.T) {
		m := valuemaptype.NewUnknownValue[types.String](ctx)
		result := StringMapToAnyMap(m)
		if result != nil {
			t.Fatalf("expected nil for unknown map, got %v", result)
		}
	})

	t.Run("returns nil for empty map", func(t *testing.T) {
		m := strmapattr.Empty() //nolint:contextcheck // Empty uses context.Background() by design
		result := StringMapToAnyMap(m)
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("converts string map values to any interface values", func(t *testing.T) {
		m := strmapattr.Value(map[string]string{"key1": "val1", "key2": "val2"}) //nolint:contextcheck // Value uses context.Background() by design
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
		result := AnyMapToStringMap(nil)
		if result != nil {
			t.Fatalf("expected nil, got %v", result)
		}
	})

	t.Run("formats string values using Sprintf", func(t *testing.T) {
		m := map[string]any{"key1": "val1", "key2": "val2"}
		result := AnyMapToStringMap(m)
		if result["key1"] != "val1" || result["key2"] != "val2" {
			t.Fatalf("unexpected values: %v", result)
		}
	})

	t.Run("formats non-string values using Sprintf", func(t *testing.T) {
		m := map[string]any{"num": 42, "bool": true}
		result := AnyMapToStringMap(m)
		if result["num"] != "42" {
			t.Fatalf("expected '42', got %q", result["num"])
		}
		if result["bool"] != "true" {
			t.Fatalf("expected 'true', got %q", result["bool"])
		}
	})

	t.Run("returns empty map for empty input", func(t *testing.T) {
		m := map[string]any{}
		result := AnyMapToStringMap(m)
		if len(result) != 0 {
			t.Fatalf("expected empty map, got %v", result)
		}
	})
}
