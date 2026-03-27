//go:build integration

package integration

import (
	"context"
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

	settings := RequireMap(t, attrs, "project_settings")
	assert.Equal(t, "3 weeks", settings["refresh_token_expiration"])
	assert.Equal(t, "1 hour", settings["session_token_expiration"])
	assert.Equal(t, true, settings["refresh_token_rotation"])

	id := StringAttr(attrs, "id")

	// Verify project exists via SDK
	assert.True(t, ProjectExistsViaSDK(t, id), "project %s should exist in API", id)

	// Import and verify settings survive
	attrs = h.ReimportResource("project/with_settings.tf", address, id, nameVar)
	settings = RequireMap(t, attrs, "project_settings")
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

	authz := RequireMap(t, attrs, "authorization")

	roles := RequireListLen(t, authz, "roles", 2)
	role0, ok := roles[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "App Developer", role0["name"])
	assert.Equal(t, "app-developer", role0["key"])
	RequireListLen(t, role0, "permissions", 3)

	permissions := RequireListLen(t, authz, "permissions", 3)
	perm0, ok := permissions[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "build-apps", perm0["name"])
	assert.Equal(t, "Allowed to build and sign applications", perm0["description"])

	// Verify project exists via SDK
	id := StringAttr(attrs, "id")
	assert.True(t, ProjectExistsViaSDK(t, id), "project %s should exist in API", id)

	// Verify roles and permissions exist in the created project via SDK
	projectClient := newSDKClientWithProject(t, id)
	ctx := context.Background()
	sdkRoles, err := projectClient.Management.Role().LoadAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sdkRoles), 2, "project should have at least 2 roles")

	sdkPerms, err := projectClient.Management.Permission().LoadAll(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(sdkPerms), 3, "project should have at least 3 permissions")

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
