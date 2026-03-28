package thirdpartyapp

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelToRequest(t *testing.T) {
	ctx := context.Background()
	diags := diag.Diagnostics{}

	model := &Model{
		ID:                   types.StringValue("app-123"),
		Name:                 types.StringValue("Test App"),
		Description:          types.StringValue("A test application"),
		Logo:                 types.StringValue("https://example.com/logo.png"),
		LoginPageURL:         types.StringValue("https://example.com/login"),
		ApprovedCallbackUrls: strsetattr.Value([]string{"https://example.com/callback", "https://example.com/auth"}),
	}

	req := ModelToRequest(ctx, model, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "app-123", req.ID)
	assert.Equal(t, "Test App", req.Name)
	assert.Equal(t, "A test application", req.Description)
	assert.Equal(t, "https://example.com/logo.png", req.Logo)
	assert.Equal(t, "https://example.com/login", req.LoginPageURL)
	assert.ElementsMatch(t, []string{"https://example.com/callback", "https://example.com/auth"}, req.ApprovedCallbackUrls)
}

func TestModelToRequestEmpty(t *testing.T) {
	ctx := context.Background()
	diags := diag.Diagnostics{}

	model := &Model{
		ID:                   types.StringValue(""),
		Name:                 types.StringValue("Minimal App"),
		Description:          types.StringValue(""),
		Logo:                 types.StringValue(""),
		LoginPageURL:         types.StringValue(""),
		ApprovedCallbackUrls: strsetattr.Value([]string{}),
	}

	req := ModelToRequest(ctx, model, &diags)

	require.False(t, diags.HasError())
	assert.Equal(t, "Minimal App", req.Name)
	assert.Empty(t, req.Description)
	assert.Empty(t, req.Logo)
	assert.Empty(t, req.LoginPageURL)
	assert.Empty(t, req.ApprovedCallbackUrls)
}

func TestRefreshModelFromResponse(t *testing.T) {
	ctx := context.Background()

	model := &Model{
		ClientSecret: types.StringValue("original-secret"),
	}

	app := &descope.ThirdPartyApplication{
		ID:                   "app-456",
		Name:                 "Updated App",
		Description:          "Updated description",
		Logo:                 "https://example.com/new-logo.png",
		LoginPageURL:         "https://example.com/new-login",
		ClientID:             "client-789",
		ApprovedCallbackUrls: []string{"https://example.com/new-callback"},
	}

	RefreshModelFromResponse(ctx, model, app)

	assert.Equal(t, "app-456", model.ID.ValueString())
	assert.Equal(t, "Updated App", model.Name.ValueString())
	assert.Equal(t, "Updated description", model.Description.ValueString())
	assert.Equal(t, "https://example.com/new-logo.png", model.Logo.ValueString())
	assert.Equal(t, "https://example.com/new-login", model.LoginPageURL.ValueString())
	assert.Equal(t, "client-789", model.ClientID.ValueString())
	// client_secret must be preserved from previous state
	assert.Equal(t, "original-secret", model.ClientSecret.ValueString())
}

func TestRefreshModelFromResponseNilCallbackUrls(t *testing.T) {
	ctx := context.Background()
	model := &Model{}

	app := &descope.ThirdPartyApplication{
		ID:                   "app-000",
		Name:                 "No Callbacks",
		ApprovedCallbackUrls: nil,
	}

	RefreshModelFromResponse(ctx, model, app)

	assert.Equal(t, "app-000", model.ID.ValueString())
	assert.Equal(t, "No Callbacks", model.Name.ValueString())
	assert.False(t, model.ApprovedCallbackUrls.IsNull())
}
