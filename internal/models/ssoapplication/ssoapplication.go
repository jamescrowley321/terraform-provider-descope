package ssoapplication

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var Attributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"enabled":     boolattr.Default(true),
	"logo":        stringattr.Default(""),
	"app_type":    stringattr.Identifier(),
	"oidc": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: oidcAttributes,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplaceIf(requiresReplaceIfBlockToggled,
				"Requires replace when switching between OIDC and SAML application types.",
				"Requires replace when switching between OIDC and SAML application types.",
			),
		},
	},
	"saml": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: samlAttributes,
		PlanModifiers: []planmodifier.Object{
			objectplanmodifier.RequiresReplaceIf(requiresReplaceIfBlockToggled,
				"Requires replace when switching between OIDC and SAML application types.",
				"Requires replace when switching between OIDC and SAML application types.",
			),
		},
	},
}

var oidcAttributes = map[string]schema.Attribute{
	"login_page_url":          stringattr.Default(""),
	"force_authentication":    boolattr.Default(false),
	"back_channel_logout_url": stringattr.Default(""),
}

var samlAttributes = map[string]schema.Attribute{
	"login_page_url":       stringattr.Default(""),
	"use_metadata_info":    boolattr.Default(false),
	"metadata_url":         stringattr.Default(""),
	"entity_id":            stringattr.Default(""),
	"acs_url":              stringattr.Default(""),
	"certificate":          stringattr.Default(""),
	"default_relay_state":  stringattr.Default(""),
	"force_authentication": boolattr.Default(false),
	"logout_redirect_url":  stringattr.Default(""),
}

type Model struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Enabled     boolattr.Type   `tfsdk:"enabled"`
	Logo        stringattr.Type `tfsdk:"logo"`
	AppType     stringattr.Type `tfsdk:"app_type"`
	OIDC        *OIDCModel      `tfsdk:"oidc"`
	SAML        *SAMLModel      `tfsdk:"saml"`
}

type OIDCModel struct {
	LoginPageURL         stringattr.Type `tfsdk:"login_page_url"`
	ForceAuthentication  boolattr.Type   `tfsdk:"force_authentication"`
	BackChannelLogoutURL stringattr.Type `tfsdk:"back_channel_logout_url"`
}

// requiresReplaceIfBlockToggled triggers resource replacement when the block
// transitions between null and non-null (i.e., switching from OIDC to SAML or vice versa).
// Updates within the same block type do not trigger replacement.
func requiresReplaceIfBlockToggled(_ context.Context, req planmodifier.ObjectRequest, resp *objectplanmodifier.RequiresReplaceIfFuncResponse) {
	resp.RequiresReplace = req.StateValue.IsNull() != req.PlanValue.IsNull()
}

type SAMLModel struct {
	LoginPageURL        stringattr.Type `tfsdk:"login_page_url"`
	UseMetadataInfo     boolattr.Type   `tfsdk:"use_metadata_info"`
	MetadataURL         stringattr.Type `tfsdk:"metadata_url"`
	EntityID            stringattr.Type `tfsdk:"entity_id"`
	AcsURL              stringattr.Type `tfsdk:"acs_url"`
	Certificate         stringattr.Type `tfsdk:"certificate"`
	DefaultRelayState   stringattr.Type `tfsdk:"default_relay_state"`
	ForceAuthentication boolattr.Type   `tfsdk:"force_authentication"`
	LogoutRedirectURL   stringattr.Type `tfsdk:"logout_redirect_url"`
}
