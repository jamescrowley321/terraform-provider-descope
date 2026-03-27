//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutboundAppCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_outbound_application.test"

	// Create
	attrs := h.ApplyFixture("outbound_application/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test outbound application", attrs["description"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	// Verify via SDK
	sdkApp := LoadOutboundAppViaSDK(t, id)
	assert.Equal(t, name, sdkApp.Name)
	assert.Equal(t, "Test outbound application", sdkApp.Description)

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
