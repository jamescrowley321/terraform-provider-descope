//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_access_key.test"

	// Create
	attrs := h.ApplyFixture("access_key/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	require.NotEmpty(t, attrs["id"])
	require.NotEmpty(t, attrs["cleartext"])
	require.NotEmpty(t, attrs["client_id"])

	id := StringAttr(attrs, "id")

	// Update (set status to inactive, add description)
	attrs = h.ApplyFixture("access_key/update.tf", address, nameVar)
	assert.Equal(t, "inactive", attrs["status"])
	assert.Equal(t, "Updated via integration test", attrs["description"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Import
	attrs = h.ReimportResource("access_key/update.tf", address, id, nameVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))
	assert.Equal(t, name, attrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}

func TestAccessKeyWithOptions(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_access_key.test"

	// Create with description, permitted_ips, custom_claims
	attrs := h.ApplyFixture("access_key/with_options.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	assert.Equal(t, "Test access key", attrs["description"])

	// Verify permitted_ips
	ips, ok := attrs["permitted_ips"].([]any)
	require.True(t, ok, "permitted_ips should be a list")
	require.Len(t, ips, 1)
	assert.Equal(t, "192.168.1.0/24", ips[0])

	// Verify custom_claims
	claims, ok := attrs["custom_claims"].(map[string]any)
	require.True(t, ok, "custom_claims should be a map")
	assert.Equal(t, "value1", claims["claim1"])

	// Verify role_names
	roles, ok := attrs["role_names"].([]any)
	require.True(t, ok, "role_names should be a list")
	require.Len(t, roles, 1)
	assert.Equal(t, "Tenant Admin", roles[0])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
