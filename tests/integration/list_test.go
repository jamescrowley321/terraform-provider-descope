//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_list.test"

	// Create
	attrs := h.ApplyFixture("list/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test IP list", attrs["description"])
	assert.Equal(t, "ips", attrs["type"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
