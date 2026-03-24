//go:build integration || fork

package integration

import (
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
	assert.NotEmpty(t, files)
	assert.Contains(t, files, "{")
}
