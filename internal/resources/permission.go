package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/permission"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &permissionResource{}
	_ resource.ResourceWithConfigure   = &permissionResource{}
	_ resource.ResourceWithImportState = &permissionResource{}
)

func NewPermissionResource() resource.Resource {
	return &permissionResource{}
}

type permissionResource struct {
	management sdk.Management
}

func (r *permissionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *permissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *permissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a standalone Descope permission.",
		Attributes:  permission.Attributes,
	}
}

func (r *permissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating permission resource")

	var model permission.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name := model.Name.ValueString()
	desc := model.Description.ValueString()

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Permission().Create(ctx, name, desc)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating permission", err.Error())
		return
	}

	model.ID = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Permission resource created")
}

func (r *permissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading permission resource")

	var model permission.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name := model.ID.ValueString()
	p, err := r.findPermission(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error reading permission", err.Error())
		return
	}
	if p == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	model.Name = types.StringValue(p.Name)
	model.Description = types.StringValue(p.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Permission resource read")
}

func (r *permissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating permission resource")

	var plan permission.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state permission.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	oldName := state.ID.ValueString()
	newName := plan.Name.ValueString()
	desc := plan.Description.ValueString()

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Permission().Update(ctx, oldName, newName, desc)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating permission", err.Error())
		return
	}

	plan.ID = types.StringValue(newName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Permission resource updated")
}

func (r *permissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting permission resource")

	var model permission.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name := model.ID.ValueString()
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Permission().Delete(ctx, name)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting permission", err.Error())
		return
	}

	tflog.Info(ctx, "Permission resource deleted")
}

func (r *permissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing permission resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *permissionResource) findPermission(ctx context.Context, name string) (*descope.Permission, error) {
	perms, err := infra.RetryOnRateLimit(ctx, func() ([]*descope.Permission, error) {
		return r.management.Permission().LoadAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, p := range perms {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, nil
}
