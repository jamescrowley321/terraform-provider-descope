package datasources

import (
	"context"
	"fmt"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/fga"
)

var (
	_ datasource.DataSource              = &fgaCheckDataSource{}
	_ datasource.DataSourceWithConfigure = &fgaCheckDataSource{}
)

func NewFGACheckDataSource() datasource.DataSource {
	return &fgaCheckDataSource{}
}

type fgaCheckDataSource struct {
	management sdk.Management
}

func (d *fgaCheckDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		d.management = data.Management
	}
}

func (d *fgaCheckDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fga_check"
}

func (d *fgaCheckDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Checks whether a given FGA relation is authorized. Returns whether the target has the specified relation to the resource.",
		Attributes: map[string]dsschema.Attribute{
			"id":            dsschema.StringAttribute{Computed: true},
			"resource":      dsschema.StringAttribute{Required: true},
			"resource_type": dsschema.StringAttribute{Required: true},
			"relation":      dsschema.StringAttribute{Required: true},
			"target":        dsschema.StringAttribute{Required: true},
			"target_type":   dsschema.StringAttribute{Required: true},
			"allowed":       dsschema.BoolAttribute{Computed: true},
		},
	}
}

func (d *fgaCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading FGA check data source")

	var model fga.CheckModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	relation := &descope.FGARelation{
		Resource:     model.Resource.ValueString(),
		ResourceType: model.ResourceType.ValueString(),
		Relation:     model.Relation.ValueString(),
		Target:       model.Target.ValueString(),
		TargetType:   model.TargetType.ValueString(),
	}

	checks, err := infra.RetryOnRateLimit(ctx, func() ([]*descope.FGACheck, error) {
		return d.management.FGA().Check(ctx, []*descope.FGARelation{relation})
	})
	if err != nil {
		resp.Diagnostics.AddError("Error checking FGA relation", err.Error())
		return
	}

	if len(checks) == 0 || checks[0] == nil {
		resp.Diagnostics.AddError("Error checking FGA relation", "no check result returned")
		return
	}

	model.Allowed = types.BoolValue(checks[0].Allowed)
	model.ID = types.StringValue(fmt.Sprintf("%s:%s#%s@%s:%s",
		model.ResourceType.ValueString(), model.Resource.ValueString(),
		model.Relation.ValueString(),
		model.TargetType.ValueString(), model.Target.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "FGA check data source read")
}
