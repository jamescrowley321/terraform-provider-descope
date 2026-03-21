package resources

import (
	"context"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/convert"
	"github.com/descope/terraform-provider-descope/internal/models/sso"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &ssoResource{}
	_ resource.ResourceWithConfigure   = &ssoResource{}
	_ resource.ResourceWithImportState = &ssoResource{}
)

func NewSSOResource() resource.Resource {
	return &ssoResource{}
}

type ssoResource struct {
	management sdk.Management
}

func (r *ssoResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *ssoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso"
}

func (r *ssoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages SSO configuration for a Descope tenant. Supports OIDC and SAML.",
		Attributes:  sso.Attributes,
	}
}

func (r *ssoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating SSO resource")

	var model sso.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	tenantID := model.TenantID.ValueString()
	ssoID := model.SSOID.ValueString()
	displayName := model.DisplayName.ValueString()

	// Create the SSO configuration slot
	result, err := infra.RetryOnRateLimit(ctx, func() (*descope.SSOTenantSettingsResponse, error) {
		return r.management.SSO().NewSettings(ctx, tenantID, ssoID, displayName)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSO configuration", err.Error())
		return
	}

	ssoID = result.SSOID
	model.SSOID = types.StringValue(ssoID)
	model.ID = types.StringValue(ssoCompositeID(tenantID, ssoID))

	// Configure the SSO type
	r.configureSSOType(ctx, &model, tenantID, ssoID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read back to populate computed fields
	if found := r.refreshModel(ctx, &model, tenantID, ssoID, &resp.Diagnostics); !found && !resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Error reading SSO after create", "SSO configuration not found after creation")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "SSO resource created")
}

func (r *ssoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading SSO resource")

	var model sso.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	tenantID, ssoID := parseSSOCompositeID(model.ID.ValueString())

	found := r.refreshModel(ctx, &model, tenantID, ssoID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "SSO resource read")
}

func (r *ssoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating SSO resource")

	var plan sso.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	tenantID, ssoID := parseSSOCompositeID(plan.ID.ValueString())

	r.configureSSOType(ctx, &plan, tenantID, ssoID, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	if found := r.refreshModel(ctx, &plan, tenantID, ssoID, &resp.Diagnostics); !found && !resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Error reading SSO after update", "SSO configuration not found after update")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "SSO resource updated")
}

func (r *ssoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting SSO resource")

	var model sso.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	tenantID, ssoID := parseSSOCompositeID(model.ID.ValueString())

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.SSO().DeleteSettings(ctx, tenantID, ssoID)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting SSO configuration", err.Error())
		return
	}

	tflog.Info(ctx, "SSO resource deleted")
}

func (r *ssoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing SSO resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *ssoResource) configureSSOType(ctx context.Context, model *sso.Model, tenantID, ssoID string, diags *diag.Diagnostics) {
	domains := convert.StringSetToSlice(ctx, model.Domains, diags)
	if diags.HasError() {
		return
	}

	if model.OIDC != nil {
		settings := sso.ModelToOIDCSettings(ctx, model.OIDC, diags)
		if diags.HasError() {
			return
		}
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.SSO().ConfigureOIDCSettings(ctx, tenantID, settings, domains, ssoID)
		})
		if err != nil {
			diags.AddError("Error configuring OIDC SSO", err.Error())
		}
	} else if model.SAML != nil {
		settings, redirectURL := sso.ModelToSAMLSettings(model.SAML)
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.SSO().ConfigureSAMLSettings(ctx, tenantID, settings, redirectURL, domains, ssoID)
		})
		if err != nil {
			diags.AddError("Error configuring SAML SSO", err.Error())
		}
	} else if model.SAMLMetadata != nil {
		settings, redirectURL := sso.ModelToSAMLMetadataSettings(model.SAMLMetadata)
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.SSO().ConfigureSAMLSettingsByMetadata(ctx, tenantID, settings, redirectURL, domains, ssoID)
		})
		if err != nil {
			diags.AddError("Error configuring SAML SSO by metadata", err.Error())
		}
	}
}

// refreshModel loads SSO settings from the API and updates the model.
// Returns false if the SSO configuration was not found.
func (r *ssoResource) refreshModel(ctx context.Context, model *sso.Model, tenantID, ssoID string, diags *diag.Diagnostics) bool {
	result, err := infra.RetryOnRateLimit(ctx, func() (*descope.SSOTenantSettingsResponse, error) {
		return r.management.SSO().LoadSettings(ctx, tenantID, ssoID)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return false
		}
		diags.AddError("Error reading SSO configuration", err.Error())
		return false
	}

	model.TenantID = types.StringValue(tenantID)
	model.SSOID = types.StringValue(result.SSOID)
	model.ID = types.StringValue(ssoCompositeID(tenantID, result.SSOID))

	// Auto-initialize type blocks from API response (needed for import)
	if result.Oidc != nil {
		if model.OIDC == nil {
			model.OIDC = &sso.OIDCModel{}
		}
		sso.RefreshOIDCFromResponse(ctx, model.OIDC, result.Oidc)
	}
	if result.Saml != nil {
		if model.SAML == nil && model.SAMLMetadata == nil {
			model.SAML = &sso.SAMLModel{}
		}
		if model.SAML != nil {
			sso.RefreshSAMLFromResponse(model.SAML, result.Saml)
		}
		if model.SAMLMetadata != nil {
			sso.RefreshSAMLMetaFromResponse(model.SAMLMetadata, result.Saml)
		}
	}

	return true
}

func ssoCompositeID(tenantID, ssoID string) string {
	return tenantID + "/" + ssoID
}

func parseSSOCompositeID(id string) (tenantID, ssoID string) {
	if i := strings.Index(id, "/"); i >= 0 {
		return id[:i], id[i+1:]
	}
	return id, ""
}
