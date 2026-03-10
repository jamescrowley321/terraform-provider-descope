//go:build integration

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessKeyCRUD(t *testing.T) {
	// Bootstrap: create a management key to use as the SDK bearer token,
	// since the access key resource authenticates via the Descope SDK which
	// requires a valid management key (the infra API key may not suffice).
	bootstrap := NewHarness(t)
	bootstrapName := GenerateName(t) + "-bootstrap"
	bootstrap.LoadFixture("access_key/bootstrap.tf")
	bootstrap.Apply("name=" + bootstrapName)
	mgmtKey := bootstrap.Output("cleartext")
	require.NotEmpty(t, mgmtKey, "bootstrap management key cleartext must not be empty")

	// Create a harness that uses the bootstrapped management key
	h := NewHarnessWithManagementKey(t, mgmtKey)
	name := GenerateName(t)
	nameVar := "name=" + name

	// Create
	h.LoadFixture("access_key/create.tf")
	out := h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs := h.StateResource("descope_access_key.test")
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "active", attrs["status"])
	require.NotEmpty(t, attrs["id"])
	require.NotEmpty(t, attrs["cleartext"])

	id := fmt.Sprintf("%v", attrs["id"])

	// Update (set status to inactive, add description)
	h.LoadFixture("access_key/update.tf")
	out = h.Apply(nameVar)
	assert.Contains(t, out, "Apply complete!")

	attrs = h.StateResource("descope_access_key.test")
	assert.Equal(t, "inactive", attrs["status"])
	assert.Equal(t, "Updated via integration test", attrs["description"])
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))

	// Import (remove from state, then import by ID using update fixture to match current state)
	h.StateRM("descope_access_key.test")
	h.LoadFixture("access_key/update.tf")
	h.Import("descope_access_key.test", id, nameVar)

	attrs = h.StateResource("descope_access_key.test")
	assert.Equal(t, id, fmt.Sprintf("%v", attrs["id"]))
	assert.Equal(t, name, attrs["name"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
