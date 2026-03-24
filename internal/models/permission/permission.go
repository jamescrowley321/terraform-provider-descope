package permission

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var Attributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
}

type Model struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
}
