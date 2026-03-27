//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFGASchemaCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_fga_schema.test"

	// Create
	attrs := h.ApplyFixture("fgaschema/create.tf", address, nameVar)

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	schema := StringAttr(attrs, "schema")
	assert.Contains(t, schema, "document")
	assert.Contains(t, schema, "owner")

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
