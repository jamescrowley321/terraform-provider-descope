//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDescoperCRUD(t *testing.T) {
	h := NewHarness(t)
	name := GenerateName(t)
	email := fmt.Sprintf("%s@test.descope.com", name)
	nameVar := "name=" + name
	emailVar := "email=" + email
	address := "descope_descoper.test"

	// Create
	attrs := h.ApplyFixture("descoper/create.tf", address, nameVar, emailVar)
	assert.Equal(t, email, attrs["email"])
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	id := StringAttr(attrs, "id")

	// Update (add phone number)
	attrs = h.ApplyFixture("descoper/update.tf", address, nameVar, emailVar)
	assert.Equal(t, "+15551234567", attrs["phone"])
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Import
	attrs = h.ReimportResource("descoper/create.tf", address, id, nameVar, emailVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}
