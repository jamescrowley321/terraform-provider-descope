//go:build integration || fork

package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDataSource(t *testing.T) {
	h := NewHarness(t)
	projectID := os.Getenv("DESCOPE_PROJECT_ID")
	require.NotEmpty(t, projectID, "DESCOPE_PROJECT_ID must be set")
	idVar := "project_id=" + projectID

	// Read the existing project via data source
	h.LoadFixture("project_data/read.tf")
	h.Apply(idVar)

	// Verify the data source reads the project
	dataAttrs := h.StateResource("data.descope_project.test")
	assert.Equal(t, projectID, dataAttrs["id"])
	require.NotEmpty(t, dataAttrs["name"])
}
