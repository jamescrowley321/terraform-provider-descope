package ssoapplication

import (
	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ModelToOIDCRequest(model *Model) *descope.OIDCApplicationRequest {
	req := &descope.OIDCApplicationRequest{
		ID:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Enabled:     model.Enabled.ValueBool(),
		Logo:        model.Logo.ValueString(),
	}
	if model.OIDC != nil {
		req.LoginPageURL = model.OIDC.LoginPageURL.ValueString()
		req.ForceAuthentication = model.OIDC.ForceAuthentication.ValueBool()
		req.BackChannelLogoutURL = model.OIDC.BackChannelLogoutURL.ValueString()
	}
	return req
}

func ModelToSAMLRequest(model *Model) *descope.SAMLApplicationRequest {
	req := &descope.SAMLApplicationRequest{
		ID:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Enabled:     model.Enabled.ValueBool(),
		Logo:        model.Logo.ValueString(),
	}
	if model.SAML != nil {
		req.LoginPageURL = model.SAML.LoginPageURL.ValueString()
		req.UseMetadataInfo = model.SAML.UseMetadataInfo.ValueBool()
		req.MetadataURL = model.SAML.MetadataURL.ValueString()
		req.EntityID = model.SAML.EntityID.ValueString()
		req.AcsURL = model.SAML.AcsURL.ValueString()
		req.Certificate = model.SAML.Certificate.ValueString()
		req.DefaultRelayState = model.SAML.DefaultRelayState.ValueString()
		req.ForceAuthentication = model.SAML.ForceAuthentication.ValueBool()
		req.LogoutRedirectURL = model.SAML.LogoutRedirectURL.ValueString()
	}
	return req
}

func RefreshModelFromResponse(model *Model, app *descope.SSOApplication) {
	model.ID = types.StringValue(app.ID)
	model.Name = types.StringValue(app.Name)
	model.Description = types.StringValue(app.Description)
	model.Enabled = types.BoolValue(app.Enabled)
	model.Logo = types.StringValue(app.Logo)
	model.AppType = types.StringValue(app.AppType)

	if app.OIDCSettings != nil && model.OIDC != nil {
		model.OIDC.LoginPageURL = types.StringValue(app.OIDCSettings.LoginPageURL)
		model.OIDC.ForceAuthentication = types.BoolValue(app.OIDCSettings.ForceAuthentication)
		model.OIDC.BackChannelLogoutURL = types.StringValue(app.OIDCSettings.BackChannelLogoutURL)
	}

	if app.SAMLSettings != nil && model.SAML != nil {
		refreshSAMLFromResponse(model.SAML, app.SAMLSettings)
	}
}

func RefreshModelFromResponseForImport(model *Model, app *descope.SSOApplication) {
	model.ID = types.StringValue(app.ID)
	model.Name = types.StringValue(app.Name)
	model.Description = types.StringValue(app.Description)
	model.Enabled = types.BoolValue(app.Enabled)
	model.Logo = types.StringValue(app.Logo)
	model.AppType = types.StringValue(app.AppType)

	if app.OIDCSettings != nil {
		if model.OIDC == nil {
			model.OIDC = &OIDCModel{}
		}
		model.OIDC.LoginPageURL = types.StringValue(app.OIDCSettings.LoginPageURL)
		model.OIDC.ForceAuthentication = types.BoolValue(app.OIDCSettings.ForceAuthentication)
		model.OIDC.BackChannelLogoutURL = types.StringValue(app.OIDCSettings.BackChannelLogoutURL)
	}

	if app.SAMLSettings != nil {
		if model.SAML == nil {
			model.SAML = &SAMLModel{}
		}
		refreshSAMLFromResponse(model.SAML, app.SAMLSettings)
	}
}

func refreshSAMLFromResponse(model *SAMLModel, s *descope.SSOApplicationSAMLSettings) {
	model.LoginPageURL = types.StringValue(s.LoginPageURL)
	model.UseMetadataInfo = types.BoolValue(s.UseMetadataInfo)
	model.MetadataURL = types.StringValue(s.MetadataURL)
	model.EntityID = types.StringValue(s.EntityID)
	model.AcsURL = types.StringValue(s.AcsURL)
	model.Certificate = types.StringValue(s.Certificate)
	model.DefaultRelayState = types.StringValue(s.DefaultRelayState)
	model.ForceAuthentication = types.BoolValue(s.ForceAuthentication)
	model.LogoutRedirectURL = types.StringValue(s.LogoutRedirectURL)
}
