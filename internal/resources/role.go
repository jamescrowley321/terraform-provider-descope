package resources

import (
	"context"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/descope/terraform-provider-descope/internal/models/convert"
	"github.com/descope/terraform-provider-descope/internal/models/role"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &roleResource{}
	_ resource.ResourceWithConfigure   = &roleResource{}
	_ resource.ResourceWithImportState = &roleResource{}
)

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

type roleResource struct {
	management sdk.Management
}

func (r *roleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if data, ok := req.ProviderData.(*infra.ProviderData); ok {
		r.management = data.Management
	}
}

func (r *roleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a standalone Descope role.",
		Attributes:  role.Attributes,
	}
}

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating role resource")

	var model role.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name := model.Name.ValueString()
	desc := model.Description.ValueString()
	tenantID := model.TenantID.ValueString()
	defaultRole := model.DefaultRole.ValueBool()
	private := model.Private.ValueBool()
	permNames := convert.StringSetToSlice(ctx, model.PermissionNames, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Role().Create(ctx, name, desc, permNames, tenantID, defaultRole, private)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating role", err.Error())
		return
	}

	model.ID = types.StringValue(roleID(name, tenantID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Role resource created")
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading role resource")

	var model role.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name, tenantID := parseRoleID(model.ID.ValueString())
	found, err := r.findRole(ctx, name, tenantID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading role", err.Error())
		return
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	refreshRoleModel(ctx, &model, found)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
	tflog.Info(ctx, "Role resource read")
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Updating role resource")

	var plan role.Model
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...); resp.Diagnostics.HasError() {
		return
	}

	var state role.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &state)...); resp.Diagnostics.HasError() {
		return
	}

	oldName, tenantID := parseRoleID(state.ID.ValueString())
	newName := plan.Name.ValueString()
	desc := plan.Description.ValueString()
	defaultRole := plan.DefaultRole.ValueBool()
	private := plan.Private.ValueBool()
	permNames := convert.StringSetToSlice(ctx, plan.PermissionNames, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Role().Update(ctx, oldName, tenantID, newName, desc, permNames, defaultRole, private)
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating role", err.Error())
		return
	}

	plan.ID = types.StringValue(roleID(newName, tenantID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Info(ctx, "Role resource updated")
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Deleting role resource")

	var model role.Model
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	name, tenantID := parseRoleID(model.ID.ValueString())
	err := infra.RetryOnRateLimitNoResult(ctx, func() error {
		return r.management.Role().Delete(ctx, name, tenantID)
	})
	if err != nil {
		if infra.IsNotFoundError(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting role", err.Error())
		return
	}

	tflog.Info(ctx, "Role resource deleted")
}

func (r *roleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, "Importing role resource")
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *roleResource) findRole(ctx context.Context, name, tenantID string) (*descope.Role, error) {
	roles, err := infra.RetryOnRateLimit(ctx, func() ([]*descope.Role, error) {
		return r.management.Role().LoadAll(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, rl := range roles {
		if rl.Name == name && rl.TenantID == tenantID {
			return rl, nil
		}
	}
	return nil, nil
}

// roleID creates the composite ID stored in Terraform state.
// For global roles: "name", for tenant-scoped roles: "tenantID/name".
func roleID(name, tenantID string) string {
	if tenantID == "" {
		return name
	}
	return tenantID + "/" + name
}

// parseRoleID splits a composite role ID back into name and tenantID.
// Format: "name" for global roles, "tenantID/name" for tenant-scoped roles.
// Uses the first "/" as separator since tenant IDs are UUIDs (no slashes),
// while role names may contain slashes.
func parseRoleID(id string) (name, tenantID string) {
	if i := strings.Index(id, "/"); i >= 0 {
		return id[i+1:], id[:i]
	}
	return id, ""
}

func refreshRoleModel(ctx context.Context, model *role.Model, r *descope.Role) {
	model.Name = types.StringValue(r.Name)
	model.Description = types.StringValue(r.Description)
	if r.TenantID != "" {
		model.TenantID = types.StringValue(r.TenantID)
	} else {
		model.TenantID = types.StringNull()
	}
	model.DefaultRole = types.BoolValue(r.Default)
	model.Private = types.BoolValue(r.Private)
	model.ID = types.StringValue(roleID(r.Name, r.TenantID))

	model.PermissionNames = strsetattr.ValueCtx(ctx, r.PermissionNames)
}
