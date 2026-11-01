package oauthprovider

import (
	"context"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestClaimMappingRefresh(t *testing.T) {
	for _, imported := range []bool{false, true} {
		t.Run(map[bool]string{false: "existing", true: "import"}[imported], func(t *testing.T) {
			ctx := context.Background()
			if imported {
				ctx = helpers.MarkImportContext(ctx)
			}
			var diags diag.Diagnostics
			h := helpers.NewHandler(ctx, &diags)
			m := OAuthProviderModel{ID: stringattr.Value("custom"), ClaimMapping: strmapattr.Value(map[string]string{"old": "old"})}
			m.SetValues(h, map[string]any{"userDataClaimsMapping": map[string]any{"loginId": "sub", "email": "", "customAttributes": map[string]any{"department": "dept"}}})
			require.False(t, diags.HasError())
			got, d := m.ClaimMapping.ToMap(ctx)
			require.False(t, d.HasError())
			require.Equal(t, "sub", got["loginId"].ValueString())
			require.Equal(t, "dept", got["department"].ValueString())
			require.Len(t, got, 2)
			m.SetValues(h, map[string]any{"userDataClaimsMapping": map[string]any{"loginId": "", "customAttributes": map[string]any{}}})
			require.Empty(t, m.ClaimMapping.Elements())
		})
	}
}

func TestClaimMappingKeepsConfiguredNulls(t *testing.T) {
	ctx := context.Background()
	var diags diag.Diagnostics
	h := helpers.NewHandler(ctx, &diags)
	m := OAuthProviderModel{ID: stringattr.Value("custom")}
	m.ClaimMapping.MapValue = types.MapValueMust(types.StringType, map[string]attr.Value{"loginId": types.StringValue("sub"), "email": types.StringNull(), "department": types.StringNull()})
	m.SetValues(h, map[string]any{"userDataClaimsMapping": map[string]any{"loginId": "sub", "email": "email", "familyName": "family_name", "customAttributes": map[string]any{}}})
	require.False(t, diags.HasError())
	require.Len(t, m.ClaimMapping.Elements(), 3)
	require.True(t, m.ClaimMapping.Elements()["email"].IsNull())
	require.True(t, m.ClaimMapping.Elements()["department"].IsNull())
}

func TestClaimMappingDefaultsAndConsoleChanges(t *testing.T) {
	ctx := helpers.MarkImportContext(context.Background())
	var diags diag.Diagnostics
	h := helpers.NewHandler(ctx, &diags)
	m := OAuthProviderModel{ID: stringattr.Value("custom")}
	defaults := map[string]any{"loginId": "sub", "email": "email", "name": "name", "customAttributes": map[string]any{}}
	m.SetValues(h, map[string]any{"userDataClaimsMapping": defaults})
	require.Empty(t, m.ClaimMapping.Elements())
	defaults["customAttributes"] = map[string]any{"department": "dept"}
	m.SetValues(h, map[string]any{"userDataClaimsMapping": defaults})
	require.Equal(t, types.StringValue("sub"), m.ClaimMapping.Elements()["loginId"])
	require.Equal(t, types.StringValue("dept"), m.ClaimMapping.Elements()["department"])
	defaults["email"] = "custom_email"
	m.SetValues(h, map[string]any{"userDataClaimsMapping": defaults})
	require.Equal(t, types.StringValue("custom_email"), m.ClaimMapping.Elements()["email"])
	defaults["email"] = "email"
	m.SetValues(h, map[string]any{"userDataClaimsMapping": defaults})
	_, present := m.ClaimMapping.Elements()["email"]
	require.False(t, present)
	require.False(t, diags.HasError())
}
