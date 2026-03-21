//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagementKeyCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_management_key.test"

	// Create
	attrs := h.ApplyFixture("management_key/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	require.NotEmpty(t, attrs["id"])
	require.NotEmpty(t, attrs["cleartext"])

	id := StringAttr(attrs, "id")

	// Update (set status to inactive, add description)
	attrs = h.ApplyFixture("management_key/update.tf", address, nameVar)
	assert.Equal(t, "inactive", attrs["status"])
	assert.Equal(t, "Updated via integration test", attrs["description"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Import
	attrs = h.ReimportResource("management_key/create.tf", address, id, nameVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))
	assert.Equal(t, name, attrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}

func TestManagementKeyPermittedIPs(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_management_key.test"

	attrs := h.ApplyFixture("management_key/with_permitted_ips.tf", address, nameVar)
	assert.Equal(t, "With permitted IPs", attrs["description"])

	ips := RequireListLen(t, attrs, "permitted_ips", 2)
	assert.Equal(t, "192.168.1.0/24", ips[0])
	assert.Equal(t, "10.0.0.1", ips[1])

	rebac := RequireMap(t, attrs, "rebac")
	companyRoles := RequireListLen(t, rebac, "company_roles", 1)
	assert.Equal(t, "company-full-access", companyRoles[0])

	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}

func TestManagementKeyTagRoles(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_management_key.test"

	attrs := h.ApplyFixture("management_key/with_tag_roles.tf", address, nameVar)
	assert.Equal(t, "With tag roles", attrs["description"])

	rebac := RequireMap(t, attrs, "rebac")
	tagRoles := RequireListLen(t, rebac, "tag_roles", 1)

	tagRole, ok := tagRoles[0].(map[string]any)
	require.True(t, ok)

	tags := RequireListLen(t, tagRole, "tags", 2)
	assert.Contains(t, tags, "production")
	assert.Contains(t, tags, "staging")

	roles := RequireListLen(t, tagRole, "roles", 1)
	assert.Equal(t, "tag-infra-read-write", roles[0])

	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
