package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// Used as a sentinel value when the JSON values for an object are at the root of the map.
const RootKey string = ""

// Require is a helper function that panics if the provided diagnostics contain errors.
func Require[T any](v T, diags diag.Diagnostics) T {
	if errs := diags.Errors(); len(errs) > 0 {
		panic(fmt.Sprintf("%s: %s", errs[0].Summary(), errs[0].Detail()))
	}
	return v
}

// Checks if any of the provided values are in an Unknown state,
// including unknown elements within maps, lists, and objects.
func HasUnknownValues(values ...any) bool {
	for _, v := range values {
		if u, ok := v.(interface{ IsUnknown() bool }); ok && u.IsUnknown() {
			return true
		}
		if m, ok := v.(interface{ Elements() map[string]attr.Value }); ok {
			for _, elem := range m.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		} else if l, ok := v.(interface{ Elements() []attr.Value }); ok {
			for _, elem := range l.Elements() {
				if elem.IsUnknown() {
					return true
				}
			}
		}
		if o, ok := v.(interface{ Attributes() map[string]attr.Value }); ok {
			for _, elem := range o.Attributes() {
				if elem.IsUnknown() {
					return true
				}
			}
		}
	}
	return false
}
