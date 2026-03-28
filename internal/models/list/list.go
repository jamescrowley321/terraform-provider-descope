package list

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/strsetattr"
)

var Attributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(), // Must be unique per project; duplicate names will cause an API error on create.
	"description": stringattr.Default(""),
	// Validator intentionally excludes "json" type — JSON lists use map data
	// which is incompatible with the string set data model.
	"type": stringattr.Required(stringvalidator.OneOf("texts", "ips"), stringplanmodifier.RequiresReplace()),
	"data": strsetattr.Default(),
}

type Model struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Type        stringattr.Type `tfsdk:"type"`
	Data        strsetattr.Type `tfsdk:"data"`
}
