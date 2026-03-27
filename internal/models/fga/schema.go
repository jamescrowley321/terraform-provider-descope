package fga

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var SchemaAttributes = map[string]schema.Attribute{
	"id":     stringattr.Identifier(),
	"schema": stringattr.Required(),
}

type SchemaModel struct {
	ID     stringattr.Type `tfsdk:"id"`
	Schema stringattr.Type `tfsdk:"schema"`
}
