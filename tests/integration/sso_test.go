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

	ssoID := StringAttr(attrs, "sso_id")
	require.NotEmpty(t, ssoID)

	oidc := RequireMap(t, attrs, "oidc")
	assert.Equal(t, "Test OIDC", oidc["name"])
	assert.Equal(t, "test-client-id", oidc["client_id"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
