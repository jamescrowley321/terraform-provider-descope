//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenantCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_tenant.test"

	// Create
	attrs := h.ApplyFixture("tenant/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])
	id := StringAttr(attrs, "id")

	// Update (add self_provisioning_domains, enable enforce_sso)
	attrs = h.ApplyFixture("tenant/update.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, true, attrs["enforce_sso"])
	domains := RequireListLen(t, attrs, "self_provisioning_domains", 1)
	assert.Equal(t, "example.com", domains[0])

	// Import
	attrs = h.ReimportResource("tenant/update.tf", address, id, nameVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))
	assert.Equal(t, name, attrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}

func TestTenantWithCustomID(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	customID := "tid-" + name
	address := "descope_tenant.test"

	attrs := h.ApplyFixture("tenant/with_custom_id.tf", address, "name="+name, "tenant_id="+customID)
	assert.Equal(t, customID, StringAttr(attrs, "id"))
	assert.Equal(t, customID, StringAttr(attrs, "tenant_id"))
	assert.Equal(t, name, attrs["name"])

	h.Destroy("name="+name, "tenant_id="+customID)
	assert.False(t, h.HasState())
}

func TestTenantWithSettings(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	address := "descope_tenant.test"

	attrs := h.ApplyFixture("tenant/with_settings.tf", address, "name="+name)
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	settings := RequireMap(t, attrs, "settings")
	assert.Equal(t, true, settings["session_settings_enabled"])
	assert.Equal(t, float64(30), settings["refresh_token_expiration"])
	assert.Equal(t, "days", settings["refresh_token_expiration_unit"])
	assert.Equal(t, float64(10), settings["session_token_expiration"])
	assert.Equal(t, "minutes", settings["session_token_expiration_unit"])

	h.Destroy("name=" + name)
	assert.False(t, h.HasState())
}
