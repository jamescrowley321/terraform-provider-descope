//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordSettingsDataSource(t *testing.T) {
	h := NewHarness(t)
	name := "name=" + GenerateName(t)
	resourceAddr := "descope_password_settings.test"
	dataAddr := "data.descope_password_settings.test"

	// Apply resource and data source together
	h.LoadFixture("password_settings/datasource.tf")
	h.Apply(name)

	// Verify the data source reads the same values as the resource
	rAttrs := h.StateResource(resourceAddr)
	dAttrs := h.StateResource(dataAddr)

	assert.Equal(t, rAttrs["disabled"], dAttrs["disabled"])
	assert.Equal(t, rAttrs["min_length"], dAttrs["min_length"])
	assert.Equal(t, rAttrs["lowercase"], dAttrs["lowercase"])
	assert.Equal(t, rAttrs["uppercase"], dAttrs["uppercase"])
	assert.Equal(t, rAttrs["number"], dAttrs["number"])
	assert.Equal(t, rAttrs["non_alphanumeric"], dAttrs["non_alphanumeric"])
	assert.Equal(t, rAttrs["expiration"], dAttrs["expiration"])
	assert.Equal(t, rAttrs["expiration_weeks"], dAttrs["expiration_weeks"])
	assert.Equal(t, rAttrs["reuse"], dAttrs["reuse"])
	assert.Equal(t, rAttrs["reuse_amount"], dAttrs["reuse_amount"])
	assert.Equal(t, rAttrs["lock"], dAttrs["lock"])
	assert.Equal(t, rAttrs["lock_attempts"], dAttrs["lock_attempts"])

	h.Destroy(name)
}
