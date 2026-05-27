package datasources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/passwordsettings"
)

var (
	_ datasource.DataSource              = &passwordSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordSettingsDataSource{}
)

func NewPasswordSettingsDataSource() datasource.DataSource {
	return &passwordSettingsDataSource{}
}

type passwordSettingsDataSource struct {
	management sdk.Management
}

func (d *passwordSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		d.management = data.Management
	}
}

func (d *passwordSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_settings"
}

func (d *passwordSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Reads password authentication settings for a Descope project.",
		Attributes:  toComputedAttributes(passwordsettings.PasswordSettingsAttributes),
	}
}

func (d *passwordSettingsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading password settings data source")

	settings, err := infra.RetryOnRateLimit(ctx, func() (*descope.PasswordSettings, error) {
		return d.management.Password().GetSettings(ctx, "")
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading password settings", err.Error())
		return
	}

	var model passwordsettings.PasswordSettingsModel
	model.SetFromSDK(settings)
	model.ID = types.StringValue("password_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Password settings data source read")
}
