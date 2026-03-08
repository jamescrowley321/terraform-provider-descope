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

	// Create
	h.LoadFixture("descoper/create.tf")
	out := h.Apply(nameVar, emailVar)
	assert.Contains(t, out, "Apply complete!")

	attrs := h.StateResource("descope_descoper.test")
	assert.Equal(t, email, attrs["email"])
	assert.Equal(t, name, attrs["name"])
	require.NotEmpty(t, attrs["id"])

	id := fmt.Sprintf("%v", attrs["id"])

	// Update (add phone number)
	h.LoadFixture("descoper/update.tf")
	out = h.Apply(nameVar, emailVar)
	assert.Contains(t, out, "Apply complete!")

	attrs = h.StateResource("descope_descoper.test")
	assert.Equal(t, "+15551234567", attrs["phone"])
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Import (remove from state, then import by ID)
	h.StateRM("descope_descoper.test")
	h.LoadFixture("descoper/create.tf")
	h.Import("descope_descoper.test", id, nameVar, emailVar)

	attrs = h.StateResource("descope_descoper.test")
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Destroy
	h.Destroy(nameVar, emailVar)
	assert.False(t, h.HasState())
}
