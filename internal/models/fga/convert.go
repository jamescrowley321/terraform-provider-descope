package fga

import (
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SchemaModelToSDK(model *SchemaModel) *descope.FGASchema {
	return &descope.FGASchema{
		Schema: model.Schema.ValueString(),
	}
}

func SchemaModelFromSDK(model *SchemaModel, fgaSchema *descope.FGASchema) {
	// Normalize trailing whitespace to prevent false diffs when the API
	// returns a normalized version of the DSL string.
	model.Schema = types.StringValue(strings.TrimRight(fgaSchema.Schema, " \t\n\r"))
}
