package accesskey

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var TenantAttributes = map[string]schema.Attribute{
	"tenant_id":   stringattr.Required(),
	"tenant_name": stringattr.Identifier(),
	"roles":       strsetattr.Default(),
}

type TenantModel struct {
	TenantID   stringattr.Type `tfsdk:"tenant_id"`
	TenantName stringattr.Type `tfsdk:"tenant_name"`
	Roles      strsetattr.Type `tfsdk:"roles"`
}
