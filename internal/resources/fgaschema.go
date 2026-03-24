package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/fga"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const fgaSchemaID = "fga_schema"

var (
	_ resource.Resource              = &fgaSchemaResource{}
	_ resource.ResourceWithConfigure = &fgaSchemaResource{}
)

func NewFGASchemaResource() resource.Resource {
	return &fgaSchemaResource{}
}

type fgaSchemaResource struct {
	management sdk.Management
}

func (r *fgaSchemaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *fgaSchemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fga_schema"
}

func (r *fgaSchemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the Fine-Grained Authorization (FGA) schema for a Descope project. The schema defines object types and their relations for relationship-based access control.",
		Attributes:  fga.SchemaAttributes,
	}
}

func (r *fgaSchemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating FGA schema resource")

	var model fga.SchemaModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	if r.saveSchema(ctx, &model, &resp.Diagnostics) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	}

	tflog.Info(ctx, "FGA schema resource created")
}

func (r *fgaSchemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading FGA schema resource")

	var model fga.SchemaModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	fgaSchema, err := infra.RetryOnRateLimit(ctx, func() (*descope.FGASchema, error) {
		return r.management.FGA().LoadSchema(ctx)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading FGA schema", err.Error())
		return
	}

	fga.SchemaModelFromSDK(&model, fgaSchema)
	model.ID = types.StringValue(fgaSchemaID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "FGA schema resource read")
}

func (r *fgaSchemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating FGA schema resource")

	var model fga.SchemaModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	if r.saveSchema(ctx, &model, &resp.Diagnostics) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	}

	tflog.Info(ctx, "FGA schema resource updated")
}

func (r *fgaSchemaResource) saveSchema(ctx context.Context, model *fga.SchemaModel, diags *diag.Diagnostics) bool {
	sdkSchema := fga.SchemaModelToSDK(model)
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.FGA().SaveSchema(ctx, sdkSchema)
	})
	if err != nil {
		diags.AddError("Error saving FGA schema", err.Error())
		return false
	}
	model.ID = types.StringValue(fgaSchemaID)
	return true
}

func (r *fgaSchemaResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting FGA schema resource")

	// Clear the schema by saving an empty one
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.FGA().SaveSchema(ctx, &descope.FGASchema{Schema: ""})
	})
	if err != nil {
		if !infra.IsNotFoundError(err) {
			resp.Diagnostics.AddError("Error clearing FGA schema", err.Error())
			return
		}
	}

	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "FGA schema resource deleted")
}
