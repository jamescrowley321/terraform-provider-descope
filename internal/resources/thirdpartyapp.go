package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/thirdpartyapp"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &thirdPartyAppResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyAppResource{}
	_ resource.ResourceWithImportState = &thirdPartyAppResource{}
)

func NewThirdPartyAppResource() resource.Resource {
	return &thirdPartyAppResource{}
}

type thirdPartyAppResource struct {
	management sdk.Management
}

func (r *thirdPartyAppResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *thirdPartyAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_application"
}

func (r *thirdPartyAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Descope third-party application that authenticates against Descope as an OAuth/OIDC provider.",
		Attributes:  thirdpartyapp.Attributes,
	}
}

func (r *thirdPartyAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating third-party application resource")

	var model thirdpartyapp.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	appReq := thirdpartyapp.ModelToRequest(ctx, &model, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	type createResult struct {
		id     string
		secret string
	}

	result, err := infra.RetryOnRateLimit(ctx, func() (*createResult, error) {
		id, secret, err := r.management.ThirdPartyApplication().CreateApplication(ctx, appReq)
		if err != nil {
			return nil, err
		}
		return &createResult{id: id, secret: secret}, nil
	})
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid third-party application configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating third-party application", err.Error())
		return
	}

	// Write partial state immediately so the resource is tracked even if the
	// subsequent Load call fails (prevents orphaned resources in Descope).
	model.ID = types.StringValue(result.id)
	model.ClientSecret = types.StringValue(result.secret)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Load the full application to populate computed fields
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.ThirdPartyApplication, error) {
		return r.management.ThirdPartyApplication().LoadApplication(ctx, result.id)
	})
	if err != nil {
		resp.Diagnostics.AddWarning("Error reading third-party application after creation",
			"The application was created successfully but could not be read back. Run 'terraform refresh' to sync state. Error: "+err.Error())
		return
	}

	thirdpartyapp.RefreshModelFromResponse(ctx, &model, app)
	model.ClientSecret = types.StringValue(result.secret)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Third-party application resource created")
}

func (r *thirdPartyAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading third-party application resource")

	var model thirdpartyapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.ThirdPartyApplication, error) {
		return r.management.ThirdPartyApplication().LoadApplication(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading third-party application", err.Error())
		return
	}

	clientSecret := model.ClientSecret
	thirdpartyapp.RefreshModelFromResponse(ctx, &model, app)
	model.ClientSecret = clientSecret
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Third-party application resource read")
}

func (r *thirdPartyAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating third-party application resource")

	var plan thirdpartyapp.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state thirdpartyapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	appReq := thirdpartyapp.ModelToRequest(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.ThirdPartyApplication().PatchApplication(ctx, appReq)
	})
	if failure, ok := infra.AsValidationError(err); ok {
		resp.Diagnostics.AddError("Invalid third-party application configuration", failure)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error updating third-party application", err.Error())
		return
	}

	// Load the updated application to refresh computed fields
	app, err := infra.RetryOnRateLimit(ctx, func() (*descope.ThirdPartyApplication, error) {
		return r.management.ThirdPartyApplication().LoadApplication(ctx, plan.ID.ValueString())
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading third-party application after update", err.Error())
		return
	}

	thirdpartyapp.RefreshModelFromResponse(ctx, &plan, app)
	plan.ClientSecret = state.ClientSecret
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Third-party application resource updated")
}

func (r *thirdPartyAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting third-party application resource")

	var model thirdpartyapp.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.ThirdPartyApplication().DeleteApplication(ctx, id)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting third-party application", err.Error())
		return
	}

	tflog.Info(ctx, "Third-party application resource deleted")
}

func (r *thirdPartyAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing third-party application resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
