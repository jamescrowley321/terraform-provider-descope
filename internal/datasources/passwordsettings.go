package datasources

import (
	"context"

	"github.com/descope/terraform-provider-descope/internal/helpers"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/settings"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type passwordSettingsDataSource struct{ client *infra.Client }

func NewPasswordSettingsDataSource() datasource.DataSource { return &passwordSettingsDataSource{} }
func (d *passwordSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*infra.Client)
}
func (d *passwordSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_settings"
}
func (d *passwordSettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := toComputedAttributes(settings.PasswordSettingsAttributes)
	attrs["project_id"] = dsschema.StringAttribute{Required: true}
	resp.Schema = dsschema.Schema{Description: "Reads password authentication settings for a Descope project.", Attributes: attrs}
}
func (d *passwordSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var pid types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("project_id"), &pid)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data, err := d.client.Get(ctx, pid.ValueString(), "/v2/mgmt/password/settings", nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading password settings", err.Error())
		return
	}
	model := settings.PasswordSettingsModel{ProjectID: pid, ID: pid}
	model.SetValues(helpers.NewHandler(helpers.MarkImportContext(ctx), &resp.Diagnostics), data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
