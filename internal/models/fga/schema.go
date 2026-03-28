package fga

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/attrs/stringattr"
)

var SchemaAttributes = map[string]schema.Attribute{
	"id":     stringattr.Identifier(),
	"schema": stringattr.Required(),
}

type SchemaModel struct {
	ID     stringattr.Type `tfsdk:"id"`
	Schema stringattr.Type `tfsdk:"schema"`
}
