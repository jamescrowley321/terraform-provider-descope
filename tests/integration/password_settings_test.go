//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordSettingsCRUD(t *testing.T) {
	h := NewHarness(t)
	address := "descope_password_settings.test"

	// Create — apply initial password policy
	attrs := h.ApplyFixture("password_settings/create.tf", address)
	assert.Equal(t, "password_settings", attrs["id"])
	assert.Equal(t, true, attrs["enabled"])
	assert.Equal(t, float64(10), attrs["min_length"])
	assert.Equal(t, true, attrs["lowercase"])
	assert.Equal(t, true, attrs["uppercase"])
	assert.Equal(t, true, attrs["number"])
	assert.Equal(t, false, attrs["non_alphanumeric"])
	assert.Equal(t, false, attrs["expiration"])
	assert.Equal(t, false, attrs["lock"])

	// Update — tighten the policy
	attrs = h.ApplyFixture("password_settings/update.tf", address)
	assert.Equal(t, float64(12), attrs["min_length"])
	assert.Equal(t, true, attrs["non_alphanumeric"])
	assert.Equal(t, true, attrs["expiration"])
	assert.Equal(t, float64(26), attrs["expiration_weeks"])
	assert.Equal(t, true, attrs["reuse"])
	assert.Equal(t, float64(5), attrs["reuse_amount"])
	assert.Equal(t, true, attrs["lock"])
	assert.Equal(t, float64(5), attrs["lock_attempts"])

	// Destroy — removes from state only
	h.Destroy()
	assert.False(t, h.HasState())
}

func TestPasswordSettingsDataSource(t *testing.T) {
	h := NewHarness(t)
	resourceAddr := "descope_password_settings.test"
	dataAddr := "data.descope_password_settings.test"

	// Apply resource and data source together
	h.LoadFixture("password_settings/datasource.tf")
	h.Apply()

	// Verify the data source reads the same values as the resource
	rAttrs := h.StateResource(resourceAddr)
	dAttrs := h.StateResource(dataAddr)

	assert.Equal(t, rAttrs["enabled"], dAttrs["enabled"])
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

	h.Destroy()
}
