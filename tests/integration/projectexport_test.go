//go:build integration || fork

package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectExportDataSource(t *testing.T) {
	h := NewHarness(t)
	address := "data.descope_project_export.test"

	attrs := h.ApplyFixture("projectexport/read.tf", address)

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	files := StringAttr(attrs, "files")
	require.NotEmpty(t, files)

	// Validate the response is valid JSON, not just a string containing "{"
	var parsed map[string]any
	err := json.Unmarshal([]byte(files), &parsed)
	require.NoError(t, err, "files attribute must be valid JSON")
	assert.NotEmpty(t, parsed, "exported project files must not be empty")
}
