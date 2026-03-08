//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name

	// Create
	h.LoadFixture("access_key/create.tf")
	out := h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs := h.StateResource("descope_access_key.test")
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	require.NotEmpty(t, attrs["id"])
	require.NotEmpty(t, attrs["cleartext"])

	id := fmt.Sprintf("%v", attrs["id"])

	// Update (set status to inactive, add description)
	h.LoadFixture("access_key/update.tf")
	out = h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs = h.StateResource("descope_access_key.test")
	assert.Equal(t, "inactive", attrs["status"])
	assert.Equal(t, "Updated via integration test", attrs["description"])
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Import (remove from state, then import by ID)
	h.StateRM("descope_access_key.test")
	h.LoadFixture("access_key/create.tf")
	h.Import("descope_access_key.test", id, nameVar)

	attrs = h.StateResource("descope_access_key.test")
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))
	assert.Equal(t, name, attrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
