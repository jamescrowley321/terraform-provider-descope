package helpers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequire(t *testing.T) {
	t.Run("returns value when no errors", func(t *testing.T) {
		var diags diag.Diagnostics
		result := Require("hello", diags)
		if result != "hello" {
			t.Fatalf("expected 'hello', got %q", result)
		}
	})

	t.Run("panics when diagnostics have errors", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic, got none")
			}
		}()
		var diags diag.Diagnostics
		diags.AddError("test error", "some detail")
		Require("value", diags)
	})
}

func TestHasUnknownValues(t *testing.T) {
	t.Run("returns false when called with no arguments", func(t *testing.T) {
		if HasUnknownValues() {
			t.Fatal("expected false for empty args")
		}
	})

	t.Run("returns false for known string", func(t *testing.T) {
		v := types.StringValue("hello")
		if HasUnknownValues(v) {
			t.Fatal("expected false for known string")
		}
	})

	t.Run("returns true for unknown string", func(t *testing.T) {
		v := types.StringUnknown()
		if !HasUnknownValues(v) {
			t.Fatal("expected true for unknown string")
		}
	})

	t.Run("returns false for null string", func(t *testing.T) {
		v := types.StringNull()
		if HasUnknownValues(v) {
			t.Fatal("expected false for null string")
		}
	})

	t.Run("returns true when at least one of multiple values is unknown", func(t *testing.T) {
		v1 := types.StringValue("known")
		v2 := types.StringUnknown()
		if !HasUnknownValues(v1, v2) {
			t.Fatal("expected true when one value is unknown")
		}
	})

	t.Run("returns false for known int", func(t *testing.T) {
		v := types.Int64Value(42)
		if HasUnknownValues(v) {
			t.Fatal("expected false for known int")
		}
	})

	t.Run("returns true for unknown int", func(t *testing.T) {
		v := types.Int64Unknown()
		if !HasUnknownValues(v) {
			t.Fatal("expected true for unknown int")
		}
	})

	t.Run("returns true for unknown bool", func(t *testing.T) {
		v := types.BoolUnknown()
		if !HasUnknownValues(v) {
			t.Fatal("expected true for unknown bool")
		}
	})

	t.Run("returns true for map with unknown element", func(t *testing.T) {
		m := types.MapValueMust(types.StringType, map[string]attr.Value{
			"known": types.StringValue("ok"),
			"bad":   types.StringUnknown(),
		})
		if !HasUnknownValues(m) {
			t.Fatal("expected true for map containing unknown element")
		}
	})

	t.Run("returns false for map with all known elements", func(t *testing.T) {
		m := types.MapValueMust(types.StringType, map[string]attr.Value{
			"a": types.StringValue("ok"),
			"b": types.StringValue("fine"),
		})
		if HasUnknownValues(m) {
			t.Fatal("expected false for map with all known elements")
		}
	})

	t.Run("returns true for list with unknown element", func(t *testing.T) {
		l := types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("ok"),
			types.StringUnknown(),
		})
		if !HasUnknownValues(l) {
			t.Fatal("expected true for list containing unknown element")
		}
	})

	t.Run("returns false for list with all known elements", func(t *testing.T) {
		l := types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
		})
		if HasUnknownValues(l) {
			t.Fatal("expected false for list with all known elements")
		}
	})

	t.Run("returns true for object with unknown attribute", func(t *testing.T) {
		o := types.ObjectValueMust(
			map[string]attr.Type{
				"name": types.StringType,
				"age":  types.Int64Type,
			},
			map[string]attr.Value{
				"name": types.StringValue("ok"),
				"age":  types.Int64Unknown(),
			},
		)
		if !HasUnknownValues(o) {
			t.Fatal("expected true for object containing unknown attribute")
		}
	})

	t.Run("returns false for object with all known attributes", func(t *testing.T) {
		o := types.ObjectValueMust(
			map[string]attr.Type{
				"name": types.StringType,
			},
			map[string]attr.Value{
				"name": types.StringValue("ok"),
			},
		)
		if HasUnknownValues(o) {
			t.Fatal("expected false for object with all known attributes")
		}
	})
}
