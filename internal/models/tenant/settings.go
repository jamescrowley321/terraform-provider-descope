package tenant

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
)

var timeUnitValidator = stringvalidator.OneOf("", "seconds", "minutes", "hours", "days", "weeks")

var SettingsAttributes = map[string]schema.Attribute{
	"session_settings_enabled":      boolattr.Default(false),
	"refresh_token_expiration":      intattr.Default(0),
	"refresh_token_expiration_unit": stringattr.Default("", timeUnitValidator),
	"session_token_expiration":      intattr.Default(0),
	"session_token_expiration_unit": stringattr.Default("", timeUnitValidator),
	"stepup_token_expiration":       intattr.Default(0),
	"stepup_token_expiration_unit":  stringattr.Default("", timeUnitValidator),
	"enable_inactivity":             boolattr.Default(false),
	"inactivity_time":               intattr.Default(0),
	"inactivity_time_unit":          stringattr.Default("", timeUnitValidator),
	"jit_disabled":                  boolattr.Default(false),
}

type SettingsModel struct {
	SessionSettingsEnabled     boolattr.Type   `tfsdk:"session_settings_enabled"`
	RefreshTokenExpiration     intattr.Type    `tfsdk:"refresh_token_expiration"`
	RefreshTokenExpirationUnit stringattr.Type `tfsdk:"refresh_token_expiration_unit"`
	SessionTokenExpiration     intattr.Type    `tfsdk:"session_token_expiration"`
	SessionTokenExpirationUnit stringattr.Type `tfsdk:"session_token_expiration_unit"`
	StepupTokenExpiration      intattr.Type    `tfsdk:"stepup_token_expiration"`
	StepupTokenExpirationUnit  stringattr.Type `tfsdk:"stepup_token_expiration_unit"`
	EnableInactivity           boolattr.Type   `tfsdk:"enable_inactivity"`
	InactivityTime             intattr.Type    `tfsdk:"inactivity_time"`
	InactivityTimeUnit         stringattr.Type `tfsdk:"inactivity_time_unit"`
	JITDisabled                boolattr.Type   `tfsdk:"jit_disabled"`
}
