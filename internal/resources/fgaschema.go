package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jamescrowley321/terraform-provider-descope/internal/infra"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/fga"
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
		Description: "Manages the Fine-Grained Authorization (FGA) schema for a Descope project. The schema defines object types and their relations for relationship-based access control. Note: destroying this resource only removes it from Terraform state; the schema remains active on the project.",
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
		if infra.IsNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading FGA schema", err.Error())
		return
	}

	if fgaSchema == nil {
		resp.State.RemoveResource(ctx)
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
	tflog.Info(ctx, "Removing FGA schema resource from state")
	// FGA schema is a singleton that cannot be deleted via the API.
	// Removing from state is sufficient; the schema remains on the project.
	resp.State.RemoveResource(ctx)
	tflog.Info(ctx, "FGA schema resource removed from state")
}
