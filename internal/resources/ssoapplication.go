package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/ssoapplication"
)

var (
	_ resource.Resource                   = &ssoApplicationResource{}
	_ resource.ResourceWithConfigure      = &ssoApplicationResource{}
	_ resource.ResourceWithImportState    = &ssoApplicationResource{}
	_ resource.ResourceWithValidateConfig = &ssoApplicationResource{}
)

func NewSSOApplicationResource() resource.Resource {
	return &ssoApplicationResource{}
}

type ssoApplicationResource struct {
	management sdk.Management
}

func (r *ssoApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *ssoApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_application"
}

func (r *ssoApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Descope SSO application (OIDC or SAML).",
		Attributes:  ssoapplication.Attributes,
	}
}

func (r *ssoApplicationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model ssoapplication.Model
	if resp.Diagnostics.Append(req.Config.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	hasOIDC := model.OIDC != nil
	hasSAML := model.SAML != nil

	if hasOIDC && hasSAML {
		resp.Diagnostics.AddError(
			"Invalid SSO application configuration",
			"Only one of 'oidc' or 'saml' block can be specified, not both.",
		)
	}

	if !hasOIDC && !hasSAML {
		resp.Diagnostics.AddError(
			"Invalid SSO application configuration",
			"Either 'oidc' or 'saml' block must be specified.",
		)
	}
}

func (r *ssoApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating SSO application resource")

	var model ssoapplication.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	var id string
	var err error

	if model.OIDC != nil {
		id, err = infra.RetryOnRateLimit(ctx, func() (string, error) {
			return r.management.SSOApplication().CreateOIDCApplication(ctx, ssoapplication.ModelToOIDCRequest(&model))
		})
	} else if model.SAML != nil {
		id, err = infra.RetryOnRateLimit(ctx, func() (string, error) {
			return r.management.SSOApplication().CreateSAMLApplication(ctx, ssoapplication.ModelToSAMLRequest(&model))
		})
	} else {
		resp.Diagnostics.AddError("Invalid SSO application", "Either oidc or saml block must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating SSO application", err.Error())
		return
	}

	model.ID = types.StringValue(id)

	// Read back to populate computed fields
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.SSOApplication, error) {
		return r.management.SSOApplication().Load(ctx, id)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO application after create", err.Error())
		return
	}

	ssoapplication.RefreshModelFromResponse(&model, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "SSO application resource created")
}

func (r *ssoApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading SSO application resource")

	var model ssoapplication.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.SSOApplication, error) {
		return r.management.SSOApplication().Load(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SSO application", err.Error())
		return
	}

	// Use import-aware refresh if type blocks are nil (import path)
	if model.OIDC == nil && model.SAML == nil {
		ssoapplication.RefreshModelFromResponseForImport(&model, app)
	} else {
		ssoapplication.RefreshModelFromResponse(&model, app)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "SSO application resource read")
}

func (r *ssoApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating SSO application resource")

	var plan ssoapplication.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state ssoapplication.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	// Preserve ID from state
	plan.ID = state.ID

	var err error
	if plan.OIDC != nil {
		err = infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.SSOApplication().UpdateOIDCApplication(ctx, ssoapplication.ModelToOIDCRequest(&plan))
		})
	} else if plan.SAML != nil {
		err = infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.SSOApplication().UpdateSAMLApplication(ctx, ssoapplication.ModelToSAMLRequest(&plan))
		})
	} else {
		resp.Diagnostics.AddError("Invalid SSO application", "Either oidc or saml block must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error updating SSO application", err.Error())
		return
	}

	// Read back
	id := plan.ID.ValueString()
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.SSOApplication, error) {
		return r.management.SSOApplication().Load(ctx, id)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSO application after update", err.Error())
		return
	}

	ssoapplication.RefreshModelFromResponse(&plan, app)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "SSO application resource updated")
}

func (r *ssoApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting SSO application resource")

	var model ssoapplication.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.SSOApplication().Delete(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting SSO application", err.Error())
		return
	}

	tflog.Info(ctx, "SSO application resource deleted")
}

func (r *ssoApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing SSO application resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
