//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescoperCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	email := fmt.Sprintf("%s@test.descope.com", name)
	nameVar := "name=" + name
	emailVar := "email=" + email
	address := "descope_descoper.test"

	// Create
	attrs := h.ApplyFixture("descoper/create.tf", address, nameVar, emailVar)
	assert.Equal(t, email, attrs["email"])
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	id := StringAttr(attrs, "id")

	// Update (add phone number)
	attrs = h.ApplyFixture("descoper/update.tf", address, nameVar, emailVar)
	assert.Equal(t, "+15551234567", attrs["phone"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Import
	attrs = h.ReimportResource("descoper/create.tf", address, id, nameVar, emailVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}

func TestDescoperTagRoles(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	email := fmt.Sprintf("%s@test.descope.com", name)
	nameVar := "name=" + name
	emailVar := "email=" + email
	address := "descope_descoper.test"

	attrs := h.ApplyFixture("descoper/with_tag_roles.tf", address, nameVar, emailVar)
	assert.Equal(t, email, attrs["email"])

	// Verify rbac
	rbac, ok := attrs["rbac"].(map[string]any)
	require.True(t, ok, "rbac should be a map")
	assert.Equal(t, false, rbac["is_company_admin"])

	tagRoles, ok := rbac["tag_roles"].([]any)
	require.True(t, ok, "tag_roles should be a list")
	require.Len(t, tagRoles, 1)

	tagRole, ok := tagRoles[0].(map[string]any)
	require.True(t, ok, "tag_role entry should be a map")
	assert.Equal(t, "admin", tagRole["role"])

	tags, ok := tagRole["tags"].([]any)
	require.True(t, ok, "tags should be a list")
	require.Len(t, tags, 2)

	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}

func TestDescoperProjectRoles(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	email := fmt.Sprintf("%s@test.descope.com", name)
	nameVar := "name=" + name
	emailVar := "email=" + email
	address := "descope_descoper.test"

	attrs := h.ApplyFixture("descoper/with_project_roles.tf", address, nameVar, emailVar)
	assert.Equal(t, email, attrs["email"])

	// Verify rbac
	rbac, ok := attrs["rbac"].(map[string]any)
	require.True(t, ok, "rbac should be a map")
	assert.Equal(t, false, rbac["is_company_admin"])

	projectRoles, ok := rbac["project_roles"].([]any)
	require.True(t, ok, "project_roles should be a list")
	require.Len(t, projectRoles, 1)

	projectRole, ok := projectRoles[0].(map[string]any)
	require.True(t, ok, "project_role entry should be a map")
	assert.Equal(t, "developer", projectRole["role"])

	projectIDs, ok := projectRole["project_ids"].([]any)
	require.True(t, ok, "project_ids should be a list")
	require.Len(t, projectIDs, 1)

	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}
