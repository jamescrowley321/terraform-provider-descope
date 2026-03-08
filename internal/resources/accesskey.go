package resources

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/accesskey"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &accessKeyResource{}
	_ resource.ResourceWithConfigure   = &accessKeyResource{}
	_ resource.ResourceWithImportState = &accessKeyResource{}
)

func NewAccessKeyResource() resource.Resource {
	return &accessKeyResource{}
}

type accessKeyResource struct {
	management sdk.Management
}

func (r *accessKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *accessKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_key"
}

func (r *accessKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: accesskey.AccessKeyAttributes,
	}
}

func (r *accessKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating access key resource")

	var model accesskey.AccessKeyModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	if model.Status.ValueString() == "inactive" {
		resp.Diagnostics.AddError("Invalid access key configuration", "Cannot set status to inactive when creating a new access key")
		return
	}

	name := model.Name.ValueString()
	description := model.Description.ValueString()
	expireTime := model.ExpireTime.ValueInt64()
	roles := accesskey.StringSetToSlice(ctx, model.RoleNames, &resp.Diagnostics)
	tenants := accesskey.TenantsToSDK(ctx, model.KeyTenants, &resp.Diagnostics)
	userID := model.UserID.ValueString()
	permittedIPs := accesskey.StringListToSlice(ctx, model.PermittedIPs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	cleartext, key, err := r.management.AccessKey().Create(ctx, name, description, expireTime, roles, tenants, userID, nil, permittedIPs, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error creating access key", err.Error())
		return
	}

	setModelFromResponse(&model, key, cleartext)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Access key resource created")
}

func (r *accessKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading access key resource")

	var model accesskey.AccessKeyModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	id := model.ID.ValueString()
	key, err := r.management.AccessKey().Load(ctx, id)
	if err != nil {
		if de := descope.AsError(err); de != nil && de.IsNotFound() {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading access key", err.Error())
		return
	}

	setModelFromResponse(&model, key, "")
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)

	tflog.Info(ctx, "Access key resource read")
}

func (r *accessKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating access key resource")

	var plan accesskey.AccessKeyModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state accesskey.AccessKeyModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()
	name := plan.Name.ValueString()
	description := plan.Description.ValueString()
	roles := accesskey.StringSetToSlice(ctx, plan.RoleNames, &resp.Diagnostics)
	tenants := accesskey.TenantsToSDK(ctx, plan.KeyTenants, &resp.Diagnostics)
	permittedIPs := accesskey.StringListToSlice(ctx, plan.PermittedIPs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.management.AccessKey().Update(ctx, id, name, &description, roles, tenants, nil, permittedIPs, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error updating access key", err.Error())
		return
	}

	// Handle status change (activate/deactivate)
	desiredStatus := plan.Status.ValueString()
	if desiredStatus != state.Status.ValueString() {
		if err := r.setAccessKeyStatus(ctx, id, desiredStatus); err != nil {
			resp.Diagnostics.AddError("Error changing access key status", err.Error())
			return
		}
	}

	setModelFromResponse(&plan, key, "")
	plan.Cleartext = state.Cleartext
	plan.Status = types.StringValue(desiredStatus)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Info(ctx, "Access key resource updated")
}

func (r *accessKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting access key resource")

	var model accesskey.AccessKeyModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	if err := r.management.AccessKey().Delete(ctx, model.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting access key", err.Error())
		return
	}

	tflog.Info(ctx, "Access key resource deleted")
}

func (r *accessKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing access key resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *accessKeyResource) setAccessKeyStatus(ctx context.Context, id, status string) error {
	if status == "inactive" {
		return r.management.AccessKey().Deactivate(ctx, id)
	}
	return r.management.AccessKey().Activate(ctx, id)
}

// setModelFromResponse populates the model from an API response, preserving
// planned values for fields the API may not echo back.
func setModelFromResponse(model *accesskey.AccessKeyModel, key *descope.AccessKeyResponse, cleartext string) {
	plannedRoles := model.RoleNames
	plannedTenants := model.KeyTenants
	plannedIPs := model.PermittedIPs

	accesskey.SetModelFromResponse(model, key, cleartext)

	if len(key.RoleNames) == 0 {
		model.RoleNames = plannedRoles
	}
	if key.KeyTenants == nil {
		model.KeyTenants = plannedTenants
	}
	if len(key.PermittedIPs) == 0 {
		model.PermittedIPs = plannedIPs
	}
}
