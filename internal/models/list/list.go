package list

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var Attributes = map[string]schema.Attribute{
	"id":          stringattr.Identifier(),
	"name":        stringattr.Required(),
	"description": stringattr.Default(""),
	"type":        stringattr.Required(stringvalidator.OneOf("texts", "ips"), stringplanmodifier.RequiresReplace()),
	"data":        strsetattr.Default(),
}

type Model struct {
	ID          stringattr.Type `tfsdk:"id"`
	Name        stringattr.Type `tfsdk:"name"`
	Description stringattr.Type `tfsdk:"description"`
	Type        stringattr.Type `tfsdk:"type"`
	Data        strsetattr.Type `tfsdk:"data"`
}
