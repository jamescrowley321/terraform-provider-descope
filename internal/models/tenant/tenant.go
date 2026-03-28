package tenant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/boolattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/intattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strmapattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/types/valuesettype"
)

var TenantAttributes = map[string]schema.Attribute{
	"id": stringattr.Identifier(),
	"tenant_id": schema.StringAttribute{
		Optional:      true,
		Computed:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace(), stringplanmodifier.UseStateForUnknown()},
	},
	"name":                      stringattr.Required(),
	"self_provisioning_domains": strsetattr.Default(),
	"custom_attributes":         strmapattr.Default(),
	"enforce_sso":               boolattr.Default(false),
	"disabled":                  boolattr.Default(false),
	"parent_tenant_id": schema.StringAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
	},
	"enforce_sso_exclusions": strsetattr.Default(),
	"role_inheritance":       stringattr.Default("", stringvalidator.OneOf("", "none", "userOnly")),
	"default_roles":          strsetattr.Default(),
	"cascade_delete":         boolattr.Default(false),
	"auth_type":              stringattr.Identifier(),
	"domains": schema.SetAttribute{
		Computed:    true,
		CustomType:  valuesettype.NewType[types.String](context.Background()),
		ElementType: types.StringType,
	},
	"created_time": schema.Int64Attribute{
		Computed:      true,
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	},
	"settings": schema.SingleNestedAttribute{
		Optional:   true,
		Attributes: SettingsAttributes,
	},
}

type TenantModel struct {
	ID                      stringattr.Type `tfsdk:"id"`
	TenantID                stringattr.Type `tfsdk:"tenant_id"`
	Name                    stringattr.Type `tfsdk:"name"`
	SelfProvisioningDomains strsetattr.Type `tfsdk:"self_provisioning_domains"`
	CustomAttributes        strmapattr.Type `tfsdk:"custom_attributes"`
	EnforceSSO              boolattr.Type   `tfsdk:"enforce_sso"`
	Disabled                boolattr.Type   `tfsdk:"disabled"`
	ParentTenantID          stringattr.Type `tfsdk:"parent_tenant_id"`
	EnforceSSOExclusions    strsetattr.Type `tfsdk:"enforce_sso_exclusions"`
	RoleInheritance         stringattr.Type `tfsdk:"role_inheritance"`
	DefaultRoles            strsetattr.Type `tfsdk:"default_roles"`
	CascadeDelete           boolattr.Type   `tfsdk:"cascade_delete"`
	AuthType                stringattr.Type `tfsdk:"auth_type"`
	Domains                 strsetattr.Type `tfsdk:"domains"`
	CreatedTime             intattr.Type    `tfsdk:"created_time"`
	Settings                *SettingsModel  `tfsdk:"settings"`
}
