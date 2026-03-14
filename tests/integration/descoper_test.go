//go:build integration || fork

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

	rbac := RequireMap(t, attrs, "rbac")
	assert.Equal(t, false, rbac["is_company_admin"])

	tagRoles := RequireListLen(t, rbac, "tag_roles", 1)
	tagRole, ok := tagRoles[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "admin", tagRole["role"])
	RequireListLen(t, tagRole, "tags", 2)

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

	rbac := RequireMap(t, attrs, "rbac")
	assert.Equal(t, false, rbac["is_company_admin"])

	projectRoles := RequireListLen(t, rbac, "project_roles", 1)
	projectRole, ok := projectRoles[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "developer", projectRole["role"])
	RequireListLen(t, projectRole, "project_ids", 1)

	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}
