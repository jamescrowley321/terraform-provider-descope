//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_role.test"

	// Create with two permissions
	attrs := h.ApplyFixture("role/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test role", attrs["description"])

	id := StringAttr(attrs, "id")
	require.Equal(t, name, id)

	permNames := RequireListLen(t, attrs, "permission_names", 2)
	require.NotEmpty(t, permNames)

	// Verify create via SDK
	sdkRole := FindRoleViaSDK(t, name)
	assert.Equal(t, name, sdkRole.Name)
	assert.Equal(t, "Test role", sdkRole.Description)
	assert.Len(t, sdkRole.PermissionNames, 2)

	// Update: change description and reduce to one permission
	attrs = h.ApplyFixture("role/update.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Updated test role", attrs["description"])
	RequireListLen(t, attrs, "permission_names", 1)

	// Verify update via SDK
	sdkRole = FindRoleViaSDK(t, name)
	assert.Equal(t, "Updated test role", sdkRole.Description)
	assert.Len(t, sdkRole.PermissionNames, 1)

	// Import
	attrs = h.ReimportResource("role/create.tf", address, name, nameVar)
	assert.Equal(t, name, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
