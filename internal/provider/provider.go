package provider

import (
	"context"
	"os"

	"github.com/descope/terraform-provider-descope/internal/datasources"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &descopeProvider{}
)

func NewDescopeProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &descopeProvider{
			version: version,
		}
	}
}

type descopeProvider struct {
	version string
}

type descopeProviderConfig struct {
	ProjectID     types.String `tfsdk:"project_id"`
	ManagementKey types.String `tfsdk:"management_key"`
	BaseURL       types.String `tfsdk:"base_url"`
}

func (p *descopeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "descope"
	resp.Version = p.version
}

func (p *descopeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use the Descope Terraform Provider to manage your Descope project's authentication methods, flows, roles, permissions, connectors, and more as infrastructure-as-code.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Optional:    true,
				Description: "The Descope project ID. Required for managing access keys without a descope_project resource. Can also be set via DESCOPE_PROJECT_ID environment variable.",
			},
			"management_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "A valid management key for your Descope company",
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "An optional base URL for the Descope API",
			},
		},
	}
}

func (p *descopeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Descope provider")

	var config descopeProviderConfig
	diags := req.Config.Get(ctx, &config)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if config.ManagementKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("management_key"), "Unknown Descope Management Key", "The provider cannot create the Descope client as there is an unknown configuration value for the Descope management key. Either target apply the source of the value first, set the value statically in the configuration, or use the DESCOPE_MANAGEMENT_KEY environment variable.")
	}
	if config.BaseURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("base_url"), "Unknown Descope Base URL", "The provider cannot create the Descope client as there is an unknown configuration value for the Descope base URL. Either target apply the source of the value first, set the value statically in the configuration, or use the DESCOPE_BASE_URL environment variable.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := os.Getenv("DESCOPE_PROJECT_ID")
	if !config.ProjectID.IsNull() {
		projectID = config.ProjectID.ValueString()
	}

	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	if !config.ManagementKey.IsNull() {
		managementKey = config.ManagementKey.ValueString()
	}

	baseURL := os.Getenv("DESCOPE_BASE_URL")
	if !config.BaseURL.IsNull() {
		baseURL = config.BaseURL.ValueString()
	}

	if managementKey == "" {
		resp.Diagnostics.AddAttributeError(path.Root("management_key"), "Missing Descope Management Key", "The provider cannot create the Descope client as there is a missing or empty value for the Descope management key. Set the management_key value in the configuration or use the DESCOPE_MANAGEMENT_KEY environment variable. If either is already set, ensure the value is not empty.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	providerData, err := infra.NewProviderData(p.version, managementKey, baseURL, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Descope client", err.Error())
		return
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData

	tflog.Info(ctx, "Configured Descope provider")
}

func (p *descopeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewProjectDataSource,
		datasources.NewPasswordSettingsDataSource,
	}
}

func (p *descopeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewProjectResource,
		resources.NewDescoperResource,
		resources.NewManagementKeyResource,
		resources.NewAccessKeyResource,
		resources.NewTenantResource,
		resources.NewInboundAppResource,
		resources.NewPasswordSettingsResource,
		resources.NewPermissionResource,
		resources.NewRoleResource,
		resources.NewSSOResource,
	}
}
