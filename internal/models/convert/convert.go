package convert

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSetToSlice converts a Terraform string set to a Go string slice.
func StringSetToSlice(_ context.Context, s strsetattr.Type, _ *diag.Diagnostics) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elems := s.Elements()
	result := make([]string, 0, len(elems))
	for _, v := range elems {
		if str, ok := v.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

// StringMapToAnyMap converts a Terraform string map to a map[string]any for SDK calls.
func StringMapToAnyMap(m strmapattr.Type) map[string]any {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	elems := m.Elements()
	if len(elems) == 0 {
		return nil
	}
	result := make(map[string]any, len(elems))
	for k, v := range elems {
		if str, ok := v.(types.String); ok {
			result[k] = str.ValueString()
		}
	}
	return result
}

// AnyMapToStringMap converts a map[string]any from the SDK to a map[string]string.
func AnyMapToStringMap(m map[string]any) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
