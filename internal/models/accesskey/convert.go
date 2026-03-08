package accesskey

import (
	"context"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringSetToSlice converts a Terraform string set to a Go string slice.
func StringSetToSlice(_ context.Context, s strsetattr.Type, _ *diag.Diagnostics) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elems := s.Elements()
	result := make([]string, 0, len(elems))
	for _, v := range elems {
		if str, ok := v.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

// StringListToSlice converts a Terraform string list to a Go string slice.
func StringListToSlice(_ context.Context, l strlistattr.Type, _ *diag.Diagnostics) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	elems := l.Elements()
	result := make([]string, 0, len(elems))
	for _, v := range elems {
		if str, ok := v.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}

// TenantsToSDK converts the Terraform key_tenants list to SDK AssociatedTenant objects.
func TenantsToSDK(ctx context.Context, tenants listattr.Type[TenantModel], diagnostics *diag.Diagnostics) []*descope.AssociatedTenant {
	if tenants.IsNull() || tenants.IsUnknown() {
		return nil
	}
	elems, diags := tenants.ToSlice(ctx)
	diagnostics.Append(diags...)
	if diags.HasError() {
		return nil
	}
	result := make([]*descope.AssociatedTenant, 0, len(elems))
	for _, t := range elems {
		tenantID := t.TenantID.ValueString()
		roles := StringSetToSlice(ctx, t.Roles, diagnostics)
		result = append(result, &descope.AssociatedTenant{
			TenantID: tenantID,
			Roles:    roles,
		})
	}
	return result
}

// SetModelFromResponse populates the Terraform model from an SDK AccessKeyResponse.
//
//nolint:contextcheck // Value helpers use context.Background() by design
func SetModelFromResponse(model *AccessKeyModel, key *descope.AccessKeyResponse, cleartext string) {
	model.ID = types.StringValue(key.ID)
	model.Name = types.StringValue(key.Name)
	model.Description = types.StringValue(key.Description)
	model.ClientID = types.StringValue(key.ClientID)
	model.CreatedBy = types.StringValue(key.CreatedBy)
	model.CreatedTime = types.Int64Value(int64(key.CreatedTime))
	model.ExpireTime = types.Int64Value(int64(key.ExpireTime))
	model.UserID = types.StringValue(key.UserID)

	if key.Status == "" {
		model.Status = types.StringValue("active")
	} else {
		model.Status = types.StringValue(key.Status)
	}

	if cleartext != "" {
		model.Cleartext = types.StringValue(cleartext)
	}

	// Set role_names
	model.RoleNames = strsetattr.Value(key.RoleNames)

	// Set permitted_ips
	model.PermittedIPs = strlistattr.Value(key.PermittedIPs)

	// Set key_tenants
	if key.KeyTenants != nil {
		tenants := make([]*TenantModel, 0, len(key.KeyTenants))
		for _, t := range key.KeyTenants {
			tenants = append(tenants, &TenantModel{
				TenantID: types.StringValue(t.TenantID),
				Roles:    strsetattr.Value(t.Roles),
			})
		}
		model.KeyTenants = listattr.Value(tenants)
	} else {
		model.KeyTenants = listattr.Empty[TenantModel]()
	}
}
