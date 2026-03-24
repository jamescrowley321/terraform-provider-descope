package outboundapp

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ModelToCreateRequest(ctx context.Context, model *Model, diags *diag.Diagnostics) *descope.CreateOutboundAppRequest {
	return &descope.CreateOutboundAppRequest{
		OutboundApp:  *modelToOutboundApp(ctx, model, diags),
		ClientSecret: model.ClientSecret.ValueString(),
	}
}

func ModelToOutboundApp(ctx context.Context, model *Model, diags *diag.Diagnostics) *descope.OutboundApp {
	return modelToOutboundApp(ctx, model, diags)
}

func ClientSecretPtr(model *Model) *string {
	if model.ClientSecret.IsNull() || model.ClientSecret.IsUnknown() || model.ClientSecret.ValueString() == "" {
		return nil
	}
	s := model.ClientSecret.ValueString()
	return &s
}

func RefreshModelFromResponse(ctx context.Context, model *Model, app *descope.OutboundApp) {
	model.ID = types.StringValue(app.ID)
	model.Name = types.StringValue(app.Name)
	model.Description = types.StringValue(app.Description)
	model.ClientID = types.StringValue(app.ClientID)
	// client_secret is write-only — not returned by API, preserve plan value
	model.Logo = types.StringValue(app.Logo)
	model.DiscoveryURL = types.StringValue(app.DiscoveryURL)
	model.AuthorizationURL = types.StringValue(app.AuthorizationURL)
	model.TokenURL = types.StringValue(app.TokenURL)
	model.RevocationURL = types.StringValue(app.RevocationURL)
	model.DefaultScopes = strsetattr.ValueCtx(ctx, app.DefaultScopes)
	model.DefaultRedirectURL = types.StringValue(app.DefaultRedirectURL)
	model.CallbackDomain = types.StringValue(app.CallbackDomain)
	model.Pkce = types.BoolValue(app.Pkce)
}

func modelToOutboundApp(ctx context.Context, model *Model, diags *diag.Diagnostics) *descope.OutboundApp {
	return &descope.OutboundApp{
		ID:                 model.ID.ValueString(),
		Name:               model.Name.ValueString(),
		Description:        model.Description.ValueString(),
		ClientID:           model.ClientID.ValueString(),
		Logo:               model.Logo.ValueString(),
		DiscoveryURL:       model.DiscoveryURL.ValueString(),
		AuthorizationURL:   model.AuthorizationURL.ValueString(),
		TokenURL:           model.TokenURL.ValueString(),
		RevocationURL:      model.RevocationURL.ValueString(),
		DefaultScopes:      convert.StringSetToSlice(ctx, model.DefaultScopes, diags),
		DefaultRedirectURL: model.DefaultRedirectURL.ValueString(),
		CallbackDomain:     model.CallbackDomain.ValueString(),
		Pkce:               model.Pkce.ValueBool(),
	}
}
