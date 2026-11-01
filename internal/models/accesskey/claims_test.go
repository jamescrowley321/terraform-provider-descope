package accesskey

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/jsonattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/require"
)

func TestClaimsRefresh(t *testing.T) {
	for _, field := range []string{"customClaims", "customAttributes"} {
		for _, tc := range []struct {
			name string
			data map[string]any
			want string
		}{
			{"drift", map[string]any{field: map[string]any{"service": "changed", "account": json.Number("9007199254740993")}}, `{"account":9007199254740993,"service":"changed"}`},
			{"empty", map[string]any{field: map[string]any{}}, `{}`},
			{"absent", map[string]any{}, `{}`},
			{"null", map[string]any{field: nil}, `{}`},
		} {
			t.Run(field+"/"+tc.name, func(t *testing.T) {
				model := AccessKeyModel{CustomClaims: jsonattr.Value(`{"service":"old"}`), CustomAttributes: jsonattr.Value(`{"service":"old"}`)}
				var diagnostics diag.Diagnostics
				model.SetValues(helpers.NewHandler(context.Background(), &diagnostics), tc.data)
				require.False(t, diagnostics.HasError())
				got := model.CustomClaims
				if field == "customAttributes" {
					got = model.CustomAttributes
				}
				require.Equal(t, tc.want, got.ValueString())
			})
		}
	}
}
