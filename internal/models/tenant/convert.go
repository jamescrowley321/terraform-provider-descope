package tenant

import (
	"context"
	"fmt"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSetToSlice converts a Terraform string set to a Go string slice.
func StringSetToSlice(_ context.Context, s strsetattr.Type, _ *diag.Diagnostics) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elems := s.Elements()
	result := make([]string, 0, len(elems))
	for _, v := range elems {
		if str, ok := v.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

// StringMapToAnyMap converts a Terraform string map to a map[string]any for SDK calls.
func StringMapToAnyMap(m strmapattr.Type) map[string]any {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	elems := m.Elements()
	if len(elems) == 0 {
		return nil
	}
	result := make(map[string]any, len(elems))
	for k, v := range elems {
		if str, ok := v.(types.String); ok {
			result[k] = str.ValueString()
		}
	}
	return result
}

// anyMapToStringMap converts a map[string]any from the SDK to a map[string]string.
func anyMapToStringMap(m map[string]any) map[string]string {
	if m == nil {
		return nil
	}
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}

// ModelToRequest converts the TenantModel to an SDK TenantRequest for Create/Update.
func ModelToRequest(ctx context.Context, model *TenantModel, diags *diag.Diagnostics) *descope.TenantRequest {
	return &descope.TenantRequest{
		Name:                    model.Name.ValueString(),
		SelfProvisioningDomains: StringSetToSlice(ctx, model.SelfProvisioningDomains, diags),
		CustomAttributes:        StringMapToAnyMap(model.CustomAttributes),
		EnforceSSO:              model.EnforceSSO.ValueBool(),
		Disabled:                model.Disabled.ValueBool(),
		ParentTenantID:          model.ParentTenantID.ValueString(),
		EnforceSSOExclusions:    StringSetToSlice(ctx, model.EnforceSSOExclusions, diags),
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
		model.CustomAttributes = strmapattr.Value(anyMapToStringMap(t.CustomAttributes))
	}
}

// ModelToSettings converts the SettingsModel to an SDK TenantSettings.
// It also populates the overlapping fields (Domains, SelfProvisioningDomains, AuthType)
// from the top-level tenant model.
func ModelToSettings(ctx context.Context, model *TenantModel, diags *diag.Diagnostics) *descope.TenantSettings {
	s := model.Settings
	return &descope.TenantSettings{
		Domains:                    StringSetToSlice(ctx, model.Domains, diags),
		SelfProvisioningDomains:    StringSetToSlice(ctx, model.SelfProvisioningDomains, diags),
		AuthType:                   model.AuthType.ValueString(),
		SessionSettingsEnabled:     s.SessionSettingsEnabled.ValueBool(),
		RefreshTokenExpiration:     int32(s.RefreshTokenExpiration.ValueInt64()),
		RefreshTokenExpirationUnit: s.RefreshTokenExpirationUnit.ValueString(),
		SessionTokenExpiration:     int32(s.SessionTokenExpiration.ValueInt64()),
		SessionTokenExpirationUnit: s.SessionTokenExpirationUnit.ValueString(),
		StepupTokenExpiration:      int32(s.StepupTokenExpiration.ValueInt64()),
		StepupTokenExpirationUnit:  s.StepupTokenExpirationUnit.ValueString(),
		EnableInactivity:           s.EnableInactivity.ValueBool(),
		InactivityTime:             int32(s.InactivityTime.ValueInt64()),
		InactivityTimeUnit:         s.InactivityTimeUnit.ValueString(),
		JITDisabled:                s.JITDisabled.ValueBool(),
	}
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
