//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSOOIDC(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_sso.test"

	// Create OIDC SSO config
	attrs := h.ApplyFixture("sso/oidc.tf", address, nameVar)
	assert.Equal(t, "Test OIDC SSO", attrs["display_name"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	oidc := RequireMap(t, attrs, "oidc")
	assert.Equal(t, "Test OIDC", oidc["name"])
	assert.Equal(t, "test-client-id", oidc["client_id"])

	// Verify via SDK — id is composite "tenantID/ssoID", use sso_id for the SDK call
	tenantAttrs := h.StateResource("descope_tenant.test")
	tenantID := StringAttr(tenantAttrs, "id")
	ssoID := StringAttr(attrs, "sso_id")
	sdkSSO := LoadSSOSettingsViaSDK(t, tenantID, ssoID)
	require.NotNil(t, sdkSSO.Oidc, "OIDC settings should not be nil")
	assert.Equal(t, "Test OIDC", sdkSSO.Oidc.Name)
	assert.Equal(t, "test-client-id", sdkSSO.Oidc.ClientID)

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
