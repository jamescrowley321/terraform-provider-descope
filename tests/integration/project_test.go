//go:build integration

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_project.test"

	// Create
	attrs := h.ApplyFixture("project/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	id := StringAttr(attrs, "id")

	// Verify project exists via SDK
	assert.True(t, ProjectExistsViaSDK(t, id), "project %s should exist in API after create", id)

	// Update (add tags)
	attrs = h.ApplyFixture("project/update.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Import
	attrs = h.ReimportResource("project/create.tf", address, id, nameVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
