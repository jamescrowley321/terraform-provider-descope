//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPermissionCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_permission.test"

	// Create
	attrs := h.ApplyFixture("permission/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test permission", attrs["description"])

	id := StringAttr(attrs, "id")
	require.Equal(t, name, id)

	// Update description
	attrs = h.ApplyFixture("permission/update.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Updated test permission", attrs["description"])

	// Import
	attrs = h.ReimportResource("permission/create.tf", address, name, nameVar)
	assert.Equal(t, name, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
