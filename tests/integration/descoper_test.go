//go:build integration || fork

package integration

import (
	"fmt"
	"strings"
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

	// Verify create via SDK
	// Note: Descope API lowercases emails, so use case-insensitive comparison
	sdkDescoper := LoadDescoperViaSDK(t, id)
	require.NotNil(t, sdkDescoper.Attributes)
	assert.True(t, strings.EqualFold(email, sdkDescoper.Attributes.Email), "email mismatch (case-insensitive): want %q, got %q", email, sdkDescoper.Attributes.Email)
	assert.Equal(t, name, sdkDescoper.Attributes.DisplayName)

	// Update (add phone number)
	attrs = h.ApplyFixture("descoper/update.tf", address, nameVar, emailVar)
	assert.Equal(t, "+15551234567", attrs["phone"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Verify update via SDK
	sdkDescoper = LoadDescoperViaSDK(t, id)
	require.NotNil(t, sdkDescoper.Attributes)
	assert.Equal(t, "+15551234567", sdkDescoper.Attributes.Phone)

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

	// Verify via SDK
	id := StringAttr(attrs, "id")
	sdkDescoper := LoadDescoperViaSDK(t, id)
	require.NotNil(t, sdkDescoper.ReBac, "rbac should not be nil")
	assert.False(t, sdkDescoper.ReBac.IsCompanyAdmin)
	require.Len(t, sdkDescoper.ReBac.Tags, 1)
	assert.Equal(t, "admin", string(sdkDescoper.ReBac.Tags[0].Role))
	assert.Len(t, sdkDescoper.ReBac.Tags[0].Tags, 2)

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

	// Verify via SDK
	id := StringAttr(attrs, "id")
	sdkDescoper := LoadDescoperViaSDK(t, id)
	require.NotNil(t, sdkDescoper.ReBac, "rbac should not be nil")
	assert.False(t, sdkDescoper.ReBac.IsCompanyAdmin)
	require.Len(t, sdkDescoper.ReBac.Projects, 1)
	assert.Equal(t, "developer", string(sdkDescoper.ReBac.Projects[0].Role))
	assert.Len(t, sdkDescoper.ReBac.Projects[0].ProjectIDs, 1)

	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}
