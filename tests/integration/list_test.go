//go:build integration || fork

package integration

import (
	"testing"

	"github.com/descope/go-sdk/descope"
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

	// Verify create via SDK
	sdkList := LoadListViaSDK(t, id)
	assert.Equal(t, name, sdkList.Name)
	assert.Equal(t, "Test IP list", sdkList.Description)
	assert.Equal(t, descope.ListTypeIPs, sdkList.Type)

	// Update: change description and add a data entry
	attrs = h.ApplyFixture("list/update.tf", address, nameVar)
	assert.Equal(t, name, attrs["name"])
	assert.Equal(t, "Updated IP list", attrs["description"])
	assert.Equal(t, "ips", attrs["type"])
	RequireListLen(t, attrs, "data", 3)

	// Verify update via SDK
	sdkList = LoadListViaSDK(t, id)
	assert.Equal(t, name, sdkList.Name)
	assert.Equal(t, "Updated IP list", sdkList.Description)

	// Import
	attrs = h.ReimportResource("list/create.tf", address, id, nameVar)
	assert.Equal(t, id, StringAttr(attrs, "id"))

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
