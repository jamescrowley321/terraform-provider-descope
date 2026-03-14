package datasources

import (
	"context"
	"fmt"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/project"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &projectDataSource{}
	_ datasource.DataSourceWithConfigure = &projectDataSource{}
)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *infra.Client
}

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		d.client = data.Client
	}
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := toComputedAttributes(project.ProjectAttributes)

	// Override id to be the required lookup key
	attrs["id"] = dsschema.StringAttribute{
		Required: true,
	}

	resp.Schema = dsschema.Schema{
		Attributes: attrs,
	}
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading project data source")

	var id types.String
	if resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...); resp.Diagnostics.HasError() {
		return
	}

	projectID := id.ValueString()

	res, err := d.client.Read(ctx, projectID, "project", projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	// Mark context so ShouldSetAttributeValue allows populating null fields
	ctx = helpers.ContextForDataSource(ctx)

	model := &project.ProjectModel{}
	model.ID = types.StringValue(projectID)

	handler := helpers.NewHandler(ctx, &resp.Diagnostics)
	model.CollectReferences(handler)
	model.SetValues(handler, res.Data)
	if resp.Diagnostics.HasError() {
		return
	}
	model.CollectReferences(handler)
	model.UpdateReferences(handler)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
	tflog.Info(ctx, "Project data source read")
}

// toComputedAttributes converts resource schema attributes to data source
// schema attributes with all fields set to Computed.
func toComputedAttributes(attrs map[string]rsschema.Attribute) map[string]dsschema.Attribute {
	result := make(map[string]dsschema.Attribute, len(attrs))
	for name, attr := range attrs {
		result[name] = toComputedAttribute(attr)
	}
	return result
}

func toComputedAttribute(attr rsschema.Attribute) dsschema.Attribute {
	switch a := attr.(type) {
	case rsschema.StringAttribute:
		return dsschema.StringAttribute{
			Computed:   true,
			Sensitive:  a.Sensitive,
			CustomType: a.CustomType,
		}
	case rsschema.Int64Attribute:
		return dsschema.Int64Attribute{
			Computed:   true,
			CustomType: a.CustomType,
		}
	case rsschema.BoolAttribute:
		return dsschema.BoolAttribute{
			Computed: true,
		}
	case rsschema.Float64Attribute:
		return dsschema.Float64Attribute{
			Computed:   true,
			CustomType: a.CustomType,
		}
	case rsschema.NumberAttribute:
		return dsschema.NumberAttribute{
			Computed:   true,
			CustomType: a.CustomType,
		}
	case rsschema.SetAttribute:
		return dsschema.SetAttribute{
			Computed:    true,
			ElementType: a.ElementType,
			CustomType:  a.CustomType,
		}
	case rsschema.ListAttribute:
		return dsschema.ListAttribute{
			Computed:    true,
			ElementType: a.ElementType,
			CustomType:  a.CustomType,
		}
	case rsschema.MapAttribute:
		return dsschema.MapAttribute{
			Computed:    true,
			ElementType: a.ElementType,
			CustomType:  a.CustomType,
		}
	case rsschema.SingleNestedAttribute:
		return dsschema.SingleNestedAttribute{
			Computed:   true,
			Attributes: toComputedAttributes(a.Attributes),
			CustomType: a.CustomType,
		}
	case rsschema.ListNestedAttribute:
		return dsschema.ListNestedAttribute{
			Computed: true,
			NestedObject: dsschema.NestedAttributeObject{
				Attributes: toComputedAttributes(a.NestedObject.Attributes),
				CustomType: a.NestedObject.CustomType,
			},
			CustomType: a.CustomType,
		}
	case rsschema.MapNestedAttribute:
		return dsschema.MapNestedAttribute{
			Computed: true,
			NestedObject: dsschema.NestedAttributeObject{
				Attributes: toComputedAttributes(a.NestedObject.Attributes),
				CustomType: a.NestedObject.CustomType,
			},
			CustomType: a.CustomType,
		}
	case rsschema.SetNestedAttribute:
		return dsschema.SetNestedAttribute{
			Computed: true,
			NestedObject: dsschema.NestedAttributeObject{
				Attributes: toComputedAttributes(a.NestedObject.Attributes),
				CustomType: a.NestedObject.CustomType,
			},
			CustomType: a.CustomType,
		}
	default:
		panic(fmt.Sprintf("unsupported attribute type for data source conversion: %T", attr))
	}
}
