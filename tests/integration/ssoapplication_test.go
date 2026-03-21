//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSOApplicationOIDC(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_sso_application.test"

	// Create OIDC SSO application
	attrs := h.ApplyFixture("ssoapplication/oidc.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test OIDC SSO Application", attrs["description"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	oidc := RequireMap(t, attrs, "oidc")
	assert.Equal(t, "https://app.example.com/login", oidc["login_page_url"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
