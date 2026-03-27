package thirdpartyapp

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ModelToRequest(ctx context.Context, model *Model, diags *diag.Diagnostics) *descope.ThirdPartyApplicationRequest {
	return &descope.ThirdPartyApplicationRequest{
		ID:                   model.ID.ValueString(),
		Name:                 model.Name.ValueString(),
		Description:          model.Description.ValueString(),
		Logo:                 model.Logo.ValueString(),
		LoginPageURL:         model.LoginPageURL.ValueString(),
		ApprovedCallbackUrls: convert.StringSetToSlice(ctx, model.ApprovedCallbackUrls, diags),
	}
}

func RefreshModelFromResponse(ctx context.Context, model *Model, app *descope.ThirdPartyApplication) {
	model.ID = types.StringValue(app.ID)
	model.Name = types.StringValue(app.Name)
	model.Description = types.StringValue(app.Description)
	model.Logo = types.StringValue(app.Logo)
	model.LoginPageURL = types.StringValue(app.LoginPageURL)
	model.ClientID = types.StringValue(app.ClientID)
	// client_secret is not returned by LoadApplication — preserve existing state value
	model.ApprovedCallbackUrls = strsetattr.ValueCtx(ctx, app.ApprovedCallbackUrls)
}
