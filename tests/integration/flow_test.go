//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlowCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_flow.test"

	// Create flow
	attrs := h.ApplyFixture("flow/create.tf", address, nameVar)

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)
	assert.Equal(t, name, StringAttr(attrs, "flow_id"))
	require.NotEmpty(t, attrs["definition"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
