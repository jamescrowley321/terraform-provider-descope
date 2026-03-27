package datasources

import (
	"context"
	"encoding/json"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &projectExportDataSource{}
	_ datasource.DataSourceWithConfigure = &projectExportDataSource{}
)

func NewProjectExportDataSource() datasource.DataSource {
	return &projectExportDataSource{}
}

type projectExportDataSource struct {
	management sdk.Management
}

type projectExportModel struct {
	ID    types.String `tfsdk:"id"`
	Files types.String `tfsdk:"files"`
}

func (d *projectExportDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		d.management = data.Management
	}
}

func (d *projectExportDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_export"
}

func (d *projectExportDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Exports a snapshot of the current Descope project configuration. Returns the full project settings and configurations as a JSON string.",
		Attributes: map[string]dsschema.Attribute{
			"id": dsschema.StringAttribute{
				Computed:    true,
				Description: "Static identifier for the project export data source.",
			},
			"files": dsschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "JSON-encoded map of the exported project configuration files.",
			},
		},
	}
}

func (d *projectExportDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading project export data source")

	if d.management == nil {
		resp.Diagnostics.AddError("Provider not configured", "The provider has not been configured. Ensure the provider block is present and valid.")
		return
	}

	snapshot, err := infra.RetryOnRateLimit(ctx, func() (*descope.ExportSnapshotResponse, error) {
		return d.management.Project().ExportSnapshot(ctx, &descope.ExportSnapshotRequest{})
	})
	if err != nil {
		resp.Diagnostics.AddError("Error exporting project snapshot", err.Error())
		return
	}

	if snapshot == nil || snapshot.Files == nil {
		resp.Diagnostics.AddError("Error exporting project snapshot", "empty snapshot returned")
		return
	}

	filesJSON, err := json.Marshal(snapshot.Files)
	if err != nil {
		resp.Diagnostics.AddError("Error encoding project snapshot", err.Error())
		return
	}

	model := projectExportModel{
		ID:    types.StringValue("project_export"),
		Files: types.StringValue(string(filesJSON)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Project export data source read")
}
