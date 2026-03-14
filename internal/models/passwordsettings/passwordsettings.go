package passwordsettings

import (
	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var PasswordSettingsAttributes = map[string]schema.Attribute{
	"id":               stringattr.Identifier(),
	"enabled":          boolattr.Default(true),
	"min_length":       intattr.Default(8, int64validator.Between(4, 64)),
	"lowercase":        boolattr.Default(false),
	"uppercase":        boolattr.Default(false),
	"number":           boolattr.Default(false),
	"non_alphanumeric": boolattr.Default(false),
	"expiration":       boolattr.Default(false),
	"expiration_weeks": intattr.Default(1, int64validator.Between(1, 999)),
	"reuse":            boolattr.Default(false),
	"reuse_amount":     intattr.Default(1, int64validator.Between(1, 50)),
	"lock":             boolattr.Default(false),
	"lock_attempts":    intattr.Default(2, int64validator.Between(2, 10)),
}

type PasswordSettingsModel struct {
	ID              stringattr.Type `tfsdk:"id"`
	Enabled         boolattr.Type   `tfsdk:"enabled"`
	MinLength       intattr.Type    `tfsdk:"min_length"`
	Lowercase       boolattr.Type   `tfsdk:"lowercase"`
	Uppercase       boolattr.Type   `tfsdk:"uppercase"`
	Number          boolattr.Type   `tfsdk:"number"`
	NonAlphanumeric boolattr.Type   `tfsdk:"non_alphanumeric"`
	Expiration      boolattr.Type   `tfsdk:"expiration"`
	ExpirationWeeks intattr.Type    `tfsdk:"expiration_weeks"`
	Reuse           boolattr.Type   `tfsdk:"reuse"`
	ReuseAmount     intattr.Type    `tfsdk:"reuse_amount"`
	Lock            boolattr.Type   `tfsdk:"lock"`
	LockAttempts    intattr.Type    `tfsdk:"lock_attempts"`
}

// ToSDK converts the Terraform model to the Descope SDK PasswordSettings type.
func (m *PasswordSettingsModel) ToSDK() *descope.PasswordSettings {
	return &descope.PasswordSettings{
		Enabled:         m.Enabled.ValueBool(),
		MinLength:       int32(m.MinLength.ValueInt64()),
		Lowercase:       m.Lowercase.ValueBool(),
		Uppercase:       m.Uppercase.ValueBool(),
		Number:          m.Number.ValueBool(),
		NonAlphanumeric: m.NonAlphanumeric.ValueBool(),
		Expiration:      m.Expiration.ValueBool(),
		ExpirationWeeks: int32(m.ExpirationWeeks.ValueInt64()),
		Reuse:           m.Reuse.ValueBool(),
		ReuseAmount:     int32(m.ReuseAmount.ValueInt64()),
		Lock:            m.Lock.ValueBool(),
		LockAttempts:    int32(m.LockAttempts.ValueInt64()),
	}
}

// SetFromSDK populates the Terraform model from the Descope SDK PasswordSettings type.
func (m *PasswordSettingsModel) SetFromSDK(settings *descope.PasswordSettings) {
	m.Enabled = types.BoolValue(settings.Enabled)
	m.MinLength = types.Int64Value(int64(settings.MinLength))
	m.Lowercase = types.BoolValue(settings.Lowercase)
	m.Uppercase = types.BoolValue(settings.Uppercase)
	m.Number = types.BoolValue(settings.Number)
	m.NonAlphanumeric = types.BoolValue(settings.NonAlphanumeric)
	m.Expiration = types.BoolValue(settings.Expiration)
	m.ExpirationWeeks = types.Int64Value(int64(settings.ExpirationWeeks))
	m.Reuse = types.BoolValue(settings.Reuse)
	m.ReuseAmount = types.Int64Value(int64(settings.ReuseAmount))
	m.Lock = types.BoolValue(settings.Lock)
	m.LockAttempts = types.Int64Value(int64(settings.LockAttempts))
}
