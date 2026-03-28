package tenant

import (
	"context"
	"math"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/convert"
)

// ModelToRequest converts the TenantModel to an SDK TenantRequest for Create/Update.
func ModelToRequest(ctx context.Context, model *TenantModel, diags *diag.Diagnostics) *descope.TenantRequest {
	return &descope.TenantRequest{
		Name:                    model.Name.ValueString(),
		SelfProvisioningDomains: convert.StringSetToSlice(ctx, model.SelfProvisioningDomains, diags),
		CustomAttributes:        convert.StringMapToAnyMap(model.CustomAttributes),
		EnforceSSO:              model.EnforceSSO.ValueBool(),
		Disabled:                model.Disabled.ValueBool(),
		ParentTenantID:          model.ParentTenantID.ValueString(),
		EnforceSSOExclusions:    convert.StringSetToSlice(ctx, model.EnforceSSOExclusions, diags),
		RoleInheritance:         descope.RoleInheritance(model.RoleInheritance.ValueString()),
	}
}

// SetModelFromTenant populates the TenantModel from an SDK Tenant response.
//
//nolint:contextcheck // Value helpers use context.Background() by design
func SetModelFromTenant(model *TenantModel, t *descope.Tenant) {
	model.ID = types.StringValue(t.ID)
	model.Name = types.StringValue(t.Name)
	model.EnforceSSO = types.BoolValue(t.EnforceSSO)
	model.Disabled = types.BoolValue(t.Disabled)
	model.CreatedTime = types.Int64Value(int64(t.CreatedTime))
	model.RoleInheritance = types.StringValue(string(t.RoleInheritance))

	model.SelfProvisioningDomains = strsetattr.Value(t.SelfProvisioningDomains)
	model.EnforceSSOExclusions = strsetattr.Value(t.EnforceSSOExclusions)
	model.Domains = strsetattr.Value(t.Domains)

	if t.AuthType != "" {
		model.AuthType = types.StringValue(t.AuthType)
	} else {
		model.AuthType = types.StringValue("")
	}

	if len(t.CustomAttributes) > 0 {
		model.CustomAttributes = strmapattr.Value(convert.AnyMapToStringMap(t.CustomAttributes))
	}
}

// ModelToSettings converts the SettingsModel to an SDK TenantSettings.
// It also populates the overlapping fields (Domains, SelfProvisioningDomains, AuthType)
// from the top-level tenant model.
func ModelToSettings(ctx context.Context, model *TenantModel, diags *diag.Diagnostics) *descope.TenantSettings {
	s := model.Settings
	return &descope.TenantSettings{
		Domains:                    convert.StringSetToSlice(ctx, model.Domains, diags),
		SelfProvisioningDomains:    convert.StringSetToSlice(ctx, model.SelfProvisioningDomains, diags),
		AuthType:                   model.AuthType.ValueString(),
		SessionSettingsEnabled:     s.SessionSettingsEnabled.ValueBool(),
		RefreshTokenExpiration:     clampInt32(s.RefreshTokenExpiration.ValueInt64()),
		RefreshTokenExpirationUnit: s.RefreshTokenExpirationUnit.ValueString(),
		SessionTokenExpiration:     clampInt32(s.SessionTokenExpiration.ValueInt64()),
		SessionTokenExpirationUnit: s.SessionTokenExpirationUnit.ValueString(),
		StepupTokenExpiration:      clampInt32(s.StepupTokenExpiration.ValueInt64()),
		StepupTokenExpirationUnit:  s.StepupTokenExpirationUnit.ValueString(),
		EnableInactivity:           s.EnableInactivity.ValueBool(),
		InactivityTime:             clampInt32(s.InactivityTime.ValueInt64()),
		InactivityTimeUnit:         s.InactivityTimeUnit.ValueString(),
		JITDisabled:                s.JITDisabled.ValueBool(),
	}
}

// RefreshModelFromAPI updates the model with fresh data from the API while
// preserving fields that aren't returned by Tenant.Load (TenantID, DefaultRoles,
// CascadeDelete, ParentTenantID, and optionally CustomAttributes).
// Returns the previously saved Settings pointer so the caller can decide whether
// to load settings from the API.
func RefreshModelFromAPI(model *TenantModel, t *descope.Tenant) *SettingsModel {
	savedDefaultRoles := model.DefaultRoles
	savedCascadeDelete := model.CascadeDelete
	savedParentTenantID := model.ParentTenantID
	savedTenantID := model.TenantID
	savedCustomAttrs := model.CustomAttributes
	savedSettings := model.Settings

	SetModelFromTenant(model, t)
	model.DefaultRoles = savedDefaultRoles
	model.CascadeDelete = savedCascadeDelete
	model.ParentTenantID = savedParentTenantID
	model.TenantID = savedTenantID
	if len(t.CustomAttributes) == 0 {
		model.CustomAttributes = savedCustomAttrs
	}

	return savedSettings
}

// clampInt32 safely converts an int64 to int32, clamping to int32 bounds.
func clampInt32(v int64) int32 {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	if v < math.MinInt32 {
		return math.MinInt32
	}
	return int32(v) // #nosec G115 -- bounds checked above
}

// SetSettingsFromSDK populates the SettingsModel from an SDK TenantSettings response.
func SetSettingsFromSDK(model *SettingsModel, s *descope.TenantSettings) {
	model.SessionSettingsEnabled = types.BoolValue(s.SessionSettingsEnabled)
	model.RefreshTokenExpiration = types.Int64Value(int64(s.RefreshTokenExpiration))
	model.RefreshTokenExpirationUnit = types.StringValue(s.RefreshTokenExpirationUnit)
	model.SessionTokenExpiration = types.Int64Value(int64(s.SessionTokenExpiration))
	model.SessionTokenExpirationUnit = types.StringValue(s.SessionTokenExpirationUnit)
	model.StepupTokenExpiration = types.Int64Value(int64(s.StepupTokenExpiration))
	model.StepupTokenExpirationUnit = types.StringValue(s.StepupTokenExpirationUnit)
	model.EnableInactivity = types.BoolValue(s.EnableInactivity)
	model.InactivityTime = types.Int64Value(int64(s.InactivityTime))
	model.InactivityTimeUnit = types.StringValue(s.InactivityTimeUnit)
	model.JITDisabled = types.BoolValue(s.JITDisabled)
}
