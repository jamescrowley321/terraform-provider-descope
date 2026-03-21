package role

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Attributes = map[string]schema.Attribute{
	"id":               stringattr.Identifier(),
	"name":             stringattr.Required(),
	"description":      stringattr.Default(""),
	"permission_names": strsetattr.Default(),
	"tenant_id": schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
	},
	"default_role": boolattr.Default(false),
	"private":      boolattr.Default(false),
}

type Model struct {
	ID              stringattr.Type `tfsdk:"id"`
	Name            stringattr.Type `tfsdk:"name"`
	Description     stringattr.Type `tfsdk:"description"`
	PermissionNames strsetattr.Type `tfsdk:"permission_names"`
	TenantID        stringattr.Type `tfsdk:"tenant_id"`
	DefaultRole     boolattr.Type   `tfsdk:"default_role"`
	Private         boolattr.Type   `tfsdk:"private"`
}
