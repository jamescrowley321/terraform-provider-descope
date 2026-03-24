//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThirdPartyAppCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_third_party_application.test"

	// Create
	attrs := h.ApplyFixture("thirdpartyapp/create.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Test third-party application", attrs["description"])

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	clientID := StringAttr(attrs, "client_id")
	require.NotEmpty(t, clientID)

	clientSecret := StringAttr(attrs, "client_secret")
	require.NotEmpty(t, clientSecret)

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
