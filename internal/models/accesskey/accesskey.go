package accesskey

import (
	"math"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/listattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strlistattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
)

var AccessKeyAttributes = map[string]schema.Attribute{
	"id":                stringattr.Identifier(),
	"name":              stringattr.Required(),
	"description":       stringattr.Default("", stringattr.StandardLenValidator),
	"status":            stringattr.Default("active", stringvalidator.OneOf("active", "inactive")),
	"expire_time":       intattr.Default(0, int64planmodifier.RequiresReplace(), int64validator.AtMost(math.MaxInt32)),
	"role_names":        strsetattr.Default(),
	"key_tenants":       listattr.Default[TenantModel](TenantAttributes),
	"permitted_ips":     strlistattr.Default(),
	"custom_claims":     strmapattr.Default(),
	"custom_attributes": strmapattr.Default(),
	"user_id": schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		Default:       stringdefault.StaticString(""),
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	},
	"cleartext": stringattr.SecretComputed(),
	"client_id": stringattr.Identifier(),
	"created_time": schema.Int64Attribute{
		Computed:      true,
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	},
	"created_by": stringattr.Identifier(),
}

type AccessKeyModel struct {
	ID               stringattr.Type            `tfsdk:"id"`
	Name             stringattr.Type            `tfsdk:"name"`
	Description      stringattr.Type            `tfsdk:"description"`
	Status           stringattr.Type            `tfsdk:"status"`
	ExpireTime       intattr.Type               `tfsdk:"expire_time"`
	RoleNames        strsetattr.Type            `tfsdk:"role_names"`
	KeyTenants       listattr.Type[TenantModel] `tfsdk:"key_tenants"`
	PermittedIPs     strlistattr.Type           `tfsdk:"permitted_ips"`
	CustomClaims     strmapattr.Type            `tfsdk:"custom_claims"`
	CustomAttributes strmapattr.Type            `tfsdk:"custom_attributes"`
	UserID           stringattr.Type            `tfsdk:"user_id"`
	Cleartext        stringattr.Type            `tfsdk:"cleartext"`
	ClientID         stringattr.Type            `tfsdk:"client_id"`
	CreatedTime      intattr.Type               `tfsdk:"created_time"`
	CreatedBy        stringattr.Type            `tfsdk:"created_by"`
}
