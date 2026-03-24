package outboundapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var Attributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"client_id":   stringattr.Default(""),
	"client_secret": schema.StringAttribute{
		Optional:  true,
		Sensitive: true,
	},
	"logo":                 stringattr.Default(""),
	"discovery_url":        stringattr.Default(""),
	"authorization_url":    stringattr.Default(""),
	"token_url":            stringattr.Default(""),
	"revocation_url":       stringattr.Default(""),
	"default_scopes":       strsetattr.Default(),
	"default_redirect_url": stringattr.Default(""),
	"callback_domain":      stringattr.Default(""),
	"pkce":                 boolattr.Default(false),
}

type Model struct {
	ID                 stringattr.Type `tfsdk:"id"`
	Name               stringattr.Type `tfsdk:"name"`
	Description        stringattr.Type `tfsdk:"description"`
	ClientID           stringattr.Type `tfsdk:"client_id"`
	ClientSecret       stringattr.Type `tfsdk:"client_secret"`
	Logo               stringattr.Type `tfsdk:"logo"`
	DiscoveryURL       stringattr.Type `tfsdk:"discovery_url"`
	AuthorizationURL   stringattr.Type `tfsdk:"authorization_url"`
	TokenURL           stringattr.Type `tfsdk:"token_url"`
	RevocationURL      stringattr.Type `tfsdk:"revocation_url"`
	DefaultScopes      strsetattr.Type `tfsdk:"default_scopes"`
	DefaultRedirectURL stringattr.Type `tfsdk:"default_redirect_url"`
	CallbackDomain     stringattr.Type `tfsdk:"callback_domain"`
	Pkce               boolattr.Type   `tfsdk:"pkce"`
}
