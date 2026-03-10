//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectSettings(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_project.test"

	// Create with project_settings
	attrs := h.ApplyFixture("project/with_settings.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])

	settings, ok := attrs["project_settings"].(map[string]any)
	require.True(t, ok, "project_settings should be a map")
	assert.Equal(t, "3 weeks", settings["refresh_token_expiration"])
	assert.Equal(t, "1 hour", settings["session_token_expiration"])
	assert.Equal(t, true, settings["refresh_token_rotation"])

	id := StringAttr(attrs, "id")

	// Import and verify settings survive
	attrs = h.ReimportResource("project/with_settings.tf", address, id, nameVar)
	settings, ok = attrs["project_settings"].(map[string]any)
	require.True(t, ok, "project_settings should be a map after import")
	assert.Equal(t, "3 weeks", settings["refresh_token_expiration"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}

func TestProjectAuthorization(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_project.test"

	// Create with authorization
	attrs := h.ApplyFixture("project/with_authorization.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])

	authz, ok := attrs["authorization"].(map[string]any)
	require.True(t, ok, "authorization should be a map")

	// Verify roles
	roles, ok := authz["roles"].([]any)
	require.True(t, ok, "roles should be a list")
	require.Len(t, roles, 2)

	role0, ok := roles[0].(map[string]any)
	require.True(t, ok, "role entry should be a map")
	assert.Equal(t, "App Developer", role0["name"])
	assert.Equal(t, "app-developer", role0["key"])

	perms0, ok := role0["permissions"].([]any)
	require.True(t, ok, "permissions should be a list")
	assert.Len(t, perms0, 3)

	// Verify permissions
	permissions, ok := authz["permissions"].([]any)
	require.True(t, ok, "permissions should be a list")
	require.Len(t, permissions, 3)

	perm0, ok := permissions[0].(map[string]any)
	require.True(t, ok, "permission entry should be a map")
	assert.Equal(t, "build-apps", perm0["name"])
	assert.Equal(t, "Allowed to build and sign applications", perm0["description"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
