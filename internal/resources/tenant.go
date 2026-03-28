package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/convert"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/tenant"
)

var (
	_ resource.Resource                = &tenantResource{}
	_ resource.ResourceWithConfigure   = &tenantResource{}
	_ resource.ResourceWithImportState = &tenantResource{}
)

func NewTenantResource() resource.Resource {
	return &tenantResource{}
}

type tenantResource struct {
	management sdk.Management
}

func (r *tenantResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *tenantResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (r *tenantResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: tenant.TenantAttributes,
	}
}

func (r *tenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating tenant resource")

	var model tenant.TenantModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	tenantReq := tenant.ModelToRequest(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create tenant with optional custom ID
	var id string
	customID := model.TenantID.ValueString()
	if customID != "" {
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.Tenant().CreateWithID(ctx, customID, tenantReq)
		})
		if err != nil {
			resp.Diagnostics.AddError("Error creating tenant", err.Error())
			return
		}
		id = customID
	} else {
		var err error
		id, err = infra.RetryOnRateLimit(ctx, func() (string, error) {
			return r.management.Tenant().Create(ctx, tenantReq)
		})
		if err != nil {
			resp.Diagnostics.AddError("Error creating tenant", err.Error())
			return
		}
	}

	r.configureSettings(ctx, &model, id, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update default roles if set
	defaultRoles := convert.StringSetToSlice(ctx, model.DefaultRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(defaultRoles) > 0 {
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.Tenant().UpdateDefaultRoles(ctx, id, defaultRoles)
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating tenant default roles", err.Error())
			return
		}
	}

	// Load the created tenant
	t, err := infra.RetryOnRateLimit(ctx, func() (*descope.Tenant, error) {
		return r.management.Tenant().Load(ctx, id)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant after create", err.Error())
		return
	}

	savedSettings := tenant.RefreshModelFromAPI(&model, t)
	model.TenantID = types.StringValue(id)
	r.refreshSettings(ctx, &model, id, savedSettings, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Tenant resource created")
}

func (r *tenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading tenant resource")

	var model tenant.TenantModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	t, err := infra.RetryOnRateLimit(ctx, func() (*descope.Tenant, error) {
		return r.management.Tenant().Load(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading tenant", err.Error())
		return
	}

	savedSettings := tenant.RefreshModelFromAPI(&model, t)
	r.refreshSettings(ctx, &model, id, savedSettings, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Tenant resource read")
}

func (r *tenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating tenant resource")

	var plan tenant.TenantModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state tenant.TenantModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	tenantReq := tenant.ModelToRequest(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	// ParentTenantID is only for creation, clear it for updates
	tenantReq.ParentTenantID = ""

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Tenant().Update(ctx, id, tenantReq)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant", err.Error())
		return
	}

	r.configureSettings(ctx, &plan, id, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update default roles if changed
	planRoles := convert.StringSetToSlice(ctx, plan.DefaultRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	stateRoles := convert.StringSetToSlice(ctx, state.DefaultRoles, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if !rolesEqual(planRoles, stateRoles) {
		err := infra.RetryOnRateLimitNoResult(ctx, func() error {
			return r.management.Tenant().UpdateDefaultRoles(ctx, id, planRoles)
		})
		if err != nil {
			resp.Diagnostics.AddError("Error updating tenant default roles", err.Error())
			return
		}
	}

	// Re-load tenant
	t, err := infra.RetryOnRateLimit(ctx, func() (*descope.Tenant, error) {
		return r.management.Tenant().Load(ctx, id)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant after update", err.Error())
		return
	}

	savedSettings := tenant.RefreshModelFromAPI(&plan, t)
	plan.TenantID = state.TenantID
	r.refreshSettings(ctx, &plan, id, savedSettings, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Tenant resource updated")
}

func (r *tenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting tenant resource")

	var model tenant.TenantModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	cascade := model.CascadeDelete.ValueBool()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Tenant().Delete(ctx, model.ID.ValueString(), cascade)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return // already deleted
		}
		resp.Diagnostics.AddError("Error deleting tenant", err.Error())
		return
	}

	tflog.Info(ctx, "Tenant resource deleted")
}

func (r *tenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing tenant resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// configureSettings pushes the model's settings to the API if settings are present.
func (r *tenantResource) configureSettings(ctx context.Context, model *tenant.TenantModel, id string, diags *diag.Diagnostics) {
	if model.Settings == nil {
		return
	}
	settings := tenant.ModelToSettings(ctx, model, diags)
	if diags.HasError() {
		return
	}
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Tenant().ConfigureSettings(ctx, id, settings)
	})
	if err != nil {
		diags.AddError("Error configuring tenant settings", err.Error())
	}
}

// refreshSettings loads tenant settings from the API if the model previously had settings.
func (r *tenantResource) refreshSettings(ctx context.Context, model *tenant.TenantModel, id string, hadSettings *tenant.SettingsModel, diags *diag.Diagnostics) {
	if hadSettings == nil {
		return
	}
	settings, err := infra.RetryOnRateLimit(ctx, func() (*descope.TenantSettings, error) {
		return r.management.Tenant().GetSettings(ctx, id)
	})
	if err != nil {
		diags.AddError("Error reading tenant settings", err.Error())
		return
	}
	model.Settings = &tenant.SettingsModel{}
	tenant.SetSettingsFromSDK(model.Settings, settings)
}

// rolesEqual checks if two string slices contain the same elements regardless of order.
func rolesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	counts := make(map[string]int, len(a))
	for _, v := range a {
		counts[v]++
	}
	for _, v := range b {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	return true
}
