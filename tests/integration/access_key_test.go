//go:build integration || fork

package integration

import (
	"encoding/json"
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
	projectID := StringAttr(attrs, "project_id")

	// Verify create via SDK
	sdkKey := LoadAccessKeyViaSDK(t, projectID, id)
	assert.Equal(t, name, sdkKey.Name)
	assert.Equal(t, "active", sdkKey.Status)

	// Update (set status to inactive, add description)
	attrs = h.ApplyFixture("access_key/update.tf", address, nameVar)
	assert.Equal(t, "inactive", attrs["status"])
	assert.Equal(t, "Updated via integration test", attrs["description"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Verify update via SDK
	sdkKey = LoadAccessKeyViaSDK(t, projectID, id)
	assert.Equal(t, "inactive", sdkKey.Status)
	assert.Equal(t, "Updated via integration test", sdkKey.Description)

	// Import (project-level resource uses a composite project_id/id import ID)
	importID := StringAttr(attrs, "project_id") + "/" + id
	attrs = h.ReimportResource("access_key/update.tf", address, importID, nameVar)
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

	var claims map[string]any
	require.NoError(t, json.Unmarshal([]byte(StringAttr(attrs, "custom_claims")), &claims))
	assert.Equal(t, "value1", claims["claim1"])

	roles := RequireListLen(t, attrs, "roles", 1)
	assert.Equal(t, "Viewer", roles[0])

	// Verify via SDK
	id := StringAttr(attrs, "id")
	sdkKey := LoadAccessKeyViaSDK(t, StringAttr(attrs, "project_id"), id)
	assert.Equal(t, name, sdkKey.Name)
	assert.Equal(t, "Test access key", sdkKey.Description)
	assert.Equal(t, []string{"192.168.1.0/24"}, sdkKey.PermittedIPs)
	assert.Equal(t, "value1", sdkKey.CustomClaims["claim1"])

	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
