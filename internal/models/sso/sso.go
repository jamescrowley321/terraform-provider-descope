package sso

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
)

var samlAttributeMappingAttributes = map[string]schema.Attribute{
	"name":         stringattr.Default(""),
	"given_name":   stringattr.Default(""),
	"middle_name":  stringattr.Default(""),
	"family_name":  stringattr.Default(""),
	"picture":      stringattr.Default(""),
	"email":        stringattr.Default(""),
	"phone_number": stringattr.Default(""),
	"group":        stringattr.Default(""),
}

var Attributes = map[string]schema.Attribute{
	"id": stringattr.Identifier(),
	"tenant_id": schema.StringAttribute{
		Required:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	},
	"sso_id": schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
	},
	"display_name": stringattr.Default(""),
	"domains":      strsetattr.Default(),
	"oidc": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: oidcAttributes,
	},
	"saml": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: samlAttributes,
	},
	"saml_metadata": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: samlMetadataAttributes,
	},
}

var oidcAttributeMappingAttributes = map[string]schema.Attribute{
	"login_id":       stringattr.Default(""),
	"name":           stringattr.Default(""),
	"given_name":     stringattr.Default(""),
	"middle_name":    stringattr.Default(""),
	"family_name":    stringattr.Default(""),
	"email":          stringattr.Default(""),
	"verified_email": stringattr.Default(""),
	"username":       stringattr.Default(""),
	"phone_number":   stringattr.Default(""),
	"verified_phone": stringattr.Default(""),
	"picture":        stringattr.Default(""),
}

var oidcAttributes = map[string]schema.Attribute{
	"name":      stringattr.Required(),
	"client_id": stringattr.Required(),
	"client_secret": schema.StringAttribute{
		Optional:  true,
		Sensitive: true,
	},
	"redirect_url":           stringattr.Default(""),
	"auth_url":               stringattr.Default(""),
	"token_url":              stringattr.Default(""),
	"user_data_url":          stringattr.Default(""),
	"jwks_url":               stringattr.Default(""),
	"callback_domain":        stringattr.Default(""),
	"grant_type":             stringattr.Default(""),
	"issuer":                 stringattr.Default(""),
	"scope":                  strsetattr.Default(),
	"manage_provider_tokens": boolattr.Default(false),
	"attribute_mapping": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: oidcAttributeMappingAttributes,
	},
}

var samlAttributes = map[string]schema.Attribute{
	"idp_url":       stringattr.Required(),
	"idp_entity_id": stringattr.Required(),
	"idp_cert":      stringattr.Required(),
	"redirect_url":  stringattr.Default(""),
	"sp_entity_id":  stringattr.Identifier(),
	"sp_acs_url":    stringattr.Identifier(),
	"attribute_mapping": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: samlAttributeMappingAttributes,
	},
}

var samlMetadataAttributes = map[string]schema.Attribute{
	"idp_metadata_url": stringattr.Required(),
	"redirect_url":     stringattr.Default(""),
	"sp_entity_id":     stringattr.Identifier(),
	"sp_acs_url":       stringattr.Identifier(),
	"attribute_mapping": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: samlAttributeMappingAttributes,
	},
}

type Model struct {
	ID           stringattr.Type `tfsdk:"id"`
	TenantID     stringattr.Type `tfsdk:"tenant_id"`
	SSOID        stringattr.Type `tfsdk:"sso_id"`
	DisplayName  stringattr.Type `tfsdk:"display_name"`
	Domains      strsetattr.Type `tfsdk:"domains"`
	OIDC         *OIDCModel      `tfsdk:"oidc"`
	SAML         *SAMLModel      `tfsdk:"saml"`
	SAMLMetadata *SAMLMetaModel  `tfsdk:"saml_metadata"`
}

type OIDCModel struct {
	Name                 stringattr.Type            `tfsdk:"name"`
	ClientID             stringattr.Type            `tfsdk:"client_id"`
	ClientSecret         stringattr.Type            `tfsdk:"client_secret"`
	RedirectURL          stringattr.Type            `tfsdk:"redirect_url"`
	AuthURL              stringattr.Type            `tfsdk:"auth_url"`
	TokenURL             stringattr.Type            `tfsdk:"token_url"`
	UserDataURL          stringattr.Type            `tfsdk:"user_data_url"`
	JWKsURL              stringattr.Type            `tfsdk:"jwks_url"`
	CallbackDomain       stringattr.Type            `tfsdk:"callback_domain"`
	GrantType            stringattr.Type            `tfsdk:"grant_type"`
	Issuer               stringattr.Type            `tfsdk:"issuer"`
	Scope                strsetattr.Type            `tfsdk:"scope"`
	ManageProviderTokens boolattr.Type              `tfsdk:"manage_provider_tokens"`
	AttributeMapping     *OIDCAttributeMappingModel `tfsdk:"attribute_mapping"`
}

type OIDCAttributeMappingModel struct {
	LoginID       stringattr.Type `tfsdk:"login_id"`
	Name          stringattr.Type `tfsdk:"name"`
	GivenName     stringattr.Type `tfsdk:"given_name"`
	MiddleName    stringattr.Type `tfsdk:"middle_name"`
	FamilyName    stringattr.Type `tfsdk:"family_name"`
	Email         stringattr.Type `tfsdk:"email"`
	VerifiedEmail stringattr.Type `tfsdk:"verified_email"`
	Username      stringattr.Type `tfsdk:"username"`
	PhoneNumber   stringattr.Type `tfsdk:"phone_number"`
	VerifiedPhone stringattr.Type `tfsdk:"verified_phone"`
	Picture       stringattr.Type `tfsdk:"picture"`
}

type SAMLModel struct {
	IdpURL           stringattr.Type        `tfsdk:"idp_url"`
	IdpEntityID      stringattr.Type        `tfsdk:"idp_entity_id"`
	IdpCert          stringattr.Type        `tfsdk:"idp_cert"`
	RedirectURL      stringattr.Type        `tfsdk:"redirect_url"`
	SpEntityID       stringattr.Type        `tfsdk:"sp_entity_id"`
	SpACSUrl         stringattr.Type        `tfsdk:"sp_acs_url"`
	AttributeMapping *AttributeMappingModel `tfsdk:"attribute_mapping"`
}

type SAMLMetaModel struct {
	IdpMetadataURL   stringattr.Type        `tfsdk:"idp_metadata_url"`
	RedirectURL      stringattr.Type        `tfsdk:"redirect_url"`
	SpEntityID       stringattr.Type        `tfsdk:"sp_entity_id"`
	SpACSUrl         stringattr.Type        `tfsdk:"sp_acs_url"`
	AttributeMapping *AttributeMappingModel `tfsdk:"attribute_mapping"`
}

type AttributeMappingModel struct {
	Name        stringattr.Type `tfsdk:"name"`
	GivenName   stringattr.Type `tfsdk:"given_name"`
	MiddleName  stringattr.Type `tfsdk:"middle_name"`
	FamilyName  stringattr.Type `tfsdk:"family_name"`
	Picture     stringattr.Type `tfsdk:"picture"`
	Email       stringattr.Type `tfsdk:"email"`
	PhoneNumber stringattr.Type `tfsdk:"phone_number"`
	Group       stringattr.Type `tfsdk:"group"`
}
