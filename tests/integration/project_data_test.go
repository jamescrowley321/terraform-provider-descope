//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDataSource(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name

	// Create a project and read it back via data source
	h.LoadFixture("project_data/read.tf")
	h.Apply(nameVar)

	// Verify the resource
	resAttrs := h.StateResource("descope_project.test")
	require.NotEmpty(t, resAttrs["id"])
	assert.Equal(t, name, resAttrs["name"])

	// Verify the data source reads the same project
	dataAttrs := h.StateResource("data.descope_project.test")
	assert.Equal(t, resAttrs["id"], dataAttrs["id"])
	assert.Equal(t, name, dataAttrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
