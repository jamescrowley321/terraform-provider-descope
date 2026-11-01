package tenant

import (
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/attrs/strmapattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestRefreshClearsCustomAttributes(t *testing.T) {
	m := TenantModel{TenantID: types.StringValue("T1"), CustomAttributes: strmapattr.Value(map[string]string{"old": "value"})}
	RefreshModelFromAPI(&m, &descope.Tenant{ID: "T1", CustomAttributes: map[string]any{}})
	require.Empty(t, m.CustomAttributes.Elements())
	require.Equal(t, "T1", m.TenantID.ValueString())
}
