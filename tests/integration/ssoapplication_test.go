//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSOApplicationOIDC(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_sso_application.test"

	// Create OIDC SSO application — requires enterprise license
	h.LoadFixture("ssoapplication/oidc.tf")
	out, err := h.TryApply(nameVar)
	if err != nil && strings.Contains(out, "license") {
		t.Skip("skipping: SSO application creation requires enterprise license")
	}
	require.NoError(t, err, "terraform apply failed: %s", out)

	attrs := h.StateResource(address)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test OIDC SSO Application", attrs["description"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	oidc := RequireMap(t, attrs, "oidc")
	assert.Equal(t, "https://app.example.com/login", oidc["login_page_url"])

	// Verify via SDK
	sdkApp := LoadSSOApplicationViaSDK(t, id)
	assert.Equal(t, name, sdkApp.Name)
	assert.Equal(t, "Test OIDC SSO Application", sdkApp.Description)
	require.NotNil(t, sdkApp.OIDCSettings, "OIDC settings should not be nil")
	assert.Equal(t, "https://app.example.com/login", sdkApp.OIDCSettings.LoginPageURL)

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
