//go:build integration

package integration

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectExportDataSource(t *testing.T) {
	h := NewHarness(t)
	address := "data.descope_project_export.test"

	// Project export requires a pro license
	h.LoadFixture("projectexport/read.tf")
	out, err := h.TryApply()
	if err != nil && strings.Contains(out, "license") {
		t.Skip("skipping: project export requires pro license")
	}
	require.NoError(t, err, "terraform apply failed: %s", out)

	attrs := h.StateResource(address)

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	files := StringAttr(attrs, "files")
	require.NotEmpty(t, files)

	// Validate the response is valid JSON, not just a string containing "{"
	var parsed map[string]any
	err = json.Unmarshal([]byte(files), &parsed)
	require.NoError(t, err, "files attribute must be valid JSON")
	assert.NotEmpty(t, parsed, "exported project files must not be empty")
}
