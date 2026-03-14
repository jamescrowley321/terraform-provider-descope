package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/passwordsettings"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const passwordSettingsID = "password_settings"

var (
	_ resource.Resource              = &passwordSettingsResource{}
	_ resource.ResourceWithConfigure = &passwordSettingsResource{}
)

func NewPasswordSettingsResource() resource.Resource {
	return &passwordSettingsResource{}
}

type passwordSettingsResource struct {
	management sdk.Management
}

func (r *passwordSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *passwordSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_settings"
}

func (r *passwordSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages password authentication settings for a Descope project.",
		Attributes:  passwordsettings.PasswordSettingsAttributes,
	}
}

func (r *passwordSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating password settings resource")

	var model passwordsettings.PasswordSettingsModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	settings := model.ToSDK()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Password().ConfigureSettings(ctx, "", settings)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error configuring password settings", err.Error())
		return
	}

	model.ID = types.StringValue(passwordSettingsID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Password settings resource created")
}

func (r *passwordSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading password settings resource")

	var model passwordsettings.PasswordSettingsModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	settings, err := infra.RetryOnRateLimit(ctx, func() (*descope.PasswordSettings, error) {
		return r.management.Password().GetSettings(ctx, "")
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading password settings", err.Error())
		return
	}

	model.SetFromSDK(settings)
	model.ID = types.StringValue(passwordSettingsID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Password settings resource read")
}

func (r *passwordSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating password settings resource")

	var model passwordsettings.PasswordSettingsModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	settings := model.ToSDK()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Password().ConfigureSettings(ctx, "", settings)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating password settings", err.Error())
		return
	}

	model.ID = types.StringValue(passwordSettingsID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Password settings resource updated")
}

func (r *passwordSettingsResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Removing password settings resource from state")
	// Password settings are a singleton that cannot be deleted.
	// Removing from state is sufficient; the settings remain on the project.
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "Password settings resource removed from state")
}
