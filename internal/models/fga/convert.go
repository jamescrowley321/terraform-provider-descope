package fga

import (
	"github.com/descope/go-sdk/descope"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func SchemaModelToSDK(model *SchemaModel) *descope.FGASchema {
	return &descope.FGASchema{
		Schema: model.Schema.ValueString(),
	}
}

func SchemaModelFromSDK(model *SchemaModel, fgaSchema *descope.FGASchema) {
	model.Schema = types.StringValue(fgaSchema.Schema)
}
