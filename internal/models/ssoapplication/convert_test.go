package ssoapplication

import (
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestModelToOIDCRequest(t *testing.T) {
	t.Run("maps all model fields to OIDC request", func(t *testing.T) {
		model := &Model{
			ID:          types.StringValue("app-123"),
			Name:        types.StringValue("Test OIDC App"),
			Description: types.StringValue("A test app"),
			Enabled:     types.BoolValue(true),
			Logo:        types.StringValue("https://example.com/logo.png"),
			OIDC: &OIDCModel{
				LoginPageURL:         types.StringValue("https://app.example.com/login"),
				ForceAuthentication:  types.BoolValue(true),
				BackChannelLogoutURL: types.StringValue("https://app.example.com/logout"),
			},
		}

		req := ModelToOIDCRequest(model)

		assertEqual(t, "ID", req.ID, "app-123")
		assertEqual(t, "Name", req.Name, "Test OIDC App")
		assertEqual(t, "Description", req.Description, "A test app")
		assertBoolEqual(t, "Enabled", req.Enabled, true)
		assertEqual(t, "Logo", req.Logo, "https://example.com/logo.png")
		assertEqual(t, "LoginPageURL", req.LoginPageURL, "https://app.example.com/login")
		assertBoolEqual(t, "ForceAuthentication", req.ForceAuthentication, true)
		assertEqual(t, "BackChannelLogoutURL", req.BackChannelLogoutURL, "https://app.example.com/logout")
	})

	t.Run("handles nil OIDC block", func(t *testing.T) {
		model := &Model{
			ID:   types.StringValue("app-123"),
			Name: types.StringValue("Test App"),
			OIDC: nil,
		}

		req := ModelToOIDCRequest(model)

		assertEqual(t, "ID", req.ID, "app-123")
		assertEqual(t, "LoginPageURL", req.LoginPageURL, "")
	})
}

func TestModelToSAMLRequest(t *testing.T) {
	t.Run("maps all model fields to SAML request", func(t *testing.T) {
		model := &Model{
			ID:          types.StringValue("app-456"),
			Name:        types.StringValue("Test SAML App"),
			Description: types.StringValue("A SAML app"),
			Enabled:     types.BoolValue(false),
			Logo:        types.StringValue("https://example.com/saml-logo.png"),
			SAML: &SAMLModel{
				LoginPageURL:        types.StringValue("https://saml.example.com/login"),
				UseMetadataInfo:     types.BoolValue(true),
				MetadataURL:         types.StringValue("https://saml.example.com/metadata"),
				EntityID:            types.StringValue("https://saml.example.com/entity"),
				AcsURL:              types.StringValue("https://saml.example.com/acs"),
				Certificate:         types.StringValue("CERT_DATA"),
				DefaultRelayState:   types.StringValue("https://saml.example.com/relay"),
				ForceAuthentication: types.BoolValue(true),
				LogoutRedirectURL:   types.StringValue("https://saml.example.com/logout"),
			},
		}

		req := ModelToSAMLRequest(model)

		assertEqual(t, "ID", req.ID, "app-456")
		assertEqual(t, "Name", req.Name, "Test SAML App")
		assertBoolEqual(t, "Enabled", req.Enabled, false)
		assertEqual(t, "LoginPageURL", req.LoginPageURL, "https://saml.example.com/login")
		assertBoolEqual(t, "UseMetadataInfo", req.UseMetadataInfo, true)
		assertEqual(t, "MetadataURL", req.MetadataURL, "https://saml.example.com/metadata")
		assertEqual(t, "EntityID", req.EntityID, "https://saml.example.com/entity")
		assertEqual(t, "AcsURL", req.AcsURL, "https://saml.example.com/acs")
		assertEqual(t, "Certificate", req.Certificate, "CERT_DATA")
		assertEqual(t, "DefaultRelayState", req.DefaultRelayState, "https://saml.example.com/relay")
		assertBoolEqual(t, "ForceAuthentication", req.ForceAuthentication, true)
		assertEqual(t, "LogoutRedirectURL", req.LogoutRedirectURL, "https://saml.example.com/logout")
	})

	t.Run("handles nil SAML block", func(t *testing.T) {
		model := &Model{
			ID:   types.StringValue("app-456"),
			Name: types.StringValue("Test App"),
			SAML: nil,
		}

		req := ModelToSAMLRequest(model)

		assertEqual(t, "ID", req.ID, "app-456")
		assertEqual(t, "LoginPageURL", req.LoginPageURL, "")
	})
}

func TestRefreshModelFromResponse(t *testing.T) {
	t.Run("does not panic on nil app", func(t *testing.T) {
		model := &Model{}
		RefreshModelFromResponse(model, nil)
		// If we get here without a panic, the nil guard works
		if model.ID.ValueString() != "" {
			t.Fatal("expected empty ID for nil app")
		}
	})

	t.Run("refreshes common fields from OIDC response", func(t *testing.T) {
		model := &Model{
			OIDC: &OIDCModel{},
		}
		app := &descope.SSOApplication{
			ID:          "app-123",
			Name:        "My OIDC App",
			Description: "desc",
			Enabled:     true,
			Logo:        "logo.png",
			AppType:     "oidc",
			OIDCSettings: &descope.SSOApplicationOIDCSettings{
				LoginPageURL:         "https://login.example.com",
				ForceAuthentication:  true,
				BackChannelLogoutURL: "https://logout.example.com",
			},
		}

		RefreshModelFromResponse(model, app)

		assertEqual(t, "ID", model.ID.ValueString(), "app-123")
		assertEqual(t, "Name", model.Name.ValueString(), "My OIDC App")
		assertEqual(t, "AppType", model.AppType.ValueString(), "oidc")
		assertEqual(t, "OIDC.LoginPageURL", model.OIDC.LoginPageURL.ValueString(), "https://login.example.com")
		assertBoolEqual(t, "OIDC.ForceAuthentication", model.OIDC.ForceAuthentication.ValueBool(), true)
	})

	t.Run("refreshes SAML fields from response", func(t *testing.T) {
		model := &Model{
			SAML: &SAMLModel{},
		}
		app := &descope.SSOApplication{
			ID:      "app-456",
			Name:    "My SAML App",
			AppType: "saml",
			SAMLSettings: &descope.SSOApplicationSAMLSettings{
				LoginPageURL:    "https://saml.example.com/login",
				UseMetadataInfo: true,
				MetadataURL:     "https://saml.example.com/metadata",
				EntityID:        "entity-id",
				AcsURL:          "https://saml.example.com/acs",
				Certificate:     "CERT",
			},
		}

		RefreshModelFromResponse(model, app)

		assertEqual(t, "SAML.LoginPageURL", model.SAML.LoginPageURL.ValueString(), "https://saml.example.com/login")
		assertBoolEqual(t, "SAML.UseMetadataInfo", model.SAML.UseMetadataInfo.ValueBool(), true)
		assertEqual(t, "SAML.EntityID", model.SAML.EntityID.ValueString(), "entity-id")
	})

	t.Run("does not populate OIDC when model block is nil", func(t *testing.T) {
		model := &Model{
			OIDC: nil,
		}
		app := &descope.SSOApplication{
			ID:      "app-123",
			AppType: "oidc",
			OIDCSettings: &descope.SSOApplicationOIDCSettings{
				LoginPageURL: "https://login.example.com",
			},
		}

		RefreshModelFromResponse(model, app)

		if model.OIDC != nil {
			t.Fatal("expected OIDC block to remain nil when model has no OIDC block")
		}
	})
}

func TestRefreshModelFromResponseForImport(t *testing.T) {
	t.Run("does not panic on nil app", func(t *testing.T) {
		model := &Model{}
		RefreshModelFromResponseForImport(model, nil)
		if model.ID.ValueString() != "" {
			t.Fatal("expected empty ID for nil app")
		}
	})

	t.Run("populates OIDC block for oidc app type", func(t *testing.T) {
		model := &Model{}
		app := &descope.SSOApplication{
			ID:      "app-123",
			Name:    "OIDC App",
			AppType: "oidc",
			OIDCSettings: &descope.SSOApplicationOIDCSettings{
				LoginPageURL:        "https://login.example.com",
				ForceAuthentication: true,
			},
		}

		RefreshModelFromResponseForImport(model, app)

		if model.OIDC == nil {
			t.Fatal("expected OIDC block to be populated on import")
		}
		assertEqual(t, "OIDC.LoginPageURL", model.OIDC.LoginPageURL.ValueString(), "https://login.example.com")
		if model.SAML != nil {
			t.Fatal("expected SAML block to remain nil for OIDC app type")
		}
	})

	t.Run("populates SAML block for saml app type", func(t *testing.T) {
		model := &Model{}
		app := &descope.SSOApplication{
			ID:      "app-456",
			Name:    "SAML App",
			AppType: "saml",
			SAMLSettings: &descope.SSOApplicationSAMLSettings{
				LoginPageURL: "https://saml.example.com/login",
				EntityID:     "entity-id",
				AcsURL:       "https://saml.example.com/acs",
			},
		}

		RefreshModelFromResponseForImport(model, app)

		if model.SAML == nil {
			t.Fatal("expected SAML block to be populated on import")
		}
		assertEqual(t, "SAML.EntityID", model.SAML.EntityID.ValueString(), "entity-id")
		if model.OIDC != nil {
			t.Fatal("expected OIDC block to remain nil for SAML app type")
		}
	})

	t.Run("does not populate both blocks when API returns both settings", func(t *testing.T) {
		model := &Model{}
		app := &descope.SSOApplication{
			ID:      "app-789",
			Name:    "Hybrid App",
			AppType: "oidc",
			OIDCSettings: &descope.SSOApplicationOIDCSettings{
				LoginPageURL: "https://oidc.example.com/login",
			},
			SAMLSettings: &descope.SSOApplicationSAMLSettings{
				LoginPageURL: "https://saml.example.com/login",
			},
		}

		RefreshModelFromResponseForImport(model, app)

		if model.OIDC == nil {
			t.Fatal("expected OIDC block to be populated for oidc app type")
		}
		if model.SAML != nil {
			t.Fatal("expected SAML block to remain nil even when API returns both settings")
		}
	})
}

func assertEqual(t *testing.T, field string, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: expected %q, got %q", field, want, got)
	}
}

func assertBoolEqual(t *testing.T, field string, got, want bool) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: expected %v, got %v", field, want, got)
	}
}
