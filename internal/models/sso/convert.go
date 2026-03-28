package sso

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/convert"
)

func ModelToOIDCSettings(ctx context.Context, m *OIDCModel, diags *diag.Diagnostics) *descope.SSOOIDCSettings {
	settings := &descope.SSOOIDCSettings{
		Name:                 m.Name.ValueString(),
		ClientID:             m.ClientID.ValueString(),
		ClientSecret:         m.ClientSecret.ValueString(),
		RedirectURL:          m.RedirectURL.ValueString(),
		AuthURL:              m.AuthURL.ValueString(),
		TokenURL:             m.TokenURL.ValueString(),
		UserDataURL:          m.UserDataURL.ValueString(),
		JWKsURL:              m.JWKsURL.ValueString(),
		CallbackDomain:       m.CallbackDomain.ValueString(),
		GrantType:            m.GrantType.ValueString(),
		Issuer:               m.Issuer.ValueString(),
		Scope:                convert.StringSetToSlice(ctx, m.Scope, diags),
		ManageProviderTokens: m.ManageProviderTokens.ValueBool(),
	}
	if m.AttributeMapping != nil {
		settings.AttributeMapping = modelToOIDCAttributeMapping(m.AttributeMapping)
	}
	return settings
}

func ModelToSAMLSettings(m *SAMLModel) (*descope.SSOSAMLSettings, string) {
	settings := &descope.SSOSAMLSettings{
		IdpURL:      m.IdpURL.ValueString(),
		IdpEntityID: m.IdpEntityID.ValueString(),
		IdpCert:     m.IdpCert.ValueString(),
	}
	if m.AttributeMapping != nil {
		settings.AttributeMapping = modelToAttributeMapping(m.AttributeMapping)
	}
	return settings, m.RedirectURL.ValueString()
}

func ModelToSAMLMetadataSettings(m *SAMLMetaModel) (*descope.SSOSAMLSettingsByMetadata, string) {
	settings := &descope.SSOSAMLSettingsByMetadata{
		IdpMetadataURL: m.IdpMetadataURL.ValueString(),
	}
	if m.AttributeMapping != nil {
		settings.AttributeMapping = modelToAttributeMapping(m.AttributeMapping)
	}
	return settings, m.RedirectURL.ValueString()
}

func RefreshOIDCFromResponse(ctx context.Context, m *OIDCModel, o *descope.SSOOIDCSettings) {
	m.Name = types.StringValue(o.Name)
	m.ClientID = types.StringValue(o.ClientID)
	// client_secret is not returned by the API — preserve plan value
	m.RedirectURL = types.StringValue(o.RedirectURL)
	m.AuthURL = types.StringValue(o.AuthURL)
	m.TokenURL = types.StringValue(o.TokenURL)
	m.UserDataURL = types.StringValue(o.UserDataURL)
	m.JWKsURL = types.StringValue(o.JWKsURL)
	m.CallbackDomain = types.StringValue(o.CallbackDomain)
	m.GrantType = types.StringValue(o.GrantType)
	m.Issuer = types.StringValue(o.Issuer)
	m.Scope = strsetattr.ValueCtx(ctx, o.Scope)
	m.ManageProviderTokens = types.BoolValue(o.ManageProviderTokens)
	if o.AttributeMapping != nil {
		if m.AttributeMapping == nil {
			m.AttributeMapping = &OIDCAttributeMappingModel{}
		}
		refreshOIDCAttributeMapping(m.AttributeMapping, o.AttributeMapping)
	}
}

func RefreshSAMLFromResponse(m *SAMLModel, s *descope.SSOSAMLSettingsResponse) {
	m.IdpURL = types.StringValue(s.IdpSSOURL)
	m.IdpEntityID = types.StringValue(s.IdpEntityID)
	m.IdpCert = types.StringValue(s.IdpCertificate)
	m.RedirectURL = types.StringValue(s.RedirectURL)
	m.SpEntityID = types.StringValue(s.SpEntityID)
	m.SpACSUrl = types.StringValue(s.SpACSUrl)
	if s.AttributeMapping != nil {
		if m.AttributeMapping == nil {
			m.AttributeMapping = &AttributeMappingModel{}
		}
		refreshAttributeMapping(m.AttributeMapping, s.AttributeMapping)
	}
}

func RefreshSAMLMetaFromResponse(m *SAMLMetaModel, s *descope.SSOSAMLSettingsResponse) {
	m.IdpMetadataURL = types.StringValue(s.IdpMetadataURL)
	m.RedirectURL = types.StringValue(s.RedirectURL)
	m.SpEntityID = types.StringValue(s.SpEntityID)
	m.SpACSUrl = types.StringValue(s.SpACSUrl)
	if s.AttributeMapping != nil {
		if m.AttributeMapping == nil {
			m.AttributeMapping = &AttributeMappingModel{}
		}
		refreshAttributeMapping(m.AttributeMapping, s.AttributeMapping)
	}
}

func modelToAttributeMapping(m *AttributeMappingModel) *descope.AttributeMapping {
	return &descope.AttributeMapping{
		Name:        m.Name.ValueString(),
		GivenName:   m.GivenName.ValueString(),
		MiddleName:  m.MiddleName.ValueString(),
		FamilyName:  m.FamilyName.ValueString(),
		Picture:     m.Picture.ValueString(),
		Email:       m.Email.ValueString(),
		PhoneNumber: m.PhoneNumber.ValueString(),
		Group:       m.Group.ValueString(),
	}
}

func refreshAttributeMapping(m *AttributeMappingModel, a *descope.AttributeMapping) {
	m.Name = types.StringValue(a.Name)
	m.GivenName = types.StringValue(a.GivenName)
	m.MiddleName = types.StringValue(a.MiddleName)
	m.FamilyName = types.StringValue(a.FamilyName)
	m.Picture = types.StringValue(a.Picture)
	m.Email = types.StringValue(a.Email)
	m.PhoneNumber = types.StringValue(a.PhoneNumber)
	m.Group = types.StringValue(a.Group)
}

func modelToOIDCAttributeMapping(m *OIDCAttributeMappingModel) *descope.OIDCAttributeMapping {
	return &descope.OIDCAttributeMapping{
		LoginID:       m.LoginID.ValueString(),
		Name:          m.Name.ValueString(),
		GivenName:     m.GivenName.ValueString(),
		MiddleName:    m.MiddleName.ValueString(),
		FamilyName:    m.FamilyName.ValueString(),
		Email:         m.Email.ValueString(),
		VerifiedEmail: m.VerifiedEmail.ValueString(),
		Username:      m.Username.ValueString(),
		PhoneNumber:   m.PhoneNumber.ValueString(),
		VerifiedPhone: m.VerifiedPhone.ValueString(),
		Picture:       m.Picture.ValueString(),
	}
}

func refreshOIDCAttributeMapping(m *OIDCAttributeMappingModel, a *descope.OIDCAttributeMapping) {
	m.LoginID = types.StringValue(a.LoginID)
	m.Name = types.StringValue(a.Name)
	m.GivenName = types.StringValue(a.GivenName)
	m.MiddleName = types.StringValue(a.MiddleName)
	m.FamilyName = types.StringValue(a.FamilyName)
	m.Email = types.StringValue(a.Email)
	m.VerifiedEmail = types.StringValue(a.VerifiedEmail)
	m.Username = types.StringValue(a.Username)
	m.PhoneNumber = types.StringValue(a.PhoneNumber)
	m.VerifiedPhone = types.StringValue(a.VerifiedPhone)
	m.Picture = types.StringValue(a.Picture)
}
