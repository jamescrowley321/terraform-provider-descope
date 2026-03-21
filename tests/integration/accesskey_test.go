//go:build integration || fork

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

	attrs := h.ApplyFixture("access_key/with_options.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	assert.Equal(t, "Test access key", attrs["description"])

	ips := RequireListLen(t, attrs, "permitted_ips", 1)
	assert.Equal(t, "192.168.1.0/24", ips[0])

	claims := RequireMap(t, attrs, "custom_claims")
	assert.Equal(t, "value1", claims["claim1"])

	roles := RequireListLen(t, attrs, "role_names", 1)
	assert.Equal(t, "Tenant Admin", roles[0])

	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
