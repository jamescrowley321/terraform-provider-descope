//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name

	// Create
	h.LoadFixture("project/create.tf")
	out := h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs := h.StateResource("descope_project.test")
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	id := fmt.Sprintf("%v", attrs["id"])

	// Update (add tags)
	h.LoadFixture("project/update.tf")
	out = h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs = h.StateResource("descope_project.test")
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Import
	h.StateRM("descope_project.test")
	h.LoadFixture("project/create.tf")
	h.Import("descope_project.test", id, nameVar)

	attrs = h.StateResource("descope_project.test")
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
